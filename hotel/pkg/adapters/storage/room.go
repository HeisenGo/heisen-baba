package storage

import (
	"context"
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

func (r *roomRepo) GetByID(ctx context.Context, id uint) (*room.Room, error) {
	var roomEntity entities.Room
	if err := r.db.WithContext(ctx).First(&roomEntity, id).Error; err != nil {
		return nil, err
	}
	roomDomain := mappers.RoomEntityToDomain(roomEntity)
	return &roomDomain, nil
}

func (r *roomRepo) UpdateRoom(ctx context.Context, rm *room.Room) (*room.Room, error) {
	roomEntity := mappers.RoomDomainToEntity(*rm)
	if err := r.db.WithContext(ctx).Save(&roomEntity).Error; err != nil {
		return nil, err
	}
	return rm, nil
}

func (r *roomRepo) DeleteRoom(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&entities.Room{}, id).Error; err != nil {
		return err
	}
	return nil
}
