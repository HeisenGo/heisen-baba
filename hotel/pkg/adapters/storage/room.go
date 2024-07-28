package storage

import (
	"context"
	"hotel/internal/room"
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


func (r *roomRepo) CreateRoom(ctx context.Context, room *room.Room) (*room.Room, error) {
	roomEntity := mappers.RoomDomainToEntity(room)
	if err := r.db.WithContext(ctx).Save(&roomEntity).Error; err != nil {
		return nil , err
	}

	room.ID = roomEntity.ID
	return room , nil
}


