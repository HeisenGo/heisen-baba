package storage

import (
	"context"
	"errors"
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

func (r *hotelRepo) GetHotels(ctx context.Context, city, country string, capacity, page, pageSize int) ([]hotel.Hotel, uint, error) {
	var h []entities.Hotel
	var int64Total int64

	query := r.db.Model(&entities.Hotel{}).Preload("Rooms")

	// Filters
	if city != "" {
		query = query.Where("city = ?", city)
	}
	if country != "" {
		query = query.Where("country = ?", country)
	}
	if capacity > 0 {
		query = query.Joins("JOIN rooms ON rooms.hotel_id = hotels.id").Where("rooms.capacity >= ?", capacity)
	}

	// Count total records for pagination
	query.Count(&int64Total)

	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	if err := query.Find(&h).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, nil
		}
		return nil, 0, err
	}

	total := uint(int64Total)
	hotels := mappers.BatchHotelEntitiesToDomain(h)
	return hotels, total, nil
}
func (r *hotelRepo) GetHotelsByOwnerID(ctx context.Context, ownerID uint, page, pageSize int) ([]hotel.Hotel, int, error) {
	var hotelEntities []entities.Hotel
	var total int64

	query := r.db.WithContext(ctx).Model(&entities.Hotel{}).Where("owner_id = ?", ownerID)

	query.Count(&total)

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&hotelEntities).Error; err != nil {
		return nil, 0, err
	}

	hotels := make([]hotel.Hotel, len(hotelEntities))
	for i, hotelEntity := range hotelEntities {
		hotels[i] = mappers.HotelEntityToDomain(hotelEntity)
	}

	return hotels, int(total), nil
}
func (r *hotelRepo) GetHotelsByID(ctx context.Context, id uint) (*hotel.Hotel, error) {
	var hotel hotel.Hotel
	if err := r.db.First(&hotel, id).Error; err != nil {
		return nil, err
	}
	return &hotel, nil
}
func (r *hotelRepo) UpdateHotel(ctx context.Context, h *hotel.Hotel) error {
	hotelEntity := mappers.HotelDomainToEntity(h)
	if err := r.db.WithContext(ctx).Model(&entities.Hotel{}).Where("id = ?", h.ID).Updates(hotelEntity).Error; err != nil {
		return err
	}
	return nil
}

func (r *hotelRepo) DeleteHotel(ctx context.Context, id uint) error {
	var h entities.Hotel
	if err := r.db.WithContext(ctx).First(&h, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return hotel.ErrRecordNotFound
		}
		return err
	}

	// Cascading delete rooms
	if err := r.db.WithContext(ctx).Where("hotel_id = ?", id).Delete(&entities.Room{}).Error; err != nil {
		return err
	}

	// Delete hotel
	if err := r.db.WithContext(ctx).Delete(&h).Error; err != nil {
		return err
	}
	return nil
}
