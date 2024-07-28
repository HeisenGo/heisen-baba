package service

import (
	"context"
	"hotel/internal/room"
)

type RoomService struct {
	roomOps room.Repo
}

func NewRoomService(roomOps room.Repo) *RoomService {
	return &RoomService{
		roomOps: roomOps,
	}
}

func (s *RoomService) CreateRoom(ctx context.Context, r *room.Room) (*room.Room, error) {
	return s.roomOps.CreateRoom(ctx, r)
}

func (s *RoomService) GetRooms(ctx context.Context, page, pageSize int) ([]room.Room, int, error) {
	return s.roomOps.GetRooms(ctx, page, pageSize)
}

func (s *RoomService) UpdateRoom(ctx context.Context, r *room.Room) (*room.Room, error) {
	return s.roomOps.UpdateRoom(ctx, r)
}

func (s *RoomService) DeleteRoom(ctx context.Context, id uint) error {
	return s.roomOps.DeleteRoom(ctx, id)
}