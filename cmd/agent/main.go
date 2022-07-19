package main

import (
	"flag"
	"net/http"

	data "github.com/ElOtro/go-metrics/internal/collector"
)

// Define a config struct to hold all the configuration settings for our application.
type config struct {
	pollInterval   int
	reportInterval int
	collectorSrv   struct {
		address string
		port    int
	}
}

// Define an application struct to hold the dependencies
type application struct {
	config config
	stats  *data.Stats
	client http.Client
}

func main() {
	// Declare an instance of the config struct.
	var cfg config

	// Read the value of the port and env command-line flags into the config struct.
	flag.IntVar(&cfg.pollInterval, "pollInterval", 2, "pollInterval duration in seconds")
	flag.IntVar(&cfg.reportInterval, "reportInterval", 10, "reportInterval duration in seconds")
	flag.StringVar(&cfg.collectorSrv.address, "address", "127.0.0.1", "Collector's server address")
	flag.IntVar(&cfg.collectorSrv.port, "port", 8080, "Collector's server port")

	flag.Parse()

	// Declare an instance of the application struct, containing the config, stats structs
	// and http.Clien
	app := &application{
		config: cfg,
		stats:  data.NewStats(cfg.pollInterval),
		client: http.Client{},
	}

	// run sending metrics
	app.postMetricsHandler()

}
