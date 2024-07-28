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
	return o.repo.CreateRoom(ctx, room)
}

func (o *Ops) GetByID(ctx context.Context, id uint) (*Room, error) {
	return o.repo.GetByID(ctx, id)
}

func (o *Ops) UpdateRoom(ctx context.Context, room *Room) (*Room, error) {
	if err := ValidateRoomName(room.Name); err != nil {
		return nil, ErrInvalidName
	}
	return o.repo.UpdateRoom(ctx, room)
}

func (o *Ops) DeleteRoom(ctx context.Context, id uint) error {
	return o.repo.DeleteRoom(ctx, id)
}