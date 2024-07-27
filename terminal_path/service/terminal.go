package service

import (
	"context"
	"terminalpathservice/internal/path"
	"terminalpathservice/internal/terminal"
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

func (s *TerminalService) CreateTerminal(ctx context.Context, terminal *terminal.Terminal) error {
	return s.terminalOps.Create(ctx, terminal)
}

func (s *TerminalService) GetTerminals(ctx context.Context, country, city, terminalType string, page, pageSize uint) ([]terminal.Terminal, uint, error) {
	return s.terminalOps.CityTypeTerminals(ctx, country, city, terminalType, page, pageSize)
}

func (s *TerminalService) PatchTerminal(ctx context.Context, updatedTerminal *terminal.Terminal) (*terminal.Terminal, error) {
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

func (s *TerminalService) DeleteTerminal(ctx context.Context, terminalID uint) (*terminal.Terminal, error) {
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
