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
	ErrQuantityGreater = errors.New("quantity should not be greater than price")
	ErrWrongOrderTime  = errors.New("wrong order time")
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
