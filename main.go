package main

import (
	"fmt"
	"lending-engine/account"
	"lending-engine/blockchain"
	"lending-engine/client"
	"lending-engine/docs"
	"lending-engine/internal/database"
	"lending-engine/internal/handler"
	"lending-engine/internal/redis"
	"lending-engine/lending"
	"lending-engine/logz"
	"lending-engine/mail"
	"lending-engine/middleware"
	"lending-engine/version"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"
	_ "time/tzdata"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func init() {
	runtime.GOMAXPROCS(1)
	initTimezone()
	initViper()
}

var isReady bool

// @title Lending Financial Services
// @version 1.0
// @description Lending finance for ICFIN company.
// @termsOfService http://swagger.io/terms/
// @contact.name K.Apiwattanawong
// @contact.url http://www.swagger.io/support
// @contact.email k.apiwattanawong@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:9090
// @BasePath /lending-engine
// @schemes http https
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	timeout := viper.GetDuration("app.timeout")

	app := fiber.New(fiber.Config{
		StrictRouting: true,
		CaseSensitive: true,
		Immutable:     true,
		ReadTimeout:   timeout,
		WriteTimeout:  timeout,
		IdleTimeout:   timeout,
	})

	logger, err := logz.NewLogConfig()
	if err != nil {
		log.Fatal(err)
	}

	pool := redis.NewRedisConn()
	defer pool.Close()

	postgresDB, err := database.NewPostgresConn()
	if err != nil {
		logger.Error(err.Error())
	}
	defer postgresDB.Close()

	if err := postgresDB.Ping(); err != nil {
		logger.Error(err.Error())
	}

	httpClient := client.NewClient()

	ethClient, err := ethclient.Dial(viper.GetString("blockchain.ethereum.rpc"))
	if err != nil {
		logger.Error(err.Error())
	}

	bscClient, err := ethclient.Dial(viper.GetString("blockchain.binance.rpc"))
	if err != nil {
		logger.Error(err.Error())
	}

	middle := middleware.NewMiddleware(
		logger,
		redis.NewCheckExpireDataRedisFn(pool),
		redis.NewGetStructDataRedisFn(pool),
		redis.NewSetStructWExpireRedisFn(pool),
		redis.NewDeleteDataRedisFn(pool),
	)

	swag := app.Group("/swagger")
	swag.Use(middle.BasicAuthenicationMiddleware())
	registerSwaggerRoute(swag)

	app.Use(middle.JSONMiddleware())
	app.Use(middle.ContextLocaleMiddleware())
	app.Use(middle.LoggingMiddleware())

	baseApi := app.Group(viper.GetString("app.context"))

	accountHandler := account.NewAccountHandler(
		account.NewAccountRepositoryDB(postgresDB),
		account.NewRequestVerifyEmailClientFn(httpClient),
		redis.NewSetDataNoExpireRedisFn(pool),
		redis.NewCheckExpireDataRedisFn(pool),
		redis.NewGetDeleteIntDataRedisFn(pool),
		redis.NewSetStructWExpireRedisFn(pool),
		redis.NewGetStructDataRedisFn(pool),
		redis.NewDeleteDataRedisFn(pool),
		mail.NewRequestMailOtpClientFn(httpClient),
	)

	lendingHandler := lending.NewLendingHandler(
		lending.NewLendingRepositoryDB(postgresDB),
		blockchain.NewQueryTransactionClientFn(ethClient, bscClient),
		redis.NewGetFloatDataRedisFn(pool),
	)

	mailhandler := mail.NewMailHandler(
		mail.NewQueryAccountByIDFn(postgresDB),
		redis.NewCheckExpireDataRedisFn(pool),
		redis.NewSetDataWExpireRedisFn(pool),
		redis.NewGetIntDataRedisFn(pool),
		redis.NewSetStructWExpireRedisFn(pool),
		mail.NewRequestMailOtpClientFn(httpClient),
	)

	baseApi.Get("/price", handler.Helper(lendingHandler.GetTokenPrice, logger))
	baseApi.Post("/price/calculation", handler.Helper(lendingHandler.PreCalculationLoan, logger))

	baseApi.Post("/signup", handler.Helper(accountHandler.SignUp, logger))
	baseApi.Post("/login", handler.Helper(accountHandler.Login, logger))
	baseApi.Get("/verify/email/:ref", handler.Helper(accountHandler.ConfirmVerifyEmail, logger))
	baseApi.Post("/reset", handler.Helper(accountHandler.RequestResetPassword, logger))
	baseApi.Put("/reset", handler.Helper(accountHandler.ResetPassword, logger))

	baseApi.Get("/admin/documentInfo", handler.Helper(accountHandler.GetDocumentInfoAdmin, logger))
	baseApi.Post("/admin/documentInfo", handler.Helper(accountHandler.CreateDocumentInfoAdmin, logger))
	baseApi.Put("/admin/documentInfo", handler.Helper(accountHandler.UpdateDocumentInfoAdmin, logger))

	baseApi.Get("/admin/interest", handler.Helper(lendingHandler.GetInterestTermAdmin, logger))
	baseApi.Post("/admin/interest", handler.Helper(lendingHandler.CreateInterestTermAdmin, logger))
	baseApi.Put("/admin/interest", handler.Helper(lendingHandler.UpdateInterestTermAdmin, logger))

	baseApi.Get("/admin/account", handler.Helper(accountHandler.GetAccountAdmin, logger))
	baseApi.Post("/admin/account/confirm", handler.Helper(accountHandler.ConfirmAccountAdmin, logger))
	baseApi.Post("/admin/account/reject", handler.Helper(accountHandler.RejectAccountAdmin, logger))
	baseApi.Put("/admin/account/document", handler.Helper(accountHandler.UpdateAccountDocumentAdmin, logger))

	baseApi.Get("/admin/wallet-transaction", handler.Helper(lendingHandler.GetWalletTransactionAdmin, logger))
	baseApi.Post("/admin/deposit/confirm", handler.Helper(lendingHandler.ConfirmDepositAdmin, logger))
	baseApi.Post("/admin/deposit/reject", handler.Helper(lendingHandler.RejectDepositAdmin, logger))
	baseApi.Post("/admin/withdraw/confirm", handler.Helper(lendingHandler.ConfirmWithdrawAdmin, logger))
	baseApi.Post("/admin/withdraw/reject", handler.Helper(lendingHandler.RejectWithdrawAdmin, logger))

	baseApi.Get("/admin/contract", handler.Helper(lendingHandler.GetLoanAdmin, logger))
	baseApi.Post("/admin/contract", handler.Helper(lendingHandler.ConfirmLoanAdmin, logger))

	baseApi.Get("/admin/repay", handler.Helper(lendingHandler.GetRepayAdmin, logger))
	baseApi.Post("/admin/repay/confirm", handler.Helper(lendingHandler.ConfirmRepayAdmin, logger))
	baseApi.Post("/admin/repay/reject", handler.Helper(lendingHandler.RejectRepayAdmin, logger))

	baseApi.Use(middle.AuthorizeTokenMiddleware())

	baseApi.Get("/terms", handler.Helper(accountHandler.GetTermsCondition, logger))
	baseApi.Post("/terms", handler.Helper(accountHandler.AcceptTermsCondition, logger))

	baseApi.Get("/wallet-transaction", handler.Helper(lendingHandler.GetWalletTransaction, logger))
	baseApi.Post("/deposit", handler.Helper(lendingHandler.SubmitDeposit, logger))

	baseApi.Get("/credit", handler.Helper(lendingHandler.GetCreditAvailable, logger))
	baseApi.Get("/contract", handler.Helper(lendingHandler.GetLoan, logger))

	baseApi.Get("/repay", handler.Helper(lendingHandler.GetRepay, logger))
	baseApi.Post("/repay", handler.Helper(lendingHandler.SubmitRepay, logger))

	baseApi.Get("/otp", handler.Helper(mailhandler.Otp, logger))

	baseApi.Use(middle.VerifyOTPMiddleware())

	baseApi.Post("/borrow", handler.Helper(lendingHandler.BorrowLoan, logger))
	baseApi.Post("/withdraw", handler.Helper(lendingHandler.SubmitWithdraw, logger))

	app.Get("/version", version.VersionHandler)
	app.Get("/liveness", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	app.Get("/readiness", func(c *fiber.Ctx) error {
		if isReady {
			return c.SendStatus(fiber.StatusOK)
		}
		return c.SendStatus(fiber.StatusServiceUnavailable)
	})

	logger.Info(fmt.Sprintf("⇨ http server started on [::]:%s", viper.GetString("app.port")))

	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", viper.GetString("app.port"))); err != nil {
			logger.Info(err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	isReady = true

	select {
	case <-c:
		logger.Info("terminating: by signal")
	}

	app.Shutdown()

	logger.Info("shutting down")
	os.Exit(0)
}

func initViper() {
	viper.SetDefault("app.name", "lending-engine")
	viper.SetDefault("app.port", "9090")
	viper.SetDefault("app.timeout", "60s")
	viper.SetDefault("app.context", "/lending-engine")

	viper.SetDefault("swagger.host", "localhost:9090")
	viper.SetDefault("swagger.user", "admin")
	viper.SetDefault("swagger.password", "password")

	viper.SetDefault("log.level", "debug")
	viper.SetDefault("log.env", "dev")

	viper.SetDefault("toggle.query-txn", false)

	viper.SetDefault("postgres.type", "postgres")
	viper.SetDefault("postgres.host", "localhost")
	viper.SetDefault("postgres.port", "5432")
	viper.SetDefault("postgres.username", "postgres")
	viper.SetDefault("postgres.password", "P@ssw0rd")
	viper.SetDefault("postgres.database", "lending")
	viper.SetDefault("postgres.timeout", 100)
	viper.SetDefault("postgres.sslmode", "disable")

	viper.SetDefault("redis.max-idle", 3)
	viper.SetDefault("redis.timeout", "60s")
	viper.SetDefault("redis.host", "localhost:6379")
	viper.SetDefault("redis.password", "P@ssw0rd")
	viper.SetDefault("redis.expired-otp", 180)
	viper.SetDefault("redis.limit-otp", 3)
	viper.SetDefault("redis.limit-request", 5)

	viper.SetDefault("client.timeout", "60s")
	viper.SetDefault("client.hidebody", true)
	viper.SetDefault("client.email-api.verification.url", "http://localhost:8080/email/verification")
	viper.SetDefault("client.email-api.verification.link", "http://www.icfin.finance.com/verify-email/{ref}")
	viper.SetDefault("client.email-api.otp.url", "http://localhost:8080/email/otp")
	viper.SetDefault("client.email-api.account", "yoisak09446@gmail.com")
	viper.SetDefault("client.email-api.verify-emil-template", "verify-email.html")
	viper.SetDefault("client.email-api.otp-template", "otp.html")

	viper.SetDefault("jwt.issuer", "admin")
	viper.SetDefault("jwt.expired-at", "60m")
	viper.SetDefault("jwt.secret-key", "ICFIN")

	viper.SetDefault("loan.haircut.btc", 0.5)
	viper.SetDefault("loan.haircut.eth", 0.5)
	viper.SetDefault("loan.interest", 0.05)

	viper.SetDefault("blockchain.ethereum.rpc", "https://rinkeby.infura.io/v3/9657539221eb40a79ce550650f0530a3")
	viper.SetDefault("blockchain.ethereum.chainId", 14)
	viper.SetDefault("blockchain.binance.rpc", "https://bsc-dataseed.binance.org/")
	viper.SetDefault("blockchain.binance.chainId", 56)
	viper.SetDefault("blockchain.address", "0xa9B6D99bA92D7d691c6EF4f49A1DC909822Cee46")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func initTimezone() {
	ict, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Printf("error loading location 'Asia/Bangkok': %v\n", err)
	}
	time.Local = ict
}

func registerSwaggerRoute(swag fiber.Router) {
	swag.Get("/*", swagger.New(swagger.Config{
		URL:         fmt.Sprintf("http://%s/swagger/doc.json", viper.GetString("swagger.host")),
		DeepLinking: false,
	}))
	docs.SwaggerInfo.Host = viper.GetString("swagger.host")
	docs.SwaggerInfo.BasePath = viper.GetString("app.context")
}
