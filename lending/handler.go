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

func (s *lendingHandler) GetDepositStatus(c *handler.Ctx) error {
	bearer := c.Locals(common.JWTClaimsKey).(*jwt.Token)
	claims := bearer.Claims.(jwt.MapClaims)
	id := claims["accountId"].(float64)

	deposits, err := s.LendingRepository.QueryDepositByAccountIDRepo(c.Context(), int(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).GetDepositSuccess, &deposits))
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

	if err := s.LendingRepository.InsertDepositRepo(c.Context(), accountId, req.Address, req.ChainID, req.TxnHash, req.CollateralType, req.Volume, status); err != nil {
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
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).SubmitDepositSuccess, nil))
}

func (s *lendingHandler) GetDepositAdmin(c *handler.Ctx) error {
	var req GetDepositAdminRequest
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).GetDepositAdminRequest, err.Error()))
	}
	m := make(map[string]interface{})
	if req.ID != nil {
		m["id"] = req.ID
	}
	if req.AccountID != nil {
		m["account_id"] = req.AccountID
	}
	if req.Address != nil {
		m["address"] = req.Address
	}
	lists, err := s.LendingRepository.QueryDepositRepo(c.Context(), m)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).GetDepositAdminSuccess, &lists))
}

func (s *lendingHandler) ConfirmDepositAdmin(c *handler.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).ConfirmDepositAdminRequest, err.Error()))
	}

	deposit, err := s.LendingRepository.QueryDepositByIDRepo(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if deposit == nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).ConfirmDepositAdminRequest, "ID doesn't exist."))
	}
	if *deposit.Status != common.PendingStatus {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, "This id has already confirmed or cancelled."))
	}

	depositRows, err := s.LendingRepository.UpdateDepositRepo(c.Context(), id, common.ConfirmStatus, time.Now().Format(common.DateYYYYMMDDHHMMSSFormat))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if depositRows != 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, fmt.Sprintf("expected to affect 1 row, affected %d", depositRows)))
	}

	wallet, err := s.LendingRepository.QueryWalletRepo(c.Context(), *deposit.AccountID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if wallet == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, "Wallet doesn't exist."))
	}

	btc := *wallet.BTCVolume
	eth := *wallet.ETHVolume
	switch *deposit.CollateralType {
	case "BTC":
		btc += *deposit.Volume
	case "ETH":
		eth += *deposit.Volume
	default:
		c.Log().Info(fmt.Sprintf("AccountID: %d can't update collateral volume (%s - %f).", *deposit.AccountID, *deposit.CollateralType, *deposit.Volume))
	}
	walletRows, err := s.LendingRepository.UpdateWalletRepo(c.Context(), *deposit.AccountID, btc, eth, wallet.MarginCallDate, time.Now().Format(common.DateYYYYMMDDHHMMSSFormat))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if walletRows != 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, fmt.Sprintf("expected to affect 1 row, affected %d", walletRows)))
	}
	c.Log().Info(fmt.Sprintf("AccountID: %d | BTC: %f | ETH: %f", *deposit.AccountID, btc, eth))

	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).ConfirmDepositAdminSuccess, nil))
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

	contracts, err := s.LendingRepository.QueryContractByAccountIDRepo(c.Context(), accountId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	var totalOutstanding float64
	for _, value := range *contracts {
		totalOutstanding += *value.LoanOutstanding
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
		ContractID: int(contractId),
	}
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).BorrowLoanSuccess, &borrowLoanResponse))
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
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).ConfirmContractAdminSuccess, nil))
}

func (s *lendingHandler) GetInterestTerm(c *handler.Ctx) error {
	interestTerm, err := s.LendingRepository.QueryInterestTermRepo(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).GetInterestTermSuccess, &interestTerm))
}
