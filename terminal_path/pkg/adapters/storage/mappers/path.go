package mappers

import (
	"terminalpathservice/internal/path"
	"terminalpathservice/pkg/adapters/storage/entities"
	"terminalpathservice/pkg/fp"
)

func PathEntityToDomain(pathEntity entities.Path) path.Path {
	return path.Path{
		ID:      pathEntity.ID,
		Name:    pathEntity.Name,
		Code: pathEntity.Code,
		Distance: pathEntity.Distance,
		ToTerminalID: pathEntity.ToTerminalID,
		FromTerminalID: pathEntity.FromTerminalID,
	}
}

func PathEntitiesToDomain(pathEntities []entities.Path) []path.Path {
	return fp.Map(pathEntities, PathEntityToDomain)
}

func PathDomainToEntity(p *path.Path) *entities.Path {
	return &entities.Path{
		ToTerminalID: p.ToTerminalID,
		FromTerminalID: p.FromTerminalID,
		Name:    p.Name,
		Code: p.Code,
		Distance: p.Distance,
	}
}
