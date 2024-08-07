package handlers

import (
	"agency/api/http/handlers/presenter"
	"agency/service"
	"strconv"
	"github.com/gofiber/fiber/v2"
)

// CreateTour creates a new tour
// @Summary Create a new tour
// @Description Create a new tour
// @Tags tours
// @Accept json
// @Produce json
// @Param tour body presenter.CreateTourReq true "Tour to create"
// @Success 201 {object} presenter.FullTourResponse
// @Failure 400 {object} presenter.Response "error: bad request"
// @Failure 500 {object} presenter.Response "error: internal server error"
// @Router /tours [post]
func CreateTour(tourService *service.TourService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req presenter.CreateTourReq
		if err := c.BodyParser(&req); err != nil {
			return presenter.BadRequest(c, err)
		}

		// Validate request
		if err := BodyValidator(&req); err != nil {
			return presenter.BadRequest(c, err)
		}

		// Convert request to domain model
		t := presenter.CreateTourRequest(&req)

		// Create tour using the service
		if err := tourService.CreateTour(c.UserContext(), t); err != nil {
			return presenter.InternalServerError(c, err)
		}
		// Return response
		res := presenter.TourToCreateTourResponse(t)
		return presenter.Created(c, "Tour created successfully", res)
	}
}
// CreateTourReservation creates a new Tour reservation
// @Summary Create a new Tour reservation
// @Description Create a new Tour reservation
// @Tags reservations
// @Accept json
// @Produce json
// @Param reservation body presenter.ReservationCreateReq true "Reservation to create"
// @Success 201 {object} presenter.ReservationResp
// @Failure 400 {object} presenter.Response "error: bad request"
// @Failure 500 {object} presenter.Response "error: internal server error"
// @Router /reservations [post]
func CreateTourReservation(reservationService *service.TourService) fiber.Handler {
    return func(c *fiber.Ctx) error {
        var req presenter.ReservationCreateReq
        if err := c.BodyParser(&req); err != nil {
            return presenter.BadRequest(c, err)
        }
        res := presenter.ReservationReqToReservationDomain(&req)
        err := reservationService.CreateTourReservation(c.UserContext(), res)
        if err != nil {
            return presenter.InternalServerError(c, err)
        }
        response := presenter.ReservationToFullReservationResponse(res)
        return presenter.Created(c, "Reservation created successfully", response)
    }
}
// GetTours retrieves a paginated list of tours
// @Summary Get tours
// @Description Get a paginated list of tours
// @Tags tours
// @Produce json
// @Param agency_id query int true "Agency ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} presenter.FullTourResponse
// @Failure 500 {object} presenter.Response "error: internal server error"
// @Router /tours [get]
func GetTours(tourService *service.TourService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		agencyID := c.QueryInt("agency_id")
		page := c.QueryInt("page", 1)
		pageSize := c.QueryInt("page_size", 10)

		tours, total, err := tourService.GetTours(c.UserContext(), uint(agencyID), page, pageSize)
		if err != nil {
			return presenter.InternalServerError(c, err)
		}

		data := presenter.NewPagination(
			presenter.BatchToursToTourResponse(tours),
			uint(page),
			uint(pageSize),
			total,
		)
		if data.TotalPages == 0 {
			return presenter.NotFound(c, fiber.ErrNotFound)
		}
		return presenter.OK(c, "Tours fetched successfully", data)
	}
}

// GetToursByAgencyID retrieves tours by agency ID
// @Summary Get tours by agency ID
// @Description Get tours by agency ID
// @Tags tours
// @Produce json
// @Param agency_id query int true "Agency ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} presenter.FullTourResponse
// @Failure 500 {object} presenter.Response "error: internal server error"
// @Router /tours/agency [get]
func GetToursByAgencyID(tourService *service.TourService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		agencyID := c.QueryInt("agency_id")
		page := c.QueryInt("page", 1)
		pageSize := c.QueryInt("page_size", 10)

		tours, total, err := tourService.GetToursByAgencyID(c.UserContext(), uint(agencyID), page, pageSize)
		if err != nil {
			return presenter.InternalServerError(c, err)
		}

		data := presenter.NewPagination(
			presenter.BatchToursToTourResponse(tours),
			uint(page),
			uint(pageSize),
			uint(total),
		)
		if data.TotalPages == 0 {
			return presenter.NotFound(c, fiber.ErrNotFound)
		}
		return presenter.OK(c, "Tours fetched successfully", data)
	}
}

// UpdateTour updates a tour by ID
// @Summary Update a tour by ID
// @Description Update a tour by ID
// @Tags tours
// @Accept json
// @Produce json
// @Param id path int true "Tour ID"
// @Param tour body presenter.UpdateTourReq true "Tour to update"
// @Success 200 {object} presenter.FullTourResponse
// @Failure 400 {object} presenter.Response "error: bad request"
// @Failure 500 {object} presenter.Response "error: internal server error"
// @Router /tours/{id} [put]
func UpdateTour(tourService *service.TourService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tourID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return presenter.BadRequest(c, err)
		}

		var updateReq presenter.UpdateTourReq
		if err := c.BodyParser(&updateReq); err != nil {
			return presenter.BadRequest(c, err)
		}

		if err := BodyValidator(&updateReq); err != nil {
			return presenter.BadRequest(c, err)
		}

		tourUpdates := presenter.UpdateTourRequestToDomain(&updateReq)
		if err := tourService.UpdateTour(c.UserContext(), uint(tourID), tourUpdates); err != nil {
			return presenter.InternalServerError(c, err)
		}

		return presenter.OK(c, "Tour updated successfully", nil)
	}
}

// DeleteTour deletes a tour by ID
// @Summary Delete a tour by ID
// @Description Delete a tour by ID
// @Tags tours
// @Produce json
// @Param id path int true "Tour ID"
// @Success 204 "No Content"
// @Failure 400 {object} presenter.Response "error: bad request"
// @Failure 500 {object} presenter.Response "error: internal server error"
// @Router /tours/{id} [delete]
func DeleteTour(tourService *service.TourService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return presenter.BadRequest(c, err)
		}

		if err := tourService.DeleteTour(c.UserContext(), uint(id)); err != nil {
			return presenter.InternalServerError(c, err)
		}
		return presenter.NoContent(c)
	}
}

// ApproveTour approves a tour by ID
// @Summary Approve a tour by ID
// @Description Approve a tour by ID
// @Tags tours
// @Produce json
// @Param id path int true "Tour ID"
// @Success 200 {object} presenter.Response "Tour approved successfully"
// @Failure 400 {object} presenter.Response "Bad request"
// @Failure 500 {object} presenter.Response "Internal server error"
// @Router /tours/{id}/approve [patch]
func ApproveTour(tourService *service.TourService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tourID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return presenter.BadRequest(c, err)
		}

		if err := tourService.ApproveTour(c.UserContext(), uint(tourID)); err != nil {
			return presenter.InternalServerError(c, err)
		}

		return presenter.OK(c, "Tour approved successfully", nil)
	}
}

// SetTourStatus updates the active status of a tour
// @Summary Set the active status of a tour
// @Description Set the active status of a tour
// @Tags tours
// @Produce json
// @Param id path int true "Tour ID"
// @Param is_active query bool true "Active status"
// @Success 200 {object} presenter.Response "Tour status updated successfully"
// @Failure 400 {object} presenter.Response "Bad request"
// @Failure 500 {object} presenter.Response "Internal server error"
// @Router /tours/{id}/status [patch]
func SetTourStatus(tourService *service.TourService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tourID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return presenter.BadRequest(c, err)
		}

		isActive := c.Query("is_active") == "true"
		if err := tourService.SetTourStatus(c.UserContext(), uint(tourID), isActive); err != nil {
			return presenter.InternalServerError(c, err)
		}

		return presenter.OK(c, "Tour status updated successfully", nil)
	}
}
