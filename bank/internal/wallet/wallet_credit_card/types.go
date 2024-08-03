package wallet

import (
	"context"
	"github.com/google/uuid"
)

type Repo interface {
	Create(ctx context.Context, creditCard *WalletCreditCard) (*WalletCreditCard, error)
}

type WalletCreditCard struct {
	WalletID     uuid.UUID
	CreditCardID uuid.UUID
}

func NewWalletCreditCard(walletID uuid.UUID, creditCardID uuid.UUID) *WalletCreditCard {
	return &WalletCreditCard{WalletID: walletID, CreditCardID: creditCardID}
}
