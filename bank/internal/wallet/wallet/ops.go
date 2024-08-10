package wallet

import (
	creditCard "bankservice/internal/wallet/credit_card"
	"context"
	"github.com/google/uuid"
)

type WalletOps struct {
	repo Repo
}

func NewWalletOps(repo Repo) *WalletOps {
	return &WalletOps{
		repo: repo,
	}
}

func (o *WalletOps) Create(ctx context.Context, wallet *Wallet) (*Wallet, error) {
	return o.repo.Create(ctx, wallet)
}

func (o *WalletOps) Deposit(ctx context.Context, creditCard *creditCard.CreditCard, amount uint, userID uuid.UUID) (*Wallet, error) {
	return o.repo.Deposit(ctx, creditCard, amount, userID)
}

func (o *WalletOps) Withdraw(ctx context.Context, creditCard *creditCard.CreditCard, amount uint, userID uuid.UUID) (*Wallet, error) {
	return o.repo.Withdraw(ctx, creditCard, amount, userID)
}

func (o *WalletOps) GetWallet(ctx context.Context, userID uuid.UUID) (*Wallet, error) {
	return o.repo.GetWallet(ctx, userID)
}

func (o *WalletOps) Transfer(ctx context.Context, transaction *TransferTransaction) (*TransferTransaction, error) {
	return o.repo.Transfer(ctx, transaction)
}
