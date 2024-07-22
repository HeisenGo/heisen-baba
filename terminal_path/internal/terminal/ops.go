package terminal

import (
	"context"
	"terminalpathservice/internal"
)

type Ops struct {
	repo Repo
}

func NewOps(repo Repo) *Ops {
	return &Ops{
		repo: repo,
	}
}

func (o *Ops) Create(ctx context.Context, terminal *Terminal) error {
	if err := terminal.ValidateType(); err != nil {
		return ErrInvalidType
	}

	if err := internal.ValidateName(terminal.Name, MaxStringLength); err != nil {
		return err
	}
	if err := internal.ValidateName(terminal.City, MaxStringLength); err != nil {
		return err
	}
	if err := internal.ValidateName(terminal.Country, MaxStringLength); err != nil {
		return err
	}
	return o.repo.Insert(ctx, terminal)
}

func (o *Ops) GetTerminalByID(ctx context.Context, id uint) (*Terminal, error) {
	t, err := o.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if t == nil {
		return nil, ErrTerminalNotFound
	}

	return t, nil
}

func (o *Ops) CityTypeTerminals(ctx context.Context,country, city, terminalType string, page, pageSize uint) ([]Terminal, uint, error) {
	limit := pageSize
	offset := (page - 1) * pageSize

	terminals, total, err := o.repo.GetTerminalsByCityAndType(ctx, country, city, terminalType, limit, offset)

	return terminals, total, err
}
