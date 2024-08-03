package storage

import (
	"context"
	"fmt"
	"strings"
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

func (r *ticketRepo) GetFullTicketByID(ctx context.Context, id uint) (*ticket.Ticket, error) {
	var eT entities.Ticket

	if err := r.db.WithContext(ctx).
		Preload("Trip").                  // Preload related Trip
		Preload("Trip.TransportCompany").
		Preload("Trip.TripCancellingPenalty"). // Preload TransportCompany within Trip
		Preload("Ticket").                // Preload related Invoice
		First(&eT, id).Error; err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, nil
		}
		return nil, fmt.Errorf("%w %w %d", ticket.ErrFailedToGetTicket, err, id)
	}
	dT := mappers.TicketEntityToDomainWithTripWithCompanyWithPenalty(eT)

	return &dT, nil
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
	dTickets := mappers.TicketEntitiesToDomainWithTrips(tickets)
	return dTickets, uint(total), nil
}

func (r *ticketRepo) UpdateTicket(ctx context.Context, id uint, updates map[string]interface{}) error {
	var t entities.Ticket

	if err := r.db.WithContext(ctx).Model(&t).Updates(updates).Error; err != nil {
		return fmt.Errorf("%w %w", ticket.ErrFailedToUpdate, err)
	}

	return nil
}