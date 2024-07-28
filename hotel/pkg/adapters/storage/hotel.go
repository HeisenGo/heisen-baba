package storage

import (
	"context"
	"hotel/internal/hotel"
	"hotel/pkg/adapters/storage/mappers"

	"gorm.io/gorm"
)

type hotelRepo struct {
	db *gorm.DB
}

func NewHotelRepo(db *gorm.DB) hotel.Repo {
	return &hotelRepo{
		db: db,
	}
}


func (r *hotelRepo) CreateHotel(ctx context.Context, h *hotel.Hotel) error {
	hotelEntity := mappers.HotelDomainToEntity(h)
	if err := r.db.WithContext(ctx).Save(&hotelEntity).Error; err != nil {
		return err
	}

	h.ID = hotelEntity.ID
	return nil
}


