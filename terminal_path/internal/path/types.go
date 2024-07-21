package path

import "terminalpathservice/internal/terminal"

type Path struct {
	ID	uint
	FromTerminalID uint         
	ToTerminalID   uint          
	FromTerminal   terminal.Terminal     
	ToTerminal     terminal.Terminal
	Distance       float64        // in kilometers
	Code           string       
	Name           string       
}