package http

import (
	"fmt"
	"hotel/api/http/handlers"
	"hotel/api/http/middlewares"
	"hotel/config"
	_ "hotel/docs"
	"hotel/pkg/adapters"
	"hotel/service"
	"log"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func Run(cfg config.Config, app *service.AppContainer) {
	fiberApp := fiber.New()
	api := fiberApp.Group("/api/v1", middlewares.SetUserContext())
	registerGlobalRoutes(api)
	api.Use(middlewares.Auth(app.AuthClient()))
	registerHotelRoutes(api, app)
	registerRoomRoutes(api, app)
	registerReservationRoutes(api, app)
	fiberApp.Get("/swagger/*", fiberSwagger.WrapHandler)

	log.Fatal(fiberApp.Listen(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.HTTPPort)))
}

func registerGlobalRoutes(router fiber.Router) {
	// Setup a simple health check route
	router.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Service is up and running")
	})
}

func registerHotelRoutes(router fiber.Router, app *service.AppContainer) {
	router = router.Group("/hotels")
	router.Post("",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		handlers.CreateHotel(app.HotelService()),
	)
	router.Get("",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		handlers.GetHotels(app.HotelService()),
	)
	router.Get("/owner",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		handlers.GetHotelsByOwnerID(app.HotelService()),
	)
	router.Put("/:id",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		handlers.UpdateHotel(app.HotelService()),
	)
	router.Delete("/:id",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		handlers.DeleteHotel(app.HotelService()),
	)
	router.Patch("/:id/block",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		handlers.BlockHotel(app.HotelService()),
	)
}

func registerRoomRoutes(router fiber.Router, app *service.AppContainer) {
	router = router.Group("/hotels/:hotel_id/rooms")
	router.Post("",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		handlers.CreateRoom(app.RoomService()),
	)
	router.Get("",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		handlers.GetRooms(app.RoomService()),
	)
	// router.Get("/:id",
	// 	middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
	// 	handlers.GetRoomByID(app.RoomService()),
	// )
	router.Put("/:id",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		handlers.UpdateRoom(app.RoomService()),
	)
	router.Delete("/:id",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		handlers.DeleteRoom(app.RoomService()),
	)
}

func registerReservationRoutes(router fiber.Router, app *service.AppContainer) {
	router = router.Group("/reservations")
	router.Post("",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		handlers.CreateReservation(app.ReservationService()),
	)
	router.Get("/hotel",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		handlers.GetReservationsByHotelOwner(app.ReservationService()),
	)
	router.Get("/user",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		handlers.GetReservationByUserID(app.ReservationService()),
	)
	router.Get("/:id",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		handlers.GetReservationByID(app.ReservationService()),
	)
	router.Put("/:id",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		handlers.UpdateReservation(app.ReservationService()),
	)
	router.Delete("/:id",
		middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		handlers.DeleteReservation(app.ReservationService()),
	)
}
