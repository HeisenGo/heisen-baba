package storage

import (
	creditCard "bankservice/internal/wallet/credit_card"
	wallet "bankservice/internal/wallet/wallet_credit_card"
	"bankservice/pkg/adapters/storage/entities"
	"bankservice/pkg/adapters/storage/mappers"
	"context"
	"github.com/google/uuid"
	"strings"

	"gorm.io/gorm"
)

type creditCardRepo struct {
	db *gorm.DB
}

func NewCreditCardRepo(db *gorm.DB) creditCard.Repo {
	return &creditCardRepo{
		db: db,
	}
}
func (r *creditCardRepo) CreateCardAndAddToWallet(ctx context.Context, card *creditCard.CreditCard, userID uuid.UUID) (*creditCard.CreditCard, error) {
	var userWalletEntity *entities.Wallet
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&userWalletEntity).Error; err != nil {
		return nil, err
	}
	newCreditCard := mappers.CreditCardDomainToEntity(card)
	if err := r.db.WithContext(ctx).Where("number = ?", newCreditCard.Number).First(&newCreditCard).Error; err != nil {
		if err = r.db.Create(&newCreditCard).Error; err != nil {
			return nil, err
		}
	}
	walletCreditCardEntity := wallet.NewWalletCreditCard(userWalletEntity.ID, newCreditCard.ID)
	if err := r.db.WithContext(ctx).Create(&walletCreditCardEntity).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil, creditCard.ErrCardAlreadyExists
		}
		return nil, err
	}
	createdCreditCard := mappers.CreditCardEntityToDomain(newCreditCard)
	return createdCreditCard, nil
}

func (r *creditCardRepo) GetUserWalletCards(ctx context.Context, userID uuid.UUID) ([]creditCard.CreditCard, error) {
	var creditCardEntities []*entities.CreditCard

	err := r.db.WithContext(ctx).Joins("JOIN wallet_credit_cards ON wallet_credit_cards.credit_card_id = credit_cards.id").
		Joins("JOIN wallets ON wallets.id = wallet_credit_cards.wallet_id").
		Where("wallets.user_id = ?", userID).
		Find(&creditCardEntities).Error

	if err != nil {
		return nil, err
	}
	allDomainCards := mappers.BatchCreditCardEntityToDomain(creditCardEntities)
	return allDomainCards, nil
}
