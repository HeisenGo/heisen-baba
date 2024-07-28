package handlers

import (
	"hotel/api/http/handlers/presenter"
	"hotel/service"
	"strconv"

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
// @Failure 400 {object} presenter.Response "error: bad request"
// @Failure 500 {object} presenter.Response "error: internal server error"
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
// GetHotels gets a paginated list of hotels
// @Summary Get hotels
// @Description Get paginated list of hotels with filters
// @Tags hotels
// @Produce json
// @Param city query string false "City"
// @Param country query string false "Country"
// @Param capacity query int false "Room capacity"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} presenter.FullHotelResponse
// @Failure 500 {object} presenter.Response
// @Router /hotels [get]
func GetHotels(hotelService *service.HotelService) fiber.Handler {
    return func(c *fiber.Ctx) error {
        city := c.Query("city")
        country := c.Query("country")
        capacity, _ := strconv.Atoi(c.Query("capacity"))
        // Pagination parameters
        page, _ := strconv.Atoi(c.Query("page", "1"))
        pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

        hotels, total, err := hotelService.GetHotels(c.Context(), city, country, capacity, page, pageSize)
        if err != nil {
            return presenter.InternalServerError(c, err)
        }

        response := presenter.NewPagination(hotels, uint(page), uint(pageSize), total)
        return presenter.OK(c, "Hotels fetched successfully", response)
    }
}

// GetHotelsByOwnerID gets hotels by owner ID
// @Summary Get hotels by owner ID
// @Description Get hotels by owner ID
// @Tags hotels
// @Produce json
// @Param owner_id query int true "Owner ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} presenter.FullHotelResponse
// @Failure 500 {object} presenter.Response
// @Router /hotels/owner [get]
func GetHotelsByOwnerID(hotelService *service.HotelService) fiber.Handler {
    return func(c *fiber.Ctx) error {
        ownerID := c.QueryInt("owner_id")
        page := c.QueryInt("page", 1)
        pageSize := c.QueryInt("page_size", 10)
		
        hotels, total, err := hotelService.GetHotelsByOwnerID(c.UserContext(), uint(ownerID), page, pageSize)
        if err != nil {
            return presenter.InternalServerError(c, err)
        }
        res := make([]presenter.FullHotelResponse, len(hotels))
        for i, hotel := range hotels {
            res[i] = *presenter.HotelToFullHotelResponse(&hotel) // Dereference the pointer here
        }
		
        pagination := presenter.NewPagination(res, uint(page), uint(pageSize), uint(total))
        return presenter.OK(c, "Hotels retrieved successfully", pagination)
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
// @Failure 400 {object} presenter.Response "error: bad request"
// @Failure 400 {object} presenter.Response "error: bad request"
// @Failure 500 {object} presenter.Response "error: internal server error"
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
// @Failure 400 {object} presenter.Response "error: bad request"
// @Failure 500 {object} presenter.Response "error: internal server error"
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