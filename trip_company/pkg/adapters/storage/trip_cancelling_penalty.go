package storage

import (
	"context"
	"strings"
	tripcancellingpenalty "tripcompanyservice/internal/trip_cancelling_penalty"
	"tripcompanyservice/pkg/adapters/storage/entities"
	"tripcompanyservice/pkg/adapters/storage/mappers"

	"gorm.io/gorm"
)

type tripCancellingPenaltyRepo struct {
	db *gorm.DB
}

func NewTripCancellingPenaltyRepo(db *gorm.DB) tripcancellingpenalty.Repo {
	return &tripCancellingPenaltyRepo{db}
}

func (r *tripCancellingPenaltyRepo) GetByID(ctx context.Context, id uint) (*tripcancellingpenalty.TripCancelingPenalty, error) {
	var t entities.TripCancellingPenalty

	err := r.db.WithContext(ctx).Model(&entities.TripCancellingPenalty{}).Where("id = ?", id).First(&t).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, nil
		}
		return nil, err
	}
	dC := mappers.PenaltyEntityToDomain(t)
	return &dC, nil
}
