package handlers

import (
	"errors"
	"terminalpathservice/api/http/handlers/presenter"
	"terminalpathservice/internal"
	"terminalpathservice/internal/path"
	"terminalpathservice/internal/terminal"
	"terminalpathservice/internal/user"
	"terminalpathservice/service"

	"github.com/gofiber/fiber/v2"
)

func CreatePath(pathService *service.PathService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var req presenter.PathRequest

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

		p := presenter.PathRequestToPath(&req)

		if err := pathService.CreatePath(c.UserContext(), p, userClaims.IsAdmin); err != nil {
			if errors.Is(err, service.ErrForbidden) {
				return presenter.Forbidden(c, err)
			}
			if errors.Is(err, path.ErrCodeIsImpossibleToUse) || errors.Is(err, path.ErrDuplication) || errors.Is(err, path.ErrMisMatchStartEndTerminalType) || errors.Is(err, path.ErrSameCitiesTerminals) || errors.Is(err, terminal.ErrTerminalNotFound) || errors.Is(err, internal.ErrEmptyString) || errors.Is(err, internal.ErrConsecutiveSpaces) || errors.Is(err, internal.ErrExceedsMaxLength) || errors.Is(err, internal.ErrInvalidCharacters) {
				return presenter.BadRequest(c, err)
			}
			if errors.Is(err, path.ErrFailedToRestore) {
				return presenter.BadRequest(c, err)
			}
			// apply trace ID here .... TODO
			return presenter.InternalServerError(c, err)
		}

		res := presenter.PathToPathResponse(*p)
		return presenter.Created(c, "Path created successfully", res)
	}
}

func GetPathsByOriginDestinationType(pathService *service.PathService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		// if !ok {
		// 	return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		// }
		//query parameter
		page, pageSize := PageAndPageSize(c)
		originCity := c.Query("from")
		destinationCity := c.Query("to")
		pathType := c.Query("type")

		paths, total, err := pathService.GetPathsByOriginDestinationType(c.UserContext(), originCity, destinationCity, pathType, uint(page), uint(pageSize))
		if err != nil {
			if errors.Is(err, terminal.ErrRecordsNotFound) {
				presenter.BadRequest(c, err)
			}
			err := errors.New("Error")
			return presenter.InternalServerError(c, err)
		}
		data := presenter.NewPagination(
			presenter.BatchPathsToPathResponse(paths),
			uint(page),
			uint(pageSize),
			total,
		)
		return presenter.OK(c, "Paths fetched successfully", data)
	}
}

func GetFullPathByID(pathService *service.PathService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		//query parameter
		pathID, err := c.ParamsInt("pathID")
		if err != nil {
			return presenter.BadRequest(c, errWrongIDType)
		}

		if pathID < 0 {
			return presenter.BadRequest(c, errWrongIDType)
		}

		path, err := pathService.GetFullPathByID(c.UserContext(), uint(pathID))
		if err != nil {
			if errors.Is(err, service.ErrForbidden) {
				return presenter.Forbidden(c, err)
			}
			if errors.Is(err, terminal.ErrRecordsNotFound) {
				return presenter.BadRequest(c, err)
			}
			return presenter.InternalServerError(c, err)
		}
		data := presenter.PathToPathResponse(*path)
		return presenter.OK(c, "Path fetched successfully", data)
	}
}

func PatchPath(pathService *service.PathService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var req presenter.UpdatePathRequest

		if err := c.BodyParser(&req); err != nil {
			return presenter.BadRequest(c, err)
		}

		userClaims, ok := c.Locals(UserClaimKey).(*user.User)
		if !ok {
			return presenter.BadRequest(c, errWrongClaimType)
		}
		pathID, err := c.ParamsInt("pathID")
		if err != nil {
			return presenter.BadRequest(c, errWrongIDType)
		}

		if pathID < 0 {
			return presenter.BadRequest(c, errWrongIDType)
		}

		newPath := presenter.UpdatePathReqToPath(&req, uint(pathID))
		changedPath, err := pathService.PatchPath(c.UserContext(), newPath, userClaims.IsAdmin)

		if err != nil {
			if errors.Is(err, service.ErrForbidden) {
				return presenter.Forbidden(c, err)
			}
			if errors.Is(err, path.ErrFailedToUpdate) || errors.Is(err, path.ErrPathNotFound) || errors.Is(err, path.ErrCanNotUpdatePath) || errors.Is(err, path.ErrDuplication) || errors.Is(err, path.ErrMisMatchStartEndTerminalType) || errors.Is(err, path.ErrSameCitiesTerminals) || errors.Is(err, terminal.ErrTerminalNotFound) || errors.Is(err, internal.ErrEmptyString) || errors.Is(err, internal.ErrConsecutiveSpaces) || errors.Is(err, internal.ErrExceedsMaxLength) || errors.Is(err, internal.ErrInvalidCharacters) {
				return presenter.BadRequest(c, err)
			}
			return presenter.InternalServerError(c, err)
		}
		res := presenter.PathToPathResponse(*changedPath)
		return presenter.OK(c, "Path updated successfully", res)
	}
}

func DeletePath(pathService *service.PathService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userClaims, ok := c.Locals(UserClaimKey).(*user.User)
		if !ok {
			return presenter.BadRequest(c, errWrongClaimType)
		}
		pathID, err := c.ParamsInt("pathID")
		if err != nil {
			return presenter.BadRequest(c, errWrongIDType)
		}

		if pathID < 0 {
			return presenter.BadRequest(c, errWrongIDType)
		}

		_, err = pathService.DeletePath(c.UserContext(), uint(pathID), userClaims.IsAdmin)

		if err != nil {
			if errors.Is(err, service.ErrForbidden) {
				return presenter.Forbidden(c, err)
			}
			if errors.Is(err, path.ErrCanNotDelete) || errors.Is(err, path.ErrPathNotFound) {
				return presenter.BadRequest(c, err)
			}
			err := errors.New("Error")
			// trace ID : TODO
			return presenter.InternalServerError(c, err)
		}
		return presenter.NoContent(c)
	}
}
