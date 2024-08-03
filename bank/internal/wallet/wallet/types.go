package wallet

import (
	creditCard "bankservice/internal/wallet/credit_card"
	"bankservice/pkg/adapters/storage/entities"
	"context"
	"errors"
	"github.com/google/uuid"
)

var (
	ErrNotEnoughBalance          = errors.New("not enough balance")
	ErrUserAlreadyHasWallet      = errors.New("user already has wallet")
	ErrUserWalletDoesNotExists   = errors.New("user wallet does not exists")
	ErrSystemWalletDoesNotExists = errors.New("system wallet does not exists")
	ErrMinTrans                  = errors.New("minimum value of transaction is 100")
)

type Repo interface {
	Create(ctx context.Context, user *Wallet) (*Wallet, error)
	Deposit(ctx context.Context, creditCard *creditCard.CreditCard, amount uint, userID uuid.UUID) (*Wallet, error)
	Withdraw(ctx context.Context, creditCard *creditCard.CreditCard, amount uint, userID uuid.UUID) (*Wallet, error)
	GetWallet(ctx context.Context, userID uuid.UUID) (*Wallet, error)
	Transfer(ctx context.Context, tr *TransferTransaction) (*TransferTransaction, error)
}

type Wallet struct {
	ID             *uuid.UUID `json:"id"`
	IsSystemWallet bool       `json:"is_system_wallet"`
	UserID         uuid.UUID  `json:"user_id"`
	Balance        uint       `json:"balance"`
}

type TransferTransaction struct {
	Amount         uint
	Status         entities.TransferTransactionStatus
	FromWallet     *Wallet
	ToWallet       *Wallet
	IsPaidToSystem bool
}
