package handlers

import (
	"hotel/api/http/handlers/presenter"
	"hotel/service"
	"github.com/gofiber/fiber/v2"
)




func CreateHotel(hotelService *service.HotelService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req presenter.CreateHotelReq

		if err := c.BodyParser(&req); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		h:= presenter.CreateHotelRequest(&req)
		if err := hotelService.CreateHotel(c.UserContext(), h); err != nil {
			if err != nil{
				return presenter.BadRequest(c, err)
			}

			return presenter.InternalServerError(c, err)
		}
		res := presenter.HotelToCreateHotelResponse(h)
		return presenter.Created(c, "hotel created successfully", res)
	}
}