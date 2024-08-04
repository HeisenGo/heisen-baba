package http

import (
	"fmt"
	"log"
	"vehicle/api/http/handlers"
	"vehicle/api/http/middlewares"
	"vehicle/config"
	"vehicle/service"
	_ "vehicle/docs"
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// Run initializes the Fiber application and sets up routes
func Run(cfg config.Config, app *service.AppContainer) {
	fiberApp := fiber.New()

	// Swagger UI for API documentation
	fiberApp.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Apply global middleware and setup routes
	api := fiberApp.Group("/api/v1", middlewares.SetUserContext())
	registerVehicleRoutes(api, app)

	// Start the server
	log.Fatal(fiberApp.Listen(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.HttpPort)))
}

// registerVehicleRoutes sets up routes for vehicle-related operations
func registerVehicleRoutes(router fiber.Router, app *service.AppContainer) {
	vehicles := router.Group("/vehicles")

	// Create a new vehicle
	vehicles.Post("", handlers.CreateVehicle(app.VehicleService()))

	// Get a list of vehicles with optional filters
	vehicles.Get("", handlers.GetVehicles(app.VehicleService()))

	// Get vehicles by owner ID
	vehicles.Get("/owner", handlers.GetVehiclesByOwnerID(app.VehicleService()))

	// Update a vehicle by ID
	vehicles.Put("/:id", handlers.UpdateVehicle(app.VehicleService()))

	// Delete a vehicle by ID
	vehicles.Delete("/:id", handlers.DeleteVehicle(app.VehicleService()))

	// Approve a vehicle by ID
	vehicles.Post("/:id/approve", handlers.ApproveVehicle(app.VehicleService()))

	// Set vehicle status by ID
	vehicles.Patch("/:id/status", handlers.SetVehicleStatus(app.VehicleService()))

	// Select vehicles based on passenger count and cost
	vehicles.Get("/select", handlers.SelectVehicles(app.VehicleService()))
}