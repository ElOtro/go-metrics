package main

import (
	"log"
	"net/http"

	data "github.com/ElOtro/go-metrics/internal/collector"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address        string `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	ReportInterval int    `env:"REPORT_INTERVAL" envDefault:"10"`
	PollInterval   int    `env:"POLL_INTERVAL" envDefault:"2"`
}

// Define an application struct to hold the dependencies
type application struct {
	config Config
	stats  *data.Stats
	client http.Client
}

func main() {
	var cfg Config
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
