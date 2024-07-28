package http

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"hotel/api/http/handlers"
	"hotel/api/http/middlewares"
	"hotel/config"
	"hotel/service"
)

func Run(cfg config.Config, app *service.AppContainer) {
	fiberApp := fiber.New()
	api := fiberApp.Group("/api/v1", middlewares.SetUserContext())
	registerHotelRoutes(api, app)
	registerRoomRoutes(api, app)

	log.Fatal(fiberApp.Listen(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.HttpPort)))
}
func registerHotelRoutes(router fiber.Router, app *service.AppContainer) {
	router = router.Group("/hotels")

	router.Post("",middlewares.SetUserContext(),handlers.CreateHotel(app.HotelService()),
	)
}


func registerRoomRoutes(router fiber.Router, app *service.AppContainer) {
	router = router.Group("/rooms")
	router.Post("",middlewares.SetUserContext(),handlers.CreateRoom(app.RoomService()),
	)
}