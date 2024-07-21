package presenter

import (
	"terminalpathservice/internal/terminal"
	"terminalpathservice/pkg/fp"
)

type TerminalRequest struct {
	ID      uint   `json:"terminal_id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	City    string `json:"city"`
	Country string `json:"country"`
}

func TerminalToTerminalRequest(t terminal.Terminal) TerminalRequest {
	return TerminalRequest{
		ID:      t.ID,
		Name:    t.Name,
		Type:    string(t.Type),
		City:    t.City,
		Country: t.Country,
	}
}

func TerminalsToTerminalRequest(terminals []terminal.Terminal) []TerminalRequest {
	return fp.Map(terminals, TerminalToTerminalRequest)
}

func TerminalRequestToTerminal(terminalReq *TerminalRequest) *terminal.Terminal {
	return &terminal.Terminal{
		ID:      terminalReq.ID,
		Name:    terminalReq.Name,
		Type:    terminal.TerminalType(terminalReq.Type),
		City:    terminalReq.City,
		Country: terminalReq.Country,
	}
}
