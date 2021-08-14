package mail

import (
	"fmt"
	"lending-engine/common"
	"lending-engine/internal/handler"
	"lending-engine/internal/redis"
	"lending-engine/response"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

type mailHandler struct {
	QueryAccountByIDFn      QueryAccountByIDFn
	CheckExpireDataRedisFn  redis.CheckExpireDataRedisFn
	SetDataWExpireRedisFn   redis.SetDataWExpireRedisFn
	GetIntDataRedisFn       redis.GetIntDataRedisFn
	SetStructWExpireRedisFn redis.SetStructWExpireRedisFn
	RequestMailOtpClientFn  RequestMailOtpClientFn
}

func NewMailHandler(queryAccountByIDFn QueryAccountByIDFn, checkExpireDataRedisFn redis.CheckExpireDataRedisFn, setDataWExpireRedisFn redis.SetDataWExpireRedisFn, getIntDataRedisFn redis.GetIntDataRedisFn, setStructWExpireRedisFn redis.SetStructWExpireRedisFn, requestMailOtpClientFn RequestMailOtpClientFn) *mailHandler {
	return &mailHandler{
		QueryAccountByIDFn:      queryAccountByIDFn,
		CheckExpireDataRedisFn:  checkExpireDataRedisFn,
		SetDataWExpireRedisFn:   setDataWExpireRedisFn,
		GetIntDataRedisFn:       getIntDataRedisFn,
		SetStructWExpireRedisFn: setStructWExpireRedisFn,
		RequestMailOtpClientFn:  requestMailOtpClientFn,
	}
}

// RequestOTP
// @Summary Request OTP Email
// @Description request otp from email
// @Tags Mail
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=mail.OtpMailResponse} "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /otp [get]
func (s *mailHandler) Otp(c *handler.Ctx) error {
	bearer := c.Locals(common.JWTClaimsKey).(*jwt.Token)
	claims := bearer.Claims.(jwt.MapClaims)
	id := claims["accountId"].(float64)
	accountId := int(id)

	account, err := s.QueryAccountByIDFn(c.Context(), accountId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalDatabase, err.Error()))
	}
	if account == nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).GetOTPRequest, "ID doesn't exist."))
	}

	penaltyTime, err := s.CheckExpireDataRedisFn(fmt.Sprintf("%d-%s", accountId, common.PenaltyRedis))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalRedis, err.Error()))
	}
	requestTime := time.Now()
	if penaltyTime != -2 {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).OTPRequestDuplicate, fmt.Sprintf("AccountID: %d has already requested otp. It will expire at %s", accountId, requestTime.Add(time.Duration(penaltyTime)*time.Second).Format(common.DateYYYYMMDDHHMMSSFormat))))
	}

	accountIdStr := strconv.Itoa(accountId)
	expiredTime, err := s.CheckExpireDataRedisFn(accountIdStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalRedis, err.Error()))
	}
	if expiredTime == -2 {
		if err := s.SetDataWExpireRedisFn(accountIdStr, common.TimeDateSecondLeft(), 1); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalRedis, err.Error()))
		}
	} else {
		maxCount, err := s.GetIntDataRedisFn(accountIdStr)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalRedis, err.Error()))
		}
		if maxCount >= viper.GetInt("redis.limit-request") {
			return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).OTPRequestMaxLimit, fmt.Sprintf("AccountID: %d has requested otp exceeding limit per day.", accountId)))
		}
		if err := s.SetDataWExpireRedisFn(accountIdStr, common.TimeDateSecondLeft(), maxCount+1); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalRedis, err.Error()))
		}
	}

	if err := s.SetDataWExpireRedisFn(fmt.Sprintf("%d-%s", accountId, common.PenaltyRedis), viper.GetInt("redis.expired-otp"), requestTime.Format(common.DateYYYYMMDDHHMMSSFormat)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalRedis, err.Error()))
	}

	refNo := common.RandStringBytes(6)
	otp := common.RandIntBytes(6)
	referenceData := ReferenceData{
		Otp:       otp,
		FailCount: 0,
	}
	if err := s.SetStructWExpireRedisFn(refNo, viper.GetInt("redis.expired-otp"), &referenceData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalRedis, err.Error()))
	}

	sendMailOtpClientRequest := SendMailOtpClientRequest{
		From: viper.GetString("client.email-api.account"),
		To: []string{
			*account.Email,
		},
		Subject:  "You have requested OTP",
		Template: viper.GetString("client.email-api.otp.template"),
		Body: BodySendMailOtpClient{
			Name:  fmt.Sprintf("%s %s", *account.FirstName, *account.LastName),
			RefNo: refNo,
			Otp:   otp,
		},
		Auth: true,
	}

	if err := s.RequestMailOtpClientFn(c.Log(), string(c.Request().Header.Peek(common.XRequestID)), &sendMailOtpClientRequest); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).GetOTPThirdParty, err.Error()))
	}
	otpMailResponse := OtpMailResponse{
		ReferenceNo: refNo,
		ExpiredTime: requestTime.Add(time.Duration(viper.GetInt64("redis.expired-otp")) * time.Second).Format(common.DateYYYYMMDDHHMMSSFormat),
	}
	return c.Status(fiber.StatusOK).JSON(response.NewResponse(response.ResponseContextLocale(c.Context()).GetOTPSuccess, &otpMailResponse))
}
