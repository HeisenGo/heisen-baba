package presenter

import (
	"bankservice/internal/wallet/wallet"
	"bankservice/protobufs"
	"github.com/google/uuid"
)

func CreateWalletReqToWalletDomain(w *protobufs.CreateWalletRequest) (*wallet.Wallet, error) {
	uid, err := uuid.Parse(w.UserID)
	if err != nil {
		return nil, err
	}
	return &wallet.Wallet{
		UserID: uid,
	}, nil
}
func TransferReqToTransferTransactionDomain(t *protobufs.TransferRequest) (*wallet.TransferTransaction, error) {
	var receiverUserUUID uuid.UUID
	senderUserUUID, err := uuid.Parse(t.SenderOwnerID)
	if err != nil {
		return nil, err
	}
	if !t.IsPaidToSystem {
		receiverUserUUID, err = uuid.Parse(t.ReceiverOwnerID)
		if err != nil {
			return nil, err
		}
	}
	fromWl := &wallet.Wallet{
		UserID: senderUserUUID,
	}
	toWl := &wallet.Wallet{
		UserID: receiverUserUUID,
	}
	return &wallet.TransferTransaction{
		Amount:         uint(t.Amount),
		FromWallet:     fromWl,
		ToWallet:       toWl,
		IsPaidToSystem: t.IsPaidToSystem,
	}, nil
}
