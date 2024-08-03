package mappers

import (
	"tripcompanyservice/internal/company"
	"tripcompanyservice/pkg/adapters/storage/entities"
	"tripcompanyservice/pkg/fp"
)

func CompanyEntityToDomain(companyEntity entities.TransportCompany) company.TransportCompany {
	return company.TransportCompany{
		ID:          companyEntity.ID,
		Name:        companyEntity.Name,
		Description: companyEntity.Description,
		OwnerID:     companyEntity.OwnerID,
		Address:     companyEntity.Address,
		IsBlocked:   companyEntity.IsBlocked,
		//PhoneNumber: companyEntity.PhoneNumber,
		//Email:       companyEntity.Email,
	}
}

func CompanyEntitiesToDomain(companyEntities []entities.TransportCompany) []company.TransportCompany {
	return fp.Map(companyEntities, CompanyEntityToDomain)
}

func CompanyDomainToEntity(c *company.TransportCompany) *entities.TransportCompany {
	return &entities.TransportCompany{
		Name:        c.Name,
		Description: c.Description,
		OwnerID:     c.OwnerID,
		Address:     c.Address,
		IsBlocked: c.IsBlocked,
		//PhoneNumber: c.PhoneNumber,
		//Email:       c.Email,
	}
}