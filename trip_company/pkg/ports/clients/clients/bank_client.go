package clients

type IBankClient interface {
	Transfer(senderOwnerID, receiverOwnerID string, isPaidToSystem bool, amount uint64) (bool, error)
}
