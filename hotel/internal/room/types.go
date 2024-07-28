package room

import (
	"context"
	"errors"
	"regexp"
)

type Repo interface {
	CreateRoom(ctx context.Context, room *Room) (*Room, error)
	GetByID(ctx context.Context, id uint) (*Room, error)
	UpdateRoom(ctx context.Context, room *Room) (*Room, error)
	DeleteRoom(ctx context.Context, id uint) error
}

type Room struct {
	ID          uint
	Name        string
	HotelID     uint
	AgencyPrice uint64
	UserPrice   uint64
	Facilities  string
	Capacity    uint8
	IsAvailable bool
}

var (
	ErrInvalidName = errors.New("invalid room name: must be 1-100 characters long and can only contain alphanumeric characters, spaces, hyphens, underscores, and periods")
)

func ValidateRoomName(name string) error {
	var validRoomName = regexp.MustCompile(`^[a-zA-Z0-9 ._-]{1,100}$`)
	if !validRoomName.MatchString(name) {
		return ErrInvalidName
	}
	return nil
}