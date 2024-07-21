package handlers

import (
	"errors"
	"terminalpathservice/api/http/handlers/presenter"
	"terminalpathservice/internal"
	"terminalpathservice/internal/terminal"
	"terminalpathservice/service"

	"github.com/gofiber/fiber/v2"
)

func CreateTerminal(terminalService *service.TerminalService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var req presenter.TerminalRequest

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

		t := presenter.TerminalRequestToTerminal(&req)

		if err := terminalService.CreateTerminal(c.UserContext(), t); err != nil {
			if errors.Is(err, terminal.ErrInvalidType) || errors.Is(err, internal.ErrEmptyString) || errors.Is(err, internal.ErrConsecutiveSpaces) || errors.Is(err, internal.ErrExceedsMaxLength) || errors.Is(err, internal.ErrInvalidCharacters) {
				return presenter.BadRequest(c, err)
			}
			err := errors.New("Error")
			// apply trace ID here .... TODO
			return presenter.InternalServerError(c, err)
		}

		res := presenter.TerminalToTerminalRequest(*t)
		return presenter.Created(c, "Terminal created successfully", res)
	}
}
