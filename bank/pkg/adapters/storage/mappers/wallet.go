package mappers

import (
	creditCard "bankservice/internal/wallet/credit_card"
	"bankservice/internal/wallet/wallet"
	"bankservice/pkg/adapters/storage/entities"
)

func WalletEntityToDomain(entity *entities.Wallet) *wallet.Wallet {
	return &wallet.Wallet{
		ID:      &entity.ID,
		UserID:  *entity.UserID,
		Balance: entity.Balance,
	}
}

func WalletDomainToEntity(domainWallet *wallet.Wallet) *entities.Wallet {
	return &entities.Wallet{
		UserID:  &domainWallet.UserID,
		Balance: domainWallet.Balance,
	}
}

func CreditCardEntityToDomain(entity *entities.CreditCard) *creditCard.CreditCard {
	return &creditCard.CreditCard{
		ID:     entity.ID,
		Number: entity.Number,
	}
}

func BatchCreditCardEntityToDomain(entities []*entities.CreditCard) []creditCard.CreditCard {
	var domainCreditCards []creditCard.CreditCard
	for _, e := range entities {
		domainCreditCards = append(domainCreditCards, creditCard.CreditCard{ID: e.ID, Number: e.Number})
	}
	return domainCreditCards
}

func CreditCardDomainToEntity(domainWallet *creditCard.CreditCard) *entities.CreditCard {
	return &entities.CreditCard{
		Number: domainWallet.Number,
	}
}

func DomainTransactionToTransactionEntity(domainTr *wallet.TransferTransaction) *entities.TransferTransaction {
	var toWl *entities.Wallet
	fromWl := WalletDomainToEntity(domainTr.FromWallet)
	if !domainTr.IsPaidToSystem {
		toWl = WalletDomainToEntity(domainTr.ToWallet)
	}
	return &entities.TransferTransaction{
		Amount:         domainTr.Amount,
		FromWallet:     fromWl,
		ToWallet:       toWl,
		IsPaidToSystem: domainTr.IsPaidToSystem,
	}
}

func TransactionEntityToDomain(entity *entities.TransferTransaction) *wallet.TransferTransaction {
	var toWalDomain *wallet.Wallet
	fromWalDomain := WalletEntityToDomain(entity.FromWallet)
	if !entity.IsPaidToSystem {
		toWalDomain = WalletEntityToDomain(entity.ToWallet)
	}
	return &wallet.TransferTransaction{
		Amount:         entity.Amount,
		Status:         entity.Status,
		FromWallet:     fromWalDomain,
		ToWallet:       toWalDomain,
		IsPaidToSystem: entity.IsPaidToSystem,
	}
}
