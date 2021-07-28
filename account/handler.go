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
			Link: strings.Replace(viper.GetString("client.email-api.link"), "{ref}", ref, 1),
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

	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).SignUpAccountSuccess, id))
}

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

	var documents []Document
	for _, value := range *lists {
		document := Document{
			DocumentID:   *value.DocumentID,
			DocumentType: *value.DocumentType,
			FileName:     *value.FileName,
			FileContext:  *value.FileContext,
			Tag:          *value.Tag,
		}
		documents = append(documents, document)
	}

	getAccountAdminResponse := GetAccountAdminResponse{
		AccountID:     *(*lists)[0].AccountID,
		FirstName:     *(*lists)[0].FirstName,
		LastName:      *(*lists)[0].LastName,
		Phone:         *(*lists)[0].Phone,
		Password:      *(*lists)[0].Password,
		AccountNumber: *(*lists)[0].AccountNumber,
		IsVerify:      *(*lists)[0].IsVerify,
		Status:        *(*lists)[0].Status,
		TermCondition: *(*lists)[0].CurrentAcceptVersion,
		Document:      documents,
	}
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).GetAccountAdminSuccess, &getAccountAdminResponse))
}

func (s *accountHandler) ConfirmAccountAdmin(c *handler.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).ConfirmAccountAdminRequest, err.Error()))
	}

	account, err := s.AccountRepository.GetAccountByIDRepo(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if account == nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).ConfirmAccountAdminRequest, "ID doesn't exist."))
	}
	if *account.Status != common.PendingStatus {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, "This id has already confirmed."))
	}

	accountRows, err := s.AccountRepository.UpdateAccountRepo(c.Context(), id, common.ConfirmStatus)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if accountRows != 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, fmt.Sprintf("expected to affect 1 row, affected %d", accountRows)))
	}
	c.Log().Info(fmt.Sprintf("AccountID: %d | Status: %s", id, *account.Status))
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).ConfirmAccountAdminSuccess, nil))
}

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

func (s *accountHandler) AcceptTermsCondition(c *handler.Ctx) error {
	bearer := c.Locals(common.JWTClaimsKey).(*jwt.Token)
	claims := bearer.Claims.(jwt.MapClaims)
	id := claims["accountId"].(float64)

	version := c.Params("version")

	affect, err := s.AccountRepository.AcceptTermsConditionRepo(c.Context(), int(id), version)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if affect == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalOperation, "There is no affected row."))
	}

	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).AcceptTermsConditionSuccess, nil))
}

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
