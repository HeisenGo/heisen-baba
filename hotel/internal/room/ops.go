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

func (o *Ops) CreateRoom(ctx context.Context, room *Room) (*Room,error) {
	if err := ValidateRoomName(room.Facilities); err != nil {
		return nil , ErrInvalidName
	}
	
	return o.repo.CreateRoom(ctx, room)
}