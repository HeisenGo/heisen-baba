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

func (r *ticketRepo) UpdateTicketStatus(ctx context.Context, ticketID uint, status string) error {
	return r.db.WithContext(ctx).Model(&entities.Ticket{}).Where("id = ?", ticketID).Update("status", status).Error
}

func (r *ticketRepo) GetTicketsByUserOrAgency(ctx context.Context, userID *uint, agencyID *uint, limit, offset uint) ([]ticket.Ticket, uint, error) {
	query := r.db.WithContext(ctx).Model(&entities.Ticket{}).
		Preload("Trip")

	if userID != nil {
		query = query.Where("user_id = ?", userID)
	}
	if agencyID != nil {
		query = query.Where("agency_id = ?", agencyID)
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

	var tickets []entities.Ticket

	if err := query.Find(&tickets).Error; err != nil {
		return nil, 0, err
	}
	dTickets := mappers.TicketEntitiesToDomain(tickets)
	return dTickets, uint(total), nil
}
