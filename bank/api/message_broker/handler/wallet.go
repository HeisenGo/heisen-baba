package handler

import (
	"bankservice/internal/wallet/wallet"
	services "bankservice/service"
	"context"
	"github.com/google/uuid"
	"log"
)

type WalletHandler struct {
	walletService *services.WalletService
}

func NewWalletHandler(walletService *services.WalletService) *WalletHandler {
	return &WalletHandler{walletService: walletService}
}

func (h *WalletHandler) CreateWallet(userID string) {
	uid, _ := uuid.Parse(userID)
	log.Printf("Creating wallet for user ID: %s", userID)
	h.walletService.CreateWallet(context.Background(), &wallet.Wallet{
		UserID: uid,
	})
}
