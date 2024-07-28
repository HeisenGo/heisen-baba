package http

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"hotel/api/http/handlers"
	"hotel/api/http/middlewares"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"hotel/config"
	"hotel/service"
	_ "hotel/docs"
)

func Run(cfg config.Config, app *service.AppContainer) {
	fiberApp := fiber.New()
	api := fiberApp.Group("/api/v1", middlewares.SetUserContext())
	registerHotelRoutes(api, app)
	registerRoomRoutes(api, app)
	fiberApp.Get("/swagger/*", fiberSwagger.WrapHandler)

	log.Fatal(fiberApp.Listen(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.HttpPort)))
}

func registerHotelRoutes(router fiber.Router, app *service.AppContainer) {
	router = router.Group("/hotels")
	router.Post("", middlewares.SetUserContext(), handlers.CreateHotel(app.HotelService()))
	router.Get("/:id", middlewares.SetUserContext(), handlers.GetHotel(app.HotelService()))
	router.Put("/:id", middlewares.SetUserContext(), handlers.UpdateHotel(app.HotelService()))
	router.Delete("/:id", middlewares.SetUserContext(), handlers.DeleteHotel(app.HotelService()))
}

func registerRoomRoutes(router fiber.Router, app *service.AppContainer) {
	router = router.Group("/rooms")
	router.Post("", middlewares.SetUserContext(), handlers.CreateRoom(app.RoomService()))
	router.Get("/:id", middlewares.SetUserContext(), handlers.GetRoom(app.RoomService()))
	router.Put("/:id", middlewares.SetUserContext(), handlers.UpdateRoom(app.RoomService()))
	router.Delete("/:id", middlewares.SetUserContext(), handlers.DeleteRoom(app.RoomService()))
}