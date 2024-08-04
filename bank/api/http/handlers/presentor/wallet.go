package presenter

import (
	creditCard "bankservice/internal/wallet/credit_card"
	"bankservice/internal/wallet/wallet"
)

type AddCardToWalletReq struct {
	CardNumber string `json:"card_number" validate:"required"`
}

type DepositReq struct {
	CardNumber string `json:"card_number" validate:"required"`
	Amount     uint   `json:"amount" validate:"required"`
}
type WithdrawReq struct {
	CardNumber string `json:"card_number" validate:"required"`
	Amount     uint   `json:"amount" validate:"required"`
}

type AddCardToWalletResp struct {
	Card *creditCard.CreditCard `json:"card"`
}

type WalletCardsResp struct {
	Cards []creditCard.CreditCard `json:"cards"`
}
type DepositResp struct {
	Wallet *wallet.Wallet `json:"wallet"`
}
type WithdrawResp struct {
	Message string         `json:"message"`
	Wallet  *wallet.Wallet `json:"wallet"`
}

type GetWalletResp struct {
	Wallet *wallet.Wallet `json:"wallet"`
}

func AddCardToWalletReqToCard(c *AddCardToWalletReq) *creditCard.CreditCard {
	return &creditCard.CreditCard{
		Number: c.CardNumber,
	}
}

func CardToAddCardToWalletResp(c creditCard.CreditCard) AddCardToWalletResp {
	return AddCardToWalletResp{Card: &c}
}

func CardsToWalletCardsResp(cards []creditCard.CreditCard) WalletCardsResp {
	return WalletCardsResp{Cards: cards}
}

func DepositReqNumToCard(cardNum string) *creditCard.CreditCard {
	return &creditCard.CreditCard{
		Number: cardNum,
	}
}
func WithdrawReqNumToCard(cardNum string) *creditCard.CreditCard {
	return &creditCard.CreditCard{
		Number: cardNum,
	}
}

func WalletToDepositResp(wl wallet.Wallet) DepositResp {
	return DepositResp{Wallet: &wl}
}

func WalletToWithdrawResp(wl wallet.Wallet) DepositResp {
	return DepositResp{Wallet: &wl}
}

func WalletToGetWalletResp(wl wallet.Wallet) GetWalletResp {
	return GetWalletResp{
		Wallet: &wl,
	}
}
