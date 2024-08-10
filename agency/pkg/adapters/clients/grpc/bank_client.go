package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"agency/internal/bank"
	"agency/pkg/ports"
	"agency/protobufs"
)

type GRPCBankClient struct {
	ServiceRegistry ports.IServiceRegistry
	BankServiceName string
}

func NewGRPCBankClient(serviceRegistry ports.IServiceRegistry, bankServiceName string) *GRPCBankClient {
	return &GRPCBankClient{ServiceRegistry: serviceRegistry, BankServiceName: bankServiceName}
}

func (g *GRPCBankClient) Transfer(senderOwnerID, receiverOwnerID string, isPaidToSystem bool, amount uint64) (bool, error) {
	port, ip, err := g.ServiceRegistry.DiscoverService(g.BankServiceName)
	if err != nil {
		return false, err
	}

	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", ip, port), grpc.WithInsecure())
	if err != nil {
		return false, err
	}

	defer conn.Close()

	// Create a new AuthService client
	client := protobufs.NewBankServiceClient(conn)

	// Create a context
	ctx := context.Background()

	// Prepare the request
	request := &protobufs.TransferRequest{
		SenderOwnerID:   "ecfe4a77-a1c3-4dd9-abfc-3349bd4b9db2",
		ReceiverOwnerID: "ecfe4a77-a1c3-4dd9-abfc-3349bd4b9db2",
		IsPaidToSystem:  false,
		Amount:          50000,
	}

	// Call the GetUserByToken method
	// Call the Transfer method
	response, err := client.Transfer(ctx, request)
	if err != nil {
		return false, err
	}

	// Extract the gRPC status code
	st, ok := status.FromError(err)
	if !ok {
		return false, err
	}
	switch st.Code() {
	case codes.OK:
		if response.Status == "success" {
			return true, nil
		}
	default:
		return false, bank.ErrNotEnoughMoney
	}
	return false, nil
}
