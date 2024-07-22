package handlers

import (
	"errors"
	"fmt"
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
			if errors.Is(err, terminal.ErrDuplication) || errors.Is(err, terminal.ErrInvalidType) || errors.Is(err, internal.ErrEmptyString) || errors.Is(err, internal.ErrConsecutiveSpaces) || errors.Is(err, internal.ErrExceedsMaxLength) || errors.Is(err, internal.ErrInvalidCharacters) {
				return presenter.BadRequest(c, err)
			}
			fmt.Print(err)
			err := errors.New("Error")
			// apply trace ID here .... TODO
			return presenter.InternalServerError(c, err)
		}

		res := presenter.TerminalToTerminalRequest(*t)
		return presenter.Created(c, "Terminal created successfully", res)
	}
}

func CityTerminals(terminalService *service.TerminalService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		// if !ok {
		// 	return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		// }
		//query parameter
		page, pageSize := PageAndPageSize(c)
		city := c.Query("city")
		terminalType := c.Query("type")
		country := c.Query("country")
		if country == "" {
			return SendError(c, errors.New("country is required"), fiber.StatusBadRequest)
		}
		terminals, total, err := terminalService.GetTerminalsOfCity(c.UserContext(), country, city, terminalType, uint(page), uint(pageSize))
		if err != nil {
			status := fiber.StatusInternalServerError
			if errors.Is(err, terminal.ErrRecordsNotFound) {
				status = fiber.StatusBadRequest
			}
			err := errors.New("Error")
			return SendError(c, err, status)
		}
		fmt.Println(terminals)
		data := presenter.NewPagination(
			presenter.TerminalsToTerminalResponse(terminals),
			uint(page),
			uint(pageSize),
			total,
		)
		fmt.Println(data)
		return presenter.OK(c, "Terminals fetched successfully", data)
	}
}
