package clients


type IBankClient interface {
	Transfer(senderOwnerID, receiverOwnerID, isPaidToSystem string, amount uint64) (bool,error)
}