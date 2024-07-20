package routes

import (
	"api-gateway/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, handler *handlers.Handler) {
	app.All("/api/:service/*", handler.ProxyRequest)
}
