package storage

import (
	"context"
	"fmt"
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
		if strings.Contains(err.Error(), "duplicate key") {
			var existingPath entities.Path

			// Check for an existing soft-deleted path with the same Code
			if err := r.db.Unscoped().
				Where("code = ?", pathEntity.Code).
				First(&existingPath).Error; err == nil {

				// Check if the record is soft-deleted
				if existingPath.DeletedAt.Valid {
					// Compare other fields
					if existingPath.FromTerminalID == pathEntity.FromTerminalID &&
						existingPath.ToTerminalID == pathEntity.ToTerminalID &&
						existingPath.Type == pathEntity.Type {

						// Restore the soft-deleted record
						existingPath.DeletedAt = gorm.DeletedAt{}
						if err := r.db.Save(&existingPath).Error; err != nil {
							return fmt.Errorf("%w %w", path.ErrFailedToRestore, err)
						}
						p.ID = existingPath.ID
						return nil
					} else {
						// Fields do not match, return a duplication error
						return path.ErrCodeIsImpossibleToUse
					}
				}
			}
			return path.ErrDuplication
		}
		return err
	}

	p.ID = pathEntity.ID
	return nil
}

func (r *pathRepo) GetByID(ctx context.Context, id uint) (*path.Path, error) {
	var p entities.Path

	err := r.db.WithContext(ctx).Model(&entities.Path{}).Where("id = ?", id).First(&p).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, nil
		}
		return nil, err
	}
	dPath := mappers.PathEntityToDomain(p)
	return &dPath, nil
}

func (r *pathRepo) GetFullPathByID(ctx context.Context, id uint) (*path.Path, error) {
	var p entities.Path
	if err := r.db.WithContext(ctx).
		Preload("FromTerminal").
		Preload("ToTerminal").
		First(&p, id).Error; err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, nil
		}
		return nil, fmt.Errorf("%w %w", path.ErrFailedToGetPath, err)
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

func (r *pathRepo) PatchPath(ctx context.Context, updatedPath, originalPath *path.Path) error {
	// Prepare a map to hold the fields to be updated
	updates := make(map[string]interface{})
	if updatedPath.FromTerminalID != uint(0) {
		updates["from_terminal_id"] = updatedPath.FromTerminalID
	}
	if updatedPath.ToTerminalID != uint(0) {
		updates["to_terminal_id"] = updatedPath.ToTerminalID
	}
	if updatedPath.Type != "" {
		updates["type"] = updatedPath.Type
		originalPath.Type = updatedPath.Type
	}
	// Add fields to the map if they are provided in the request
	if updatedPath.Name != "" {
		updates["name"] = updatedPath.Name
		originalPath.Name = updatedPath.Name
	}

	if updatedPath.Code != "" {
		updates["code"] = updatedPath.Code
		originalPath.Code = updatedPath.Code
	}

	if updatedPath.DistanceKM != 0 {
		updates["distance_km"] = updatedPath.DistanceKM
		originalPath.DistanceKM = updatedPath.DistanceKM
	}
	if err := r.db.Model(&entities.Path{}).Where("id = ?", updatedPath.ID).Updates(updates).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return path.ErrDuplication
		}
		return path.ErrFailedToUpdate
	}

	return nil
}

func (r *pathRepo) Delete(ctx context.Context, pathID uint) error {
	// if has unfinished trip can not be deleted this logic checked in ops layer

	// Delete the path
	if err := r.db.WithContext(ctx).Delete(&entities.Path{}, pathID).Error; err != nil {
		return fmt.Errorf("%w %w", path.ErrDeletePath, err)

	}
	return nil

}
