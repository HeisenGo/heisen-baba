package mappers

import (
	"agency/internal/tour"
	"agency/pkg/adapters/storage/entities"
	"agency/pkg/fp"
)

func TourEntityToDomain(tourEntity entities.Tour) tour.Tour {
	return tour.Tour{
		ID:          tourEntity.ID,
		AgencyID:    tourEntity.AgencyID,
		GoTicketID:  tourEntity.GoTicketID,
		BackTicketID: tourEntity.BackTicketID,
		HotelID:     tourEntity.HotelID,
		Capacity:    tourEntity.Capacity,
		ReleaseDate: tourEntity.ReleaseDate,
	}
}

func BatchTourEntitiesToDomain(tourEntities []entities.Tour) []tour.Tour {
	return fp.Map(tourEntities, TourEntityToDomain)
}

func TourDomainToEntity(t *tour.Tour) *entities.Tour {
	return &entities.Tour{
		AgencyID:    t.AgencyID,
		GoTicketID:  t.GoTicketID,
		BackTicketID: t.BackTicketID,
		HotelID:     t.HotelID,
		Capacity:    t.Capacity,
	}
}