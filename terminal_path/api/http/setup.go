package http

import (
	"fmt"
	"log"
	"terminalpathservice/api/http/handlers"
	"terminalpathservice/api/http/middlewares"
	"terminalpathservice/config"
	"terminalpathservice/service"

	"github.com/gofiber/fiber/v2"
)

func Run(cfg config.Config, app *service.AppContainer) {
	fiberApp := fiber.New()
	fiberApp.Use(middlewares.SetupLimiterMiddleware(1, 1, cfg.Redis))
	api := fiberApp.Group("/api/v1") //, middlerwares.SetUserContext())

	registerGlobalRoutes(api)
	registerTerminalRouts(api, app)
	registerPathRouts(api, app)
	// run server
	log.Fatal(fiberApp.Listen(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.HttpPort)))
}

func registerGlobalRoutes(router fiber.Router) {
	// Setup a simple health check route
	router.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Service is up and running")
	})
}

func registerTerminalRouts(router fiber.Router, app *service.AppContainer) {
	terminalGroup := router.Group("/terminals") //, middlerwares.Auth(secret), middlerwares.RoleChecker("user", "admin"))
	terminalGroup.Post("",
		middlewares.Auth(app.AuthClient()),
		handlers.CreateTerminal(app.TerminalService()),
	)

	terminalGroup.Get("",
		middlewares.SetupCacheMiddleware(2),
		handlers.CityTerminals(app.TerminalService()))

	terminalGroup.Patch(":terminalID", handlers.PatchTerminal(app.TerminalService()))
	terminalGroup.Delete(":terminalID", handlers.DeleteTerminal(app.TerminalService()))
}

func registerPathRouts(router fiber.Router, app *service.AppContainer) {
	pathGroup := router.Group("terminals/paths") //, middlerwares.Auth(secret), middlerwares.RoleChecker("user", "admin"))
	pathGroup.Post("",
		handlers.CreatePath(app.PathService()),
	)
	pathGroup.Get(":pathID", handlers.GetFullPathByID(app.PathService()))
	pathGroup.Get("",
		middlewares.SetupCacheMiddleware(2),
		handlers.GetPathsByOriginDestinationType(app.PathService()))
	pathGroup.Patch(":pathID", handlers.PatchPath(app.PathService()))
	pathGroup.Delete(":pathID", handlers.DeletePath(app.PathService()))
}
