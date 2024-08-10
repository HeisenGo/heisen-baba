package mappers

import (
	"agency/internal/invoice"
	"agency/pkg/adapters/storage/entities"
	"agency/pkg/fp"

	"gorm.io/gorm"
)

// Invoice Mappers
func InvoiceEntityToDomain(invoiceEntity entities.Invoice) invoice.Invoice {
	return invoice.Invoice{
		ID:            invoiceEntity.ID,
		ReservationID: invoiceEntity.ReservationID,
		IssueDate:     invoiceEntity.IssueDate,
		Amount:        invoiceEntity.Amount,
		Paid:          invoiceEntity.Paid,
		UserID:        invoiceEntity.UserID,
		OwnerID:       invoiceEntity.OwnerID,
	}
}

func BatchInvoiceEntitiesToDomain(invoiceEntities []entities.Invoice) []invoice.Invoice {
	return fp.Map(invoiceEntities, InvoiceEntityToDomain)
}

func InvoiceDomainToEntity(inv *invoice.Invoice) entities.Invoice {
	return entities.Invoice{
		Model: gorm.Model{
			ID: inv.ID,
		},
		ReservationID: inv.ReservationID,
		IssueDate:     inv.IssueDate,
		Amount:        inv.Amount,
		Paid:          inv.Paid,
		UserID:        inv.UserID,
		OwnerID:       inv.OwnerID,
	}
}
