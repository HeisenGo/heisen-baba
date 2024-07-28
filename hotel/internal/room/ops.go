package room

import (
	"context"
)

type Ops struct {
	repo Repo
}

func NewOps(repo Repo) *Ops {
	return &Ops{repo}
}

func (o *Ops) CreateRoom(ctx context.Context, room *Room) (*Room, error) {
	if err := ValidateRoomName(room.Name); err != nil {
		return nil, ErrInvalidName
	}
	if err := ValidatePrice(uint(room.UserPrice),uint(room.AgencyPrice)); err != nil {
		return nil, ErrPrice
	}
	return o.repo.CreateRoom(ctx, room)
}

func (o *Ops) GetRooms(ctx context.Context, page, pageSize int) ([]Room, int, error) {
	return o.repo.GetRooms(ctx, page, pageSize)
}

func (o *Ops) UpdateRoom(ctx context.Context, room *Room) (*Room, error) {
	if err := ValidateRoomName(room.Name); err != nil {
		return nil, ErrInvalidName
	}
	if err := ValidatePrice(uint(room.UserPrice),uint(room.AgencyPrice)); err != nil {
		return nil, ErrPrice
	}
	return o.repo.UpdateRoom(ctx, room)
}

func (o *Ops) DeleteRoom(ctx context.Context, id uint) error {
	return o.repo.DeleteRoom(ctx, id)
}