package grpc

import (
	"fmt"
	"log"
	"net"
	"terminalpathservice/api/grpc/handlers"
	"terminalpathservice/config"
	"terminalpathservice/protobufs"
	"terminalpathservice/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func Run(cfg config.Config, app *service.AppContainer) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", cfg.Server.GRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	protobufs.RegisterPathServiceServer(s, handlers.NewPathServiceServer(app.PathService()))
	// Register the Health Service server
	healthServer := &handlers.HealthServer{}
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	// Register reflection service on gRPC server
	reflection.Register(s)

	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		fmt.Printf("Failed to serve: %v", err)
	}
}
