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
