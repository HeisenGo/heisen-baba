package http

import (
	"fmt"
	"hotel/api/http/handlers"
	"hotel/api/http/middlewares"
	"hotel/config"
	_ "hotel/docs"
	"hotel/service"
	"log"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
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
	router.Get("", middlewares.SetUserContext(), handlers.GetHotels(app.HotelService()))
	router.Get("/owner", middlewares.SetUserContext(), handlers.GetHotelsByOwnerID(app.HotelService()))
	router.Put("/:id", middlewares.SetUserContext(), handlers.UpdateHotel(app.HotelService()))
	router.Delete("/:id", middlewares.SetUserContext(), handlers.DeleteHotel(app.HotelService()))
	router.Patch("/:id/block", middlewares.SetUserContext(), handlers.BlockHotel(app.HotelService()))
}

func registerRoomRoutes(router fiber.Router, app *service.AppContainer) {
	router = router.Group("/rooms")
	router.Post("", middlewares.SetUserContext(), handlers.CreateRoom(app.RoomService()))
	router.Get("", middlewares.SetUserContext(), handlers.GetRooms(app.RoomService()))
	router.Get("/:id", middlewares.SetUserContext(), handlers.GetRooms(app.RoomService()))
	router.Put("/:id", middlewares.SetUserContext(), handlers.UpdateRoom(app.RoomService()))
	router.Delete("/:id", middlewares.SetUserContext(), handlers.DeleteRoom(app.RoomService()))
}
