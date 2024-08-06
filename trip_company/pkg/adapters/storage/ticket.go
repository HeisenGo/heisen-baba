package storage

import (
	"context"
	"fmt"
	"strings"
	"tripcompanyservice/internal/ticket"
	"tripcompanyservice/pkg/adapters/storage/entities"
	"tripcompanyservice/pkg/adapters/storage/mappers"

	"github.com/google/uuid"
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
		Preload("Trip.TripCancelingPenalty"). // Preload TransportCompany within Trip
		Preload("Invoice").                // Preload related Invoice
		First(&eT, id).Error; err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, nil
		}
		return nil, fmt.Errorf("%w %w %d", ticket.ErrFailedToGetTicket, err, id)
	}
	dT := mappers.TicketEntityToDomainWithTripWithCompanyWithPenaltyWithInvoice(eT)

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

func (r *ticketRepo) GetTicketsByUserOrAgency(ctx context.Context, userID *uuid.UUID, agencyID *uint, limit, offset uint) ([]ticket.Ticket, uint, error) {
	query := r.db.WithContext(ctx).Model(&entities.Ticket{}).
		Preload("Trip").Preload("Trip.TripCancelingPenalty").Preload("Invoice")

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

	if err := r.db.WithContext(ctx).Model(&t).Where("id = ?", id).Updates(updates).Error; err != nil {
		return fmt.Errorf("%w %w", ticket.ErrFailedToUpdate, err)
	}

	return nil
}

func (r *ticketRepo) GetTicketsWithInvoicesByTripID(ctx context.Context, tripID uint) ([]ticket.Ticket, error) {
	var tickets []entities.Ticket

	err := r.db.WithContext(ctx).
		Preload("Invoice").       // Preload the related Invoice
		Where("trip_id = ?", tripID). // Filter by trip_id
		Find(&tickets).Error

	if err != nil {
		return nil, err
	}
	dTickets := mappers.BatchTicketEntitiesToTickets(tickets)
	return dTickets, nil
}


