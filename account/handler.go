package account

import (
	"fmt"
	"lending-engine/common"
	"lending-engine/internal/handler"
	"lending-engine/internal/redis"
	"lending-engine/response"
	"lending-engine/token"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type accountHandler struct {
	AccountRepository          AccountRepository
	RequestVerifyEmailClientFn RequestVerifyEmailClientFn
	SetDataNoExpireRedisFn     redis.SetDataNoExpireRedisFn
	GetDeleteIntDataRedisFn    redis.GetDeleteIntDataRedisFn
}

func NewAccountHandler(accountRepository AccountRepository, requestVerifyEmailClientFn RequestVerifyEmailClientFn, setDataNoExpireRedisFn redis.SetDataNoExpireRedisFn, getDeleteIntDataRedisFn redis.GetDeleteIntDataRedisFn) *accountHandler {
	return &accountHandler{
		AccountRepository:          accountRepository,
		RequestVerifyEmailClientFn: requestVerifyEmailClientFn,
		SetDataNoExpireRedisFn:     setDataNoExpireRedisFn,
		GetDeleteIntDataRedisFn:    getDeleteIntDataRedisFn,
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
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
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

	account, err := s.AccountRepository.GetAccountRepo(c.Context(), req.Email)
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
