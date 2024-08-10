package services

import (
	creditCard "bankservice/internal/wallet/credit_card"
	"bankservice/internal/wallet/wallet"
	"context"
	"github.com/google/uuid"
)

type WalletService struct {
	walletOps     *wallet.WalletOps
	creditCardOps *creditCard.CreditCardOps
}

func NewWalletService(walletOps *wallet.WalletOps, creditCardOps *creditCard.CreditCardOps) *WalletService {
	return &WalletService{
		walletOps:     walletOps,
		creditCardOps: creditCardOps,
	}
}

func (s *WalletService) CreateWallet(ctx context.Context, wl *wallet.Wallet) (*wallet.Wallet, error) {
	createdWallet, err := s.walletOps.Create(ctx, wl)
	if err != nil {
		return nil, err
	}
	return createdWallet, nil
}

func (s *WalletService) AddCardToWalletByUserID(ctx context.Context, card *creditCard.CreditCard, userID uuid.UUID) (*creditCard.CreditCard, error) {
	createdCard, err := s.creditCardOps.CreateCardAndAddToWallet(ctx, card, userID)
	if err != nil {
		return nil, err
	}
	return createdCard, nil
}

func (s *WalletService) GetUserWalletCards(ctx context.Context, userID uuid.UUID) ([]creditCard.CreditCard, error) {
	userWalletCards, err := s.creditCardOps.GetUserWalletCards(ctx, userID)
	if err != nil {
		return nil, err
	}
	return userWalletCards, nil
}

func (s *WalletService) Deposit(ctx context.Context, card *creditCard.CreditCard, amount uint, userID uuid.UUID) (*wallet.Wallet, error) {
	userWallet, err := s.walletOps.Deposit(ctx, card, amount, userID)
	if err != nil {
		return nil, err
	}
	return userWallet, nil
}

func (s *WalletService) Withdraw(ctx context.Context, card *creditCard.CreditCard, amount uint, userID uuid.UUID) (*wallet.Wallet, error) {
	userWallet, err := s.walletOps.Withdraw(ctx, card, amount, userID)
	if err != nil {
		return nil, err
	}
	return userWallet, nil
}

func (s *WalletService) GetWallet(ctx context.Context, userID uuid.UUID) (*wallet.Wallet, error) {
	userWallet, err := s.walletOps.GetWallet(ctx, userID)
	if err != nil {
		return nil, err
	}
	return userWallet, nil
}

func (s *WalletService) Transfer(ctx context.Context, tr *wallet.TransferTransaction) (*wallet.TransferTransaction, error) {
	createdTransaction, err := s.walletOps.Transfer(ctx, tr)
	if err != nil {
		return nil, err
	}
	return createdTransaction, nil
}
