package config

type Config struct {
	Address        string `env:"ADDRESS,required" envDefault:"127.0.0.1:8080"`
	ReportInterval int    `env:"REPORT_INTERVAL,required" envDefault:"10"`
	PollInterval   int    `env:"POLL_INTERVAL,required" envDefault:"2"`
	Enviroment     string `env:"ENVIROMENT" envDefault:"debug"`
}
