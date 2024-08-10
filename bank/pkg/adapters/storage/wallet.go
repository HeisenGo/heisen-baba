package storage

import (
	creditCard "bankservice/internal/wallet/credit_card"
	"bankservice/internal/wallet/wallet"
	"bankservice/pkg/adapters/storage/entities"
	"bankservice/pkg/adapters/storage/mappers"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
)

type walletRepo struct {
	db *gorm.DB
}

func NewWalletRepo(db *gorm.DB) wallet.Repo {
	return &walletRepo{
		db: db,
	}
}
func (r *walletRepo) Create(ctx context.Context, wl *wallet.Wallet) (*wallet.Wallet, error) {
	newWallet := mappers.WalletDomainToEntity(wl)
	err := r.db.WithContext(ctx).Create(&newWallet).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil, wallet.ErrUserAlreadyHasWallet
		}
		return nil, err
	}
	createdWallet := mappers.WalletEntityToDomain(newWallet)
	return createdWallet, nil
}

func (r *walletRepo) Deposit(ctx context.Context, card *creditCard.CreditCard, amount uint, userID uuid.UUID) (*wallet.Wallet, error) {
	var userWalletEntity *entities.Wallet
	var cardEntity *entities.CreditCard

	if err := r.db.Where("user_id = ?", userID).First(&userWalletEntity).Error; err != nil {
		return nil, err
	}

	// Check if the credit card exists and belongs to the user's wallet
	if err := r.db.WithContext(ctx).Joins("JOIN wallet_credit_cards ON wallet_credit_cards.credit_card_id = credit_cards.id").
		Joins("JOIN wallets ON wallets.id = wallet_credit_cards.wallet_id").
		Where("credit_cards.number = ? AND wallets.user_id = ?", card.Number, userID).
		First(&cardEntity).Error; err != nil {
		return nil, err
	}

	// Increase the wallet balance
	userWalletEntity.Balance += amount
	if err := r.db.WithContext(ctx).Save(&userWalletEntity).Error; err != nil {
		return nil, err
	}

	createdWallet := mappers.WalletEntityToDomain(userWalletEntity)
	return createdWallet, nil
}

func (r *walletRepo) Withdraw(ctx context.Context, card *creditCard.CreditCard, amount uint, userID uuid.UUID) (*wallet.Wallet, error) {
	var userWalletEntity *entities.Wallet
	var cardEntity *entities.CreditCard

	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&userWalletEntity).Error; err != nil {
		return nil, err
	}

	// Check if the credit card exists and belongs to the user's wallet
	if err := r.db.WithContext(ctx).Joins("JOIN wallet_credit_cards ON wallet_credit_cards.credit_card_id = credit_cards.id").
		Joins("JOIN wallets ON wallets.id = wallet_credit_cards.wallet_id").
		Where("credit_cards.number = ? AND wallets.user_id = ?", card.Number, userID).
		First(&cardEntity).Error; err != nil {
		return nil, err
	}

	if userWalletEntity.Balance < amount {
		return nil, wallet.ErrNotEnoughBalance
	}
	userWalletEntity.Balance -= amount
	if err := r.db.WithContext(ctx).Save(&userWalletEntity).Error; err != nil {
		return nil, err
	}
	createdWallet := mappers.WalletEntityToDomain(userWalletEntity)
	return createdWallet, nil
}

func (r *walletRepo) GetWallet(ctx context.Context, userID uuid.UUID) (*wallet.Wallet, error) {
	var userWalletEntity *entities.Wallet
	err := r.db.Where("user_id = ?", userID).First(&userWalletEntity).Error
	if err != nil {
		return nil, err
	}
	var fetchedWalletEntity *entities.Wallet
	err = r.db.WithContext(ctx).Model(&entities.Wallet{}).Where("id = ?", userWalletEntity.ID).First(&fetchedWalletEntity).Error
	if err != nil {
		return nil, err
	}
	fetchedWalletDomain := mappers.WalletEntityToDomain(fetchedWalletEntity)
	return fetchedWalletDomain, nil
}

func (r *walletRepo) Transfer(ctx context.Context, tr *wallet.TransferTransaction) (*wallet.TransferTransaction, error) {
	transaction := mappers.DomainTransactionToTransactionEntity(tr)
	var wallets []entities.Wallet
	var walletIDs []uuid.UUID
	if tr.IsPaidToSystem {
		walletIDs = []uuid.UUID{*transaction.FromWallet.UserID}
	} else {
		walletIDs = []uuid.UUID{*transaction.FromWallet.UserID, *transaction.ToWallet.UserID}
	}
	err := r.db.Where("user_id IN ?", walletIDs).Find(&wallets).Error
	if err != nil {
		return nil, err
	}
	if transaction.Amount < 100 {
		return nil, wallet.ErrNotEnoughBalance
	}
	if !(len(wallets) == 2 && transaction.IsPaidToSystem == false || len(wallets) == 1 && transaction.IsPaidToSystem == true) {
		return nil, wallet.ErrUserWalletDoesNotExists
	}
	var fromWalEntity *entities.Wallet
	var toWalEntity *entities.Wallet
	for _, wal := range wallets {
		if *wal.UserID == *transaction.FromWallet.UserID {
			fromWalEntity = &wal
			transaction.FromWallet.ID = fromWalEntity.ID
		} else if *wal.UserID == *transaction.ToWallet.UserID {
			toWalEntity = &wal
			transaction.ToWallet.ID = toWalEntity.ID
		}
	}
	if fromWalEntity.Balance < transaction.Amount {
		return nil, wallet.ErrNotEnoughBalance
	}
	fromWalEntity.Balance -= transaction.Amount
	if err := r.db.WithContext(ctx).Save(&fromWalEntity).Error; err != nil {
		return nil, err
	}
	var systemWalEntity *entities.Wallet
	err = r.db.Where("is_system_wallet = ?", true).First(&systemWalEntity).Error
	if err != nil {
		return nil, wallet.ErrSystemWalletDoesNotExists
	}
	if transaction.IsPaidToSystem == true {
		systemWalEntity.Balance += transaction.Amount
	} else {
		var commissionEntity *entities.Commission
		err := r.db.WithContext(ctx).First(&commissionEntity).Error
		if err != nil {
			return nil, err
		}
		tax := transaction.Amount * commissionEntity.AppCommissionPercentage / 100
		toWalEntity.Balance += transaction.Amount - tax
		if err := r.db.WithContext(ctx).Save(&toWalEntity).Error; err != nil {
			return nil, err
		}
		systemWalEntity.Balance += tax
	}
	if err := r.db.WithContext(ctx).Save(&systemWalEntity).Error; err != nil {
		return nil, err
	}
	transaction.Status = entities.TransactionSuccess
	err = r.db.WithContext(ctx).Create(&transaction).Error
	if err != nil {
		return nil, err
	}
	createdTransaction := mappers.TransactionEntityToDomain(transaction)
	return createdTransaction, nil
}
