package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"agency/api/http/handlers/presenter"
	"agency/internal/bank"
	"agency/service"
	"strconv"
)

// CreateReservation creates a new reservation
// @Summary Create a new reservation
// @Description Create a new reservation
// @Tags reservations
// @Accept json
// @Produce json
// @Param reservation body presenter.ReservationCreateReq true "Reservation to create"
// @Success 201 {object} presenter.ReservationResp
// @Failure 400 {object} presenter.Response "error: bad request"
// @Failure 500 {object} presenter.Response "error: internal server error"
// @Router /reservations [post]
func CreateReservation(reservationService *service.ReservationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req presenter.ReservationCreateReq

		if err := c.BodyParser(&req); err != nil {
			return presenter.BadRequest(c, err)
		}
		reservation := presenter.ReservationReqToReservationDomain(&req)
		err := reservationService.CreateReservation(c.UserContext(), reservation)
		if err != nil {
			if errors.Is(err, bank.ErrNotEnoughMoney) {
				return presenter.BadRequest(c, err)
			}
			return presenter.InternalServerError(c, err)
		}

		resp := presenter.ReservationToReservationResp(reservation)
		return presenter.Created(c, "Reservation created successfully", resp)
	}
}

// GetReservationsByHotelOwner gets reservations by agency owner ID
// @Summary Get reservations by agency owner ID
// @Description Get reservations by agency owner ID
// @Tags reservations
// @Produce json
// @Param owner_id query string true "Owner ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} presenter.FullReservationResponse
// @Failure 400 {object} presenter.Response "error: bad request"
// @Failure 500 {object} presenter.Response "error: internal server error"
// @Router /reservations/agency [get]
func GetReservationsByHotelOwner(reservationService *service.ReservationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ownerID, err := uuid.Parse(c.Query("owner_id"))
		if err != nil {
			return presenter.BadRequest(c, err)
		}

		page := c.QueryInt("page", 1)
		pageSize := c.QueryInt("page_size", 10)

		reservations, total, err := reservationService.GetReservationsByHotelOwner(c.UserContext(), ownerID, page, pageSize)
		if err != nil {
			return presenter.InternalServerError(c, err)
		}

		resp := presenter.BatchReservationsToReservationResponse(reservations)
		pagination := presenter.NewPagination(resp, uint(page), uint(pageSize), uint(total))
		return presenter.OK(c, "Hotels retrieved successfully", pagination)
	}
}

// GetReservationByUserID gets reservations by user ID
// @Summary Get reservations by user ID
// @Description Get reservations by user ID
// @Tags reservations
// @Produce json
// @Param user_id query string true "User ID"
// @Success 200 {object} presenter.ReservationResp
// @Failure 400 {object} presenter.Response "error: bad request"
// @Failure 500 {object} presenter.Response "error: internal server error"
// @Router /reservations/user [get]
func GetReservationByUserID(reservationService *service.ReservationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, err := uuid.Parse(c.Query("user_id"))
		if err != nil {
			return presenter.BadRequest(c, err)
		}

		reservations, err := reservationService.GetReservationByUserID(c.UserContext(), userID)
		if err != nil {
			return presenter.InternalServerError(c, err)
		}

		resp := presenter.BatchReservationsToReservationResponse(reservations)
		return presenter.OK(c, "Reservations retrieved successfully", resp)
	}
}

// GetReservationByID gets a reservation by its ID
// @Summary Get a reservation by ID
// @Description Get a reservation by ID
// @Tags reservations
// @Produce json
// @Param id path int true "Reservation ID"
// @Success 200 {object} presenter.ReservationResp
// @Failure 400 {object} presenter.Response "error: bad request"
// @Failure 500 {object} presenter.Response "error: internal server error"
// @Router /reservations/{id} [get]
func GetReservationByID(reservationService *service.ReservationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return presenter.BadRequest(c, err)
		}

		reservation, err := reservationService.GetReservationByID(c.UserContext(), uint(id))
		if err != nil {
			return presenter.InternalServerError(c, err)
		}

		resp := presenter.ReservationToFullReservationResponse(reservation)
		return presenter.OK(c, "Reservation retrieved successfully", resp)
	}
}

// UpdateReservation updates a reservation by its ID
// @Summary Update a reservation by ID
// @Description Update a reservation by ID
// @Tags reservations
// @Accept json
// @Produce json
// @Param id path int true "Reservation ID"
// @Param reservation body presenter.ReservationCreateReq true "Reservation to update"
// @Success 200 {object} presenter.ReservationResp
// @Failure 400 {object} presenter.Response "error: bad request"
// @Failure 500 {object} presenter.Response "error: internal server error"
// @Router /reservations/{id} [put]
func UpdateReservation(reservationService *service.ReservationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return presenter.BadRequest(c, err)
		}

		var req presenter.ReservationCreateReq
		if err := c.BodyParser(&req); err != nil {
			return presenter.BadRequest(c, err)
		}

		reservation := presenter.ReservationReqToReservationDomain(&req)
		err = reservationService.UpdateReservation(c.UserContext(), uint(id), reservation)
		if err != nil {
			return presenter.InternalServerError(c, err)
		}

		resp := presenter.ReservationToReservationResp(reservation)
		return presenter.OK(c, "Reservation updated successfully", resp)
	}
}

// DeleteReservation deletes a reservation by its ID
// @Summary Delete a reservation by ID
// @Description Delete a reservation by ID
// @Tags reservations
// @Produce json
// @Param id path int true "Reservation ID"
// @Success 204 "No Content"
// @Failure 400 {object} presenter.Response "error: bad request"
// @Failure 500 {object} presenter.Response "error: internal server error"
// @Router /reservations/{id} [delete]
func DeleteReservation(reservationService *service.ReservationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return presenter.BadRequest(c, err)
		}

		err = reservationService.DeleteReservation(c.UserContext(), uint(id))
		if err != nil {
			return presenter.InternalServerError(c, err)
		}

		return presenter.NoContent(c)
	}
}
