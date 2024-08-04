package tour

import (
	"context"
	"errors"
	"time"
)

type Repo interface {
	CreateTour(ctx context.Context, tour *Tour) error
	GetTours(ctx context.Context, agencyID uint, page, pageSize int) ([]Tour, uint, error)
	GetToursByAgencyID(ctx context.Context, agencyID uint, page, pageSize int) ([]Tour, int, error)
	GetTourByID(ctx context.Context, id uint) (*Tour, error)
	UpdateTour(ctx context.Context, tour *Tour) error
	DeleteTour(ctx context.Context, id uint) error
	ApproveTour(ctx context.Context, tourID uint) error
	SetTourStatus(ctx context.Context, tourID uint, isActive bool) error
}

type Tour struct {
	ID           uint
	AgencyID     uint
	GoTicketID   uint
	BackTicketID uint
	HotelID      uint // Changed from 'HotelID' to 'AgencyID'
	Capacity     uint
	ReleaseDate  time.Time
	IsApproved   bool
	IsActive     bool
}

var (
	ErrInvalidTourCapacity = errors.New("invalid tour capacity")
	ErrRecordNotFound      = errors.New("record not found")
)

func ValidateTourCapacity(capacity uint) error {
	if capacity <= 0 {
		return ErrInvalidTourCapacity
	}
	return nil
}
