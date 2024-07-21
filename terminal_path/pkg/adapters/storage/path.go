package storage

import (
	"context"
	"errors"
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
