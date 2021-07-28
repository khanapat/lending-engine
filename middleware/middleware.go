package middleware

import (
	"fmt"
	"lending-engine/common"
	"lending-engine/internal/redis"
	"lending-engine/mail"
	"lending-engine/response"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	jwtware "github.com/gofiber/jwt/v2"
)

type middleware struct {
	ZapLogger               *zap.Logger
	CheckExpireDataRedisFn  redis.CheckExpireDataRedisFn
	GetStructDataRedisFn    redis.GetStructDataRedisFn
	SetStructWExpireRedisFn redis.SetStructWExpireRedisFn
	DeleteDataRedisFn       redis.DeleteDataRedisFn
}

func NewMiddleware(zapLogger *zap.Logger, checkExpireDataRedisFn redis.CheckExpireDataRedisFn, getStructDataRedisFn redis.GetStructDataRedisFn, setStructWExpireRedisFn redis.SetStructWExpireRedisFn, deleteDataRedisFn redis.DeleteDataRedisFn) *middleware {
	return &middleware{
		ZapLogger:               zapLogger,
		CheckExpireDataRedisFn:  checkExpireDataRedisFn,
		GetStructDataRedisFn:    getStructDataRedisFn,
		SetStructWExpireRedisFn: setStructWExpireRedisFn,
		DeleteDataRedisFn:       deleteDataRedisFn,
	}
}

func (m *middleware) BasicAuthenicationMiddleware() fiber.Handler {
	return basicauth.New(basicauth.Config{
		Users: map[string]string{
			viper.GetString("swagger.user"): viper.GetString("swagger.password"),
		},
		Realm: "Restricted",
		Authorizer: func(user, pass string) bool {
			if user == viper.GetString("swagger.user") && pass == viper.GetString("swagger.password") {
				return true
			}
			return false
		},
		Unauthorized:    nil,
		ContextUsername: "_user",
		ContextPassword: "_password",
	})
}

func (m *middleware) AuthorizeTokenMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:    []byte(viper.GetString("jwt.secret-key")),
		SigningMethod: "HS256",
		SuccessHandler: func(c *fiber.Ctx) error {
			return c.Next()
		},
		ErrorHandler: func(c *fiber.Ctx, e error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).AuthorizationToken, "Access Token is unauthorized."))
		},
		ContextKey: common.JWTClaimsKey,
	})
}

func (m *middleware) JSONMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts(common.ApplicationJSON)
		return c.Next()
	}
}

func (m *middleware) ContextLocaleMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals(common.LocaleKey, c.Query(common.LocaleKey))
		return c.Next()
	}
}

func (m *middleware) LoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		if c.Request().Header.Peek(common.XRequestID) == nil {
			c.Request().Header.Add(common.XRequestID, uuid.New().String())
		}

		logger := m.ZapLogger.With(zap.String(common.XRequestID, string(c.Request().Header.Peek(common.XRequestID))))

		logger.Debug(common.RequestInfoMsg,
			zap.String("method", string(c.Request().Header.Method())),
			zap.String("host", string(c.Request().Header.Host())),
			zap.String("path_uri", c.Request().URI().String()),
			zap.String("remote_addr", c.Context().RemoteAddr().String()),
			zap.String("body", string(c.Request().Body())),
		)

		if err := c.Next(); err != nil {
			return err
		}
		logger.Debug(common.ResponseInfoMsg,
			zap.String("body", string(c.Response().Body())),
		)
		logger.Info("Summary Information",
			zap.String("method", string(c.Request().Header.Method())),
			zap.String("path_uri", c.Request().URI().String()),
			zap.Duration("duration", time.Since(start)),
			zap.Int("status_code", c.Response().StatusCode()),
		)
		return nil
	}
}

func (m *middleware) VerifyOTPMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		bearer := c.Locals(common.JWTClaimsKey).(*jwt.Token)
		claims := bearer.Claims.(jwt.MapClaims)
		id := claims["accountId"].(float64)
		accountId := int(id)

		refNo := string(c.Request().Header.Peek(common.ReferenceOTPKey))
		otp := string(c.Request().Header.Peek(common.OTPKey))

		m.ZapLogger.Info(fmt.Sprintf("RefNo. - %s | OTP - %s", refNo, otp))

		expiredTime, err := m.CheckExpireDataRedisFn(refNo)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalRedis, err.Error()))
		}
		if expiredTime == -2 {
			return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).OTPRequestInvalid, fmt.Sprintf("RefNo. - %s has already expired.", refNo)))
		}

		var referenceData mail.ReferenceData
		if err := m.GetStructDataRedisFn(refNo, &referenceData); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalRedis, err.Error()))
		}

		m.ZapLogger.Debug(fmt.Sprintf("IN REDIS [ OTP - %s | FailCount - %d | ExpiredTime - %d ]", referenceData.Otp, referenceData.FailCount, expiredTime))

		if referenceData.FailCount >= viper.GetInt("redis.limit-otp") {
			return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).OTPRequestFailLimit, fmt.Sprintf("RefNo. - %s has %d failed attempts.", refNo, referenceData.FailCount)))
		}

		if referenceData.Otp != otp {
			referenceData.FailCount++
			if err := m.SetStructWExpireRedisFn(refNo, expiredTime, &referenceData); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalRedis, err.Error()))
			}
			return c.Status(fiber.StatusBadRequest).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).OTPRequestInvalid, fmt.Sprintf("RefNo. - %s has %d failed attempts.", refNo, referenceData.FailCount)))
		} else {
			if err := m.DeleteDataRedisFn(refNo); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalRedis, err.Error()))
			}
			if err := m.DeleteDataRedisFn(fmt.Sprintf("%d-%s", accountId, common.PenaltyRedis)); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrResponse(response.ResponseContextLocale(c.Context()).InternalRedis, err.Error()))
			}
		}

		m.ZapLogger.Info("OTP is valid")
		return c.Next()
	}
}
