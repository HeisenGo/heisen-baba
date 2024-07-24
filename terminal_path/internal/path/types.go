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
	ErrSameCitiesTerminals          = errors.New("same city terminals with less than 90 km DistanceKM")
)

type Repo interface {
	GetByID(ctx context.Context, id uint) (*Path, error)
	Insert(ctx context.Context, p *Path) error
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
