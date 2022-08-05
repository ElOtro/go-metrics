package config

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

type config struct {
	address        string
	reportInterval time.Duration
	pollInterval   time.Duration
	storeInterval  time.Duration
	storeFile      string
	restore        bool
}

type EnvConfig struct {
	Address        string        `env:"ADDRESS,required" envDefault:"127.0.0.1:8080"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
	StoreInterval  time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFile      string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore        bool          `env:"RESTORE" envDefault:"true"`
}

// NewConfig returns app config.
func NewConfig() (*EnvConfig, error) {

	// Declare an instance of the environment config struct.
	envCfg := &EnvConfig{}
	err := env.Parse(envCfg)
	if err != nil {
		log.Fatal(err)
	}

	// Declare an instance of the config struct for values
	// from a command line.
	var cfg config

	// Read the value of the port and env command-line flags into the config struct.
	flag.StringVar(&cfg.address, "ADDRESS", "127.0.0.1:8080", "Metrics server address")
	flag.DurationVar(&cfg.reportInterval, "REPORT_INTERVAL", time.Duration(10), "reportInterval duration in seconds")
	flag.DurationVar(&cfg.pollInterval, "POLL_INTERVAL", time.Duration(2), "pollInterval duration in seconds")
	flag.DurationVar(&cfg.storeInterval, "STORE_INTERVAL", time.Duration(300), "pollInterval duration in seconds")
	flag.StringVar(&cfg.storeFile, "STORE_FILE", "/tmp/devops-metrics-db.json", "json filename to store metrics")
	flag.BoolVar(&cfg.restore, "RESTORE", true, "restore from json file")

	flag.Parse()

	if envCfg.Address == "" {
		envCfg.Address = cfg.address
	}

	if envCfg.ReportInterval == 0 {
		envCfg.ReportInterval = cfg.reportInterval
	}

	if envCfg.PollInterval == 0 {
		envCfg.PollInterval = cfg.pollInterval
	}

	if envCfg.StoreInterval == 0 {
		envCfg.StoreInterval = cfg.storeInterval
	}

	if envCfg.StoreFile == "" {
		envCfg.StoreFile = cfg.storeFile
	}

	return envCfg, nil
}
