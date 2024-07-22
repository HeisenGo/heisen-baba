package mappers

import (
	"hotel/internal/hotel"
	"hotel/pkg/adapters/storage/entities"
)

func HotelEntityToDomain(hotelEntity entities.Hotel) hotel.Hotel {
	return hotel.Hotel{
		ID:        hotelEntity.ID,
		Name:      hotelEntity.Name,
		City:      hotelEntity.City,
		Country:   hotelEntity.Country,
		Details:   hotelEntity.Details,
		IsBlocked: hotelEntity.IsBlocked,
	}
}

func HotelDomainToEntity(h *hotel.Hotel) *entities.Hotel {
	return &entities.Hotel{
		Name:      h.Name,
		City:      h.City,
		Country:   h.Country,
		Details:   h.Details,
		IsBlocked: h.IsBlocked,
	}
}
