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
	if err := internal.ValidateName(path.Code, MaxCodeLength); err != nil {
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

func (o *Ops) GetPathsByOriginDestinationType(ctx context.Context, originCity, destinationCity, pathType string, page, pageSize uint) ([]Path, uint, error) {
	limit := pageSize
	offset := (page - 1) * pageSize

	terminals, total, err := o.repo.GetPathsByOriginDestinationType(ctx, originCity, destinationCity, pathType, limit, offset)

	return terminals, total, err

}

func (o *Ops) GetFullPathByID(ctx context.Context, id uint) (*Path, error) {
	p, err := o.repo.GetFullPathByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if p == nil {
		return nil, ErrPathNotFound
	}

	return p, nil

}

func (o *Ops) PatchPath(ctx context.Context, updatedPath, originalPath *Path, hasUnFinishedTrip bool) error {
	if hasUnFinishedTrip {
		if updatedPath.FromTerminalID != uint(0) || updatedPath.ToTerminalID != uint(0) {
			return ErrCanNotUpdatePath
		}
	} else {
		if updatedPath.FromTerminalID != uint(0) && updatedPath.ToTerminalID == uint(0) {
			updatedPath.ToTerminal = originalPath.ToTerminal
			updatedPath.ToTerminalID = originalPath.ToTerminalID
			err := updatedPath.ValidateStartEndTerminalTypes()
			if err != nil {
				return err
			}
			err = updatedPath.ValidateStartAndEndCities()
			if err != nil {
				return err
			}
			originalPath.FromTerminal = updatedPath.FromTerminal
			originalPath.FromTerminalID = updatedPath.FromTerminalID
		}

		if updatedPath.FromTerminalID == uint(0) && updatedPath.ToTerminalID != uint(0) {
			updatedPath.FromTerminal = originalPath.FromTerminal
			updatedPath.FromTerminalID = originalPath.FromTerminalID
			err := updatedPath.ValidateStartEndTerminalTypes()
			if err != nil {
				return err
			}
			err = updatedPath.ValidateStartAndEndCities()
			if err != nil {
				return err
			}
			originalPath.ToTerminal = updatedPath.ToTerminal
			originalPath.ToTerminalID = updatedPath.ToTerminalID
		}
		if updatedPath.FromTerminalID != uint(0) && updatedPath.ToTerminalID != uint(0) {
			err := updatedPath.ValidateStartEndTerminalTypes()
			if err != nil {
				return err
			}
			err = updatedPath.ValidateStartAndEndCities()
			if err != nil {
				return err
			}
			updatedPath.Type = updatedPath.FromTerminal.Type
			originalPath.FromTerminal = updatedPath.FromTerminal
			originalPath.FromTerminalID = updatedPath.FromTerminalID
			originalPath.ToTerminal = updatedPath.ToTerminal
			originalPath.ToTerminalID = updatedPath.ToTerminalID

		}
	}
	if updatedPath.Name != "" {
		if err := internal.ValidateName(updatedPath.Name, MaxStringLength); err != nil {
			return err
		}
	}
	if updatedPath.Code != "" {
		if err := internal.ValidateName(updatedPath.Code, MaxCodeLength); err != nil {
			return err
		}
	}
	return o.repo.PatchPath(ctx, updatedPath, originalPath)
}

func (o *Ops) Delete(ctx context.Context, pathID uint, hasUnfinishedTrip bool) error {
	if hasUnfinishedTrip {
		return ErrCanNotDelete
	}
	return o.repo.Delete(ctx, pathID)
}

func (o *Ops) AreTherePathRelatedToTerminalID(ctx context.Context, terminalID uint) (bool, error) {
	return o.repo.AreTherePathRelatedToTerminalID(ctx, terminalID)
}

func (o *Ops) IsTherePathRelatedToTerminalID(ctx context.Context, terminalID uint) (bool, error) {
	return o.repo.IsTherePathRelatedToTerminalID(ctx, terminalID)
}
