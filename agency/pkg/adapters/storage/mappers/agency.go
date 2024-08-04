package mappers

import (
	"agency/internal/agency"
	"agency/pkg/adapters/storage/entities"
	"agency/pkg/fp"
)

func AgencyEntityToDomain(agencyEntity entities.Agency) agency.Agency {
	return agency.Agency{
		ID:        agencyEntity.ID,
		OwnerID:   agencyEntity.OwnerID,
		Name:      agencyEntity.Name,
		IsBlocked: agencyEntity.IsBlocked,
	}
}

func BatchAgencyEntitiesToDomain(agencyEntities []entities.Agency) []agency.Agency {
	return fp.Map(agencyEntities, AgencyEntityToDomain)
}

func AgencyDomainToEntity(a *agency.Agency) *entities.Agency {
	return &entities.Agency{
		Name:      a.Name,
		OwnerID:   a.OwnerID,
		IsBlocked: a.IsBlocked,
	}
}