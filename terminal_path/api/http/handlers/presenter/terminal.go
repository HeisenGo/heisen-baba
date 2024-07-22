package presenter

import (
	"terminalpathservice/internal/terminal"
	"terminalpathservice/pkg/fp"
)

type TerminalRequest struct {
	ID      uint   `json:"terminal_id"`
	Name    string `json:"name" validate:"required"`
	Type    string `json:"type" validate:"required"`
	City    string `json:"city" validate:"required"`
	Country string `json:"country" validate:"required"`
}


func TerminalToTerminalGetResponse(t terminal.Terminal) TerminalGetResponse {
	return TerminalGetResponse{
		ID:      t.ID,
		Name:    t.Name,
		Type:    string(t.Type),
		City:    t.City,
	}
}

func TerminalsToTerminalResponse(terminals []terminal.Terminal) []TerminalGetResponse {
	return fp.Map(terminals, TerminalToTerminalGetResponse)
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

type SearchTerminalByCityTypeReq struct {
	Type string `json:"type"`
	City string `json:"city" validate:"required"` 
}

type SearchTerminalByCityTypeRes struct{
	ID uint `json:"terminal_id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type TerminalGetResponse struct {
	ID      uint   `json:"terminal_id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	City    string `json:"city"`
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