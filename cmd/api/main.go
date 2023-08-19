package main

import (
	"log"

	"github.com/athunlal/bookNowBooking-svc/pkg/config"
	"github.com/athunlal/bookNowBooking-svc/pkg/di"
)

func main() {
	cfg, cfgErr := config.LoadConfig()
	if cfgErr != nil {
		log.Fatal("Could not load the config file:", cfgErr)
		return
	}

	// Initialize the API server
	server, err := di.InitApi(cfg)
	if err != nil {
		log.Fatalln("Error in initializing the API:", err)
	}

	// Start the API server
	server.Start()
}
