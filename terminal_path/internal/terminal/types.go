package terminal

import (
	"context"
	"errors"
)

type TerminalType string

const (
	Air     TerminalType = "ait"
	Rail    TerminalType = "rail"
	Road    TerminalType = "road"
	Sailing TerminalType = "sailing" // port
)

var (
	ErrInvalidType      = errors.New("terminal type is not valid. It should be one of the following: road, rail, air, sailing")
	ErrTerminalNotFound = errors.New("terminal not found")
)

type Repo interface {
	Insert(ctx context.Context, t *Terminal) error
	GetByID(ctx context.Context, id uint) (*Terminal, error)
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
