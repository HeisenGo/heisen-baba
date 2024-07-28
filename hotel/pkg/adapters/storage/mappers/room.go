package mappers

import (
	"hotel/internal/room"
	"hotel/pkg/adapters/storage/entities"

	"gorm.io/gorm"
)

func RoomEntityToDomain(roomEntity entities.Room) room.Room {
	return room.Room{
		ID:          roomEntity.ID,
		HotelID:     roomEntity.HotelID,
		Name:        roomEntity.Name,
		AgencyPrice: roomEntity.AgencyPrice,
		UserPrice:   roomEntity.UserPrice,
		Facilities:  roomEntity.Facilities,
		Capacity:    roomEntity.Capacity,
		IsAvailable: roomEntity.IsAvailable,
	}
}

func RoomDomainToEntity(r *room.Room) *entities.Room {
	return &entities.Room{
		Model: gorm.Model{
			ID: r.ID,
		},
		Name:        r.Name,
		HotelID:     r.HotelID,
		AgencyPrice: r.AgencyPrice,
		UserPrice:   r.UserPrice,
		Facilities:  r.Facilities,
		Capacity:    r.Capacity,
		IsAvailable: r.IsAvailable,
	}
}
