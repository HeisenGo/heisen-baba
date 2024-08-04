package mappers

import (
	"tripcompanyservice/internal/invoice"
	"tripcompanyservice/pkg/adapters/storage/entities"
	"tripcompanyservice/pkg/fp"
)

func InvoiceEntityToDomain(invoiceEntity entities.Invoice) invoice.Invoice {
	//ticket := TicketEntityToDomainWithTrip(*invoiceEntity.Ticket)
	return invoice.Invoice{
		ID:             invoiceEntity.ID,
		TicketID:       invoiceEntity.TicketID,
		//Ticket:         *ticket,
		IssuedDate:     invoiceEntity.IssuedDate,
		Info:           invoiceEntity.Info,
		PerAmountPrice: invoiceEntity.PerAmountPrice,
		TotalPrice:     invoiceEntity.TotalPrice,
		Status:         invoiceEntity.Status,
		Penalty: invoiceEntity.Penalty,
	}
}

func InvoiceEntitiesToDomain(invoiceEntities []entities.Invoice) []invoice.Invoice {
	return fp.Map(invoiceEntities, InvoiceEntityToDomain)
}

func InvoiceDomainToEntity(i *invoice.Invoice) *entities.Invoice {
	//ticket := TicketDomainToEntity(i.Ticket)
	return &entities.Invoice{
		TicketID:       i.TicketID,
		IssuedDate:     i.IssuedDate,
		Info:           i.Info,
		PerAmountPrice: i.PerAmountPrice,
		TotalPrice:     i.TotalPrice,
		Penalty: i.Penalty,
		//Ticket: ticket,
	}
}

func SimpleInvoiceEntityToDomain(invoiceEntity entities.Invoice) invoice.Invoice {
	return invoice.Invoice{
		ID:             invoiceEntity.ID,
		TicketID:       invoiceEntity.TicketID,
		IssuedDate:     invoiceEntity.IssuedDate,
		Info:           invoiceEntity.Info,
		PerAmountPrice: invoiceEntity.PerAmountPrice,
		TotalPrice:     invoiceEntity.TotalPrice,
		Status:         invoiceEntity.Status,
		Penalty: invoiceEntity.Penalty,

	}
}