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
	if err := ValidatePrice(uint(room.UserPrice), uint(room.AgencyPrice)); err != nil {
		return nil, ErrPrice
	}
	return o.repo.CreateRoom(ctx, room)
}

func (o *Ops) GetRooms(ctx context.Context, page, pageSize int) ([]Room, int, error) {
	return o.repo.GetRooms(ctx, page, pageSize)
}
func (o *Ops) GetRoomByID(ctx context.Context, id uint) (*Room, error) {
	return o.repo.GetRoomByID(ctx, id)
}
func (o *Ops) UpdateRoom(ctx context.Context, room *Room) error {
	// Ensure room exists before updating
	existingRoom, err := o.repo.GetRoomByID(ctx, room.ID)
	if err != nil {
		return err
	}
	if existingRoom == nil {
		return ErrRecordNotFound
	}
	if err := ValidateRoomName(room.Name); err != nil {
		return ErrInvalidName
	}
	if err := ValidatePrice(uint(room.UserPrice), uint(room.AgencyPrice)); err != nil {
		return ErrPrice
	}
	return o.repo.UpdateRoom(ctx, room)
}

func (o *Ops) DeleteRoom(ctx context.Context, id uint) error {
	// Ensure room exists before deleting
	existingRoom, err := o.repo.GetRoomByID(ctx, id)
	if err != nil {
		return err
	}
	if existingRoom == nil {
		return ErrRecordNotFound
	}

	return o.repo.DeleteRoom(ctx, id)
}
