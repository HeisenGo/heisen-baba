package http

import (
	"bankservice/api/http/handlers"
	"bankservice/api/http/middlewares"
	"bankservice/config"
	"bankservice/pkg/adapters"
	services "bankservice/service"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"os"
	"path/filepath"
)

func Run(cfg config.Config, app *services.AppContainer) {
	fiberApp := fiber.New()

	api := fiberApp.Group("/api/v1", middlewares.SetUserContext())

	createGroupLogger := loggerSetup(fiberApp)

	registerGlobalRoutes(api)

	api.Use(middlewares.Auth(app.AuthClient()))
	registerWalletRoutes(api, app, createGroupLogger("wallets"))

	log.Fatal(fiberApp.Listen(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.HTTPPort)))
}

func registerGlobalRoutes(router fiber.Router) {
	// Setup a simple health check route
	router.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Service is up and running")
	})
}

func registerWalletRoutes(router fiber.Router, app *services.AppContainer, loggerMiddleWare fiber.Handler) {
	router = router.Group("/bank")
	router.Use(loggerMiddleWare)

	router.Post("/add-card",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		handlers.AddCardToWallet(app.WalletServiceFromCtx),
	)
	router.Get("/cards",
		handlers.WalletCards(app.WalletService()),
	)
	router.Post("/deposit",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		handlers.Deposit(app.WalletServiceFromCtx),
	)
	router.Post("/withdraw",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		handlers.Withdraw(app.WalletServiceFromCtx),
	)
	router.Get("/my-wallet",
		handlers.GetWallet(app.WalletService()),
	)
}

func loggerSetup(app *fiber.App) func(groupName string) fiber.Handler {

	// Create the logs directory if it does not exist
	logDir := "./logs"
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		log.Fatalf("error creating logs directory: %v", err)
	}

	// Common format for console logging
	consoleLoggerConfig := logger.Config{
		Format:     "${time} [${ip}]:${port} ${status} - ${method} ${path} - ${latency} ${bytesSent} ${bytesReceived} ${userAgent}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "Local",
	}
	app.Use(logger.New(consoleLoggerConfig))

	// Function to create a logger middleware with separate log file
	createGroupLogger := func(groupName string) fiber.Handler {
		logFilePath := filepath.Join(logDir, groupName+".log")
		file, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}

		return logger.New(logger.Config{
			Format:     "${time} [${ip}]:${port} ${status} - ${method} ${path} - ${latency} ${bytesSent} ${bytesReceived} ${userAgent}\n",
			TimeFormat: "02-Jan-2006 15:04:05",
			TimeZone:   "Local",
			Output:     file,
		})
	}
	return createGroupLogger
}
