package main

import (
	"auth/internal/service"
	"auth/pkg/db"
	"auth/pkg/logger"
	"net"

	pb "auth/api"
	"auth/internal/config"

	"google.golang.org/grpc"
)

func main() {
	log := logger.New()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dbConn := db.MustInitDB(cfg.DSN)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	lis, err := net.Listen("tcp", cfg.ServerAddress)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, auth.NewService(dbConn, []byte(cfg.Secret)))

	log.Infof("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
