package mappers

import (
	"hotel/internal/hotel"
	"hotel/pkg/adapters/storage/entities"
	"hotel/pkg/fp"
)

func HotelEntityToDomain(hotelEntity entities.Hotel) hotel.Hotel {
	domainRooms := BatchRoomEntitiesToDomain(hotelEntity.Rooms)
	return hotel.Hotel{
		ID:        hotelEntity.ID,
		OwnerID:   hotelEntity.OwnerID,
		Name:      hotelEntity.Name,
		City:      hotelEntity.City,
		Country:   hotelEntity.Country,
		Details:   hotelEntity.Details,
		IsBlocked: hotelEntity.IsBlocked,
		Rooms:     domainRooms,
	}
}
func BatchHotelEntitiesToDomain(hotelEntities []entities.Hotel) []hotel.Hotel {
	return fp.Map(hotelEntities, HotelEntityToDomain)
}

func HotelDomainToEntity(h *hotel.Hotel) *entities.Hotel {
	return &entities.Hotel{
		Name:      h.Name,
		City:      h.City,
		Country:   h.Country,
		Details:   h.Details,
		OwnerID:   h.OwnerID,
		IsBlocked: h.IsBlocked,
	}
}
