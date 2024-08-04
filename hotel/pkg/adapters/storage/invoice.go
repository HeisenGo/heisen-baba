package storage

import (
	"context"
	"hotel/internal/invoice"
	"hotel/pkg/adapters/storage/entities"
	"hotel/pkg/adapters/storage/mappers"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type invoiceRepo struct {
	db *gorm.DB
}

func NewInvoiceRepo(db *gorm.DB) invoice.Repo {
	return &invoiceRepo{
		db: db,
	}
}

func (r *invoiceRepo) CreateInvoice(ctx context.Context, inv *invoice.Invoice) error {
	invoiceEntity := mappers.InvoiceDomainToEntity(inv)
	if err := r.db.WithContext(ctx).Create(&invoiceEntity).Error; err != nil {
		return err
	}
	inv.ID = invoiceEntity.ID
	return nil
}

func (r *invoiceRepo) GetInvoicesByHotelID(ctx context.Context, hotelID uint, page, pageSize int) ([]invoice.Invoice, int, error) {
	var invoiceEntities []entities.Invoice
	var total int64

	query := r.db.WithContext(ctx).Model(&entities.Invoice{}).Joins("JOIN reservations ON invoices.reservation_id = reservations.id").Joins("JOIN rooms ON reservations.room_id = rooms.id").Where("rooms.hotel_id = ?", hotelID)

	query.Count(&total)

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&invoiceEntities).Error; err != nil {
		return nil, 0, err
	}

	invoices := mappers.BatchInvoiceEntitiesToDomain(invoiceEntities)
	return invoices, int(total), nil
}

func (r *invoiceRepo) GetInvoicesByUserID(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]invoice.Invoice, int, error) {
	var invoiceEntities []entities.Invoice
	var total int64

	query := r.db.WithContext(ctx).Model(&entities.Invoice{}).Joins("JOIN reservations ON invoices.reservation_id = reservations.id").Joins("JOIN rooms ON reservations.room_id = rooms.id").Where("rooms.hotel_id = ?", userID)

	query.Count(&total)

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&invoiceEntities).Error; err != nil {
		return nil, 0, err
	}

	invoices := mappers.BatchInvoiceEntitiesToDomain(invoiceEntities)
	return invoices, int(total), nil
}

func (r *invoiceRepo) GetInvoiceByID(ctx context.Context, id uint) (*invoice.Invoice, error) {
	var invoiceEntity entities.Invoice
	if err := r.db.WithContext(ctx).First(&invoiceEntity, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, invoice.ErrRecordNotFound
		}
		return nil, err
	}
	in :=mappers.InvoiceEntityToDomain(invoiceEntity)
	return &in , nil
}

func (r *invoiceRepo) UpdateInvoice(ctx context.Context, inv *invoice.Invoice) error {
	invoiceEntity := mappers.InvoiceDomainToEntity(inv)
	if err := r.db.WithContext(ctx).Model(&entities.Invoice{}).Where("id = ?", inv.ID).Updates(invoiceEntity).Error; err != nil {
		return err
	}
	return nil
}

func (r *invoiceRepo) DeleteInvoice(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&entities.Invoice{}, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return invoice.ErrRecordNotFound
		}
		return err
	}
	return nil
}