package wallet

import (
	"context"
	"github.com/google/uuid"
)

type CreditCardOps struct {
	repo Repo
}

func NewCreditCardOps(repo Repo) *CreditCardOps {
	return &CreditCardOps{
		repo: repo,
	}
}

func (o *CreditCardOps) CreateCardAndAddToWallet(ctx context.Context, creditCard *CreditCard, userID uuid.UUID) (*CreditCard, error) {
	if !isValidCardNumber(creditCard.Number) {
		return nil, ErrInvalidCardNumber
	}
	return o.repo.CreateCardAndAddToWallet(ctx, creditCard, userID)
}

func (o *CreditCardOps) GetUserWalletCards(ctx context.Context, userID uuid.UUID) ([]CreditCard, error) {
	return o.repo.GetUserWalletCards(ctx, userID)
}
