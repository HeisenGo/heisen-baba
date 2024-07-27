package mappers

import (
	"terminalpathservice/internal/path"
	"terminalpathservice/internal/terminal"
	"terminalpathservice/pkg/adapters/storage/entities"
	"terminalpathservice/pkg/fp"
)

func PathEntityToDomain(pathEntity entities.Path) path.Path {
	fromTerminal := TerminalEntityToDomain(pathEntity.FromTerminal)
	toTerminal := TerminalEntityToDomain(pathEntity.ToTerminal)
	return path.Path{
		ID:             pathEntity.ID,
		Name:           pathEntity.Name,
		Code:           pathEntity.Code,
		DistanceKM:     pathEntity.DistanceKM,
		ToTerminalID:   pathEntity.ToTerminalID,
		ToTerminal:     &toTerminal,
		FromTerminalID: pathEntity.FromTerminalID,
		FromTerminal:   &fromTerminal,
		Type:           terminal.TerminalType(pathEntity.Type),
	}
}

func PathEntitiesToDomain(pathEntities []entities.Path) []path.Path {

	return fp.Map(pathEntities, PathEntityToDomain)
}

func PathDomainToEntity(p *path.Path) *entities.Path {
	return &entities.Path{
		ToTerminalID:   p.ToTerminalID,
		FromTerminalID: p.FromTerminalID,
		Name:           p.Name,
		Code:           p.Code,
		DistanceKM:     p.DistanceKM,
		Type:           string(p.Type),
	}
}
