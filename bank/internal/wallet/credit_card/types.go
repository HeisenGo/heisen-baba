package wallet

import (
	"context"
	"errors"
	"github.com/google/uuid"
)

var (
	ErrInvalidCardNumber = errors.New("invalid card number")
	ErrCardAlreadyExists = errors.New("card already exists in wallet")
)

type Repo interface {
	CreateCardAndAddToWallet(ctx context.Context, creditCard *CreditCard, userID uuid.UUID) (*CreditCard, error)
	GetUserWalletCards(ctx context.Context, userID uuid.UUID) ([]CreditCard, error)
}

type CreditCard struct {
	ID     uuid.UUID `json:"id"`
	Number string    `json:"number"`
}

func NewCreditCard(number string) *CreditCard {
	return &CreditCard{
		Number: number,
	}
}
