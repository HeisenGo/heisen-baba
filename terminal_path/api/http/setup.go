package http

import (
	"fmt"
	"log"
	"terminalpathservice/api/http/handlers"
	"terminalpathservice/config"
	"terminalpathservice/service"

	"github.com/gofiber/fiber/v2"
)

func Run(cfg config.Server, app *service.AppContainer) {
	fiberApp := fiber.New()
	api := fiberApp.Group("/api/v1") //, middlerwares.SetUserContext())

	secret := []byte(cfg.Secret)
	fmt.Println(api, secret)
	registerTerminalRouts(api, app, secret)
	// run server
	log.Fatal(fiberApp.Listen(fmt.Sprintf("%s:%d", cfg.Host, cfg.HttpPort)))
}

func registerTerminalRouts(router fiber.Router, app *service.AppContainer, secret []byte) {
	terminalGroup := router.Group("/terminals") //, middlerwares.Auth(secret), middlerwares.RoleChecker("user", "admin"))
	fmt.Print(secret)
	terminalGroup.Post("",
		handlers.CreateTerminal(app.TerminalService()),
	)
}
