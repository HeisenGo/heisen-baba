package mappers

import (
	"agency/internal/agency"
	"agency/pkg/adapters/storage/entities"
	"agency/pkg/fp"

	"gorm.io/gorm"
)

// Convert a single Agency entity to domain model
func AgencyEntityToDomain(agencyEntity entities.Agency) agency.Agency {
	// Batch convert Tours if necessary
	domainTours := BatchTourEntitiesToDomain(agencyEntity.Tours)
	return agency.Agency{
		ID:        agencyEntity.ID,
		OwnerID:   agencyEntity.OwnerID,
		Name:      agencyEntity.Name,
		IsBlocked: agencyEntity.IsBlocked,
		Tours:     domainTours,
	}
}

// Convert a batch of Agency entities to domain models
func BatchAgencyEntitiesToDomain(agencyEntities []entities.Agency) []agency.Agency {
	return fp.Map(agencyEntities, AgencyEntityToDomain)
}

// Convert a domain model Agency to entity
func AgencyDomainToEntity(a *agency.Agency) *entities.Agency {
	return &entities.Agency{
		Model:     gorm.Model{ID: a.ID}, // Set the ID for existing records
		Name:      a.Name,
		OwnerID:   a.OwnerID,
		IsBlocked: a.IsBlocked,
	}
}
