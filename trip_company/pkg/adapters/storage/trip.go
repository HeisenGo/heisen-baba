package storage

import (
	"context"
	"fmt"
	"strings"
	"tripcompanyservice/internal/trip"
	"tripcompanyservice/pkg/adapters/storage/entities"
	"tripcompanyservice/pkg/adapters/storage/mappers"

	"gorm.io/gorm"
)

type tripRepo struct {
	db *gorm.DB
}

func NewTripRepo(db *gorm.DB) trip.Repo {
	return &tripRepo{db}
}

func (r *tripRepo) Insert(ctx context.Context, t *trip.Trip) error {
	tripEntity := mappers.TripDomainToEntity(t)

	result := r.db.WithContext(ctx).Create(&tripEntity)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {

			var existingTrip entities.Trip
			// Search for the soft-deleted record with the same unique constraints
			if r.db.WithContext(ctx).Unscoped().Where("path_id = ? AND transport_company_id= ? AND start_date = ?", t.PathID, t.TransportCompanyID, t.StartDate).First(&existingTrip).Error == nil {
				// Check if the record is soft-deleted
				if existingTrip.DeletedAt.Valid {
					// Restore the soft-deleted record
					existingTrip.DeletedAt = gorm.DeletedAt{}
					if err := r.db.WithContext(ctx).Save(&existingTrip).Error; err != nil {
						return fmt.Errorf("%w %w", trip.ErrFailedToRestore, err)
					}
					t.ID = existingTrip.ID
					return nil
				}
			}

			return trip.ErrDuplication
		}
		return result.Error

	}
	t.ID = tripEntity.ID

	return nil
}

func (r *tripRepo) GetCompanyTrips(ctx context.Context, companyID uint, limit, offset uint) ([]trip.Trip, uint, error) {
	return nil, 0, nil
}
