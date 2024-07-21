package path

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

func (o *Ops) Create(ctx context.Context, path *Path) error {
	if err := path.ValidateStartEndTerminalTypes(); err != nil {
		return err
	}
	if err := internal.ValidateName(path.Name, MaxStringLength); err != nil {
		return err
	}
	if err := path.ValidateStartAndEndCities(); err != nil {
		return err
	}
	path.Type = path.FromTerminal.Type
	return o.repo.Insert(ctx, path)
}

func (o *Ops) GetPathByID(ctx context.Context, id uint) (*Path, error) {
	p, err := o.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if p == nil {
		return nil, ErrPathNotFound
	}

	return p, nil
}
