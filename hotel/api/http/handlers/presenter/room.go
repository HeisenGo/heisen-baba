package presenter

import (
	"hotel/pkg/adapters/storage/entities"

)

type CreateRoomRequest struct {
	HotelID uint               `json:"hotel_id"`
	Rooms []CreateRoomItem `json:"rooms"`
}

type CreateRoomsItem struct {
	Name string `json:"name"`
}

type CreateRoomsResponse struct {
	Data    []RoomResponseItem `json:"data"`
	Message string               `json:"message"`
}

func CreateRoomsRequestToEntities(req CreateRoomRequest, maxOrder uint) []entities.Room {
	rooms := make([]entities.Room, len(req.Rooms))
	for i, room := range req.Rooms {
		rooms[i] = entities.Room{
			Name:     room.Name,
			HotelID:  req.HotelID,
		}
	}
	return rooms
}

func EntitiesToCreateColumnsResponse(columns []entities.Column) CreateColumnsResponse {
	respItems := make([]ColumnResponseItem, len(columns))
	for i, col := range columns {
		respItems[i] = ColumnResponseItem{
			ID:    col.ID,
			Name:  col.Name,
			Order: col.OrderNum,
		}
	}
	return CreateColumnsResponse{
		Data:    respItems,
		Message: "Columns successfully created.",
	}
}

func EntityToColumnResponse(c column.Column) ColumnResponseItem {
	return ColumnResponseItem{
		ID:    c.ID,
		Name:  c.Name,
		Order: c.OrderNum,
	}
}
