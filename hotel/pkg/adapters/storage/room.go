package storage

import (
	"context"
	"errors"
	"hotel/internal/room"
	"hotel/pkg/adapters/storage/entities"
	"hotel/pkg/adapters/storage/mappers"

	"gorm.io/gorm"
)

type roomRepo struct {
	db *gorm.DB
}

func NewRoomRepo(db *gorm.DB) room.Repo {
	return &roomRepo{
		db: db,
	}
}

func (r *roomRepo) CreateRoom(ctx context.Context, rm *room.Room) (*room.Room, error) {
	roomEntity := mappers.RoomDomainToEntity(*rm)
	if err := r.db.WithContext(ctx).Create(&roomEntity).Error; err != nil {
		return nil, err
	}
	rm.ID = roomEntity.ID
	return rm, nil
}

func (r *roomRepo) GetRooms(ctx context.Context, page, pageSize int) ([]room.Room, int, error) {
	var roomEntities []entities.Room
	var total int64

	query := r.db.WithContext(ctx).Model(&entities.Room{})

	query.Count(&total)

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&roomEntities).Error; err != nil {
		return nil, 0, err
	}

	rooms := make([]room.Room, len(roomEntities))
	for i, roomEntity := range roomEntities {
		rooms[i] = mappers.RoomEntityToDomain(roomEntity)
	}

	return rooms, int(total), nil
}

func (r *roomRepo) GetRoomByID(ctx context.Context, id uint) (*room.Room, error) {
	var ro entities.Room
	if err := r.db.WithContext(ctx).First(&ro, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, room.ErrRecordNotFound
		}
		return nil, err
	}
	selectedroom := mappers.RoomEntityToDomain(ro)
	return &selectedroom, nil
}

func (r *roomRepo) UpdateRoom(ctx context.Context, rm *room.Room) error {
	roomEntity := mappers.RoomDomainToEntity(*rm)
	if err := r.db.WithContext(ctx).Model(&entities.Room{}).Where("id = ?", rm.ID).Updates(roomEntity).Error; err != nil {
		return err
	}
	return nil
}
func (r *roomRepo) DeleteRoom(ctx context.Context, id uint) error {
	var selectedroom entities.Room
	if err := r.db.WithContext(ctx).First(&selectedroom, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return room.ErrRecordNotFound
		}
		return err
	}

	if err := r.db.WithContext(ctx).Delete(&selectedroom).Error; err != nil {
		return err
	}
	return nil
}

func (r *roomRepo) GetRoomsByHotelID(ctx context.Context, hotelID uint) ([]room.Room, error) {
	var roomEntities []entities.Room
	if err := r.db.WithContext(ctx).Where("hotel_id = ?", hotelID).Find(&roomEntities).Error; err != nil {
		return nil, err
	}
	rooms := make([]room.Room, len(roomEntities))
	for i, roomEntity := range roomEntities {
		rooms[i] = mappers.RoomEntityToDomain(roomEntity)
	}
	return rooms, nil
}
