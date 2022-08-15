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
	dsn           string
}

type ServerEnvConfig struct {
	Address       string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
	Key           string        `env:"KEY"`
	Dsn           string        `env:"DATABASE_DSN"` // PostgreSQL DSN
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
	flag.StringVar(&cfg.storeFile, "f", "", "json filename to store metrics")
	flag.BoolVar(&cfg.restore, "r", true, "restore from json file")
	// Key for sha256
	// test key is "2bb80d537b1da3e38bd30361aa855686bde0eacd7162fef6a25fe97bf527a25b"
	flag.StringVar(&cfg.key, "k", "", "key")
	// Read the DSN value from the d command-line flag into the config struct. We
	// default to using our development DSN if no flag is provided.
	flag.StringVar(&cfg.dsn, "d", "", "PostgreSQL DSN")

	flag.Parse()

	log.Println("-------------envCfg-----------------")
	log.Printf("%+v", envCfg)
	log.Println("-------------cfg-----------------")
	log.Printf("%+v", cfg)

	if envCfg.Address == "" {
		envCfg.Address = addr.String()
	}

	if envCfg.StoreInterval == 0 {
		envCfg.StoreInterval = cfg.storeInterval
	}

	if envCfg.StoreFile == "" {
		envCfg.StoreFile = cfg.storeFile
	}

	if !cfg.restore {
		envCfg.Restore = cfg.restore
	}

	if envCfg.Key == "" {
		envCfg.Key = cfg.key
	}

	if envCfg.Dsn == "" {
		envCfg.Dsn = cfg.dsn
	}

	// if envCfg.Dsn != "" {
	// 	envCfg.Restore = false
	// 	envCfg.StoreFile = ""
	// }

	return envCfg, nil
}
