package storage

import (
	"context"
	"errors"
	"agency/internal/tour"
	"agency/pkg/adapters/storage/entities"
	"agency/pkg/adapters/storage/mappers"

	"gorm.io/gorm"
)

type tourRepo struct {
	db *gorm.DB
}

func NewTourRepo(db *gorm.DB) tour.Repo {
	return &tourRepo{
		db: db,
	}
}

func (r *tourRepo) CreateTour(ctx context.Context, t *tour.Tour) error {
	tourEntity := mappers.TourDomainToEntity(t)
	if err := r.db.WithContext(ctx).Create(&tourEntity).Error; err != nil {
		return err
	}
	t.ID = tourEntity.ID
	return nil
}

func (r *tourRepo) GetTours(ctx context.Context, agencyID uint, page, pageSize int) ([]tour.Tour, uint, error) {
	var t []entities.Tour
	var int64Total int64

	query := r.db.Model(&entities.Tour{}).Where("agency_id = ?", agencyID)

	// Count total records for pagination
	query.Count(&int64Total)

	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	if err := query.Find(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, nil
		}
		return nil, 0, err
	}

	total := uint(int64Total)
	tours := mappers.BatchTourEntitiesToDomain(t)
	return tours, total, nil
}

func (r *tourRepo) GetToursByAgencyID(ctx context.Context, agencyID uint, page, pageSize int) ([]tour.Tour, int, error) {
	var tourEntities []entities.Tour
	var total int64

	query := r.db.WithContext(ctx).Model(&entities.Tour{}).Where("agency_id = ?", agencyID)

	query.Count(&total)

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&tourEntities).Error; err != nil {
		return nil, 0, err
	}

	tours := make([]tour.Tour, len(tourEntities))
	for i, tourEntity := range tourEntities {
		tours[i] = mappers.TourEntityToDomain(tourEntity)
	}

	return tours, int(total), nil
}

func (r *tourRepo) GetTourByID(ctx context.Context, id uint) (*tour.Tour, error) {
	var tourEntity entities.Tour
	if err := r.db.First(&tourEntity, id).Error; err != nil {
		return nil, err
	}
	return mappers.TourEntityToDomain(tourEntity), nil
}

func (r *tourRepo) UpdateTour(ctx context.Context, t *tour.Tour) error {
	tourEntity := mappers.TourDomainToEntity(t)
	if err := r.db.WithContext(ctx).Model(&entities.Tour{}).Where("id = ?", t.ID).Updates(tourEntity).Error; err != nil {
		return err
	}
	return nil
}

func (r *tourRepo) DeleteTour(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&entities.Tour{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tour.ErrRecordNotFound
		}
		return err
	}
	return nil
}

func (r *tourRepo) ApproveTour(ctx context.Context, tourID uint) error {
	if err := r.db.WithContext(ctx).Model(&entities.Tour{}).Where("id = ?", tourID).Update("is_approved", true).Error; err != nil {
		return err
	}
	return nil
}

func (r *tourRepo) SetTourStatus(ctx context.Context, tourID uint, isActive bool) error {
	if err := r.db.WithContext(ctx).Model(&entities.Tour{}).Where("id = ?", tourID).Update("is_active", isActive).Error; err != nil {
		return err
	}
	return nil
}
