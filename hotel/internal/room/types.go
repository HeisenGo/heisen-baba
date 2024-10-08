package room

import (
	"context"
	"errors"
	"regexp"
)

type Repo interface {
	CreateRoom(ctx context.Context, room *Room) (*Room, error)
	GetRooms(ctx context.Context, page, pageSize int) ([]Room, int, error)
	UpdateRoom(ctx context.Context, room *Room) error
	DeleteRoom(ctx context.Context, id uint) error
	GetRoomByID(ctx context.Context, id uint) (*Room, error)
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
	ErrInvalidName    = errors.New("invalid room name: must be 1-100 characters long and can only contain alphanumeric characters, spaces, hyphens, underscores, and periods")
	ErrPrice          = errors.New("user price should be bigger than agency price")
	ErrRecordNotFound = errors.New("record not found")
	ErrEmptyHotelID   = errors.New("hotel_id must be entered")
)

func ValidatePrice(userprice, agencyprice uint) error {
	if userprice <= agencyprice {
		return ErrPrice
	}
	return nil
}
func ValidateRoomName(name string) error {
	var validRoomName = regexp.MustCompile(`^[a-zA-Z0-9 ._-]{1,100}$`)
	if !validRoomName.MatchString(name) {
		return ErrInvalidName
	}
	return nil
}
