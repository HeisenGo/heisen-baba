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
	//api := fiberApp.Group("/api/v1", middlewares.SetUserContext())
	fiberApp.Get("/swagger/*", fiberSwagger.WrapHandler)

	log.Fatal(fiberApp.Listen(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.HttpPort)))
}
