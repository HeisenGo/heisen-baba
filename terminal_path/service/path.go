package service

import (
	"context"
	"errors"
	"fmt"
	"terminalpathservice/internal/path"
	"terminalpathservice/internal/terminal"
)

type PathService struct {
	pathOps     *path.Ops
	terminalOps *terminal.Ops
}

func NewPathService(pathOps *path.Ops, terminalOps *terminal.Ops) *PathService {
	return &PathService{
		pathOps:     pathOps,
		terminalOps: terminalOps,
	}
}

func (s *PathService) CreatePath(ctx context.Context, path *path.Path) error {
	var err error
	path.FromTerminal, err = s.terminalOps.GetTerminalByID(ctx, path.FromTerminalID)
	if err != nil {
		if errors.Is(err, terminal.ErrTerminalNotFound) {
			return fmt.Errorf("from %w", err)
		}
		return err
	}
	path.ToTerminal, err = s.terminalOps.GetTerminalByID(ctx, path.ToTerminalID)
	if err != nil {
		if errors.Is(err, terminal.ErrTerminalNotFound) {
			return fmt.Errorf("to %w", err)
		}
		return err
	}
	return s.pathOps.Create(ctx, path)
}

func (s *PathService) GetPathsByOriginDestinationType(ctx context.Context, originCity, destinationCity, pathType string, page, pageSize uint) ([]path.Path, uint, error) {
	return s.pathOps.GetPathsByOriginDestinationType(ctx, originCity, destinationCity, pathType, page, pageSize)
}

func (s *PathService) GetFullPathByID(ctx context.Context, id uint) (*path.Path, error) {
	return s.pathOps.GetFullPathByID(ctx, id)
}

func (s *PathService) PatchPath(ctx context.Context, updatedPath *path.Path, hasUnfinishedTrips bool) (*path.Path, error) {
	originalPath, err := s.GetFullPathByID(ctx, updatedPath.ID)
	if err != nil {
		return nil, err
	}
	if updatedPath.FromTerminalID != 0 {
		updatedPath.FromTerminal, err = s.terminalOps.GetTerminalByID(ctx, updatedPath.FromTerminalID)
		if err != nil {
			return nil, err
		}
	}

	if updatedPath.ToTerminalID != 0 {
		updatedPath.ToTerminal, err = s.terminalOps.GetTerminalByID(ctx, updatedPath.ToTerminalID)
		if err != nil {
			return nil, err
		}
	}
	err = s.pathOps.PatchPath(ctx, updatedPath, originalPath, hasUnfinishedTrips)
	if err != nil {
		return nil, err
	}
	return originalPath, nil
}

func (s *PathService) DeletePath(ctx context.Context, pathID uint, hasUnFinishedTrips bool) (*path.Path, error) {
	p, err := s.pathOps.GetFullPathByID(ctx, pathID)
	if err != nil {
		return nil, err
	}
	err = s.pathOps.Delete(ctx, pathID, hasUnFinishedTrips)
	if err != nil {
		return nil, err
	}
	return p, nil
}
