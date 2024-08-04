package service

import (
	"context"
	"hotel/internal/invoice"

	"github.com/google/uuid"
)


type InvoiceService struct {
	invoiceOps *invoice.Ops
}

func NewInvoiceService(invoiceOps *invoice.Ops) *InvoiceService {
	return &InvoiceService{
		invoiceOps: invoiceOps,
	}
}

func (s *InvoiceService) CreateInvoice(ctx context.Context, inv *invoice.Invoice) error {
	return s.invoiceOps.Create(ctx, inv)
}

func (s *InvoiceService) GetInvoiceByID(ctx context.Context, id uint) (*invoice.Invoice, error) {
	inv, err := s.invoiceOps.GetInvoiceByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if inv == nil {
		return nil, invoice.ErrRecordNotFound
	}
	return inv, nil
}

func (s *InvoiceService) GetInvoicesByHotelID(ctx context.Context, hotelID uint, page, pageSize int) ([]invoice.Invoice, int, error) {
	invoices, total, err := s.invoiceOps.GetInvoicesByHotelID(ctx, hotelID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	if len(invoices) == 0 {
		return nil, 0, invoice.ErrRecordNotFound
	}
	return invoices, total, nil
}

func (s *InvoiceService) GetInvoicesByUserID(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]invoice.Invoice, int, error) {
	invoices, total, err := s.invoiceOps.GetInvoicesByUserID(ctx, userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	if len(invoices) == 0 {
		return nil, 0, invoice.ErrRecordNotFound
	}
	return invoices, total, nil
}

func (s *InvoiceService) UpdateInvoice(ctx context.Context, id uint, updates *invoice.Invoice) error {
	existingInvoice, err := s.invoiceOps.GetInvoiceByID(ctx, id)
	if err != nil {
		return err
	}
	if existingInvoice == nil {
		return invoice.ErrRecordNotFound
	}

	// Update only the fields that are provided
	if updates.ReservationID != 0 {
		existingInvoice.ReservationID = updates.ReservationID
	}
	if !updates.IssueDate.IsZero() {
		existingInvoice.IssueDate = updates.IssueDate
	}
	if updates.Amount != 0 {
		existingInvoice.Amount = updates.Amount
	}
	existingInvoice.Paid = updates.Paid

	return s.invoiceOps.Update(ctx, existingInvoice)
}

func (s *InvoiceService) DeleteInvoice(ctx context.Context, id uint) error {
	_, err := s.invoiceOps.GetInvoiceByID(ctx, id)
	if err != nil {
		return err
	}
	return s.invoiceOps.Delete(ctx, id)
}
