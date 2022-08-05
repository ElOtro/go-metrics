package main

import (
	"log"
	"net/http"

	data "github.com/ElOtro/go-metrics/internal/collector"
	"github.com/ElOtro/go-metrics/internal/config"
)

// Define an application struct to hold the dependencies
type application struct {
	config *config.EnvConfig
	stats  *data.Stats
	client http.Client
}

func main() {
	// Declare an instance of the config struct.
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(cfg)

	app := &application{
		config: cfg,
		stats:  data.NewStats(cfg.PollInterval),
		client: http.Client{},
	}

	// run sending metrics
	app.postMetrics()

}
