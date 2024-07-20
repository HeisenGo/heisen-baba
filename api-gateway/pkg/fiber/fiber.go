package fiber

import (
	"github.com/gofiber/fiber/v2"
)

func SetupFiber() *fiber.App {
	app := fiber.New()

	// Middleware can be added here

	return app
}
