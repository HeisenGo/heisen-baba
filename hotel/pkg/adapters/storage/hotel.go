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

func (r *hotelRepo) GetHotels(ctx context.Context, city, country string, capacity,page,pageSize int) ([]hotel.Hotel, uint, error) {
    var hotels []hotel.Hotel
    query := r.db.Model(&hotel.Hotel{})

    if city != "" {
        query = query.Where("city = ?", city)
    }
    if country != "" {
        query = query.Where("country = ?", country)
    }
    if capacity > 0 {
        query = query.Joins("JOIN rooms ON rooms.hotel_id = hotels.id").Where("rooms.capacity >= ?", capacity)
    }


    var total int64
    query.Count(&total)

    offset := (page - 1) * pageSize
    if err := query.Limit(pageSize).Offset(offset).Preload("Rooms").Find(&hotels).Error; err != nil {
        return nil, 0, err
    }

    return hotels, uint(total), nil
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


