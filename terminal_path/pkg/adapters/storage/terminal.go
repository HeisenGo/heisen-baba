package storage

import (
	"context"
	"errors"
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
	if err := r.db.WithContext(ctx).Save(&terminalEntity).Error; err != nil {
		return err
	}

	t.ID = terminalEntity.ID

	return nil
}

func (r *terminalRepo) GetByID(ctx context.Context, id uint) (*terminal.Terminal, error) {
	var t entities.Terminal

	err := r.db.WithContext(ctx).Model(&entities.Terminal{}).Where("id = ?", id).First(&t).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	dTerminal := mappers.TerminalEntityToDomain(t)
	return &dTerminal, nil
}
