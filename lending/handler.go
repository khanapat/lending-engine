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
	QueryTransactionClientFn   blockchain.QueryTransactionClientFn
	LendingRepository          LendingRepository
	GetFloatDataRedisFn        redis.GetFloatDataRedisFn
	RequestLiquidationClientFn RequestLiquidationClientFn
}

func NewLendingHandler(lendingRepository LendingRepository, queryTransactionClientFn blockchain.QueryTransactionClientFn, getFloatDataRedisFn redis.GetFloatDataRedisFn, requestLiquidationClientFn RequestLiquidationClientFn) *lendingHandler {
	return &lendingHandler{
		QueryTransactionClientFn:   queryTransactionClientFn,
		LendingRepository:          lendingRepository,
		GetFloatDataRedisFn:        getFloatDataRedisFn,
		RequestLiquidationClientFn: requestLiquidationClientFn,
	}
}

// GetTokenPrice
// @Summary Get Token Price
// @Description get token price, haircut and interest rate
// @Tags Lending
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=lending.GetTokenPriceResponse} "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Router /price [get]
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

// GetWalletTransaction
// @Summary Get Wallet Transaction
// @Description get wallet transaction by accountId
// @Tags Lending
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]lending.WalletTransaction} "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /wallet-transaction [get]
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

// SubmitDeposit
// @Summary Submit Deposit
// @Description submit deposit transaction
// @Tags Lending
// @Accept json
// @Produce json
// @Param SubmitDeposit body lending.SubmitDepositRequest true "request body to submit deposit"
// @Success 200 {object} response.Response{data=lending.SubmitDepositResponse} "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /deposit [post]
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

// SubmitWithdraw
// @Summary Submit Withdraw
// @Description submit withdraw transaction
// @Tags Lending
// @Accept json
// @Produce json
// @Param ReferenceNo header string true "reference number."
// @Param OTP header string true "one time password."
// @Param SubmitWithdraw body lending.SubmitWithdrawRequest true "request body to submit withdraw"
// @Success 200 {object} response.Response{data=lending.SubmitWithdrawResponse} "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /withdraw [post]
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

// GetWalletTransactionAdmin
// @Summary Get Wallet Transaction Admin
// @Description get wallet transaction by id, account id, address or txn type
// @Tags Admin
// @Accept json
// @Produce json
// @Param Id query int false "Transaction ID"
// @Param accountId query int false "Account ID"
// @Param address query string false "Address"
// @Param txnType query string false "Transaction Type"
// @Success 200 {object} response.Response{data=lending.WalletTransaction} "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Router /admin/wallet-transaction [get]
func (s *lendingHandler) GetWalletTransactionAdmin(c *handler.Ctx) error {
	var req GetWalletTransactionAdminRequest
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).GetWalletTransactionAdminRequest, err.Error()))
	}
	m := make(map[string]interface{})
	if req.TxnType != nil {
		m["txn_type"] = req.TxnType
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

// ConfirmDepositAdmin
// @Summary Confirm Deposit Admin
// @Description confirm deposit transaction by account id
// @Tags Admin
// @Accept json
// @Produce json
// @Param ConfirmDepositAdmin body lending.ConfirmDepositAdminRequest true "request body to confirm deposit"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Router /admin/deposit/confirm [post]
func (s *lendingHandler) ConfirmDepositAdmin(c *handler.Ctx) error {
	var req ConfirmDepositAdminRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).ConfirmDepositAdminRequest, err.Error()))
	}
	if err := req.validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).ConfirmDepositAdminRequest, err.Error()))
	}

	txn, err := s.LendingRepository.QueryWalletTransactionByIDRepo(c.Context(), req.ID)
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

	depositRows, err := s.LendingRepository.UpdateDepositRepo(c.Context(), req.ID, common.ConfirmStatus, time.Now().Format(common.DateYYYYMMDDHHMMSSFormat))
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
	c.Log().Info(fmt.Sprintf("TxnID: %d - Status: %s | AccountID: %d - BTC: %f - ETH: %f", req.ID, common.ConfirmStatus, *txn.AccountID, btc, eth))
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).ConfirmDepositAdminSuccess, nil))
}

// RejectDepositAdmin
// @Summary Reject Deposit Admin
// @Description reject deposit transaction by account id
// @Tags Admin
// @Accept json
// @Produce json
// @Param RejectDepositAdmin body lending.RejectDepositAdminRequest true "request body to reject deposit"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Router /admin/deposit/reject [post]
func (s *lendingHandler) RejectDepositAdmin(c *handler.Ctx) error {
	var req RejectDepositAdminRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).RejectDepostiAdminRequest, err.Error()))
	}
	if err := req.validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).RejectDepostiAdminRequest, err.Error()))
	}

	txn, err := s.LendingRepository.QueryWalletTransactionByIDRepo(c.Context(), req.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if txn == nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).RejectDepostiAdminRequest, "ID doesn't exist."))
	}
	if *txn.Status != common.PendingStatus {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).RejectDepostiAdminRequest, "This id has already confirmed or cancelled."))
	}
	if *txn.TxnType != common.DepositStatus {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).RejectDepostiAdminRequest, "This id isn't deposit method."))
	}

	depositRows, err := s.LendingRepository.UpdateDepositRepo(c.Context(), req.ID, common.RejectStatus, time.Now().Format(common.DateYYYYMMDDHHMMSSFormat))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if depositRows != 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, fmt.Sprintf("expected to affect 1 row, affected %d", depositRows)))
	}
	c.Log().Info(fmt.Sprintf("TxnID: %d - Status: %s", req.ID, common.RejectStatus))
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).RejectDepositAdminSuccess, nil))
}

// ConfirmWithdrawAdmin
// @Summary Confirm Withdraw Admin
// @Description confirm withdraw transaction by account id
// @Tags Admin
// @Accept json
// @Produce json
// @Param ConfirmWithdrawAdmin body lending.ConfirmWithdrawAdminRequest true "request body to confirm withdraw"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Router /admin/withdraw/confirm [post]
func (s *lendingHandler) ConfirmWithdrawAdmin(c *handler.Ctx) error {
	var req ConfirmWithdrawAdminRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).ConfirmWithdrawAdminRequest, err.Error()))
	}
	if err := req.validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).ConfirmWithdrawAdminRequest, err.Error()))
	}

	txn, err := s.LendingRepository.QueryWalletTransactionByIDRepo(c.Context(), req.ID)
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

	withdrawRows, err := s.LendingRepository.UpdateWithdrawRepo(c.Context(), req.ID, req.TxnHash, common.ConfirmStatus, time.Now().Format(common.DateYYYYMMDDHHMMSSFormat))
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
	c.Log().Info(fmt.Sprintf("TxnID: %d - Status: %s | AccountID: %d - BTC: %f - ETH: %f", req.ID, common.ConfirmStatus, *txn.AccountID, btc, eth))
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).ConfirmWithdrawAdminSuccess, nil))
}

// RejectWithdrawAdmin
// @Summary Reject Withdraw Admin
// @Description reject withdraw transaction by account id
// @Tags Admin
// @Accept json
// @Produce json
// @Param RejectWithdrawAdmin body lending.RejectWithdrawAdminRequest true "request body to reject withdraw"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Router /admin/withdraw/reject [post]
func (s *lendingHandler) RejectWithdrawAdmin(c *handler.Ctx) error {
	var req RejectWithdrawAdminRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).RejectWithdrawAdminRequest, err.Error()))
	}
	if err := req.validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).RejectWithdrawAdminRequest, err.Error()))
	}

	txn, err := s.LendingRepository.QueryWalletTransactionByIDRepo(c.Context(), req.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if txn == nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).RejectDepostiAdminRequest, "ID doesn't exist."))
	}
	if *txn.Status != common.PendingStatus {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).RejectDepostiAdminRequest, "This id has already confirmed or cancelled."))
	}
	if *txn.TxnType != common.WithdrawStatus {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).RejectDepostiAdminRequest, "This id isn't withdraw method."))
	}

	withdrawRows, err := s.LendingRepository.UpdateWithdrawRepo(c.Context(), req.ID, "-", common.RejectStatus, time.Now().Format(common.DateYYYYMMDDHHMMSSFormat))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if withdrawRows != 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, fmt.Sprintf("expected to affect 1 row, affected %d", withdrawRows)))
	}
	c.Log().Info(fmt.Sprintf("TxnID: %d - Status: %s", req.ID, common.RejectStatus))
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).RejectWithdrawAdminSuccess, nil))
}

// GetCreditAvailable
// @Summary Get Credit Available
// @Description get user's credit available by accountId
// @Tags Lending
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=lending.GetCreditAvailableResponse} "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /credit [get]
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

// GetContract
// @Summary Get Contract Loan
// @Description get user's loan contract by accountId
// @Tags Lending
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]lending.Contract} "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /contract [get]
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

// BorrowLoan
// @Summary Borrow Loan
// @Description borrow loan
// @Tags Lending
// @Accept json
// @Produce json
// @Param ReferenceNo header string true "reference number."
// @Param OTP header string true "one time password."
// @Param BorrowLoan body lending.BorrowLoanRequest true "request body to borrow loan"
// @Success 200 {object} response.Response{data=lending.BorrowLoanResponse} "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /borrow [post]
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

// GetLoanAdmin
// @Summary Get Loan Admin
// @Description get loan by contract id or account id
// @Tags Admin
// @Accept json
// @Produce json
// @Param contractId query string false "Contract ID"
// @Param accountId query int false "Account ID"
// @Success 200 {object} response.Response{data=[]lending.Contract} "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Router /admin/contract [get]
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

// ConfirmLoanAdmin
// @Summary Confirm Loan Admin
// @Description confirm loan contract by account id
// @Tags Admin
// @Accept json
// @Produce json
// @Param ConfirmLoanAdmin body lending.ConfirmLoanAdminRequest true "request body to confirm loan contract"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Router /admin/contract [post]
func (s *lendingHandler) ConfirmLoanAdmin(c *handler.Ctx) error {
	var req ConfirmLoanAdminRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).ConfirmContractAdminRequest, err.Error()))
	}
	if err := req.validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).ConfirmContractAdminRequest, err.Error()))
	}

	contract, err := s.LendingRepository.QueryContractByIDRepo(c.Context(), req.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if contract == nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).ConfirmContractAdminRequest, "ID doesn't exist."))
	}
	if *contract.Status != common.PendingStatus {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, "This id has already run or closed."))
	}

	contractRows, err := s.LendingRepository.UpdateContractRepo(c.Context(), req.ID, common.OngoingStatus, time.Now().Format(common.DateYYYYMMDDHHMMSSFormat))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if contractRows != 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, fmt.Sprintf("expected to affect 1 row, affected %d", contractRows)))
	}
	c.Log().Info(fmt.Sprintf("ContractID: %d - Status: %s", req.ID, *contract.Status))
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).ConfirmContractAdminSuccess, nil))
}

// GetInterestTermAdmin
// @Summary Get Interest Term Admin
// @Description get all of interest term
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]lending.InterestTerm} "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Router /admin/interest [get]
func (s *lendingHandler) GetInterestTermAdmin(c *handler.Ctx) error {
	interestTerm, err := s.LendingRepository.QueryInterestTermRepo(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).GetInterestTermSuccess, &interestTerm))
}

// CreateInterestTermAdmin
// @Summary Create Interest Term Admin
// @Description create new interest rate
// @Tags Admin
// @Accept json
// @Produce json
// @Param CreateInterestTerm body lending.CreateInterestTermAdminRequest true "request body to create interest term"
// @Success 200 {object} response.Response{data=lending.CreateInterestTermAdminResponse} "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Router /admin/interest [post]
func (s *lendingHandler) CreateInterestTermAdmin(c *handler.Ctx) error {
	var req CreateInterestTermAdminRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).CreateInterestTermAdminRequest, err.Error()))
	}
	if err := req.validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).CreateInterestTermAdminRequest, err.Error()))
	}

	interestCode, err := s.LendingRepository.InsertInterestTermRepo(c.Context(), req.InterestRate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}

	createInterestTermAdminResponse := CreateInterestTermAdminResponse{
		InterestCode: interestCode,
	}
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).CreateInterestTermAdminSuccess, &createInterestTermAdminResponse))
}

// UpdateInterestTermAdmin
// @Summary Update Interest Term Admin
// @Description update interest rate by interest code
// @Tags Admin
// @Accept json
// @Produce json
// @Param UpdateInterestTerm body lending.UpdateInterestTermAdminRequest true "request body to update interest term"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Router /admin/interest [put]
func (s *lendingHandler) UpdateInterestTermAdmin(c *handler.Ctx) error {
	var req UpdateInterestTermAdminRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).UpdateInterestTermAdminRequest, err.Error()))
	}
	if err := req.validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).UpdateInterestTermAdminRequest, err.Error()))
	}

	rows, err := s.LendingRepository.UpdateInterestTermRepo(c.Context(), req.InterestCode, req.InterestRate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, fmt.Sprintf("expected to affect 1 row, affected %d", rows)))
	}
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).UpdateInterestTermAdminSuccess, nil))
}

// GetRepay
// @Summary Get Repay
// @Description get repayment by account id
// @Tags Lending
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]lending.RepayTransaction} "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /repay [get]
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

// SubmitRepay
// @Summary Submit Repay
// @Description submit repayment
// @Tags Lending
// @Accept json
// @Produce json
// @Param SubmitRepay body lending.SubmitRepayRequest true "request body to submit repay"
// @Success 200 {object} response.Response{data=lending.SubmitRepayResponse} "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /repay [post]
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

// GetRepayAdmin
// @Summary Get Repay Admin
// @Description get repayment by id, contract id and account id
// @Tags Admin
// @Accept json
// @Produce json
// @Param id query int false "ID"
// @Param contractId query int false "Contract ID"
// @Param accountId query int false "Account ID"
// @Success 200 {object} response.Response{data=[]lending.RepayTransaction} "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Router /admin/repay [get]
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

// ConfirmRepayAdmin
// @Summary Confirm Repay Admin
// @Description confirm repayment by account id
// @Tags Admin
// @Accept json
// @Produce json
// @Param ConfirmRepayAdmin body lending.ConfirmRepayAdminRequest true "request body to confirm repay"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Router /admin/repay/confirm [post]
func (s *lendingHandler) ConfirmRepayAdmin(c *handler.Ctx) error {
	var req ConfirmRepayAdminRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).ConfirmRepaymentAdminRequest, err.Error()))
	}
	if err := req.validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).ConfirmRepaymentAdminRequest, err.Error()))
	}

	repay, err := s.LendingRepository.QueryRepayTransactionByIDRepo(c.Context(), req.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if repay == nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).ConfirmRepaymentAdminRequest, "ID doesn't exist."))
	}
	if *repay.Status != common.PendingStatus {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).ConfirmRepaymentAdminRequest, "This id has already confirmed."))
	}

	repayRows, err := s.LendingRepository.UpdateRepayTransactionRepo(c.Context(), req.ID, common.ConfirmStatus, time.Now().Format(common.DateYYYYMMDDHHMMSSFormat))
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
	c.Log().Info(fmt.Sprintf("RepayID: %d - Status: %s | ContractID: %d - Status: %s", req.ID, common.ConfirmStatus, *repay.ContractID, common.ClosedStatus))
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).ConfirmRepaymentAdminSuccess, nil))
}

// RejectRepayAdmin
// @Summary Reject Repay Admin
// @Description reject repayment by account id
// @Tags Admin
// @Accept json
// @Produce json
// @Param RejectRepayAdmin body lending.RejectRepayAdminRequest true "request body to reject repay"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Router /admin/repay/reject [post]
func (s *lendingHandler) RejectRepayAdmin(c *handler.Ctx) error {
	var req RejectRepayAdminRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).RejectRepaymentAdminRequest, err.Error()))
	}
	if err := req.validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).RejectRepaymentAdminRequest, err.Error()))
	}

	repay, err := s.LendingRepository.QueryRepayTransactionByIDRepo(c.Context(), req.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if repay == nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).RejectRepaymentAdminRequest, "ID doesn't exist."))
	}
	if *repay.Status != common.PendingStatus {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).RejectRepaymentAdminRequest, "This id has already confirmed."))
	}

	repayRows, err := s.LendingRepository.UpdateRepayTransactionRepo(c.Context(), req.ID, common.RejectStatus, time.Now().Format(common.DateYYYYMMDDHHMMSSFormat))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if repayRows != 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, fmt.Sprintf("expected to affect 1 row, affected %d", repayRows)))
	}
	c.Log().Info(fmt.Sprintf("RepayID: %d - Status: %s", req.ID, common.ConfirmStatus))
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).RejectRepaymentAdminSuccess, nil))
}

// LiquidateFundAdmin
// @Summary Liquidate Fund Admin
// @Description liquidate all fund by account id and contract id
// @Tags Admin
// @Accept json
// @Produce json
// @Param LiquidateFundAdmin body lending.LiquidateFundRequest true "request body to liquidate fund"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Router /admin/liquidation [post]
func (s *lendingHandler) LiquidateFundAdmin(c *handler.Ctx) error {
	var req LiquidateFundRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).LiquidateFundAdminRequest, err.Error()))
	}
	if err := req.validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).LiquidateFundAdminRequest, err.Error()))
	}

	liq, err := s.LendingRepository.LiquidationRepo(c.Context(), req.AccountID, req.ContractID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if liq == nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).LiquidateFundAdminRequest, "AccountID or ContractID doesn't exist."))
	}

	margin, err := time.ParseInLocation(time.RFC3339, *liq.MarginCallDate, time.Local)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, err.Error()))
	}
	count := int(time.Since(margin).Hours() / 24)
	c.Log().Info(fmt.Sprintf("Margin Count: %d", count))
	if count <= viper.GetInt("loan.liquidate-limit") {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).LiquidateFundAdminRequest, "Margin Call doesn't reach limit."))
	}

	if *liq.Status != common.OngoingStatus {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).LiquidateFundAdminRequest, "ContractID is inactive."))
	}

	rows, err := s.LendingRepository.UpdateWalletRepo(c.Context(), req.AccountID, 0.0, 0.0, nil, time.Now().Format(common.DateYYYYMMDDHHMMSSFormat))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if rows != 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, fmt.Sprintf("expected to affect 1 row, affected %d", rows)))
	}
	c.Log().Info(fmt.Sprintf("AccountID: %d | BTC: %f | ETH: %f", req.AccountID, 0.0, 0.0))

	contractRows, err := s.LendingRepository.UpdateContractRepo(c.Context(), req.ContractID, common.ClosedStatus, time.Now().Format(common.DateYYYYMMDDHHMMSSFormat))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if contractRows != 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, fmt.Sprintf("expected to affect 1 row, affected %d", contractRows)))
	}
	c.Log().Info(fmt.Sprintf("ContractID: %d - Status: %s", req.ContractID, common.ClosedStatus))

	sendLiquidationClientRequest := SendLiquidationClientRequest{
		From:     viper.GetString("client.email-api.account"),
		To:       []string{},
		Subject:  "Asset Liquidation Notice",
		Template: viper.GetString("client.email-api.liquidation.template"),
		Body: BodySendLiquidationClient{
			Name:       fmt.Sprintf("%s %s", *liq.FirstName, *liq.LastName),
			BTCAmount:  *liq.BTCVolume,
			ETHAmount:  *liq.ETHVolume,
			ContractID: req.ContractID,
		},
		Auth: true,
	}
	if err := s.RequestLiquidationClientFn(c.Log(), string(c.Request().Header.Peek(common.XRequestID)), &sendLiquidationClientRequest); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).LiquidateFundAdminThirdParty, err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).LiquidateFundAdminSuccess, nil))
}
