package service

import (
	"context"
	"tripcompanyservice/internal/invoice"
	"tripcompanyservice/internal/ticket"
	"tripcompanyservice/internal/trip"

	"github.com/google/uuid"
)

type InvoiceService struct {
	invoiceOps *invoice.Ops
	ticketOps  *ticket.Ops
	tripOps    *trip.Ops
}

func NewInvoiceService(invoiceOps *invoice.Ops, ticketOps *ticket.Ops, tripOps *trip.Ops) *InvoiceService {
	return &InvoiceService{
		invoiceOps: invoiceOps,
		ticketOps:  ticketOps,
		tripOps:    tripOps,
	}
}

func (s *InvoiceService) GetInvoicesByUserOrAgency(ctx context.Context, userID *uuid.UUID, agencyID *uint, page, pageSize uint) ([]invoice.Invoice,uint, error) {
	// check one of them should be nill !!!
	return s.invoiceOps.GetInvoicesByUserOrAgency(ctx, userID, agencyID, page, pageSize)
}
