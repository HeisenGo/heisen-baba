package http

import (
	"authservice/api/http/handlers"
	"authservice/api/http/middlewares"
	"authservice/config"
	"authservice/service"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/template/html/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func Run(cfg config.Config, app *service.AppContainer) {
	fiberApp := fiber.New(fiber.Config{
		Views: html.New("./templates", ".html"),
	})
	// Serve static files from the "assets" directory
	fiberApp.Static("/assets", "./assets")

	api := fiberApp.Group("/api/v1", middlewares.SetUserContext())

	createGroupLogger := loggerSetup(fiberApp)

	registerUIRoutes(fiberApp, app, createGroupLogger("templates"))
	// register global routes
	registerGlobalRoutes(api, app,
		createGroupLogger("global"),
		middlewares.SetupLimiterMiddleware(1, 1, cfg.Redis),
	)
	//secret := []byte(cfg.Server.TokenSecret)
	
	log.Fatal(fiberApp.Listen(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.HTTPPort)))
}

func registerUIRoutes(router fiber.Router, app *service.AppContainer, loggerMiddleWare fiber.Handler) {
	router.Use(loggerMiddleWare)
	router.Get("/", func(c *fiber.Ctx) error {
		// Render index
		return c.Render("index", fiber.Map{
			"Title": "HeisenFlow",
		})
	})
	router.Get("/swagger/*", fiberSwagger.WrapHandler)
	router.Get("/metrics", monitor.New(monitor.Config{Title: "HeisenFlow Metrics Page"}))
}

func registerGlobalRoutes(router fiber.Router, app *service.AppContainer, loggerMiddleWare fiber.Handler, limiterMiddleWare fiber.Handler) {
	router.Use(loggerMiddleWare)
	router.Post("/register", limiterMiddleWare, handlers.RegisterUser(app.AuthService()))
	router.Post("/login", handlers.LoginUser(app.AuthService()))
	router.Get("/refresh", handlers.RefreshToken(app.AuthService()))
}

func userRoleChecker() fiber.Handler {
	return middlewares.RoleChecker("user")
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
