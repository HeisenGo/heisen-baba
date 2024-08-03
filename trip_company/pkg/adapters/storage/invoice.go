package storage

import (
	"context"
	"tripcompanyservice/internal/invoice"
	"tripcompanyservice/pkg/adapters/storage/entities"
	"tripcompanyservice/pkg/adapters/storage/mappers"

	"gorm.io/gorm"
)

type invoiceRepo struct {
	db *gorm.DB
}

func NewInvoiceRepo(db *gorm.DB) invoice.Repo {
	return &invoiceRepo{db}
}

func (r *invoiceRepo) Insert(ctx context.Context, i *invoice.Invoice) error {
	iEntity := mappers.InvoiceDomainToEntity(i)

	// Create the new company record
	result := r.db.WithContext(ctx).Create(&iEntity)
	if result.Error != nil {
		return result.Error
	}

	i.ID = iEntity.ID

	return nil

}

func (r *invoiceRepo) UpdateInvoiceStatus(ctx context.Context, invoiceID uint, status string) error {
	// Ensure status is valid, e.g., 'paid', 'pending', etc.
	return r.db.WithContext(ctx).Model(&entities.Invoice{}).Where("id = ?", invoiceID).Update("status", status).Error
}

func (r *invoiceRepo) GetInvoicesByUserOrAgency(ctx context.Context, userID *uint, agencyID *uint, limit, offset uint) ([]invoice.Invoice, uint,error) {
	query := r.db.WithContext(ctx).Model(&entities.Invoice{}).
		Preload("Ticket").
		Preload("Ticket.Trip")

		if userID != nil {
			query = query.Joins("JOIN tickets ON tickets.id = invoices.ticket_id").
				Where("tickets.user_id = ?", userID)
		}
		if agencyID != nil {
			query = query.Joins("JOIN tickets ON tickets.id = invoices.ticket_id").
				Where("tickets.agency_id = ?", agencyID)
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
	
		var invoices []entities.Invoice
	
		if err := query.Find(&invoices).Error; err != nil {
			return nil, 0, err
		}
	dInvoices := mappers.InvoiceEntitiesToDomain(invoices)
	return dInvoices, uint(total), nil
}
