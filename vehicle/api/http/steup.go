package http

import (
	"fmt"
	"vehicle/api/http/handlers"
	"vehicle/api/http/middlewares"
	"vehicle/config"
	_ "vehicle/docs"
	"vehicle/service"
	"log"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func Run(cfg config.Config, app *service.AppContainer) {
	fiberApp := fiber.New()
	api := fiberApp.Group("/api/v1", middlewares.SetUserContext())
	fiberApp.Get("/swagger/*", fiberSwagger.WrapHandler)
	registerVehicleRoutes(api, app)
	log.Fatal(fiberApp.Listen(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.HttpPort)))
}


func registerVehicleRoutes(router fiber.Router, app *service.AppContainer) {
	router = router.Group("/vehicles")
	
	// Create a new vehicle
	router.Post("", middlewares.SetUserContext(), handlers.CreateVehicle(app.VehicleService()))
	
	// Get a list of vehicles with optional filters
	router.Get("", middlewares.SetUserContext(), handlers.GetVehicles(app.VehicleService()))
	
	// Get vehicles by owner ID
	router.Get("/owner", middlewares.SetUserContext(), handlers.GetVehiclesByOwnerID(app.VehicleService()))
	
	// Update a vehicle by ID
	router.Put("/:id", middlewares.SetUserContext(), handlers.UpdateVehicle(app.VehicleService()))
	
	// Delete a vehicle by ID
	router.Delete("/:id", middlewares.SetUserContext(), handlers.DeleteVehicle(app.VehicleService()))
}