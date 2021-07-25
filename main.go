package main

import (
	"fmt"
	"lending-engine/account"
	"lending-engine/blockchain"
	"lending-engine/client"
	"lending-engine/internal/database"
	"lending-engine/internal/handler"
	"lending-engine/internal/redis"
	"lending-engine/lending"
	"lending-engine/logz"
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

	_ "lending-engine/docs"

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

// @title Template Fiber API
// @version 1.0
// @description Template api with fiber framework.
// @termsOfService http://swagger.io/terms/
// @contact.name K.apiwattanawong
// @contact.url http://www.swagger.io/support
// @contact.email k.apiwattanawong@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:9090
// @BasePath /template-api
// @schemes http https
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

	middle := middleware.NewMiddleware(logger)

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
		redis.NewGetDeleteIntDataRedisFn(pool),
	)

	lendingHandler := lending.NewLendingHandler(
		lending.NewLendingRepositoryDB(postgresDB),
		blockchain.NewQueryTransactionClientFn(ethClient, bscClient),
		redis.NewGetFloatDataRedisFn(pool),
	)

	baseApi.Get("/price", handler.Helper(lendingHandler.GetTokenPrice, logger))
	baseApi.Post("/price/calculation", handler.Helper(lendingHandler.PreCalculationLoan, logger))

	baseApi.Post("/signup", handler.Helper(accountHandler.SignUp, logger))
	baseApi.Post("/login", handler.Helper(accountHandler.Login, logger))
	baseApi.Get("/verify/email/:ref", handler.Helper(accountHandler.ConfirmVerifyEmail, logger))

	baseApi.Get("/interest", handler.Helper(lendingHandler.GetInterestTerm, logger))

	baseApi.Get("/admin/deposit", handler.Helper(lendingHandler.GetDepositAdmin, logger))
	baseApi.Get("/admin/deposit/:id", handler.Helper(lendingHandler.ConfirmDepositAdmin, logger))

	baseApi.Get("/admin/contract", handler.Helper(lendingHandler.GetLoanAdmin, logger))
	baseApi.Get("/admin/contract/:id", handler.Helper(lendingHandler.ConfirmLoanAdmin, logger))

	baseApi.Use(middle.AuthorizeTokenMiddleware())

	baseApi.Get("/terms", handler.Helper(accountHandler.GetTermsCondition, logger))
	baseApi.Get("/terms/:version", handler.Helper(accountHandler.AcceptTermsCondition, logger))

	baseApi.Get("/deposit", handler.Helper(lendingHandler.GetDepositStatus, logger))
	baseApi.Post("/deposit", handler.Helper(lendingHandler.SubmitDeposit, logger))

	baseApi.Get("/credit", handler.Helper(lendingHandler.GetCreditAvailable, logger))
	baseApi.Post("/borrow", handler.Helper(lendingHandler.BorrowLoan, logger))

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

	viper.SetDefault("client.timeout", "60s")
	viper.SetDefault("client.hidebody", true)
	viper.SetDefault("client.email-api.url", "http://localhost:8080/email/verify")
	viper.SetDefault("client.email-api.link", "http://www.icfin.finance.com/verify-email/{ref}")
	viper.SetDefault("client.email-api.account", "yoisak09446@gmail.com")
	viper.SetDefault("client.email-api.verify-emil-template", "verify-email.html")

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
}
