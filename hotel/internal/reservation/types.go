package reservation

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"time"
)

type Repo interface {
	CreateReservation(ctx context.Context, reservation *Reservation) error
	GetReservationsByHotelOwner(ctx context.Context, ownerID uuid.UUID, page, pageSize int) ([]Reservation, int, error)
	GetReservationByUserID(ctx context.Context, userid uuid.UUID) ([]Reservation, error)
	GetReservationByID(ctx context.Context, id uint) (*Reservation, error)
	UpdateReservation(ctx context.Context, reservation *Reservation) error
	DeleteReservation(ctx context.Context, id uint) error
}

type Reservation struct {
	ID         uint
	OwnerID    uuid.UUID
	RoomID     uint
	UserID     uuid.UUID
	CheckIn    time.Time
	CheckOut   time.Time
	TotalPrice uint64
	Status     string // e.g., "booked", "checked_in", "checked_out", "canceled"
}

var (
	ErrInvalidReservationStatus = errors.New("invalid reservation status: must be one of 'booked', 'checked_in', 'checked_out', 'canceled'")
	ErrInvalidAmount            = errors.New("invalid amount: must be a positive number")
	ErrInvalidDates             = errors.New("invalid dates: check-in date must be before check-out date")
	ErrRecordNotFound           = errors.New("record not found")
)

func ValidateReservationStatus(status string) error {
	var validStatuses = map[string]bool{
		"booked":      true,
		"checked_in":  true,
		"checked_out": true,
		"canceled":    true,
		"pending" : true,
	}
	_, ok := validStatuses[status]
	if !ok {
		return ErrInvalidReservationStatus
	}
	return nil
}

func ValidateAmount(amount uint64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	return nil
}

func ValidateDates(checkIn, checkOut time.Time) error {
	if checkIn.After(checkOut) {
		return ErrInvalidDates
	}
	return nil
}
