package handlers

import (
	"hotel/api/http/handlers/presenter"
	"hotel/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// CreateRoom creates a new room
// @Summary Create a new room
// @Description Create a new room
// @Tags rooms
// @Accept json
// @Produce json
// @Param room body presenter.CreateRoomReq true "Room to create"
// @Success 201 {object} presenter.RoomResp
// @Failure 400 {object} presenter.Response "error: bad request"
// @Failure 500 {object} presenter.Response "error: internal server error"
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

// GetRooms gets a paginated list of rooms
// @Summary Get rooms
// @Description Get paginated list of rooms
// @Tags rooms
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} presenter.RoomResp
// @Failure 500 {object} presenter.Response
// @Router /rooms [get]
func GetRooms(roomService *service.RoomService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page := c.QueryInt("page", 1)
		pageSize := c.QueryInt("page_size", 10)

		rooms, total, err := roomService.GetRooms(c.UserContext(), page, pageSize)
		if err != nil {
			return presenter.InternalServerError(c, err)
		}

		res := make([]presenter.RoomResp, len(rooms))
		for i, room := range rooms {
			res[i] = *presenter.RoomToFullRoomResponse(&room) // Dereference the pointer here
		}

		pagination := presenter.NewPagination(res, uint(page), uint(pageSize), uint(total))
		return presenter.OK(c, "Rooms retrieved successfully", pagination)
	}
}
// GetRoomsByID gets a room by its ID
// @Summary Get a room by ID
// @Description Get a room by ID
// @Tags rooms
// @Produce json
// @Param id path int true "room ID"
// @Success 200 {object} presenter.roomResp
// @Failure 400 {object} presenter.Response "error: bad request"
// @Failure 500 {object} presenter.Response "error: internal server error"
// @Router /rooms/{id} [get]
func GetRoomsByID(roomService *service.RoomService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return presenter.BadRequest(c, err)
		}

		room, err := roomService.GetRoomByID(c.UserContext(), uint(id))
		if err != nil {
			return presenter.InternalServerError(c, err)
		}

		resp := presenter.RoomToFullRoomResponse(room)
		return presenter.OK(c, "room retrieved successfully", resp)
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
// @Failure 400 {object} presenter.Response "error: bad request"
// @Failure 404 {object} presenter.Response "error: bad request"
// @Failure 500 {object} presenter.Response "error: internal server error"
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
		err = roomService.UpdateRoom(c.UserContext(), r)
		if err != nil {
			return presenter.InternalServerError(c, err)
		}

		return presenter.OK(c, "Room updated successfully", nil)
	}
}

// DeleteRoom deletes a room by ID
// @Summary Delete a room by ID
// @Description Delete a room by ID
// @Tags rooms
// @Produce json
// @Param id path int true "Room ID"
// @Success 204 "No Content"
// @Failure 404 {object} presenter.Response "error: bad request"
// @Failure 500 {object} presenter.Response "error: internal server error"
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
