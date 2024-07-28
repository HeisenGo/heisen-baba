package handlers


import (
	"hotel/api/http/handlers/presenter"
	"hotel/service"
	"github.com/gofiber/fiber/v2"
	"strconv"
)
// CreateRoom creates a new room
// @Summary Create a new room
// @Description Create a new room
// @Tags rooms
// @Accept json
// @Produce json
// @Param room body presenter.CreateRoomReq true "Room to create"
// @Success 201 {object} presenter.RoomResp
// @Failure 400 {object} map[string]interface{} "error: bad request"
// @Failure 500 {object} map[string]interface{} "error: internal server error"
// @Router /rooms [post]
func CreateRoom(roomService *service.RoomService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req presenter.CreateRoomReq
		if err := c.BodyParser(&req); err != nil {
			return presenter.BadRequest(c, err)
		}

		r := presenter.CreateRoomRequest(&req)
		createdRoom, err := roomService.CreateRoom(c.UserContext(), r)
		if err != nil {
			return presenter.InternalServerError(c, err)
		}

		res := presenter.RoomToCreateRoomResponse(createdRoom)
		return presenter.Created(c, "Room created successfully", res)
	}
}
// GetRoom gets a room by ID
// @Summary Get a room by ID
// @Description Get a room by ID
// @Tags rooms
// @Produce json
// @Param id path int true "Room ID"
// @Success 200 {object} presenter.RoomResp
// @Failure 400 {object} map[string]interface{} "error: bad request"
// @Failure 500 {object} map[string]interface{} "error: internal server error"
// @Router /rooms/{id} [get]
func GetRoom(roomService *service.RoomService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roomID, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil {
			return presenter.BadRequest(c, err)
		}

		r, err := roomService.GetRoom(c.UserContext(), uint(roomID))
		if err != nil {
			return presenter.NotFound(c, err)
		}

		res := presenter.RoomToFullRoomResponse(r)
		return presenter.OK(c, "Room Fetched Successfully", res)
	}
}
// UpdateRoom updates a room by ID
// @Summary Update a room by ID
// @Description Update a room by ID
// @Tags rooms
// @Accept json
// @Produce json
// @Param id path int true "Room ID"
// @Param room body presenter.CreateRoomReq true "Room to update"
// @Success 200 {object} presenter.RoomResp
// @Failure 400 {object} map[string]interface{} "error: bad request"
// @Failure 404 {object} map[string]interface{} "error: bad request"
// @Failure 500 {object} map[string]interface{} "error: internal server error"
// @Router /rooms/{id} [put]
func UpdateRoom(roomService *service.RoomService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roomID, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil {
			return presenter.BadRequest(c, err)
		}

		var req presenter.CreateRoomReq
		if err := c.BodyParser(&req); err != nil {
			return presenter.BadRequest(c, err)
		}

		r := presenter.CreateRoomRequest(&req)
		r.ID = uint(roomID)
		updatedRoom, err := roomService.UpdateRoom(c.UserContext(), r)
		if err != nil {
			return presenter.InternalServerError(c, err)
		}

		res := presenter.RoomToCreateRoomResponse(updatedRoom)
		return presenter.OK(c, "Room updated successfully", res)
	}
}
// DeleteRoom deletes a room by ID
// @Summary Delete a room by ID
// @Description Delete a room by ID
// @Tags rooms
// @Produce json
// @Param id path int true "Room ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]interface{} "error: bad request"
// @Failure 500 {object} map[string]interface{} "error: internal server error"
// @Router /rooms/{id} [delete]
func DeleteRoom(roomService *service.RoomService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roomID, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil {
			return presenter.BadRequest(c, err)
		}

		if err := roomService.DeleteRoom(c.UserContext(), uint(roomID)); err != nil {
			return presenter.InternalServerError(c, err)
		}

		return presenter.NoContent(c)
	}
}
