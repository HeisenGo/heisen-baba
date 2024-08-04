package storage

import (
	"agency/internal/agency"
	"agency/pkg/adapters/storage/entities"
	"agency/pkg/adapters/storage/mappers"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type agencyRepo struct {
	db *gorm.DB
}

func NewAgencyRepo(db *gorm.DB) agency.Repo {
	return &agencyRepo{
		db: db,
	}
}

func (r *agencyRepo) CreateAgency(ctx context.Context, a *agency.Agency) error {
	agencyEntity := mappers.AgencyDomainToEntity(a)
	if err := r.db.WithContext(ctx).Create(&agencyEntity).Error; err != nil {
		return err
	}
	a.ID = agencyEntity.ID
	return nil
}

func (r *agencyRepo) GetAgencies(ctx context.Context, name string, page, pageSize int) ([]agency.Agency, uint, error) {
	var a []entities.Agency
	var int64Total int64

	query := r.db.Model(&entities.Agency{}).Preload("Tours")

	// Filters
	if name != "" {
		query = query.Where("name = ?", name)
	}
	// Filter out blocked agencies
	query = query.Where("is_blocked = ?", false)
	// Count total records for pagination
	query.Count(&int64Total)

	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	if err := query.Find(&a).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, nil
		}
		return nil, 0, err
	}

	total := uint(int64Total)
	agencies := mappers.BatchAgencyEntitiesToDomain(a)
	return agencies, total, nil
}

func (r *agencyRepo) GetAgenciesByOwnerID(ctx context.Context, ownerID uuid.UUID, page, pageSize int) ([]agency.Agency, int, error) {
	var agencyEntities []entities.Agency
	var total int64

	query := r.db.WithContext(ctx).Model(&entities.Agency{}).Where("owner_id = ?", ownerID).Preload("Tours")

	query.Count(&total)

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&agencyEntities).Error; err != nil {
		return nil, 0, err
	}

	agencies := make([]agency.Agency, len(agencyEntities))
	for i, agencyEntity := range agencyEntities {
		agencies[i] = mappers.AgencyEntityToDomain(agencyEntity)
	}

	return agencies, int(total), nil
}

func (r *agencyRepo) GetAgencyByID(ctx context.Context, id uint) (*agency.Agency, error) {
	var agencyEntity entities.Agency
	if err := r.db.Preload("Tours").First(&agencyEntity, id).Error; err != nil {
		return nil, err
	}
	ag := mappers.AgencyEntityToDomain(agencyEntity)
	return &ag, nil
}

func (r *agencyRepo) UpdateAgency(ctx context.Context, a *agency.Agency) error {
	agencyEntity := mappers.AgencyDomainToEntity(a)
	if err := r.db.WithContext(ctx).Model(&entities.Agency{}).Where("id = ?", a.ID).Updates(agencyEntity).Error; err != nil {
		return err
	}
	return nil
}

func (r *agencyRepo) DeleteAgency(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&entities.Agency{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return agency.ErrRecordNotFound
		}
		return err
	}
	return nil
}

func (r *agencyRepo) BlockAgency(ctx context.Context, agencyID uint) error {
	if err := r.db.WithContext(ctx).Model(&entities.Agency{}).Where("id = ?", agencyID).Update("is_blocked", true).Error; err != nil {
		return err
	}
	return nil
}
