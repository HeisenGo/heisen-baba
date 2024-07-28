package handlers

import (
	presenter "authservice/api/grpc/handlers/presentor"
	"authservice/internal/user"
	"authservice/protobufs"
	"authservice/service"
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

type GRPCAuthHandler struct {
	protobufs.UnimplementedAuthServiceServer
	authService *service.AuthService
}

func NewGRPCAuthHandler(authService *service.AuthService) *GRPCAuthHandler {
	return &GRPCAuthHandler{authService: authService}
}

func (a *GRPCAuthHandler) Register(ctx context.Context, req *protobufs.RegisterRequest) (*protobufs.RegisterResponse, error) {
	u := presenter.UserRegisterToUserDomain(req)
	newUser, err := a.authService.CreateUser(ctx, u)
	if err != nil {
		if errors.Is(err, user.ErrInvalidEmail) || errors.Is(err, user.ErrInvalidPassword) {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, user.ErrEmailAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, err.Error())
		}

		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &protobufs.RegisterResponse{UserId: newUser.ID.String()}, nil
}

func (a *GRPCAuthHandler) Login(ctx context.Context, req *protobufs.LoginRequest) (*protobufs.LoginResponse, error) {
	authToken, err := a.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &protobufs.LoginResponse{Token: authToken.AuthorizationToken, RefreshToken: authToken.RefreshToken}, nil
}

func (a *GRPCAuthHandler) GetUserByToken(ctx context.Context, req *protobufs.GetUserByTokenRequest) (*protobufs.GetUserByTokenResponse, error) {
	claim, err := a.authService.GetUserByToken(ctx, req.Token)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	return &protobufs.GetUserByTokenResponse{
		UserId:  claim.UserID.String(),
		IsAdmin: claim.IsAdmin,
	}, nil
}

func (a *GRPCAuthHandler) RefreshToken(ctx context.Context, req *protobufs.RefreshTokenRequest) (*protobufs.LoginResponse, error) {
	refToken := req.RefreshToken
	if len(refToken) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "token should not be empty")
	}
	authToken, err := a.authService.RefreshAuth(ctx, refToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}
	return &protobufs.LoginResponse{Token: authToken.AuthorizationToken, RefreshToken: req.RefreshToken}, nil
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
