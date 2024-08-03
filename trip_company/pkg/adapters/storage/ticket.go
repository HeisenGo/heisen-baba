package storage

import (
	"context"
	"tripcompanyservice/internal/ticket"
	"tripcompanyservice/pkg/adapters/storage/entities"
	"tripcompanyservice/pkg/adapters/storage/mappers"

	"gorm.io/gorm"
)

type ticketRepo struct {
	db *gorm.DB
}

func NewTicketRepo(db *gorm.DB) ticket.Repo {
	return &ticketRepo{db}
}

func (r *ticketRepo) Insert(ctx context.Context, t *ticket.Ticket) error {
	ticketEntity := mappers.TicketDomainToEntity(t)
	result := r.db.WithContext(ctx).Create(&ticketEntity)
	if result.Error != nil {
		return result.Error
	}

	t.ID = ticketEntity.ID

	return nil

}

func (r *ticketRepo)  UpdateTicketStatus(ctx context.Context, ticketID uint, status string) error {
    return r.db.WithContext(ctx).Model(&entities.Ticket{}).Where("id = ?", ticketID).Update("status", status).Error
}
