package mappers

import (
	"agency/internal/tour"
	"agency/pkg/adapters/storage/entities"
	"agency/pkg/fp"

	"gorm.io/gorm"
)

// Convert a single Tour entity to domain model
func TourEntityToDomain(tourEntity entities.Tour) tour.Tour {
	return tour.Tour{
		ID:           tourEntity.ID,
		AgencyID:     tourEntity.AgencyID,
		GoTicketID:   tourEntity.GoTicketID,
		BackTicketID: tourEntity.BackTicketID,
		HotelID:      tourEntity.HotelID,
		Capacity:     tourEntity.Capacity,
		ReleaseDate:  tourEntity.ReleaseDate,
		IsActive:     tourEntity.IsActive,
		UserPrice:    tourEntity.UserPrice,
		IsApproved:   tourEntity.IsApproved,
	}
}

// Convert a batch of Tour entities to domain models
func BatchTourEntitiesToDomain(tourEntities []entities.Tour) []tour.Tour {
	return fp.Map(tourEntities, TourEntityToDomain)
}

// Convert a domain model Tour to entity
func TourDomainToEntity(t *tour.Tour) *entities.Tour {
	return &entities.Tour{
		Model:        gorm.Model{ID: t.ID}, // Set the ID for existing records
		AgencyID:     t.AgencyID,
		GoTicketID:   t.GoTicketID,
		BackTicketID: t.BackTicketID,
		HotelID:      t.HotelID,
		Capacity:     t.Capacity,
		ReleaseDate:  t.ReleaseDate,
		IsActive:     t.IsActive,
		IsApproved:   t.IsApproved,
		UserPrice:    t.UserPrice,
	}
}
