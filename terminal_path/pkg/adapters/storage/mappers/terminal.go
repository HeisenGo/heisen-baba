package mappers

import (
	"terminalpathservice/internal/terminal"
	"terminalpathservice/pkg/adapters/storage/entities"
	"terminalpathservice/pkg/fp"
)

func TerminalEntityToDomain(terminalEntity entities.Terminal) terminal.Terminal {
	return terminal.Terminal{
		ID:      terminalEntity.ID,
		Name:    terminalEntity.Name,
		Type:    terminal.TerminalType(terminalEntity.Type),
		City:    terminalEntity.City,
		Country: terminalEntity.Country,
	}
}

func TerminalEntitiesToDomain(terminalEntities []entities.Terminal) []terminal.Terminal {
	return fp.Map(terminalEntities, TerminalEntityToDomain)
}

func TerminalDomainToEntity(t *terminal.Terminal) *entities.Terminal {
	return &entities.Terminal{
		City:    t.City,
		Name:    t.Name,
		Type:    string(t.Type),
		Country: t.Country,
	}
}
