package storage

import (
	"context"
	"hotel/internal/hotel"
	"hotel/pkg/adapters/storage/entities"
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
	if err := r.db.WithContext(ctx).Create(&hotelEntity).Error; err != nil {
		return err
	}
	h.ID = hotelEntity.ID
	return nil
}

func (r *hotelRepo) GetByID(ctx context.Context, id uint) (*hotel.Hotel, error) {
	var hotelEntity entities.Hotel
	if err := r.db.WithContext(ctx).Preload("Rooms").First(&hotelEntity, id).Error; err != nil {
		return nil, err
	}
	domainHotel := mappers.HotelEntityToDomain(hotelEntity)
	return &domainHotel, nil
}

func (r *hotelRepo) UpdateHotel(ctx context.Context, h *hotel.Hotel) error {
	hotelEntity := mappers.HotelDomainToEntity(h)
	if err := r.db.WithContext(ctx).Save(hotelEntity).Error; err != nil {
		return err
	}
	return nil
}

func (r *hotelRepo) DeleteHotel(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&entities.Hotel{}, id).Error; err != nil {
		return err
	}
	return nil
}


