package lending

import (
	"fmt"
	"lending-engine/blockchain"
	"lending-engine/common"
	"lending-engine/internal/handler"
	"lending-engine/internal/redis"
	"lending-engine/response"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

type lendingHandler struct {
	QueryTransactionClientFn blockchain.QueryTransactionClientFn
	LendingRepository        LendingRepository
	GetFloatDataRedisFn      redis.GetFloatDataRedisFn
}

func NewLendingHandler(lendingRepository LendingRepository, queryTransactionClientFn blockchain.QueryTransactionClientFn, getFloatDataRedisFn redis.GetFloatDataRedisFn) *lendingHandler {
	return &lendingHandler{
		QueryTransactionClientFn: queryTransactionClientFn,
		LendingRepository:        lendingRepository,
		GetFloatDataRedisFn:      getFloatDataRedisFn,
	}
}

func (s *lendingHandler) GetTokenPrice(c *handler.Ctx) error {
	thbbtc, err := s.GetFloatDataRedisFn(common.THBBTCRedis)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalRedis, err.Error()))
	}
	thbeth, err := s.GetFloatDataRedisFn(common.THBETHRedis)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalRedis, err.Error()))
	}
	getTokenPriceResponse := GetTokenPriceResponse{
		BTC: TokenPrice{
			Price:   thbbtc,
			Haircut: viper.GetFloat64("loan.haircut.btc"),
		},
		ETH: TokenPrice{
			Price:   thbeth,
			Haircut: viper.GetFloat64("loan.haircut.eth"),
		},
		InterestRate: viper.GetFloat64("loan.interest"),
	}
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).GetTokenPriceSuccess, &getTokenPriceResponse))
}

func (s *lendingHandler) PreCalculationLoan(c *handler.Ctx) error {
	var req PreCalculationLoanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).PreCalculationLoanRequest, err.Error()))
	}
	if err := req.validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).PreCalculationLoanRequest, err.Error()))
	}

	thbbtc, err := s.GetFloatDataRedisFn(common.THBBTCRedis)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalRedis, err.Error()))
	}
	thbeth, err := s.GetFloatDataRedisFn(common.THBETHRedis)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalRedis, err.Error()))
	}

	btcLoan := req.BTCAmount * thbbtc * viper.GetFloat64("loan.haircut.btc")
	ethLoan := req.ETHAmount * thbeth * viper.GetFloat64("loan.haircut.eth")

	totalLoanAmount := btcLoan + ethLoan
	monthlyInterest := totalLoanAmount * viper.GetFloat64("loan.interest") / 12
	totalInterest := monthlyInterest * float64(req.Period)

	preCalculationLoanResponse := PreCalculationLoanResponse{
		BTC: TokenPriceRate{
			Volume:     req.BTCAmount,
			Haircut:    viper.GetFloat64("loan.haircut.btc"),
			LoanAmount: btcLoan,
		},
		ETH: TokenPriceRate{
			Volume:     req.ETHAmount,
			Haircut:    viper.GetFloat64("loan.haircut.btc"),
			LoanAmount: ethLoan,
		},
		Summary: SummaryLoan{
			TotalLoanAmount: totalLoanAmount,
			InterestRate:    viper.GetFloat64("loan.interest"),
			MonthlyInterest: monthlyInterest,
			Period:          req.Period,
			TotalInterest:   totalInterest,
		},
	}
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).PreCalculationLoanSuccess, &preCalculationLoanResponse))
}

func (s *lendingHandler) GetWalletTransaction(c *handler.Ctx) error {
	bearer := c.Locals(common.JWTClaimsKey).(*jwt.Token)
	claims := bearer.Claims.(jwt.MapClaims)
	id := claims["accountId"].(float64)

	deposits, err := s.LendingRepository.QueryWalletTransactionRepo(c.Context(), map[string]interface{}{"account_id": id})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).GetWalletTransactionSuccess, &deposits))
}

func (s *lendingHandler) SubmitDeposit(c *handler.Ctx) error {
	bearer := c.Locals(common.JWTClaimsKey).(*jwt.Token)
	claims := bearer.Claims.(jwt.MapClaims)
	id := claims["accountId"].(float64)
	accountId := int(id)

	var req SubmitDepositRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).SubmitDepositRequest, err.Error()))
	}
	if err := req.validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).SubmitDepositRequest, err.Error()))
	}

	status := common.PendingStatus

	if viper.GetBool("toggle.query-txn") {
		result, isPending, err := s.QueryTransactionClientFn(c.Context(), req.ChainID, req.TxnHash)
		if err != nil {
			c.Log().Error(err.Error())
		}
		if isPending {
			c.Log().Info(fmt.Sprintf("Txn Hash: %s | Txn Status: %t", req.TxnHash, isPending))
		}
		if result != nil {
			if result.TokenTransfer.From == req.Address && result.TokenTransfer.To == viper.GetString("blockchain.address") && result.TokenTransfer.Amount == req.Volume {
				status = common.ConfirmStatus
			}
			c.Log().Info(fmt.Sprintf("From: %s | Interacted With(To): %s | To: %s | Amount: %f", result.From, result.InteractedWith, result.TokenTransfer.To, result.TokenTransfer.Amount))
		}
	}

	c.Log().Info(fmt.Sprintf("Deposit Status: %s", status))

	depositId, err := s.LendingRepository.InsertDepositRepo(c.Context(), accountId, req.Address, req.ChainID, req.TxnHash, req.CollateralType, req.Volume, common.DepositStatus, status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}

	if status == common.ConfirmStatus {
		wallet, err := s.LendingRepository.QueryWalletRepo(c.Context(), accountId)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
		}
		if wallet == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, "Wallet doesn't exist."))
		}

		btc := *wallet.BTCVolume
		eth := *wallet.ETHVolume
		switch req.CollateralType {
		case "BTC":
			btc += req.Volume
		case "ETH":
			eth += req.Volume
		default:
			c.Log().Info(fmt.Sprintf("AccountID: %d can't update collateral volume (%s - %f).", accountId, req.CollateralType, req.Volume))
		}
		rows, err := s.LendingRepository.UpdateWalletRepo(c.Context(), accountId, btc, eth, wallet.MarginCallDate, time.Now().Format(common.DateYYYYMMDDHHMMSSFormat))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
		}
		if rows != 1 {
			return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, fmt.Sprintf("expected to affect 1 row, affected %d", rows)))
		}
		c.Log().Info(fmt.Sprintf("AccountID: %d | BTC: %f | ETH: %f", accountId, btc, eth))
	}
	submitDepositResponse := SubmitDepositResponse{
		DepositID: depositId,
	}
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).SubmitDepositSuccess, &submitDepositResponse))
}

func (s *lendingHandler) SubmitWithdraw(c *handler.Ctx) error {
	bearer := c.Locals(common.JWTClaimsKey).(*jwt.Token)
	claims := bearer.Claims.(jwt.MapClaims)
	id := claims["accountId"].(float64)
	accountId := int(id)

	var req SubmitWithdrawRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).SubmitWithdrawRequest, err.Error()))
	}
	if err := req.validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).SubmitWithdrawRequest, err.Error()))
	}

	withdrawId, err := s.LendingRepository.InsertWithdrawRepo(c.Context(), accountId, req.Address, req.ChainID, req.CollateralType, req.Volume, common.WithdrawStatus, common.PendingStatus)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	submitWithdrawResponse := SubmitWithdrawResponse{
		WithdrawID: withdrawId,
	}
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).SubmitWithdrawSuccess, &submitWithdrawResponse))
}

func (s *lendingHandler) GetWalletTransactionAdmin(c *handler.Ctx) error {
	var req GetWalletTransactionAdminRequest
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).GetWalletTransactionAdminRequest, err.Error()))
	}
	m := map[string]interface{}{
		"txn_type": req.TxnType,
	}
	if req.ID != nil {
		m["id"] = req.ID
	}
	if req.AccountID != nil {
		m["account_id"] = req.AccountID
	}
	if req.Address != nil {
		m["address"] = req.Address
	}
	lists, err := s.LendingRepository.QueryWalletTransactionRepo(c.Context(), m)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).GetWalletTransactionAdminSuccess, &lists))
}

func (s *lendingHandler) ConfirmDepositAdmin(c *handler.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).ConfirmDepositAdminRequest, err.Error()))
	}

	txn, err := s.LendingRepository.QueryWalletTransactionByIDRepo(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if txn == nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).ConfirmDepositAdminRequest, "ID doesn't exist."))
	}
	if *txn.Status != common.PendingStatus {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).ConfirmDepositAdminRequest, "This id has already confirmed or cancelled."))
	}
	if *txn.TxnType != common.DepositStatus {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).ConfirmDepositAdminRequest, "This id isn't deposit method."))
	}

	depositRows, err := s.LendingRepository.UpdateDepositRepo(c.Context(), id, common.ConfirmStatus, time.Now().Format(common.DateYYYYMMDDHHMMSSFormat))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if depositRows != 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, fmt.Sprintf("expected to affect 1 row, affected %d", depositRows)))
	}

	wallet, err := s.LendingRepository.QueryWalletRepo(c.Context(), *txn.AccountID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if wallet == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, "Wallet doesn't exist."))
	}

	btc := *wallet.BTCVolume
	eth := *wallet.ETHVolume
	switch *txn.CollateralType {
	case "BTC":
		btc += *txn.Volume
	case "ETH":
		eth += *txn.Volume
	default:
		c.Log().Info(fmt.Sprintf("AccountID: %d can't update collateral volume (%s - %f).", *txn.AccountID, *txn.CollateralType, *txn.Volume))
	}
	walletRows, err := s.LendingRepository.UpdateWalletRepo(c.Context(), *txn.AccountID, btc, eth, wallet.MarginCallDate, time.Now().Format(common.DateYYYYMMDDHHMMSSFormat))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if walletRows != 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, fmt.Sprintf("expected to affect 1 row, affected %d", walletRows)))
	}
	c.Log().Info(fmt.Sprintf("AccountID: %d | BTC: %f | ETH: %f", *txn.AccountID, btc, eth))

	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).ConfirmDepositAdminSuccess, nil))
}

func (s *lendingHandler) ConfirmWithdrawAdmin(c *handler.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).ConfirmWithdrawAdminRequest, err.Error()))
	}
	txnHash := c.Params("txnHash")
	if txnHash == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).ConfirmWithdrawAdminRequest, fmt.Sprintf("'txnHash' must be REQUIRED field but the input is '%v'.", txnHash)))
	}

	txn, err := s.LendingRepository.QueryWalletTransactionByIDRepo(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if txn == nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).ConfirmWithdrawAdminRequest, "ID doesn't exist."))
	}
	if *txn.Status != common.PendingStatus {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).ConfirmWithdrawAdminRequest, "This id has already confirmed or cancelled."))
	}
	if *txn.TxnType != common.WithdrawStatus {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).ConfirmWithdrawAdminRequest, "This id isn't withdraw method."))
	}

	withdrawRows, err := s.LendingRepository.UpdateWithdrawRepo(c.Context(), id, txnHash, common.ConfirmStatus, time.Now().Format(common.DateYYYYMMDDHHMMSSFormat))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if withdrawRows != 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, fmt.Sprintf("expected to affect 1 row, affected %d", withdrawRows)))
	}

	wallet, err := s.LendingRepository.QueryWalletRepo(c.Context(), *txn.AccountID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if wallet == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, "Wallet doesn't exist."))
	}

	btc := *wallet.BTCVolume
	eth := *wallet.ETHVolume
	switch *txn.CollateralType {
	case "BTC":
		btc -= *txn.Volume
		if btc < 0 {
			btc = 0
		}
	case "ETH":
		eth -= *txn.Volume
		if eth < 0 {
			eth = 0
		}
	default:
		c.Log().Info(fmt.Sprintf("AccountID: %d can't update collateral volume (%s - %f).", *txn.AccountID, *txn.CollateralType, *txn.Volume))
	}
	walletRows, err := s.LendingRepository.UpdateWalletRepo(c.Context(), *txn.AccountID, btc, eth, wallet.MarginCallDate, time.Now().Format(common.DateYYYYMMDDHHMMSSFormat))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if walletRows != 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, fmt.Sprintf("expected to affect 1 row, affected %d", walletRows)))
	}
	c.Log().Info(fmt.Sprintf("AccountID: %d | BTC: %f | ETH: %f", *txn.AccountID, btc, eth))

	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).ConfirmWithdrawAdminSuccess, nil))
}

func (s *lendingHandler) GetCreditAvailable(c *handler.Ctx) error {
	bearer := c.Locals(common.JWTClaimsKey).(*jwt.Token)
	claims := bearer.Claims.(jwt.MapClaims)
	id := claims["accountId"].(float64)
	accountId := int(id)

	wallet, err := s.LendingRepository.QueryWalletRepo(c.Context(), accountId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if wallet == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, "Wallet doesn't exist."))
	}

	thbbtc, err := s.GetFloatDataRedisFn(common.THBBTCRedis)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalRedis, err.Error()))
	}
	thbeth, err := s.GetFloatDataRedisFn(common.THBETHRedis)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalRedis, err.Error()))
	}

	btcLoan := *wallet.BTCVolume * thbbtc * viper.GetFloat64("loan.haircut.btc")
	ethLoan := *wallet.ETHVolume * thbeth * viper.GetFloat64("loan.haircut.eth")

	totalCollateralValue := btcLoan + ethLoan

	contracts, err := s.LendingRepository.QueryContractRepo(c.Context(), map[string]interface{}{"account_id": accountId})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	var totalOutstanding float64
	for _, value := range *contracts {
		if *value.Status != common.ClosedStatus {
			totalOutstanding += *value.LoanOutstanding
		}
	}

	getCreditAvailableResponse := GetCreditAvailableResponse{
		BTCVolume:       *wallet.BTCVolume,
		ETHVolume:       *wallet.ETHVolume,
		CollateralValue: totalCollateralValue,
		LoanOutstanding: totalOutstanding,
		CreditAvailable: totalCollateralValue - totalOutstanding,
	}
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).GetCreditAvailableSuccess, &getCreditAvailableResponse))
}

func (s *lendingHandler) GetLoan(c *handler.Ctx) error {
	bearer := c.Locals(common.JWTClaimsKey).(*jwt.Token)
	claims := bearer.Claims.(jwt.MapClaims)
	id := claims["accountId"].(float64)

	lists, err := s.LendingRepository.QueryContractRepo(c.Context(), map[string]interface{}{"account_id": id})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).GetLoanSuccess, &lists))
}

func (s *lendingHandler) BorrowLoan(c *handler.Ctx) error {
	bearer := c.Locals(common.JWTClaimsKey).(*jwt.Token)
	claims := bearer.Claims.(jwt.MapClaims)
	id := claims["accountId"].(float64)
	accountId := int(id)

	var req BorrowLoanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).BorrowLoanRequest, err.Error()))
	}
	if err := req.validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).BorrowLoanRequest, err.Error()))
	}

	contractId, err := s.LendingRepository.InsertContractRepo(c.Context(), accountId, req.InterestCode, req.Loan, req.Term)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	borrowLoanResponse := BorrowLoanResponse{
		ContractID: contractId,
	}
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).BorrowLoanSuccess, &borrowLoanResponse))
}

func (s *lendingHandler) GetLoanAdmin(c *handler.Ctx) error {
	var req GetLoanAdminRequest
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).GetContractAdminRequest, err.Error()))
	}
	m := make(map[string]interface{})
	if req.ContractID != nil {
		m["contract_id"] = req.ContractID
	}
	if req.AccountID != nil {
		m["account_id"] = req.AccountID
	}
	lists, err := s.LendingRepository.QueryContractRepo(c.Context(), m)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).GetContractAdminSuccess, &lists))
}

func (s *lendingHandler) ConfirmLoanAdmin(c *handler.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).ConfirmContractAdminRequest, err.Error()))
	}

	contract, err := s.LendingRepository.QueryContractByIDRepo(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if contract == nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).ConfirmContractAdminRequest, "ID doesn't exist."))
	}
	if *contract.Status != common.PendingStatus {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, "This id has already run or closed."))
	}

	contractRows, err := s.LendingRepository.UpdateContractRepo(c.Context(), id, common.OngoingStatus, time.Now().Format(common.DateYYYYMMDDHHMMSSFormat))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if contractRows != 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, fmt.Sprintf("expected to affect 1 row, affected %d", contractRows)))
	}
	c.Log().Info(fmt.Sprintf("ContractID: %d | Status: %s", id, *contract.Status))
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).ConfirmContractAdminSuccess, nil))
}

func (s *lendingHandler) GetInterestTerm(c *handler.Ctx) error {
	interestTerm, err := s.LendingRepository.QueryInterestTermRepo(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).GetInterestTermSuccess, &interestTerm))
}

func (s *lendingHandler) GetRepay(c *handler.Ctx) error {
	bearer := c.Locals(common.JWTClaimsKey).(*jwt.Token)
	claims := bearer.Claims.(jwt.MapClaims)
	id := claims["accountId"].(float64)

	lists, err := s.LendingRepository.QueryRepayTransactionRepo(c.Context(), map[string]interface{}{"account_id": id})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).GetRepaymentSuccess, &lists))
}

func (s *lendingHandler) SubmitRepay(c *handler.Ctx) error {
	bearer := c.Locals(common.JWTClaimsKey).(*jwt.Token)
	claims := bearer.Claims.(jwt.MapClaims)
	id := claims["accountId"].(float64)

	var req SubmitRepayRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).SubmitRepaymentRequest, err.Error()))
	}
	if err := req.validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).SubmitRepaymentRequest, err.Error()))
	}

	contract, err := s.LendingRepository.QueryContractByIDRepo(c.Context(), req.ContractID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if contract == nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).SubmitRepaymentRequest, "ContractID doesn't exist."))
	}

	repayId, err := s.LendingRepository.InsertRepayTransactionRepo(c.Context(), req.ContractID, int(id), float64(req.Amount), req.Slip)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	submitRepayResponse := SubmitRepayResponse{
		RepayID: repayId,
	}
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).SubmitRepaymentSuccess, &submitRepayResponse))
}

func (s *lendingHandler) GetRepayAdmin(c *handler.Ctx) error {
	var req GetRepayAdminRequest
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).GetRepaymentAdminRequest, err.Error()))
	}
	m := make(map[string]interface{})
	if req.ID != nil {
		m["id"] = req.ID
	}
	if req.ContractID != nil {
		m["contract_id"] = req.ContractID
	}
	if req.AccountID != nil {
		m["account_id"] = req.AccountID
	}
	lists, err := s.LendingRepository.QueryRepayTransactionRepo(c.Context(), m)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).GetRepaymentAdminSuccess, &lists))
}

func (s *lendingHandler) ConfirmRepayAdmin(c *handler.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).ConfirmRepaymentAdminRequest, err.Error()))
	}

	repay, err := s.LendingRepository.QueryRepayTransactionByIDRepo(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if repay == nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).ConfirmRepaymentAdminRequest, "ID doesn't exist."))
	}
	if *repay.Status != common.PendingStatus {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, "This id has already confirmed."))
	}

	repayRows, err := s.LendingRepository.UpdateRepayTransactionRepo(c.Context(), id, common.ConfirmStatus, time.Now().Format(common.DateYYYYMMDDHHMMSSFormat))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if repayRows != 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, fmt.Sprintf("expected to affect 1 row, affected %d", repayRows)))
	}

	contractRows, err := s.LendingRepository.UpdateContractRepo(c.Context(), *repay.ContractID, common.ClosedStatus, time.Now().Format(common.DateYYYYMMDDHHMMSSFormat))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if contractRows != 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, fmt.Sprintf("expected to affect 1 row, affected %d", contractRows)))
	}
	c.Log().Info(fmt.Sprintf("RepayID: %d - Status: %s | ContractID: %d Status: %s", id, common.ConfirmStatus, *repay.ContractID, common.ClosedStatus))
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).ConfirmRepaymentAdminSuccess, nil))
}
