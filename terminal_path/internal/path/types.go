package path

import (
	"context"
	"errors"
	"terminalpathservice/internal/terminal"
)

const (
	MaxStringLength = 100
	MaxCodeLength   = 50
)

var (
	ErrPathNotFound                 = errors.New("path not found")
	ErrMisMatchStartEndTerminalType = errors.New("terminal types for starting and ending a path should be the same")
	ErrSameCitiesTerminals          = errors.New("same city terminals")
	ErrRecordsNotFound              = errors.New("any path exists")
	ErrDuplication                  = errors.New("a path with this code already exists")
	ErrCanNotUpdatePath             = errors.New("path can not be updated deu to existing unfinished trips")
	ErrFailedToUpdate               = errors.New("updating failed try again later")
	ErrDeletePath                   = errors.New("path was not deleted try later")
	ErrCanNotDelete                 = errors.New("can not delete path due to existing unfinished trips")
	ErrFailedToGetPath              = errors.New("failed to get path")
	ErrFailedToRestore              = errors.New("failed to restore the soft-deleted path")
	ErrCodeIsImpossibleToUse = errors.New("please change the entered code this code was used in a deleted path with different type and start and end terminals")

)

type Repo interface {
	GetByID(ctx context.Context, id uint) (*Path, error)
	Insert(ctx context.Context, p *Path) error
	GetPathsByOriginDestinationType(ctx context.Context, originCity, destinationCity, pathType string, limit, offset uint) ([]Path, uint, error)
	PatchPath(ctx context.Context, updatedPath, originalPath *Path) error
	Delete(ctx context.Context, pathID uint) error
	GetFullPathByID(ctx context.Context, id uint) (*Path, error)
	AreTherePathRelatedToTerminalID(ctx context.Context, terminalID uint) (bool, error)
	IsTherePathRelatedToTerminalID(ctx context.Context, terminalID uint) (bool, error)
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
