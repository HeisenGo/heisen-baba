package middlewares

import (
	"errors"
	"strings"
	"terminalpathservice/api/http/handlers/presenter"
	"terminalpathservice/pkg/jwt"
	"terminalpathservice/pkg/ports/clients/clients"

	"github.com/gofiber/fiber/v2"
)

func Auth(GRPCAuthClient clients.IAuthClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorization := c.Get("Authorization")

		if authorization == "" {
			return presenter.Unauthorized(c, errors.New("authorization header missing"))
		}

		// Split the Authorization header value
		parts := strings.Split(authorization, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return presenter.Unauthorized(c, errors.New("invalid authorization token format"))
		}

		//pureToken := parts[1]
		pureToken := parts[1]
		user, err := GRPCAuthClient.GetUserByToken(pureToken)
		if err != nil {
			return presenter.Unauthorized(c, err)
		}

		c.Locals(jwt.UserClaimKey, user)

		return c.Next()
	}
}
