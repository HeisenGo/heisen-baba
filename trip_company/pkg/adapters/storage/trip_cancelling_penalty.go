package storage

// import (
// 	"context"
// 	"strings"
// 	"tripcompanyservice/internal/trip"
// 	tripcancellingpenalty "tripcompanyservice/internal/trip_cancelling_penalty"
// 	"tripcompanyservice/pkg/adapters/storage/mappers"

// 	"gorm.io/gorm"
// )

// type tripCancellingPenaltyRepo struct {
// 	db *gorm.DB
// }

// func NewTripCancellingPenaltyRepo(db *gorm.DB) tripcancellingpenalty.Repo {
// 	return &tripCancellingPenaltyRepo{db}
// }

// func (r *tripCancellingPenaltyRepo) Insert(ctx context.Context, p *tripcancellingpenalty.TripCancelingPenalty) error {
// 	penaltyEntity := mappers.PenaltyDomainToEntity(p)

// 	result := r.db.WithContext(ctx).Save(&penaltyEntity)
// 	if result.Error != nil {
// 		if strings.Contains(result.Error.Error(), "duplicate key") {

// 			// var existingTerminal entities.Trip
// 			// // Search for the soft-deleted record with the same unique constraints
// 			// if r.db.WithContext(ctx).Unscoped().Where("name = ? AND type = ? AND city = ? AND country = ?", t.Name, t.Type, t.City, t.Country).First(&existingTerminal).Error == nil {
// 			// 	// Check if the record is soft-deleted
// 			// 	if existingTerminal.DeletedAt.Valid {
// 			// 		// Restore the soft-deleted record
// 			// 		existingTerminal.DeletedAt = gorm.DeletedAt{}
// 			// 		if err := r.db.WithContext(ctx).Save(&existingTerminal).Error; err != nil {
// 			// 			return fmt.Errorf("%w %w", terminal.ErrFailedToRestore, err)
// 			// 		}
// 			// 		t.ID = existingTerminal.ID
// 			// 		return nil
// 			// 	}

// 			return trip.ErrDuplication
// 		}
// 		return result.Error

// 	}
// 	p.ID = penaltyEntity.ID

// 	return nil
// }
