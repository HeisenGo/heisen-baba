package handlers

import (
	"errors"
	"terminalpathservice/api/http/handlers/presenter"
	"terminalpathservice/internal"
	"terminalpathservice/internal/path"
	"terminalpathservice/internal/terminal"
	"terminalpathservice/service"

	"github.com/gofiber/fiber/v2"
)

func CreatePath(pathService *service.PathService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var req presenter.PathRequest

		if err := c.BodyParser(&req); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		err := BodyValidator(req)
		if err != nil {
			return presenter.BadRequest(c, err)
		}

		//userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		// if !ok {
		// 	return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		// }

		p := presenter.PathRequestToPath(&req)

		if err := pathService.CreatePath(c.UserContext(), p); err != nil {
			if errors.Is(err, path.ErrMisMatchStartEndTerminalType) || errors.Is(err, path.ErrSameCitiesTerminals) || errors.Is(err, terminal.ErrTerminalNotFound) || errors.Is(err, internal.ErrEmptyString) || errors.Is(err, internal.ErrConsecutiveSpaces) || errors.Is(err, internal.ErrExceedsMaxLength) || errors.Is(err, internal.ErrInvalidCharacters) {
				return presenter.BadRequest(c, err)
			}
			err := errors.New("Error")
			// apply trace ID here .... TODO
			return presenter.InternalServerError(c, err)
		}

		res := presenter.PathToPathResponse(*p)
		return presenter.Created(c, "Path created successfully", res)
	}
}
