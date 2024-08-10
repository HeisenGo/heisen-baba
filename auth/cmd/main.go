package main

import (
	grpcServer "authservice/api/grpc"
	httpServer "authservice/api/http"
	"authservice/config"
	"authservice/service"
	"flag"
	"log"
	"os"
	"sync"
)

var configPath = flag.String("config", "", "configuration path")

func main() {
	cfg := readConfig()

	app, err := service.NewAppContainer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(2) // We are running two servers

	// Start the Fiber server
	go func() {
		defer wg.Done()
		httpServer.Run(cfg, app)
	}()

	go func() {
		defer wg.Done()
		grpcServer.Run(cfg, app)
	}()

	wg.Wait()
}

func readConfig() config.Config {
	flag.Parse()

	if cfgPathEnv := os.Getenv("APP_CONFIG_PATH"); len(cfgPathEnv) > 0 {
		*configPath = cfgPathEnv
	}

	if len(*configPath) == 0 {
		log.Fatal("configuration file not found")
	}

	cfg, err := config.ReadStandard(*configPath)

	if err != nil {
		log.Fatal(err)
	}

	return cfg
}
