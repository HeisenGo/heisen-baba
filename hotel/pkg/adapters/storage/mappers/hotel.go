package mappers

import (
	"hotel/internal/hotel"
	"hotel/internal/room"
	"hotel/pkg/adapters/storage/entities"
	"hotel/pkg/fp"
)

func HotelEntityToDomain(hotelEntity entities.Hotel) hotel.Hotel {
	domainRooms := BatchRoomEntitiesToDomain(hotelEntity.Rooms)
	return hotel.Hotel{
		ID:        hotelEntity.ID,
		Name:      hotelEntity.Name,
		City:      hotelEntity.City,
		Country:   hotelEntity.Country,
		Details:   hotelEntity.Details,
		IsBlocked: hotelEntity.IsBlocked,
		Rooms:     domainRooms,
	}
}
func BatchRoomEntitiesToDomain(roomEntities []entities.Room) []room.Room {
	return fp.Map(roomEntities, RoomEntityToDomain)
}
func BatchRoomDomainToEntity(roomDomain []room.Room) []entities.Room {
	return fp.Map(roomDomain, RoomDomainToEntity)
}
func HotelDomainToEntity(h *hotel.Hotel) *entities.Hotel {
	entitiyRooms := BatchRoomDomainToEntity(h.Rooms)
	return &entities.Hotel{
		Name:      h.Name,
		City:      h.City,
		Country:   h.Country,
		Details:   h.Details,
		OwnerID:   h.OwnerID,
		IsBlocked: h.IsBlocked,
		Rooms:     entitiyRooms,
	}
}
