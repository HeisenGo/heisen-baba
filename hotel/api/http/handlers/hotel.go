package handlers

import (
	"hotel/api/http/handlers/presenter"
	"hotel/service"
	"github.com/gofiber/fiber/v2"
)
// CreateHotel creates a new hotel
// @Summary Create a new hotel
// @Description Create a new hotel
// @Tags hotels
// @Accept json
// @Produce json
// @Param hotel body presenter.CreateHotelReq true "Hotel to create"
// @Success 201 {object} presenter.CreateHotelResponse
// @Failure 400 {object} map[string]interface{} "error: bad request"
// @Failure 500 {object} map[string]interface{} "error: internal server error"
// @Router /hotels [post]
func CreateHotel(hotelService *service.HotelService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req presenter.CreateHotelReq

		if err := c.BodyParser(&req); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		h := presenter.CreateHotelRequest(&req)
		if err := hotelService.CreateHotel(c.UserContext(), h); err != nil {
			if err != nil {
				return presenter.BadRequest(c, err)
			}
			return presenter.InternalServerError(c, err)
		}
		res := presenter.HotelToCreateHotelResponse(h)
		return presenter.Created(c, "hotel created successfully", res)
	}
}
// GetHotel gets a hotel by ID
// @Summary Get a hotel by ID
// @Description Get a hotel by ID
// @Tags hotels
// @Produce json
// @Param id path int true "Hotel ID"
// @Success 200 {object} presenter.FullHotelResponse
// @Failure 400 {object} map[string]interface{} "error: bad request"
// @Failure 500 {object} map[string]interface{} "error: internal server error"
// @Router /hotels/{id} [get]
func GetHotel(hotelService *service.HotelService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return presenter.BadRequest(c, err)
		}
		h, err := hotelService.GetHotel(c.UserContext(), uint(id))
		if err != nil {
			return presenter.InternalServerError(c, err)
		}
		res := presenter.HotelToFullHotelResponse(h)
		return presenter.OK(c,"Hotel Fetched Successfully", res)
	}
}
// UpdateHotel updates a hotel by ID
// @Summary Update a hotel by ID
// @Description Update a hotel by ID
// @Tags hotels
// @Accept json
// @Produce json
// @Param id path int true "Hotel ID"
// @Param hotel body presenter.CreateHotelReq true "Hotel to update"
// @Success 200 {object} presenter.FullHotelResponse
// @Failure 400 {object} map[string]interface{} "error: bad request"
// @Failure 400 {object} map[string]interface{} "error: bad request"
// @Failure 500 {object} map[string]interface{} "error: internal server error"
// @Router /hotels/{id} [put]
func UpdateHotel(hotelService *service.HotelService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return presenter.BadRequest(c, err)
		}

		var req presenter.CreateHotelReq
		if err := c.BodyParser(&req); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		h := presenter.CreateHotelRequest(&req)
		h.ID = uint(id)
		if err := hotelService.UpdateHotel(c.UserContext(), h.ID,h); err != nil {
			return presenter.InternalServerError(c, err)
		}
		res := presenter.HotelToCreateHotelResponse(h)
		return presenter.OK(c,"Hotel Updated Succssesfully",res)
	}
}
// DeleteHotel deletes a hotel by ID
// @Summary Delete a hotel by ID
// @Description Delete a hotel by ID
// @Tags hotels
// @Produce json
// @Param id path int true "Hotel ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{} "error: bad request"
// @Failure 500 {object} map[string]interface{} "error: internal server error"
// @Router /hotels/{id} [delete]
func DeleteHotel(hotelService *service.HotelService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return presenter.BadRequest(c, err)
		}
		if err := hotelService.DeleteHotel(c.UserContext(), uint(id)); err != nil {
			return presenter.InternalServerError(c, err)
		}
		return presenter.NoContent(c)
	}
}