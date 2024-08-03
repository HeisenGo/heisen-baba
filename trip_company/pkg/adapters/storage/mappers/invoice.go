package mappers

import (
	"tripcompanyservice/internal/invoice"
	"tripcompanyservice/pkg/adapters/storage/entities"
	"tripcompanyservice/pkg/fp"
)

func InvoiceEntityToDomain(invoiceEntity entities.Invoice) invoice.Invoice {
	ticket := TicketEntityToDomain(invoiceEntity.Ticket)
	return invoice.Invoice{
		ID:             invoiceEntity.ID,
		TicketID:       invoiceEntity.TicketID,
		Ticket:         ticket,
		IssuedDate:     invoiceEntity.IssuedDate,
		Info:           invoiceEntity.Info,
		PerAmountPrice: invoiceEntity.PerAmountPrice,
		TotalPrice:     invoiceEntity.TotalPrice,
	}
}

func InvoiceEntitiesToDomain(invoiceEntities []entities.Invoice) []invoice.Invoice {
	return fp.Map(invoiceEntities, InvoiceEntityToDomain)
}

func InvoiceDomainToEntity(i invoice.Invoice) *entities.Invoice {
	//ticket := TicketDomainToEntity(i.Ticket)
	return &entities.Invoice{
		TicketID:       i.TicketID,
		IssuedDate:     i.IssuedDate,
		Info:           i.Info,
		PerAmountPrice: i.PerAmountPrice,
		TotalPrice:     i.TotalPrice,
	}
}
