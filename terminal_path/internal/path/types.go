package path

import (
	"context"
	"errors"
	"terminalpathservice/internal/terminal"
)

var (
	ErrQuantityGreater = errors.New("quantity should not be greater than price")
	ErrWrongOrderTime  = errors.New("wrong order time")
)

type Repo interface {
	GetByID(ctx context.Context, id uint) (*Path, error)
	Insert(ctx context.Context, p *Path) error
}

type Path struct {
	ID             uint
	FromTerminalID uint
	ToTerminalID   uint
	FromTerminal   terminal.Terminal
	ToTerminal     terminal.Terminal
	Distance       float64 // in kilometers
	Code           string
	Name           string
}
