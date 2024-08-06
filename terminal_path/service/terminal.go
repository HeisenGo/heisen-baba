package service

import (
	"context"
	"errors"
	"terminalpathservice/internal/path"
	"terminalpathservice/internal/terminal"
)

var (
	ErrForbidden = errors.New("you are not allowed to do this action")
)

type TerminalService struct {
	terminalOps *terminal.Ops
	pathOps     *path.Ops
}

func NewTerminalService(terminalOps *terminal.Ops, pathOps *path.Ops) *TerminalService {
	return &TerminalService{
		terminalOps: terminalOps,
		pathOps:     pathOps,
	}
}

func (s *TerminalService) CreateTerminal(ctx context.Context, terminal *terminal.Terminal, isAdmin bool) error {
	if !isAdmin {
		return ErrForbidden
	}
	return s.terminalOps.Create(ctx, terminal)
}

func (s *TerminalService) GetTerminals(ctx context.Context, country, city, terminalType string, page, pageSize uint) ([]terminal.Terminal, uint, error) {
	return s.terminalOps.CityTypeTerminals(ctx, country, city, terminalType, page, pageSize)
}

func (s *TerminalService) PatchTerminal(ctx context.Context, updatedTerminal *terminal.Terminal, isAdmin bool) (*terminal.Terminal, error) {
	
	if !isAdmin{
		return nil, ErrForbidden
	}
	// exists?
	originalTerminal, err := s.terminalOps.GetTerminalByID(ctx, updatedTerminal.ID)
	if err != nil {
		return nil, err
	}
	canUpdate, err := s.pathOps.AreTherePathRelatedToTerminalID(ctx, updatedTerminal.ID)
	if err != nil {
		return nil, err
	}

	if !canUpdate {
		if updatedTerminal.Type != "" || updatedTerminal.City != "" || updatedTerminal.Country != "" {
			return nil, terminal.ErrCanNotUpdate
		}
	}
	er := s.terminalOps.PatchTerminal(ctx, updatedTerminal, originalTerminal)
	return originalTerminal, er
}

func (s *TerminalService) DeleteTerminal(ctx context.Context, terminalID uint, isAdmin bool) (*terminal.Terminal, error) {
	if !isAdmin{
		return nil, ErrForbidden
	}

	// exists?
	t, err := s.terminalOps.GetTerminalByID(ctx, terminalID)
	if err != nil {
		return nil, err
	}
	// has any related path?
	isTherePath, err := s.pathOps.IsTherePathRelatedToTerminalID(ctx, terminalID)
	if err != nil {
		return nil, err
	}
	if isTherePath {
		return nil, terminal.ErrCanNotDelete
	}
	err = s.terminalOps.Delete(ctx, terminalID)
	if err != nil {
		return nil, err
	}
	return t, nil
}
