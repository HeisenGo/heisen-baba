package handlers

import (
	"errors"
	"terminalpathservice/api/http/handlers/presenter"
	"terminalpathservice/internal"
	"terminalpathservice/internal/terminal"
	"terminalpathservice/internal/user"
	"terminalpathservice/service"

	"github.com/gofiber/fiber/v2"
)

func CreateTerminal(terminalService *service.TerminalService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var req presenter.TerminalRequest

		if err := c.BodyParser(&req); err != nil {
			return presenter.BadRequest(c, err)
		}

		err := BodyValidator(req)
		if err != nil {
			return presenter.BadRequest(c, err)
		}

		userClaims, ok := c.Locals(UserClaimKey).(*user.User)
		if !ok {
			return presenter.BadRequest(c, errWrongClaimType)
		}

		t := presenter.TerminalRequestToTerminal(&req)

		if err := terminalService.CreateTerminal(c.UserContext(), t, userClaims.IsAdmin); err != nil {
			if errors.Is(err, service.ErrForbidden) {
				return presenter.Forbidden(c, err)
			}
			if errors.Is(err, terminal.ErrCityCountryDoNotExist) || errors.Is(err, terminal.ErrDuplication) || errors.Is(err, terminal.ErrInvalidType) || errors.Is(err, internal.ErrEmptyString) || errors.Is(err, internal.ErrConsecutiveSpaces) || errors.Is(err, internal.ErrExceedsMaxLength) || errors.Is(err, internal.ErrInvalidCharacters) {
				return presenter.BadRequest(c, err)
			}
			// apply trace ID here .... TODO
			return presenter.InternalServerError(c, err)
		}

		res := presenter.TerminalToTerminalRequest(*t)
		return presenter.Created(c, "Terminal created successfully", res)
	}
}

func CityTerminals(terminalService *service.TerminalService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		//query parameter
		page, pageSize := PageAndPageSize(c)
		city := c.Query("city")
		terminalType := c.Query("type")
		country := c.Query("country")
		if country == "" {
			return presenter.BadRequest(c, errors.New("country is required"))
		}
		terminals, total, err := terminalService.GetTerminals(c.UserContext(), country, city, terminalType, uint(page), uint(pageSize))
		if err != nil {
			if errors.Is(err, terminal.ErrRecordsNotFound) {
				return presenter.BadRequest(c, err)
			}
			return presenter.InternalServerError(c, err)
		}
		data := presenter.NewPagination(
			presenter.TerminalsToTerminalResponse(terminals),
			uint(page),
			uint(pageSize),
			total,
		)
		return presenter.OK(c, "Terminals fetched successfully", data)
	}
}

func PatchTerminal(terminalService *service.TerminalService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var req presenter.UpdateTerminalRequest

		if err := c.BodyParser(&req); err != nil {
			return presenter.BadRequest(c, err)
		}

		userClaims, ok := c.Locals(UserClaimKey).(*user.User)
		if !ok {
			return presenter.BadRequest(c, errWrongClaimType)
		}
		terminalID, err := c.ParamsInt("terminalID")
		if err != nil {
			return presenter.BadRequest(c, errWrongIDType)
		}

		if terminalID < 0 {
			return presenter.BadRequest(c, errWrongIDType)
		}

		updatedTerminal := presenter.UpdateTerminalToTerminal(&req, uint(terminalID))
		changedTerminal, err := terminalService.PatchTerminal(c.UserContext(), updatedTerminal, userClaims.IsAdmin)

		if err != nil {
			if errors.Is(err, service.ErrForbidden) {
				return presenter.Forbidden(c, err)
			}
			if errors.Is(err, terminal.ErrCityCountryDoNotExist) || errors.Is(err, terminal.ErrFailedToUpdate) || errors.Is(err, terminal.ErrTerminalNotFound) || errors.Is(err, terminal.ErrCanNotUpdate) {
				return presenter.BadRequest(c, err)
			}
			// trace ID : TODO
			return presenter.InternalServerError(c, err)
		}
		res := presenter.TerminalToTerminalRequest(*changedTerminal)
		return presenter.OK(c, "Terminal updated successfully", res)
	}
}

func DeleteTerminal(terminalService *service.TerminalService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userClaims, ok := c.Locals(UserClaimKey).(*user.User)
		if !ok {
			return presenter.BadRequest(c, errWrongClaimType)
		}
		terminalID, err := c.ParamsInt("terminalID")
		if err != nil {
			return presenter.BadRequest(c, errWrongIDType)
		}

		if terminalID < 0 {
			return presenter.BadRequest(c, errWrongIDType)
		}

		_, err = terminalService.DeleteTerminal(c.UserContext(), uint(terminalID), userClaims.IsAdmin)

		if err != nil {
			if errors.Is(err, terminal.ErrCanNotDelete) || errors.Is(err, terminal.ErrTerminalNotFound) {
				return presenter.BadRequest(c, err)
			}
			// trace ID : TODO
			return presenter.InternalServerError(c, err)
		}
		return presenter.NoContent(c)
	}
}
