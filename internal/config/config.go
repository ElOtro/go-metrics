package config

type Config struct {
	Address        string `env:"ADDRESS,required" envDefault:"localhost:8080"`
	ReportInterval int    `env:"REPORT_INTERVAL,required" envDefault:"10"`
	PollInterval   int    `env:"POLL_INTERVAL,required" envDefault:"2"`
	Environment    string `env:"ENVIRONMENT" envDefault:"debug"`
}
