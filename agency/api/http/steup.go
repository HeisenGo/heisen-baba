package http

import (
	"fmt"
	"agency/api/http/handlers"
	"agency/api/http/middlewares"
	"agency/config"
	_ "agency/docs"
	"agency/pkg/adapters"
	"agency/service"
	"log"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func Run(cfg config.Config, app *service.AppContainer) {
	fiberApp := fiber.New()
	api := fiberApp.Group("/api/v1", middlewares.SetUserContext())
	registerGlobalRoutes(api)
	api.Use(middlewares.Auth(app.AuthClient()))
	registerAgencyRoutes(api, app)
	registerTourRoutes(api, app)
	// Add any additional routes here
	fiberApp.Get("/swagger/*", fiberSwagger.WrapHandler)

	log.Fatal(fiberApp.Listen(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.HTTPPort)))
}

func registerGlobalRoutes(router fiber.Router) {
	// Setup a simple health check route
	router.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Service is up and running")
	})
}

func registerAgencyRoutes(router fiber.Router, app *service.AppContainer) {
	router = router.Group("/agencies")
	router.Post("",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		handlers.CreateAgency(app.AgencyService()),
	)
	router.Get("",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		handlers.GetAgencies(app.AgencyService()),
	)
	router.Get("/owner",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		handlers.GetAgenciesByOwnerID(app.AgencyService()),
	)
	router.Put("/:id",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		handlers.UpdateAgency(app.AgencyService()),
	)
	router.Delete("/:id",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		handlers.DeleteAgency(app.AgencyService()),
	)
	router.Patch("/:id/block",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		handlers.BlockAgency(app.AgencyService()),
	)
}

func registerTourRoutes(router fiber.Router, app *service.AppContainer) {
	router = router.Group("/tours")
	router.Post("",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		handlers.CreateTour(app.TourService()),
	)
	router.Get("",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		handlers.GetTours(app.TourService()),
	)
	router.Get("/agency",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		handlers.GetToursByAgencyID(app.TourService()),
	)
	router.Put("/:id",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		handlers.UpdateTour(app.TourService()),
	)
	router.Delete("/:id",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		handlers.DeleteTour(app.TourService()),
	)
	router.Patch("/:id/approve",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		handlers.ApproveTour(app.TourService()),
	)
	router.Patch("/:id/status",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		handlers.SetTourStatus(app.TourService()),
	)
}

// Additional route registration functions can be added here
