package handlers

import (
	presenter "bankservice/api/grpc/handlers/presentor"
	"bankservice/internal/wallet/wallet"
	"bankservice/protobufs"
	"bankservice/service"
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

type GRPCBankHandler struct {
	protobufs.UnimplementedBankServiceServer
	walletService *services.WalletService
}

func NewGRPCBankHandler(bankService *services.WalletService) *GRPCBankHandler {
	return &GRPCBankHandler{walletService: bankService}
}

func (g *GRPCBankHandler) CreateWallet(ctx context.Context, wl *protobufs.CreateWalletRequest) (*protobufs.CreateWalletRequestResponse, error) {
	domainWallet, err := presenter.CreateWalletReqToWalletDomain(wl)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	_, err = g.walletService.CreateWallet(ctx, domainWallet)
	if err != nil {
		if errors.Is(err, wallet.ErrUserAlreadyHasWallet) {
			return nil, status.Errorf(codes.AlreadyExists, "wallet already exists")
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &protobufs.CreateWalletRequestResponse{Message: "wallet created"}, nil
}
func (g *GRPCBankHandler) Transfer(ctx context.Context, tr *protobufs.TransferRequest) (*protobufs.TransferResponse, error) {
	domainTransaction, err := presenter.TransferReqToTransferTransactionDomain(tr)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	createdTransaction, err := g.walletService.Transfer(ctx, domainTransaction)
	if err != nil {
		//if errors.Is(err, wallet.ErrUserAlreadyHasWallet) {
		//	return nil, status.Errorf(codes.AlreadyExists, "wallet already exists")
		//}
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	receiverOwnerID := ""
	if !createdTransaction.IsPaidToSystem {
		receiverOwnerID = createdTransaction.ToWallet.UserID.String()
	}
	return &protobufs.TransferResponse{
		SenderOwnerID:   createdTransaction.FromWallet.UserID.String(),
		ReceiverOwnerID: receiverOwnerID,
		IsPaidToSystem:  createdTransaction.IsPaidToSystem,
		Amount:          uint64(createdTransaction.Amount),
		Status:          string(createdTransaction.Status),
	}, nil
}

type HealthServer struct {
	grpc_health_v1.HealthServer
}

// Check implements Health.Check
func (s *HealthServer) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}
