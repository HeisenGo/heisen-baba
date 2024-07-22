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
