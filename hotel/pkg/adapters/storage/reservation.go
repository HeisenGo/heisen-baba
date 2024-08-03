package storage

import (
	"context"
	"hotel/internal/reservation"
	"hotel/pkg/adapters/storage/entities"
	"hotel/pkg/adapters/storage/mappers"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type reservationRepo struct {
	db *gorm.DB
}

func NewReservationRepo(db *gorm.DB) reservation.Repo {
	return &reservationRepo{
		db: db,
	}
}

func (r *reservationRepo) CreateReservation(ctx context.Context, res *reservation.Reservation) error {
	reservationEntity := mappers.ReservationDomainToEntity(res)
	if err := r.db.WithContext(ctx).Create(&reservationEntity).Error; err != nil {
		return err
	}
	res.ID = reservationEntity.ID
	return nil
}

func (r *reservationRepo) GetReservationsByHotelOwner(ctx context.Context, ownerID uuid.UUID, page, pageSize int) ([]reservation.Reservation, int, error) {
	var reservationEntities []entities.Reservation
	var total int64

	query := r.db.WithContext(ctx).Model(&entities.Reservation{}).Joins("JOIN rooms ON reservations.room_id = rooms.id").Where("rooms.hotel_id = ?", ownerID)

	query.Count(&total)

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&reservationEntities).Error; err != nil {
		return nil, 0, err
	}

	reservations := mappers.BatchReservationEntitiesToDomain(reservationEntities)
	return reservations, int(total), nil
}
func (r *reservationRepo) GetReservationByUserID(ctx context.Context, userID uuid.UUID) ([]reservation.Reservation, error) {
	var reservationEntities []entities.Reservation
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&reservationEntities).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, reservation.ErrRecordNotFound
		}
		return nil, err
	}
	reservations := mappers.BatchReservationEntitiesToDomain(reservationEntities)
	return reservations, nil
}


func (r *reservationRepo) GetReservationByID(ctx context.Context, id uint) (*reservation.Reservation, error) {
	var reservationEntity entities.Reservation
	if err := r.db.WithContext(ctx).First(&reservationEntity, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, reservation.ErrRecordNotFound
		}
		return nil, err
	}
	re :=mappers.ReservationEntityToDomain(reservationEntity)
	return &re, nil
}

func (r *reservationRepo) UpdateReservation(ctx context.Context, res *reservation.Reservation) error {
	reservationEntity := mappers.ReservationDomainToEntity(res)
	if err := r.db.WithContext(ctx).Model(&entities.Reservation{}).Where("id = ?", res.ID).Updates(reservationEntity).Error; err != nil {
		return err
	}
	return nil
}

func (r *reservationRepo) DeleteReservation(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&entities.Reservation{}, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return reservation.ErrRecordNotFound
		}
		return err
	}
	return nil
}