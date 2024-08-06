package storage

import (
	"context"
	"fmt"
	"strings"
	"time"
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

func (r *tripRepo) GetCompanyTrips(ctx context.Context, originCity, destinationCity, pathType string, startDate *time.Time, requesterType string, companyID uint, limit, offset uint) ([]trip.Trip, uint, error) {
	query := r.db.WithContext(ctx).
		Model(&entities.Trip{}).
		Where("transport_company_id = ?", companyID).Preload("TransportCompany"). // Preload related TransportCompany
		Preload("TripCancelingPenalty").                                          // Preload related TripCancelingPenalty
		Preload("VehicleRequest").
		Preload("TechTeam").
		Preload("TechTeam.Members").
		Order("created_at DESC")

	if originCity != "" {
		query = query.Where("origin = ?", originCity)
	}

	if destinationCity != "" {
		query = query.Where("destination = ?", destinationCity)
	}

	if pathType != "" {
		query = query.Where("trip_type = ?", pathType)
	}
	if requesterType == "agency" {
		query = query.Where("tour_release_date <= ?", time.Now()).
			Where("transport_companies.is_blocked = ?", false).
			Where("is_canceled = ?", false).
			Where("is_finished = ?", false)
	} else if requesterType == "user" {
		query = query.Where("user_release_date <= ?", time.Now()).
			Where("transport_companies.is_blocked = ?", false).
			Where("is_canceled = ?", false).
			Where("is_finished = ?", false)
	}

	if startDate != nil {
		startDateStr := startDate.Format("2006-01-02")
		if startDateStr != "0001-01-01" {
			query = query.Where("DATE(start_date) = ?", startDateStr)
		}
	}

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

	var trips []entities.Trip

	if err := query.Find(&trips).Error; err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, 0, trip.ErrRecordNotFound
		}
		return nil, 0, err
	}

	return mappers.TripEntitiesToDomain(trips), uint(total), nil
}

func (r *tripRepo) GetFullTripByID(ctx context.Context, id uint) (*trip.Trip, error) {
	var t entities.Trip
	if err := r.db.WithContext(ctx).
		Preload("TransportCompany").     // Preload related TransportCompany
		Preload("TripCancelingPenalty"). // Preload related TripCancelingPenalty
		Preload("VehicleRequest").
		Preload("TechTeam").
		Preload("TechTeam.Members").
		First(&t, id).Error; err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, nil
		}
		return nil, fmt.Errorf("%w %w", trip.ErrFailedToGetTrip, err)
	}
	dPath := mappers.TripFullEntityToDomain(t)
	return &dPath, nil
}

func (r *tripRepo) GetTrips(ctx context.Context, originCity, destinationCity, pathType string, startDate *time.Time, requesterType string, limit, offset uint) ([]trip.Trip, uint, error) {
	// Start the query for trips
	query := r.db.WithContext(ctx).
		Model(&entities.Trip{}).
		Joins("JOIN transport_companies ON transport_companies.id = trips.transport_company_id").
		Preload("TransportCompany").
		Preload("TripCancelingPenalty").Preload("VehicleRequest").
		Preload("TechTeam").Preload("TechTeam.Members").
		Where("start_date > ?", time.Now()).
		Where("sold_tickets < max_tickets")
		//Where("vehicle_id IS NOT NULL").

	if originCity != "" {
		query = query.Where("origin = ?", originCity)
	}

	if destinationCity != "" {
		query = query.Where("destination = ?", destinationCity)
	}

	if pathType != "" {
		query = query.Where("trip_type = ?", pathType)
	}

	if requesterType == "agency" {
		query = query.Where("tour_release_date <= ?", time.Now()).
			Where("transport_companies.is_blocked = ?", false).
			Where("is_canceled = ?", false).
			Where("is_finished = ?", false)
		//Where("trip.is_confirmed = ?", isConfirmed)
	} else if requesterType == "user" {
		query = query.Where("user_release_date <= ?", time.Now()).
			Where("transport_companies.is_blocked = ?", false).
			Where("is_canceled = ?", false).
			Where("is_finished = ?", false)
	}

	if startDate != nil {
		startDateStr := startDate.Format("2006-01-02")
		if startDateStr != "0001-01-01" {
			query = query.Where("DATE(start_date) = ?", startDateStr)
		}
	}

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

	var trips []entities.Trip
	if err := query.Find(&trips).Error; err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, 0, trip.ErrRecordsNotFound
		}
		return nil, 0, err
	}

	return mappers.TripEntitiesToDomain(trips), uint(total), nil
}

func (r *tripRepo) UpdateTrip(ctx context.Context, id uint, updates map[string]interface{}) error {
	// Ensure updates is not empty
	if len(updates) == 0 {
		return nil
	}

	if err := r.db.WithContext(ctx).
		Model(&entities.Trip{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update trip: %w", err)
	}

	return nil
}

func (r *tripRepo) GetCountPathUnfinishedTrips(ctx context.Context, pathID uint) (uint, error) {
	var count int64

	err := r.db.WithContext(ctx).Model(&entities.Trip{}).
		Where("is_finished = ? and is_canceled=?", false, false).
		Where("path_id = ?", pathID).
		Count(&count).Error

	if err != nil {
		return 100, err
	}
	return uint(count), nil
}

func (r *tripRepo) GetUpcomingUnconfirmedTripIDsToCancel(ctx context.Context) ([]uint, error) {
	now := time.Now()
	in24Hours := now.Add(24 * time.Hour)

	var tripIDs []uint

	err := r.db.WithContext(ctx).Model(&entities.Trip{}).
		Select("id").
		Where("start_date BETWEEN ? AND ?", now, in24Hours).
		Where("is_confirmed = ? and is_canceled = ? and is_finished = ?", false, false, false).
		Pluck("id", &tripIDs).
		Error

	if err != nil {
		return nil, err
	}

	return tripIDs, nil
}

func (r *tripRepo) CheckAvailabilityTechTeam(ctx context.Context, tripID uint, techTeamID uint, startDate time.Time, endDate time.Time) (bool, error) {
	var conflictingTrips []entities.Trip

	err := r.db.WithContext(ctx).
		Where("tech_team_id = ?", techTeamID).
		Where("(start_date <= ? AND end_date >= ?)", endDate, startDate).
		Where("id <> ?", tripID).
		Find(&conflictingTrips).Error

	if err != nil {
		return false, fmt.Errorf("failed to check for conflicting trips: %w", err)
	}

	if len(conflictingTrips) > 0 {
		return false, nil
	}
	return true, nil
}

func (r *tripRepo) GetCountCompanyUnfinishedUncanceledTrips(ctx context.Context, companyID uint) (uint, error) {
	var count int64

	err := r.db.WithContext(ctx).Model(&entities.Trip{}).
		Where("is_finished = ? and is_canceled=?", false, false).
		Where("transport_company_id = ?", companyID).
		Count(&count).Error

	if err != nil {
		return 100, err
	}
	return uint(count), nil
}
