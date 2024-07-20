package main

import (
	"api-gateway/config"
	"api-gateway/internal/handlers"
	"api-gateway/internal/routes"
	"api-gateway/pkg/consul"
	"api-gateway/pkg/fiber"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	consulClient, err := consul.NewConsulClient(cfg.Consul.Address)
	if err != nil {
		log.Fatalf("Failed to connect to Consul: %v", err)
	}

	app := fiber.SetupFiber()

	handler := handlers.NewHandler(consulClient)
	routes.SetupRoutes(app, handler)

	log.Fatal(app.Listen(":" + cfg.Gateway.Port))
}
