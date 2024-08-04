package agency

import (
	"agency/internal/tour"
	"context"
	"errors"
	"regexp"

	"github.com/google/uuid"
)

type Repo interface {
	CreateAgency(ctx context.Context, agency *Agency) error
	GetAgencies(ctx context.Context, name string, page, pageSize int) ([]Agency, uint, error)
	GetAgenciesByOwnerID(ctx context.Context, ownerID uuid.UUID, page, pageSize int) ([]Agency, int, error)
	GetAgencyByID(ctx context.Context, id uint) (*Agency, error)
	UpdateAgency(ctx context.Context, agency *Agency) error
	DeleteAgency(ctx context.Context, id uint) error
	BlockAgency(ctx context.Context, agencyID uint) error
}

type Agency struct {
	OwnerID   uuid.UUID
	ID        uint
	Name      string
	IsBlocked bool
	Tours     []tour.Tour 
}

var (
	ErrInvalidAgencyName = errors.New("invalid agency name: must be 1-100 characters long and can only contain alphanumeric characters, spaces, hyphens, underscores, and periods")
	ErrRecordNotFound    = errors.New("record not found")
)

func ValidateAgencyName(name string) error {
	var validAgencyName = regexp.MustCompile(`^[a-zA-Z0-9 ._-]{1,100}$`)
	if !validAgencyName.MatchString(name) {
		return ErrInvalidAgencyName
	}
	return nil
}
