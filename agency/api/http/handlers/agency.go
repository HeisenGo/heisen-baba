package handlers

import (
	"agency/api/http/handlers/presenter"
	"agency/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"strconv"
)

// CreateAgency creates a new agency
// @Summary Create a new agency
// @Description Create a new agency
// @Tags agencies
// @Accept json
// @Produce json
// @Param agency body presenter.CreateAgencyReq true "Agency to create"
// @Success 201 {object} presenter.CreateAgencyResponse
// @Failure 400 {object} presenter.Response "error: bad request"
// @Failure 500 {object} presenter.Response "error: internal server error"
// @Router /agencies [post]
func CreateAgency(agencyService *service.AgencyService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req presenter.CreateAgencyReq

		if err := c.BodyParser(&req); err != nil {
			return presenter.BadRequest(c, err)
		}
		if err := BodyValidator(&req); err != nil {
			return presenter.BadRequest(c, err)
		}
		a := presenter.CreateAgencyRequest(&req)
		if err := agencyService.CreateAgency(c.UserContext(), a); err != nil {
			return presenter.InternalServerError(c, err)
		}
		res := presenter.AgencyToCreateAgencyResponse(a)
		return presenter.Created(c, "Agency created successfully", res)
	}
}

// GetAgencies gets a paginated list of agencies
// @Summary Get agencies
// @Description Get paginated list of agencies with filters
// @Tags agencies
// @Produce json
// @Param name query string false "Name"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} presenter.FullAgencyResponse
// @Failure 500 {object} presenter.Response
// @Router /agencies [get]
func GetAgencies(agencyService *service.AgencyService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		name := c.Query("name")
		page := c.QueryInt("page", 1)
		pageSize := c.QueryInt("page_size", 10)

		agencies, total, err := agencyService.GetAgencies(c.UserContext(), name, page, pageSize)
		if err != nil {
			return presenter.InternalServerError(c, err)
		}

		data := presenter.NewPagination(
			presenter.BatchAgenciesToAgencyResponse(agencies),
			uint(page),
			uint(pageSize),
			total,
		)
		if data.TotalPages == 0 {
			return presenter.NotFound(c, fiber.ErrNotFound)
		}
		return presenter.OK(c, "Agencies fetched successfully", data)
	}
}

// GetAgenciesByOwnerID gets agencies by owner ID
// @Summary Get agencies by owner ID
// @Description Get agencies by owner ID
// @Tags agencies
// @Produce json
// @Param owner_id query string true "Owner ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} presenter.FullAgencyResponse
// @Failure 400 {object} presenter.Response "error: bad request"
// @Failure 500 {object} presenter.Response
// @Router /agencies/owner [get]
func GetAgenciesByOwnerID(agencyService *service.AgencyService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ownerIDStr := c.Query("owner_id")
		ownerID, err := uuid.Parse(ownerIDStr)
		if err != nil {
			return presenter.BadRequest(c, err)
		}
		page := c.QueryInt("page", 1)
		pageSize := c.QueryInt("page_size", 10)

		agencies, total, err := agencyService.GetAgenciesByOwnerID(c.UserContext(), ownerID, page, pageSize)
		if err != nil {
			return presenter.InternalServerError(c, err)
		}

		res := make([]presenter.FullAgencyResponse, len(agencies))
		for i, agency := range agencies {
			res[i] = presenter.AgencyToFullAgencyResponse(agency)
		}

		pagination := presenter.NewPagination(res, uint(page), uint(pageSize), uint(total))
		return presenter.OK(c, "Agencies retrieved successfully", pagination)
	}
}

// UpdateAgency updates an agency by ID
// @Summary Update an agency by ID
// @Description Update an agency by ID
// @Tags agencies
// @Accept json
// @Produce json
// @Param id path int true "Agency ID"
// @Param agency body presenter.CreateAgencyReq true "Agency to update"
// @Success 200 {object} presenter.FullAgencyResponse
// @Failure 400 {object} presenter.Response "error: bad request"
// @Failure 500 {object} presenter.Response "error: internal server error"
// @Router /agencies/{id} [put]
func UpdateAgency(agencyService *service.AgencyService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		agencyID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return presenter.BadRequest(c, err)
		}

		var updateReq presenter.UpdateAgencyReq
		if err := c.BodyParser(&updateReq); err != nil {
			return presenter.BadRequest(c, err)
		}
		if err := BodyValidator(&updateReq); err != nil {
			return presenter.BadRequest(c, err)
		}
		if err := agencyService.UpdateAgency(c.UserContext(), uint(agencyID), presenter.UpdateAgencyRequestToDomain(&updateReq)); err != nil {
			return presenter.InternalServerError(c, err)
		}

		return presenter.OK(c, "Agency updated successfully", nil)
	}
}

// DeleteAgency deletes an agency by ID
// @Summary Delete an agency by ID
// @Description Delete an agency by ID
// @Tags agencies
// @Produce json
// @Param id path int true "Agency ID"
// @Success 204 "No Content"
// @Failure 400 {object} presenter.Response "error: bad request"
// @Failure 500 {object} presenter.Response "error: internal server error"
// @Router /agencies/{id} [delete]
func DeleteAgency(agencyService *service.AgencyService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return presenter.BadRequest(c, err)
		}
		if err := agencyService.DeleteAgency(c.UserContext(), uint(id)); err != nil {
			return presenter.InternalServerError(c, err)
		}
		return presenter.NoContent(c)
	}
}

// BlockAgency blocks an agency by its ID
// @Summary Block an agency by ID
// @Description Block an agency by its ID
// @Tags agencies
// @Produce json
// @Param id path int true "Agency ID"
// @Success 200 {object} presenter.Response "Agency blocked successfully"
// @Failure 400 {object} presenter.Response "Bad request"
// @Failure 500 {object} presenter.Response "Internal server error"
// @Router /agencies/{id}/block [patch]
func BlockAgency(agencyService *service.AgencyService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		agencyID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return presenter.BadRequest(c, err)
		}

		if err := agencyService.BlockAgency(c.UserContext(), uint(agencyID)); err != nil {
			return presenter.InternalServerError(c, err)
		}

		return presenter.OK(c, "Agency blocked successfully", nil)
	}
}
