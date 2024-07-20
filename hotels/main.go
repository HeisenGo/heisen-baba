package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/consul/api"
	"log"
)

func main() {
	// Initialize Fiber app
	app := fiber.New()

	// Define a simple route
	app.Get("/hotels", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to the Hotels Service!",
		})
	})

	// Load Consul configuration
	consulAddress := "127.0.0.1:8500"
	serviceName := "hotels"
	serviceID := "hotels-service-id"

	// Initialize Consul client
	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulAddress
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		log.Fatalf("Failed to create Consul client: %v", err)
	}

	// Register service with Consul
	registration := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Address: "127.0.0.1",
		Port:    8080,
		Tags:    []string{"hotels"},
		Check: &api.AgentServiceCheck{
			HTTP:     "http://127.0.0.1:8080/health",
			Interval: "10s",
			Timeout:  "1s",
		},
	}

	err = consulClient.Agent().ServiceRegister(registration)
	if err != nil {
		log.Fatalf("Failed to register service with Consul: %v", err)
	}

	// Setup a simple health check route
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Service is up and running")
	})

	// Start the Fiber app
	log.Printf("Hotels service listening on port 8080")
	log.Fatal(app.Listen(":8080"))
}
