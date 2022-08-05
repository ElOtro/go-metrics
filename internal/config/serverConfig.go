package config

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

type config struct {
	address       string
	storeInterval time.Duration
	storeFile     string
	restore       bool
}

type EnvConfig struct {
	Address       string        `env:"ADDRESS,required" envDefault:"127.0.0.1:8080"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFile     string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore       bool          `env:"RESTORE" envDefault:"true"`
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
	flag.StringVar(&cfg.address, "a", "127.0.0.1:8080", "Metrics server address")
	flag.DurationVar(&cfg.storeInterval, "i", time.Duration(300), "pollInterval duration in seconds")
	flag.StringVar(&cfg.storeFile, "f", "/tmp/devops-metrics-db.json", "json filename to store metrics")
	flag.BoolVar(&cfg.restore, "r", true, "restore from json file")

	flag.Parse()

	if envCfg.Address == "" {
		envCfg.Address = cfg.address
	}

	if envCfg.StoreInterval == 0 {
		envCfg.StoreInterval = cfg.storeInterval
	}

	if envCfg.StoreFile == "" {
		envCfg.StoreFile = cfg.storeFile
	}

	return envCfg, nil
}
