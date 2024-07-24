package storage

import (
	"context"
	"errors"
	"strings"
	"terminalpathservice/internal/path"
	"terminalpathservice/pkg/adapters/storage/entities"
	"terminalpathservice/pkg/adapters/storage/mappers"

	"gorm.io/gorm"
)

type pathRepo struct {
	db *gorm.DB
}

func NewPathRepo(db *gorm.DB) path.Repo {
	return &pathRepo{db}
}

func (r *pathRepo) Insert(ctx context.Context, p *path.Path) error {
	pathEntity := mappers.PathDomainToEntity(p)
	if err := r.db.WithContext(ctx).Save(&pathEntity).Error; err != nil {
		return err
	}

	p.ID = pathEntity.ID

	return nil
}

func (r *pathRepo) GetByID(ctx context.Context, id uint) (*path.Path, error) {
	var p entities.Path

	err := r.db.WithContext(ctx).Model(&entities.Path{}).Where("id = ?", id).First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	dPath := mappers.PathEntityToDomain(p)
	return &dPath, nil
}

func (r *pathRepo) GetPathsByOriginDestinationType(ctx context.Context, originCity, destinationCity, pathType string, limit, offset uint) ([]path.Path, uint, error) {
	var paths []entities.Path
	var total int64

	query := r.db.WithContext(ctx).Model(&entities.Path{}).
		Preload("FromTerminal").
		Preload("ToTerminal")

		// Build the query based on provided parameters
	if originCity != "" {
		query = query.Joins("JOIN terminals AS from_terminal ON paths.from_terminal_id = from_terminal.id").
			Where("from_terminal.city = ?", originCity)
	}
	if destinationCity != "" {
		query = query.Joins("JOIN terminals AS to_terminal ON paths.to_terminal_id = to_terminal.id").
			Where("to_terminal.city = ?", destinationCity)
	}
	if pathType != "" {
		query = query.Where("paths.type = ?", pathType)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if offset > 0 {
		query = query.Offset(int(offset))
	}
	if limit > 0 {
		query = query.Limit(int(limit))
	}

	if err := query.Find(&paths).Error; err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, 0, path.ErrRecordsNotFound
		}
		return nil, 0, err
	}
	print(paths)
	return mappers.PathEntitiesToDomain(paths), uint(total), nil
}
