package main

import (
	"log"
	"net/http"

	data "github.com/ElOtro/go-metrics/internal/collector"
	"github.com/ElOtro/go-metrics/internal/config"
	"github.com/caarlos0/env/v6"
)

// Define an application struct to hold the dependencies
type application struct {
	config config.Config
	stats  *data.Stats
	client http.Client
}

func main() {
	var cfg config.Config
	err := env.Parse(&cfg)
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
