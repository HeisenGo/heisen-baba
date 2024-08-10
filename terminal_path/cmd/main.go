package main

import (
	"flag"
	"log"
	"os"
	"sync"
	grpc_server "terminalpathservice/api/grpc"
	http_server "terminalpathservice/api/http"
	"terminalpathservice/config"
	"terminalpathservice/service"
)

var configPath = flag.String("config", "", "configuration path")

func main() {
	cfg := readConfig()

	app, err := service.NewAppContainer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		http_server.Run(cfg, app)
	}()

	go func() {
		defer wg.Done()
		grpc_server.Run(cfg, app)
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
