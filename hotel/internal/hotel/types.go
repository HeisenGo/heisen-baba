package hotel

import (
	"context"
	"errors"
	"hotel/internal/room"
	"regexp"

	"github.com/google/uuid"
)

type Repo interface {
	CreateHotel(ctx context.Context, hotel *Hotel) error
	GetHotels(ctx context.Context, city, country string, capacity, page, pageSize int) ([]Hotel, uint, error)
	GetHotelsByOwnerID(ctx context.Context, ownerID uint, page, pageSize int) ([]Hotel, int, error)
	GetHotelsByID(ctx context.Context, id uint) (*Hotel, error)
	UpdateHotel(ctx context.Context, hotel *Hotel) error
	DeleteHotel(ctx context.Context, id uint) error
	BlockHotel(ctx context.Context,hotelID uint) error
}

type Hotel struct {
	OwnerID   uuid.UUID
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
	ErrInvalidCapacity  = errors.New("invalid capacity")
	ErrRecordNotFound   = errors.New("record not found")
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
func ValidateCapacity(number int) error {
	if number <= 0 {
		return ErrInvalidCapacity
	}
	return nil
}
