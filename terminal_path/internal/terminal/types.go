package terminal

import (
	"context"
	"errors"
)

type TerminalType string

const (
	Air             TerminalType = "air"
	Rail            TerminalType = "rail"
	Road            TerminalType = "road"
	Sailing         TerminalType = "sailing" // port
	MaxStringLength int          = 100
)

var (
	ErrInvalidType      = errors.New("terminal type is not valid. It should be one of the following: road, rail, air, sailing")
	ErrTerminalNotFound = errors.New("terminal not found")
	ErrRecordsNotFound  = errors.New("any terminal exists")
	ErrDuplication      = errors.New("a terminal with this name, type, city, and country already exists")
	ErrFailedToUpdate   = errors.New("failed to update terminal ")
	ErrCanNotUpdate     = errors.New("terminal cannot be updated due to existing paths")
	ErrDeleteTerminal   = errors.New("error deleting terminal")
	ErrCanNotDelete     = errors.New("terminal cannot be deleted due to existing paths")
	ErrFailedToRestore = errors.New("failed to restore the soft-deleted terminal")
	ErrCityCountryDoNotExist = errors.New("invalid city country composite")
)

type Repo interface {
	Insert(ctx context.Context, t *Terminal) error
	GetByID(ctx context.Context, id uint) (*Terminal, error)
	GetTerminalsByCityAndType(ctx context.Context, city, terminalType, country string, limit, offset uint) ([]Terminal, uint, error)
	PatchTerminal(ctx context.Context, updatedTerminal, originalTerminal *Terminal) error
	Delete(ctx context.Context, terminalID uint) error 
}

type Terminal struct {
	ID      uint
	Name    string
	Type    TerminalType
	City    string
	Country string
}

func (t *Terminal) ValidateType() error {
	if t.Type != Air && t.Type != Sailing && t.Type != Road && t.Type != Rail {
		return ErrInvalidType
	}
	return nil
}
