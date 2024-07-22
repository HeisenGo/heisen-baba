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
	return s.terminalOps.CityTypeTerminals(ctx, country,city, terminalType, page, pageSize)
}
