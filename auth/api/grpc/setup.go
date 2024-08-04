package grpc

import (
	"authservice/api/grpc/handlers"
	"authservice/config"
	auth "authservice/protobufs"
	"authservice/service"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func Run(cfg config.Config, app *service.AppContainer) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", cfg.Server.GRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	auth.RegisterAuthServiceServer(s, handlers.NewGRPCAuthHandler(app.AuthService()))
	// Register the Health Service server
	healthServer := &handlers.HealthServer{}
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	// Register reflection service on gRPC server
	reflection.Register(s)

	log.Println("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
