package handlers

import (
	"strconv"
	"vehicle/api/http/handlers/presenter"
	"vehicle/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreateVehicle creates a new vehicle
// @Summary Create a new vehicle
// @Description Create a new vehicle
// @Tags vehicles
// @Accept json
// @Produce json
// @Param vehicle body presenter.CreateVehicleReq true "Vehicle to create"
// @Success 201 {object} presenter.CreateVehicleResponse
// @Failure 400 {object} presenter.Response "error: bad request"
// @Failure 500 {object} presenter.Response "error: internal server error"
// @Router /vehicles [post]
func CreateVehicle(vehicleService *service.VehicleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req presenter.CreateVehicleReq

		if err := c.BodyParser(&req); err != nil {
			return presenter.BadRequest(c, err)
		}

		v := presenter.CreateVehicleRequest(&req)
		if err := vehicleService.CreateVehicle(c.UserContext(), v); err != nil {
			return presenter.InternalServerError(c, err)
		}
		res := presenter.VehicleToCreateVehicleResponse(v)
		return presenter.Created(c, "Vehicle created successfully", res)
	}
}

// GetVehicles gets a paginated list of vehicles
// @Summary Get vehicles
// @Description Get paginated list of vehicles with filters
// @Tags vehicles
// @Produce json
// @Param type query string false "Vehicle type"
// @Param capacity query uint false "Capacity"
// @Param page query uint false "Page number"
// @Param page_size query uint false "Page size"
// @Success 200 {object} presenter.FullVehicleResponse
// @Failure 500 {object} presenter.Response
// @Router /vehicles [get]
func GetVehicles(vehicleService *service.VehicleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		vehicleType := c.Query("type")
		capacity, _ := strconv.ParseUint(c.Query("capacity"), 10, 32)
		page := c.QueryInt("page", 1)
		pageSize := c.QueryInt("page_size", 10)

		vehicles, total, err := vehicleService.GetVehicles(c.UserContext(), vehicleType, uint(capacity), page, pageSize)
		if err != nil {
			return presenter.InternalServerError(c, err)
		}

		data := presenter.NewPagination(
			presenter.BatchVehiclesToVehicleResponse(vehicles),
			uint(page),
			uint(pageSize),
			total,
		)
		if data.TotalPages == 0 {
			return presenter.NotFound(c, fiber.ErrNotFound)
		}
		return presenter.OK(c, "Vehicles fetched successfully", data)
	}
}

// GetVehiclesByOwnerID gets vehicles by owner ID
// @Summary Get vehicles by owner ID
// @Description Get vehicles by owner ID
// @Tags vehicles
// @Produce json
// @Param owner_id query string true "Owner ID"
// @Param page query uint false "Page number"
// @Param page_size query uint false "Page size"
// @Success 200 {object} presenter.FullVehicleResponse
// @Failure 500 {object} presenter.Response
// @Router /vehicles/owner [get]
func GetVehiclesByOwnerID(vehicleService *service.VehicleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ownerIDStr := c.Query("owner_id")
		ownerID, err := uuid.Parse(ownerIDStr)
		if err != nil {
			return presenter.BadRequest(c, err)
		}
		page := c.QueryInt("page", 1)
		pageSize := c.QueryInt("page_size", 10)

		vehicles, total, err := vehicleService.GetVehiclesByOwnerID(c.UserContext(), ownerID, page, pageSize)
		if err != nil {
			return presenter.InternalServerError(c, err)
		}

		res := make([]presenter.FullVehicleResponse, len(vehicles))
		for i, vehicle := range vehicles {
			res[i] = presenter.VehicleToFullVehicleResponse(vehicle)
		}

		pagination := presenter.NewPagination(res, uint(page), uint(pageSize), uint(total))
		return presenter.OK(c, "Vehicles retrieved successfully", pagination)
	}
}

// UpdateVehicle updates a vehicle by ID
// @Summary Update a vehicle by ID
// @Description Update a vehicle by ID
// @Tags vehicles
// @Accept json
// @Produce json
// @Param id path uint true "Vehicle ID"
// @Param vehicle body presenter.UpdateVehicleReq true "Vehicle to update"
// @Success 200 {object} presenter.FullVehicleResponse
// @Failure 400 {object} presenter.Response "error: bad request"
// @Failure 500 {object} presenter.Response "error: internal server error"
// @Router /vehicles/{id} [put]
func UpdateVehicle(vehicleService *service.VehicleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		vehicleID, err := c.ParamsInt("id")
		if err != nil {
			return presenter.BadRequest(c, err)
		}

		var updateReq presenter.UpdateVehicleReq
		if err := c.BodyParser(&updateReq); err != nil {
			return presenter.BadRequest(c, err)
		}

		if err := vehicleService.UpdateVehicle(c.UserContext(), uint(vehicleID), presenter.UpdateVehicleRequestToDomain(&updateReq)); err != nil {
			return presenter.InternalServerError(c, err)
		}

		return presenter.OK(c, "Vehicle updated successfully", nil)
	}
}

// DeleteVehicle deletes a vehicle by ID
// @Summary Delete a vehicle by ID
// @Description Delete a vehicle by ID
// @Tags vehicles
// @Produce json
// @Param id path uint true "Vehicle ID"
// @Success 204 "No Content"
// @Failure 400 {object} presenter.Response "error: bad request"
// @Failure 500 {object} presenter.Response "error: internal server error"
// @Router /vehicles/{id} [delete]
func DeleteVehicle(vehicleService *service.VehicleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return presenter.BadRequest(c, err)
		}
		if err := vehicleService.DeleteVehicle(c.UserContext(), uint(id)); err != nil {
			return presenter.InternalServerError(c, err)
		}
		return presenter.NoContent(c)
	}
}
