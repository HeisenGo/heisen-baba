package storage

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"tripcompanyservice/internal/company"
	"tripcompanyservice/pkg/adapters/storage/entities"
	"tripcompanyservice/pkg/adapters/storage/mappers"

	"gorm.io/gorm"
)

type comapnyRepo struct {
	db *gorm.DB
}

func NewTransportCompanyRepo(db *gorm.DB) company.Repo {
	return &comapnyRepo{db}
}

func (r *comapnyRepo) GetTransportCompanies(ctx context.Context, limit, offset uint) ([]company.TransportCompany, uint, error) {
	query := r.db.WithContext(ctx).Model(&entities.TransportCompany{})

	var total int64

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if offset > 0 {
		query = query.Offset(int(offset))
	}

	if limit > 0 {
		query = query.Limit(int(limit))
	}

	var companies []entities.TransportCompany

	if err := query.Find(&companies).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, nil
		}
		return nil, 0, err
	}

	return mappers.CompanyEntitiesToDomain(companies), uint(total), nil
}

func (r *comapnyRepo) Insert(ctx context.Context, c *company.TransportCompany) error {
	companyEntity := mappers.CompanyDomainToEntity(c)
	result := r.db.WithContext(ctx).Save(&companyEntity)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {

			var existingCompany entities.TransportCompany
			// Search for the soft-deleted record with the same unique constraints
			if r.db.WithContext(ctx).Unscoped().Where("name = ? AND owner_id = ? AND email = ?", c.Name, c.OwnerID, c.Email).First(&existingCompany).Error == nil {
				// Check if the record is soft-deleted
				if existingCompany.DeletedAt.Valid {
					// Restore the soft-deleted record
					existingCompany.DeletedAt = gorm.DeletedAt{}
					if err := r.db.WithContext(ctx).Save(existingCompany).Error; err != nil {
						return fmt.Errorf("%w %w", company.ErrFailedToRestore, err)
					}
					c.ID = existingCompany.ID
					return nil
				}

				return company.ErrDuplication
			}
		}
		return result.Error
	}

	c.ID = companyEntity.ID

	return nil
}

func (r *comapnyRepo) GetUserTransportCompanies(ctx context.Context, ownerID uint, limit, offset uint) ([]company.TransportCompany, uint, error) {
	query := r.db.WithContext(ctx).Model(&entities.TransportCompany{}).Where("owner_id = ?", ownerID)

	var total int64

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if offset > 0 {
		query = query.Offset(int(offset))
	}

	if limit > 0 {
		query = query.Limit(int(limit))
	}

	var companies []entities.TransportCompany

	if err := query.Find(&companies).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, nil
		}
		return nil, 0, err
	}

	return mappers.CompanyEntitiesToDomain(companies), uint(total), nil
}
