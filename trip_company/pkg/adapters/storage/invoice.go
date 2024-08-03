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

func  (r *invoiceRepo)UpdateInvoiceStatus(ctx context.Context, invoiceID uint, status string) error {
    // Ensure status is valid, e.g., 'paid', 'pending', etc.
    return r.db.WithContext(ctx).Model(&entities.Invoice{}).Where("id = ?", invoiceID).Update("status", status).Error
}
