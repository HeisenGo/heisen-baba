package storage

import (
	"context"
	"errors"
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, terminal.ErrRecordsNotFound
		}
		return nil, 0, err
	}

	return mappers.TerminalEntitiesToDomain(terminals), uint(total), nil
}
