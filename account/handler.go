package account

import (
	"fmt"
	"lending-engine/common"
	"lending-engine/internal/handler"
	"lending-engine/internal/redis"
	"lending-engine/mail"
	"lending-engine/response"
	"lending-engine/token"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type accountHandler struct {
	AccountRepository          AccountRepository
	RequestVerifyEmailClientFn RequestVerifyEmailClientFn
	SetDataNoExpireRedisFn     redis.SetDataNoExpireRedisFn
	CheckExpireDataRedisFn     redis.CheckExpireDataRedisFn
	GetDeleteIntDataRedisFn    redis.GetDeleteIntDataRedisFn
	SetStructWExpireRedisFn    redis.SetStructWExpireRedisFn
	GetStructDataRedisFn       redis.GetStructDataRedisFn
	DeleteDataRedisFn          redis.DeleteDataRedisFn
	RequestMailOtpClientFn     mail.RequestMailOtpClientFn
}

func NewAccountHandler(accountRepository AccountRepository, requestVerifyEmailClientFn RequestVerifyEmailClientFn, setDataNoExpireRedisFn redis.SetDataNoExpireRedisFn, checkExpireDataRedisFn redis.CheckExpireDataRedisFn, getDeleteIntDataRedisFn redis.GetDeleteIntDataRedisFn, setStructWExpireRedisFn redis.SetStructWExpireRedisFn, getStructDataRedisFn redis.GetStructDataRedisFn, deleteDataRedisFn redis.DeleteDataRedisFn, requestMailOtpClientFn mail.RequestMailOtpClientFn) *accountHandler {
	return &accountHandler{
		AccountRepository:          accountRepository,
		RequestVerifyEmailClientFn: requestVerifyEmailClientFn,
		SetDataNoExpireRedisFn:     setDataNoExpireRedisFn,
		CheckExpireDataRedisFn:     checkExpireDataRedisFn,
		GetDeleteIntDataRedisFn:    getDeleteIntDataRedisFn,
		SetStructWExpireRedisFn:    setStructWExpireRedisFn,
		GetStructDataRedisFn:       getStructDataRedisFn,
		DeleteDataRedisFn:          deleteDataRedisFn,
		RequestMailOtpClientFn:     requestMailOtpClientFn,
	}
}

// SignUp
// @Summary Sign up
// @Description sign up account
// @Tags Account
// @Accept json
// @Produce json
// @Param SignUp body account.SignUpRequest true "request body to sign up account"
// @Success 200 {object} response.Response{data=account.SignUpResponse} "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Router /signup [post]
func (s *accountHandler) SignUp(c *handler.Ctx) error {
	var req SignUpRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).SignUpAccountRequest, err.Error()))
	}
	if err := req.validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).SignUpAccountRequest, err.Error()))
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, err.Error()))
	}

	id, err := s.AccountRepository.SignUpAccountRepo(c.Context(), req.FirstName, req.LastName, req.Phone, req.Email, string(password), req.AccountNumber, req.CitizenName, req.CitizenCard, req.BookBankName, req.BookBankLedger)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).SignUpAccountRequest, "Duplicate email."))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	c.Log().Info(fmt.Sprintf("Sign up with email: %s success.", req.Email))

	if err := s.AccountRepository.CreateWalletRepo(c.Context(), int(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}

	ref := common.RandStringBytes(10)
	sendVerifyEmailClientRequest := SendVerifyEmailClientRequest{
		From: viper.GetString("client.email-api.account"),
		To: []string{
			req.Email,
		},
		Subject:  "Verify your email address to finish signing up for ICFin.finance",
		Template: viper.GetString("client.email-api.verify-emil-template"),
		Body: BodySendVerifyEmailClient{
			Link: strings.Replace(viper.GetString("client.email-api.verification.link"), "{ref}", ref, 1),
		},
		Auth: true,
	}
	if err := s.RequestVerifyEmailClientFn(c.Log(), string(c.Request().Header.Peek(common.XRequestID)), &sendVerifyEmailClientRequest); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).SignUpAccountThirdParty, err.Error()))
	}

	if err := s.SetDataNoExpireRedisFn(ref, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalRedis, err.Error()))
	}
	c.Log().Info(fmt.Sprintf("Verify your email with ref: %s", ref))

	signUpResponse := SignUpResponse{
		AccountID: id,
	}
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).SignUpAccountSuccess, &signUpResponse))
}

// ConfirmVerifyEmail
// @Summary Confirm Verify Email
// @Description confirm email by verifying reference
// @Tags Account
// @Accept json
// @Produce json
// @Param ref path string true "reference number"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Router /verify/email/{ref} [get]
func (s *accountHandler) ConfirmVerifyEmail(c *handler.Ctx) error {
	req := c.Params("ref")

	id, err := s.GetDeleteIntDataRedisFn(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalRedis, err.Error()))
	}
	if id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).ConfirmVerifyEmailRequest, "Ref doesn't exist."))
	}

	affect, err := s.AccountRepository.ConfirmVerifyEmailRepo(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if affect == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, "There is no affected row."))
	}

	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).ConfirmVerifyEmailSuccess, nil))
}

// GetAccountAdmin
// @Summary Get Account Admin
// @Description get account by account id or email
// @Tags Admin
// @Accept json
// @Produce json
// @Param accountId query int false "Account ID"
// @Param email query string false "Email"
// @Success 200 {object} response.Response{data=account.GetAccountAdminResponse} "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Router /admin/account [get]
func (s *accountHandler) GetAccountAdmin(c *handler.Ctx) error {
	var req GetAccountAdminRequest
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).GetAccountAdminRequest, err.Error()))
	}
	m := make(map[string]interface{})
	if req.AccountID != nil {
		m["x.account_id"] = req.AccountID
	}
	if req.Email != nil {
		m["x.email"] = req.Email
	}
	lists, err := s.AccountRepository.GetAccountRepo(c.Context(), m)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}

	var count int
	// var getAccountAdminResponses []GetAccountAdminResponse // return nil
	getAccountAdminResponses := make([]GetAccountAdminResponse, 0) // return []
	for index, value := range *lists {
		document := Document{
			DocumentID:   *value.DocumentID,
			DocumentType: *value.DocumentType,
			FileName:     *value.FileName,
			FileContext:  *value.FileContext,
			Tag:          *value.Tag,
		}
		if index != 0 {
			if *value.AccountID == getAccountAdminResponses[count].AccountID {
				getAccountAdminResponses[count].Document = append(getAccountAdminResponses[count].Document, document)
				continue
			} else {
				count++
			}
		}
		getAccountAdminResponse := GetAccountAdminResponse{
			AccountID:     *(*lists)[index].AccountID,
			FirstName:     *(*lists)[index].FirstName,
			LastName:      *(*lists)[index].LastName,
			Phone:         *(*lists)[index].Phone,
			Password:      *(*lists)[index].Password,
			AccountNumber: *(*lists)[index].AccountNumber,
			IsVerify:      *(*lists)[index].IsVerify,
			Status:        *(*lists)[index].Status,
			TermCondition: *(*lists)[index].CurrentAcceptVersion,
			Document:      []Document{document},
		}
		getAccountAdminResponses = append(getAccountAdminResponses, getAccountAdminResponse)
	}
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).GetAccountAdminSuccess, &getAccountAdminResponses))
}

// ConfirmAccountAdmin
// @Summary Confirm Account Admin
// @Description confirm account by account id
// @Tags Admin
// @Accept json
// @Produce json
// @Param ConfirmAccountAdmin body account.ConfirmAccountAdminRequest true "request body to confirm account"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Router /admin/account/confirm [post]
func (s *accountHandler) ConfirmAccountAdmin(c *handler.Ctx) error {
	var req ConfirmAccountAdminRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).ConfirmAccountAdminRequest, err.Error()))
	}
	if err := req.validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).ConfirmAccountAdminRequest, err.Error()))
	}

	account, err := s.AccountRepository.GetAccountByIDRepo(c.Context(), req.AccountID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if account == nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).ConfirmAccountAdminRequest, "ID doesn't exist."))
	}
	if *account.Status == common.ConfirmStatus {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).ConfirmAccountAdminRequest, "This id has already confirmed."))
	}

	accountRows, err := s.AccountRepository.UpdateAccountRepo(c.Context(), req.AccountID, common.ConfirmStatus)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if accountRows != 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, fmt.Sprintf("expected to affect 1 row, affected %d", accountRows)))
	}
	c.Log().Info(fmt.Sprintf("AccountID: %d - Status: %s", req.AccountID, common.ConfirmStatus))
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).ConfirmAccountAdminSuccess, nil))
}

// RejectAccountAdmin
// @Summary Reject Account Admin
// @Description reject account by account id
// @Tags Admin
// @Accept json
// @Produce json
// @Param RejectAccountAdmin body account.RejectAccountAdminRequest true "request body to reject account"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Router /admin/account/reject [post]
func (s *accountHandler) RejectAccountAdmin(c *handler.Ctx) error {
	var req RejectAccountAdminRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).RejectAccountAdminRequest, err.Error()))
	}
	if err := req.validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).RejectAccountAdminRequest, err.Error()))
	}

	account, err := s.AccountRepository.GetAccountByIDRepo(c.Context(), req.AccountID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if account == nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).RejectAccountAdminRequest, "ID doesn't exist."))
	}
	if *account.Status == common.RejectStatus {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).RejectAccountAdminRequest, "This id has already rejected."))
	}

	accountRows, err := s.AccountRepository.UpdateAccountRepo(c.Context(), req.AccountID, common.RejectStatus)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if accountRows != 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, fmt.Sprintf("expected to affect 1 row, affected %d", accountRows)))
	}
	c.Log().Info(fmt.Sprintf("AccountID: %d - Status: %s", req.AccountID, common.RejectStatus))
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).RejectAccountAdminSuccess, nil))
}

// UpdateAccountDocumentAdmin
// @Summary Update Account document Admin
// @Description update account document by account id and document id
// @Tags Admin
// @Accept json
// @Produce json
// @Param UpdateAccountDocumentAdmin body account.UpdateAccountDocumentAdminRequest true "request body to update account document"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Router /admin/account/document [put]
func (s *accountHandler) UpdateAccountDocumentAdmin(c *handler.Ctx) error {
	var req UpdateAccountDocumentAdminRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).UpdateAccountDocumentAdminRequest, err.Error()))
	}
	if err := req.validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).UpdateAccountDocumentAdminRequest, err.Error()))
	}

	docRows, err := s.AccountRepository.UpdateAccountDocumentRepo(c.Context(), req.AccountID, req.DocumentID, req.FileName, req.FileContext)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if docRows != 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, fmt.Sprintf("expected to affect 1 row, affected %d", docRows)))
	}
	c.Log().Info(fmt.Sprintf("AccountID: %d - DocumentID: %d - FileName: %s", req.AccountID, req.DocumentID, req.FileName))
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).UpdateAccountDocumentAdminSuccess, nil))
}

// GetTermsCondition
// @Summary Get Terms & Condition
// @Description get the latest terms & condition's user
// @Tags Account
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=account.TermsCondition} "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /terms [get]
func (s *accountHandler) GetTermsCondition(c *handler.Ctx) error {
	bearer := c.Locals(common.JWTClaimsKey).(*jwt.Token)
	claims := bearer.Claims.(jwt.MapClaims)
	id := claims["accountId"].(float64)

	term, err := s.AccountRepository.GetTermsConditionRepo(c.Context(), int(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).GetTermsConditionSuccess, &term))
}

// AcceptTermsCondition
// @Summary Accept Terms & Condition
// @Description user accept new version of terms & condition
// @Tags Account
// @Accept json
// @Produce json
// @Param AcceptTermsCondition body account.AcceptTermsConditionRequest true "request body to accept terms & condition"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /terms [post]
func (s *accountHandler) AcceptTermsCondition(c *handler.Ctx) error {
	bearer := c.Locals(common.JWTClaimsKey).(*jwt.Token)
	claims := bearer.Claims.(jwt.MapClaims)
	id := claims["accountId"].(float64)

	var req AcceptTermsConditionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).AcceptTermsConditionRequest, err.Error()))
	}
	if err := req.validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).AcceptTermsConditionRequest, err.Error()))
	}

	affect, err := s.AccountRepository.AcceptTermsConditionRepo(c.Context(), int(id), req.Version)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if affect == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, "There is no affected row."))
	}

	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).AcceptTermsConditionSuccess, nil))
}

// Login
// @Summary Login
// @Description login by user & password
// @Tags Account
// @Accept json
// @Produce json
// @Param Login body account.LoginRequest true "request body to login account"
// @Success 200 {object} response.Response{data=account.LoginResponse} "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Router /login [post]
func (s *accountHandler) Login(c *handler.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).LoginAccountRequest, err.Error()))
	}
	if err := req.validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).LoginAccountRequest, err.Error()))
	}

	account, err := s.AccountRepository.GetAccountByEmailRepo(c.Context(), req.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if account == nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).LoginAccountRequest, "Couldn't find this email."))
	}
	if !*account.IsVerify {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).LoginAccountRequest, "Please verify your account email address."))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*account.Password), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).LoginAccountRequest, "Wrong password."))
	}

	jwtToken, err := token.GenerateJWTToken(viper.GetString("jwt.secret-key"), *account.AccountID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, err.Error()))
	}

	loginResponse := LoginResponse{
		Token: jwtToken,
	}
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).SignUpAccountSuccess, &loginResponse))
}

// RequestResetPassword
// @Summary Request Reset Password
// @Description request to reset password
// @Tags Account
// @Accept json
// @Produce json
// @Param RequestReset body account.RequestResetPasswordRequest true "request body to request to reset password"
// @Success 200 {object} response.Response{data=account.RequestResetPasswordResponse} "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Router /reset [post]
func (s *accountHandler) RequestResetPassword(c *handler.Ctx) error {
	var req RequestResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).RequestResetPasswordRequest, err.Error()))
	}
	if err := req.validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).RequestResetPasswordRequest, err.Error()))
	}

	account, err := s.AccountRepository.GetAccountByEmailRepo(c.Context(), req.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if account == nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).RequestResetPasswordRequest, "Couldn't find this email."))
	}

	refNo := common.RandStringBytes(6)
	otp := common.RandIntBytes(6)
	sendMailOtpClientRequest := mail.SendMailOtpClientRequest{
		From: viper.GetString("client.email-api.account"),
		To: []string{
			*account.Email,
		},
		Subject:  "You have requested OTP",
		Template: viper.GetString("client.email-api.otp-template"),
		Body: mail.BodySendMailOtpClient{
			UserName: fmt.Sprintf("%s %s", *account.FirstName, *account.LastName),
			RefNo:    refNo,
			Otp:      otp,
		},
		Auth: true,
	}
	if err := s.RequestMailOtpClientFn(c.Log(), string(c.Request().Header.Peek(common.XRequestID)), &sendMailOtpClientRequest); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).RequestResetPasswordThirdParty, err.Error()))
	}

	resetPasswordData := ResetPasswordData{
		Otp:       otp,
		FailCount: 0,
		AccountID: *account.AccountID,
	}
	if err := s.SetStructWExpireRedisFn(refNo, viper.GetInt("redis.expired-otp"), &resetPasswordData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalRedis, err.Error()))
	}
	requestResetPasswordResponse := RequestResetPasswordResponse{
		ReferenceNo: refNo,
		ExpiredTime: time.Now().Add(time.Duration(viper.GetInt64("redis.expired-otp")) * time.Second).Format(common.DateYYYYMMDDHHMMSSFormat),
	}
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).RequestResetPasswordSuccess, &requestResetPasswordResponse))
}

// ResetPassword
// @Summary Reset Password
// @Description submit to reset password
// @Tags Account
// @Accept json
// @Produce json
// @Param ReferenceNo header string true "reference number."
// @Param OTP header string true "one time password."
// @Param ResetPassword body account.ResetPasswordRequest true "request body to reset password"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 400 {object} response.ErrResponse "Invalid OTP"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Router /reset [put]
func (s *accountHandler) ResetPassword(c *handler.Ctx) error {
	refNo := string(c.Request().Header.Peek(common.ReferenceOTPKey))
	otp := string(c.Request().Header.Peek(common.OTPKey))

	c.Log().Info(fmt.Sprintf("RefNo. - %s | OTP - %s", refNo, otp))

	var req ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).ResetPasswordRequest, err.Error()))
	}
	if err := req.validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).ResetPasswordRequest, err.Error()))
	}

	expiredTime, err := s.CheckExpireDataRedisFn(refNo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalRedis, err.Error()))
	}
	if expiredTime == -2 {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).OTPRequestInvalid, fmt.Sprintf("RefNo. - %s has already expired.", refNo)))
	}

	var resetPasswordData ResetPasswordData
	if err := s.GetStructDataRedisFn(refNo, &resetPasswordData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalRedis, err.Error()))
	}

	c.Log().Debug(fmt.Sprintf("IN REDIS [ OTP - %s | FailCount - %d | ExpiredTime - %d ]", resetPasswordData.Otp, resetPasswordData.FailCount, expiredTime))

	if resetPasswordData.FailCount >= viper.GetInt("redis.limit-otp") {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).OTPRequestFailLimit, fmt.Sprintf("RefNo. - %s has %d failed attempts.", refNo, resetPasswordData.FailCount)))
	}

	if resetPasswordData.Otp != otp {
		resetPasswordData.FailCount++
		if err := s.SetStructWExpireRedisFn(refNo, expiredTime, &resetPasswordData); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalRedis, err.Error()))
		}
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).OTPRequestInvalid, fmt.Sprintf("RefNo. - %s has %d failed attempts.", refNo, resetPasswordData.FailCount)))
	} else {
		if err := s.DeleteDataRedisFn(refNo); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalRedis, err.Error()))
		}
	}

	c.Log().Info("OTP is valid")

	newPass, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 10)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, err.Error()))
	}

	affect, err := s.AccountRepository.ConfirmChangePasswordRepo(c.Context(), resetPasswordData.AccountID, string(newPass))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if affect == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, "There is no affected row."))
	}
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).ResetPasswordSuccess, nil))
}
