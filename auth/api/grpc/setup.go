package grpc

import (
	"authservice/api/grpc/handlers"
	"authservice/config"
	auth "authservice/protobufs"
	"authservice/service"
	"fmt"
	"google.golang.org/grpc"
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

	log.Println("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
