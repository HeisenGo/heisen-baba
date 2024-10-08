package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"
	http_server "tripcompanyservice/api/http"
	"tripcompanyservice/config"
	"tripcompanyservice/service"
)

var configPath = flag.String("config", "", "configuration path")

func main() {
	cfg := readConfig()

	app, err := service.NewAppContainer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	http_server.Run(cfg.Server, app)

	go startTripCancellationRoutine(app.TripService())

	select {}
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

func startTripCancellationRoutine(tripService *service.TripService) {
	for {
		// Run
		err := tripService.GetUpcomingUnconfirmedTripIDsToCancel(context.Background())
		if err != nil {
			log.Println("Error getting trips to cancel:", err)
			time.Sleep(5 * time.Minute)
		}

		time.Sleep(24 * time.Hour)
	}
}
