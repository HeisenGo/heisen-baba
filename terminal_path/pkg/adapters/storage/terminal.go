package storage

import (
	"context"
	"fmt"
	"strings"
	"terminalpathservice/internal/terminal"
	"terminalpathservice/pkg/adapters/storage/entities"
	"terminalpathservice/pkg/adapters/storage/mappers"

	"gorm.io/gorm"
)

type terminalRepo struct {
	db *gorm.DB
}

func NewTerminalRepo(db *gorm.DB) terminal.Repo {
	return &terminalRepo{db}
}

func (r *terminalRepo) Insert(ctx context.Context, t *terminal.Terminal) error {
	terminalEntity := mappers.TerminalDomainToEntity(t)

	result := r.db.WithContext(ctx).Save(&terminalEntity)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {

			var existingTerminal entities.Terminal
			// Search for the soft-deleted record with the same unique constraints
			if r.db.WithContext(ctx).Unscoped().Where("name = ? AND type = ? AND city = ? AND country = ?", t.Name, t.Type, t.City, t.Country).First(&existingTerminal).Error == nil {
				// Check if the record is soft-deleted
				if existingTerminal.DeletedAt.Valid {
					// Restore the soft-deleted record
					existingTerminal.DeletedAt = gorm.DeletedAt{}
					if err := r.db.WithContext(ctx).Save(&existingTerminal).Error; err != nil {
						return fmt.Errorf("%w %w", terminal.ErrFailedToRestore, err)
					}
					t.ID = existingTerminal.ID
					return nil
				}

				return terminal.ErrDuplication
			}
		}
		return result.Error
	}
	t.ID = terminalEntity.ID

	return nil
}

func (r *terminalRepo) GetByID(ctx context.Context, id uint) (*terminal.Terminal, error) {
	var t entities.Terminal

	err := r.db.WithContext(ctx).Model(&entities.Terminal{}).Where("id = ?", id).First(&t).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, nil
		}
		return nil, err
	}
	dTerminal := mappers.TerminalEntityToDomain(t)
	return &dTerminal, nil
}

func (r *terminalRepo) GetTerminalsByCityAndType(ctx context.Context, country, city, terminalType string, limit, offset uint) ([]terminal.Terminal, uint, error) {
	var query *gorm.DB
	if city != "" && terminalType != "" {
		query = r.db.WithContext(ctx).Model(&entities.Terminal{}).Where("city = ? AND type=? AND country=?", city, terminalType, country)
	} else if city != "" && terminalType == "" {
		query = r.db.WithContext(ctx).Model(&entities.Terminal{}).Where("city = ? AND country=?", city, country)
	} else if city == "" && terminalType != "" {
		query = r.db.WithContext(ctx).Model(&entities.Terminal{}).Where("type=? AND country=?", terminalType, country)
	} else if city == "" && terminalType == "" {
		query = r.db.WithContext(ctx).Model(&entities.Terminal{}).Where("country=?", country)
	}
	var total int64

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if offset > 0 {
		query = query.Offset(int(offset))
	}

	if limit > 0 {
		query = query.Limit(int(limit))
	}

	var terminals []entities.Terminal

	if err := query.Find(&terminals).Error; err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, 0, terminal.ErrRecordsNotFound
		}
		return nil, 0, err
	}

	return mappers.TerminalEntitiesToDomain(terminals), uint(total), nil
}

func (r *terminalRepo) PatchTerminal(ctx context.Context, updatedTerminal, originalTerminal *terminal.Terminal) error {
	// Prepare a map to hold the fields to be updated
	updates := make(map[string]interface{})

	// Add fields to the map if they are provided in the request
	if updatedTerminal.Name != "" {
		updates["name"] = updatedTerminal.Name
		originalTerminal.Name = updatedTerminal.Name
		updates["normalized_name"] = entities.NormalizeName(updates["name"].(string))
	}

	if updatedTerminal.Type != "" {
		updates["type"] = updatedTerminal.Type
		originalTerminal.Type = updatedTerminal.Type
	}
	if updatedTerminal.City != "" {
		updates["city"] = updatedTerminal.City
		originalTerminal.City = updatedTerminal.City
	}
	if updatedTerminal.Country != "" {
		updates["country"] = updatedTerminal.Country
		originalTerminal.Country = updatedTerminal.Country
	}

	if err := r.db.Model(&entities.Terminal{}).Where("id = ?", updatedTerminal.ID).Updates(updates).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return terminal.ErrDuplication
		}
		return terminal.ErrFailedToUpdate
	}

	return nil
}

func (r *terminalRepo) Delete(ctx context.Context, terminalID uint) error {
	// check if there jis a path related to this terminal in business logic

	// Delete the terminal
	if err := r.db.WithContext(ctx).Delete(&entities.Terminal{}, terminalID).Error; err != nil {
		return fmt.Errorf("%w %w", terminal.ErrDeleteTerminal, err)
	} else {
		return nil
	}

}
