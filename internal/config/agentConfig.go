package config

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

type agentConfig struct {
	address        string
	reportInterval time.Duration
	pollInterval   time.Duration
}

type AgentEnvConfig struct {
	Address        string        `env:"ADDRESS,required" envDefault:"127.0.0.1:8080"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
}

// NewConfig returns app config.
func NewAgentConfig() (*AgentEnvConfig, error) {

	// Declare an instance of the environment config struct.
	envCfg := &AgentEnvConfig{}
	err := env.Parse(envCfg)
	if err != nil {
		log.Fatal(err)
	}

	// Declare an instance of the config struct for values
	// from a command line.
	var cfg agentConfig

	// Read the value of the port and env command-line flags into the config struct.
	flag.StringVar(&cfg.address, "a", "127.0.0.1:8080", "Metrics server address")
	flag.DurationVar(&cfg.reportInterval, "r", time.Duration(10), "reportInterval duration in seconds")
	flag.DurationVar(&cfg.pollInterval, "p", time.Duration(2), "pollInterval duration in seconds")

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

	return envCfg, nil
}
