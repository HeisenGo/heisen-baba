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
			return terminal.ErrDuplication
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
	canUpdate, err := r.canUpdateTerminal(ctx, updatedTerminal.ID)
	if err != nil {
		return err
	}

	if !canUpdate {
		if updatedTerminal.Type != "" || updatedTerminal.City != "" || updatedTerminal.Country != "" {
			return terminal.ErrCanNotUpdate
		}
	}

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

func (r *terminalRepo) canUpdateTerminal(ctx context.Context, terminalID uint) (bool, error) {
	var count int64

	// Count paths where the terminal is either the start or the end terminal
	if err := r.db.WithContext(ctx).Model(&entities.Path{}).
		Where("from_terminal_id = ? OR to_terminal_id = ?", terminalID, terminalID).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to count paths: %v", err)
	}

	// If there are no such paths, the terminal can be updated
	if count == 0 {
		return true, nil
	}

	// If there are paths, determine if the terminal can be updated based on business logic
	return false, nil
}
