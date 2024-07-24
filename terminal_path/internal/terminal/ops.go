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

func (o *Ops) CityTypeTerminals(ctx context.Context, country, city, terminalType string, page, pageSize uint) ([]Terminal, uint, error) {
	limit := pageSize
	offset := (page - 1) * pageSize

	terminals, total, err := o.repo.GetTerminalsByCityAndType(ctx, country, city, terminalType, limit, offset)

	return terminals, total, err
}

func (o *Ops) PatchTerminal(ctx context.Context, updatedTerminal, originalTerminal *Terminal) error {
	if updatedTerminal.Type != "" {
		if err := updatedTerminal.ValidateType(); err != nil {
			return ErrInvalidType
		}
	}
	if updatedTerminal.Name != "" {
		if err := internal.ValidateName(updatedTerminal.Name, MaxStringLength); err != nil {
			return err
		}
	}
	if updatedTerminal.City != "" {
		if err := internal.ValidateName(updatedTerminal.City, MaxStringLength); err != nil {
			return err
		}
	}
	if updatedTerminal.Country != "" {
		if err := internal.ValidateName(updatedTerminal.Country, MaxStringLength); err != nil {
			return err
		}
	}
	return o.repo.PatchTerminal(ctx, updatedTerminal, originalTerminal)
}

func (o *Ops) Delete(ctx context.Context, terminalID uint) error {
	return o.repo.Delete(ctx, terminalID)
}
