package hotel

import (
	"context"
	"errors"
	"hotel/internal/room"
	"regexp"
)

type Repo interface {
	CreateHotel(ctx context.Context, hotel *Hotel) error
	GetByID(ctx context.Context, id uint) (*Hotel, error)
	UpdateHotel(ctx context.Context, hotel *Hotel) error
	DeleteHotel(ctx context.Context, id uint) error
}

type Hotel struct {
	OwnerID   uint
	ID        uint
	Name      string
	City      string
	Country   string
	Details   string
	IsBlocked bool
	Rooms     []room.Room
}

var (
	ErrInvalidHotelName = errors.New("invalid hotel name: must be 1-100 characters long and can only contain alphanumeric characters, spaces, hyphens, underscores, and periods")
	ErrInvalidName      = errors.New("invalid city or country name : only alphabetic will be accepted")
)

func ValidateHotelName(name string) error {
	var validHotelName = regexp.MustCompile(`^[a-zA-Z0-9 ._-]{1,100}$`)
	if !validHotelName.MatchString(name) {
		return ErrInvalidHotelName
	}
	return nil
}

func ValidateName(name string) error {
	var validName = regexp.MustCompile(`^[a-zA-Z\s]+$`)
	if !validName.MatchString(name) {
		return ErrInvalidName
	}
	return nil
}