package config

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

type servConfig struct {
	storeInterval time.Duration
	storeFile     string
	restore       bool
	key           string
}

type ServerEnvConfig struct {
	Address       string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
	Key           string        `env:"KEY"`
}

// NewServerConfig returns app config.
func NewServerConfig() (*ServerEnvConfig, error) {
	// Declare an instance of the environment config struct.
	envCfg := &ServerEnvConfig{}
	err := env.Parse(envCfg)
	if err != nil {
		log.Fatal(err)
	}

	// Declare an instance of the config struct for values
	// from a command line.
	var cfg servConfig

	addr := new(NetAddress)

	// Read the value of the port and env command-line flags into the config struct.
	flag.Var(addr, "a", "Metrics server address host:port")
	flag.DurationVar(&cfg.storeInterval, "i", time.Duration(300*time.Second), "pollInterval duration in seconds")
	flag.StringVar(&cfg.storeFile, "f", "/tmp/devops-metrics-db.json", "json filename to store metrics")
	flag.BoolVar(&cfg.restore, "r", true, "restore from json file")
	flag.StringVar(&cfg.key, "k", "2bb80d537b1da3e38bd30361aa855686bde0eacd7162fef6a25fe97bf527a25b", "key")

	flag.Parse()

	if envCfg.Address == "" {
		envCfg.Address = addr.String()
	}

	if envCfg.StoreInterval == 0 {
		envCfg.StoreInterval = cfg.storeInterval
	}

	if envCfg.StoreFile == "" {
		envCfg.StoreFile = cfg.storeFile
	}

	if envCfg.Key == "" {
		envCfg.Key = cfg.key
	}

	return envCfg, nil
}
