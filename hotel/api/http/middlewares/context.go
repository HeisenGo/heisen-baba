package middlewares

import (
	"log/slog"
	"os"
	"github.com/gofiber/fiber/v2"
	"hotel/pkg/valuecontext"
)

func SetUserContext() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctxValue := &valuecontext.ContextValue{
			Logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		}

		c.SetUserContext(valuecontext.NewValueContext(c.UserContext(), ctxValue))

		return c.Next()
	}
}

