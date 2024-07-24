package path

import (
	"context"
	"errors"
	"terminalpathservice/internal/terminal"
)

const (
	MaxStringLength = 100
)

var (
	ErrPathNotFound                 = errors.New("path not found")
	ErrMisMatchStartEndTerminalType = errors.New("terminal types for starting and ending a path should be the same")
	ErrSameCitiesTerminals          = errors.New("same city terminals")
	ErrRecordsNotFound              = errors.New("any path exists")
	ErrDuplication                  = errors.New("a path with this code already exists")
)

type Repo interface {
	GetByID(ctx context.Context, id uint) (*Path, error)
	Insert(ctx context.Context, p *Path) error
	GetPathsByOriginDestinationType(ctx context.Context, originCity, destinationCity, pathType string, limit, offset uint) ([]Path, uint, error)
}

type Path struct {
	ID             uint
	FromTerminalID uint
	ToTerminalID   uint
	FromTerminal   *terminal.Terminal
	ToTerminal     *terminal.Terminal
	DistanceKM     float64 // in kilometers
	Code           string
	Name           string
	Type           terminal.TerminalType
}

func (p *Path) ValidateStartEndTerminalTypes() error {
	if p.FromTerminal.Type != p.ToTerminal.Type {
		return ErrMisMatchStartEndTerminalType
	}
	return nil
}

func (p *Path) ValidateStartAndEndCities() error {
	if p.FromTerminal.City == p.ToTerminal.City {
		return ErrSameCitiesTerminals
	}
	return nil
}
