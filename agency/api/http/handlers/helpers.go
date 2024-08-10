package handlers

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"agency/api/http/handlers/presenter"
	"agency/pkg/valuecontext"
	"strings"
)

type ServiceFactory[T any] func(context.Context) T

func SendError(c *fiber.Ctx, err error, status int) error {
	if status == 0 {
		status = fiber.StatusInternalServerError
	}

	c.Locals(valuecontext.IsTxError, err)

	return c.Status(status).JSON(map[string]any{
		"error_msg": err.Error(),
	})
}

func PageAndPageSize(c *fiber.Ctx) (int, int) {
	page, pageSize := c.QueryInt("page"), c.QueryInt("page_size")
	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 20
	}

	return page, pageSize
}

func BodyValidator[T any](req T) error {
	myValidator := presenter.GetValidator()
	if errs := myValidator.Validate(req); len(errs) > 0 {
		errMsgs := make([]string, 0)

		for _, err := range errs {
			errMsgs = append(errMsgs, err.Error)
		}

		return errors.New(strings.Join(errMsgs, "and"))
	}
	return nil
}
