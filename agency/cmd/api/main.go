package main

import (
	"flag"
	"log"
	"os"
	"agency/config"
	"agency/service"

	http_server "agency/api/http"
)

var configPath = flag.String("config", "", "configuration path")
//	@Title			heisenbaba-Agency Domain
//	@version		1.0
//	@description	Agency And Tour Reservaiton Domain

//	@contact.name	HeisenGo
//	@contact.url	https://github.com/HeisenGo

// @host			localhost:8082
// @BasePath		/api/v1
func main() {
	cfg := readConfig()

	app, err := service.NewAppContainer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	http_server.Run(cfg, app)
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
