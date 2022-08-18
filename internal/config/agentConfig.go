package config

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

type agentConfig struct {
	reportInterval time.Duration
	pollInterval   time.Duration
	key            string
}

type AgentEnvConfig struct {
	Address        string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	Key            string        `env:"KEY"`
}

// NewAgentConfig returns app config.
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

	addr := new(NetAddress)
	// Read the value of the port and env command-line flags into the config struct.
	flag.Var(addr, "a", "Metrics server address host:port")
	flag.DurationVar(&cfg.reportInterval, "r", time.Duration(10*time.Second), "reportInterval duration in seconds")
	flag.DurationVar(&cfg.pollInterval, "p", time.Duration(2*time.Second), "pollInterval duration in seconds")
	flag.StringVar(&cfg.key, "k", "", "key")
	// flag.StringVar(&cfg.key, "k", "2bb80d537b1da3e38bd30361aa855686bde0eacd7162fef6a25fe97bf527a25b", "key")

	flag.Parse()

	if envCfg.Address == "" {
		envCfg.Address = addr.String()
	}

	if envCfg.ReportInterval == 0 {
		envCfg.ReportInterval = cfg.reportInterval
	}

	if envCfg.PollInterval == 0 {
		envCfg.PollInterval = cfg.pollInterval
	}

	if envCfg.Key == "" {
		envCfg.Key = cfg.key
	}

	return envCfg, nil
}
