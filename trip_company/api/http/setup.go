package http

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"tripcompanyservice/api/http/handlers"
	"tripcompanyservice/config"
	"tripcompanyservice/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Run(cfg config.Server, app *service.AppContainer) {
	fiberApp := fiber.New()
	api := fiberApp.Group("/api/v1") //, middlerwares.SetUserContext())

	createGroupLogger := loggerSetup(fiberApp)

	secret := []byte(cfg.Secret)
	fmt.Println(api, secret)
	registerGlobalRoutes(api)
	registerTransportCompanyRoutes(api, app, secret, createGroupLogger("boards"))
	//registerPathRouts(api, app, secret)
	// run server
	log.Fatal(fiberApp.Listen(fmt.Sprintf("%s:%d", cfg.Host, cfg.HttpPort)))
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

func registerGlobalRoutes(router fiber.Router) {
	// Setup a simple health check route
	router.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Service is up and running")
	})
}

func registerTransportCompanyRoutes(router fiber.Router, app *service.AppContainer, secret []byte, loggerMiddleWare fiber.Handler) {
	router = router.Group("/companies")
	router.Use(loggerMiddleWare)

	router.Post("",
		//middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		//middlewares.Auth(secret),
		handlers.CreateTransportCompany(app.CompanyService()),
	)
	// router.Get("/my-boards",
	// 	middlewares.Auth(secret),
	// 	handlers.GetUserBoards(app.BoardService()),
	// )
	// router.Get("/publics",
	// 	middlewares.Auth(secret),
	// 	middlewares.SetupCacheMiddleware(5),
	// 	handlers.GetPublicBoards(app.BoardService()),
	// )
	// router.Get("/:boardID",
	// 	middlewares.Auth(secret),
	// 	handlers.GetFullBoardByID(app.BoardService()),
	// )

	// router.Delete("/:boardID",
	// 	middlewares.Auth(secret),
	// 	handlers.DeleteBoard(app.BoardService()),
	// )

	// router.Post("/invite", middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
	// 	middlewares.Auth(secret),
	// 	handlers.InviteToBoard(app.BoardServiceFromCtx))
}

// func registerTerminalRouts(router fiber.Router, app *service.AppContainer, secret []byte) {
// 	terminalGroup := router.Group("/terminals") //, middlerwares.Auth(secret), middlerwares.RoleChecker("user", "admin"))
// 	fmt.Print(secret)
// 	terminalGroup.Post("",
// 		handlers.CreateTerminal(app.TerminalService()),
// 	)

// 	terminalGroup.Get("", handlers.CityTerminals(app.TerminalService()))

// 	terminalGroup.Patch(":terminalID", handlers.PatchTerminal(app.TerminalService()))
// 	terminalGroup.Delete(":terminalID", handlers.DeleteTerminal(app.TerminalService()))
// }

// func registerPathRouts(router fiber.Router, app *service.AppContainer, secret []byte) {
// 	pathGroup := router.Group("/paths") //, middlerwares.Auth(secret), middlerwares.RoleChecker("user", "admin"))
// 	fmt.Print(secret)
// 	pathGroup.Post("",
// 		handlers.CreatePath(app.PathService()),
// 	)
// 	pathGroup.Get(":pathID", handlers.GetFullPathByID(app.PathService()))
// 	pathGroup.Get("", handlers.GetPathsByOriginDestinationType(app.PathService()))
// 	pathGroup.Patch(":pathID", handlers.PatchPath(app.PathService()))
// 	pathGroup.Delete(":pathID", handlers.DeletePath(app.PathService()))
// }
