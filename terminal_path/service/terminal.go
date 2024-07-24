package service

import (
	"context"
	"terminalpathservice/internal/terminal"
)

type TerminalService struct {
	terminalOps *terminal.Ops
}

func NewTerminalService(terminalOps *terminal.Ops) *TerminalService {
	return &TerminalService{
		terminalOps: terminalOps,
	}
}

func (s *TerminalService) CreateTerminal(ctx context.Context, terminal *terminal.Terminal) error {
	return s.terminalOps.Create(ctx, terminal)
}

func (s *TerminalService) GetTerminalsOfCity(ctx context.Context, country, city, terminalType string, page, pageSize uint) ([]terminal.Terminal, uint, error) {
	return s.terminalOps.CityTypeTerminals(ctx, country, city, terminalType, page, pageSize)
}

func (s *TerminalService) PatchTerminal(ctx context.Context, updatedTerminal *terminal.Terminal) (*terminal.Terminal, error) {
	originalTerminal, err := s.terminalOps.GetTerminalByID(ctx, updatedTerminal.ID)
	if err != nil {
		return nil, err
	}
	er := s.terminalOps.PatchTerminal(ctx, updatedTerminal, originalTerminal)
	return originalTerminal, er
}

func (s *TerminalService) DeleteTerminal(ctx context.Context, terminalID uint) (*terminal.Terminal, error) {

	terminal, err := s.terminalOps.GetTerminalByID(ctx, terminalID)
	if err != nil {
		return nil, err
	}
	err = s.terminalOps.Delete(ctx, terminalID)
	if err != nil {
		return nil, err
	}
	return terminal, nil
}
