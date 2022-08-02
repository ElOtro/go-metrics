package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ElOtro/go-metrics/internal/handlers"
	"github.com/ElOtro/go-metrics/internal/repo"
	"github.com/caarlos0/env/v6"
)

// Define a config struct to hold all the configuration settings for our application.
type Config struct {
	Address    string `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	Enviroment string `env:"ENVIROMENT" envDefault:"debug"`
}

func main() {
	// Declare an instance of the config struct.
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(cfg)

	rep, err := repo.NewMemStorage(&repo.Options{Environment: cfg.Enviroment})
	if err != nil {
		log.Fatalln(err)
	}

	// Initialize a new Handlers struct
	h := handlers.NewHandlers(rep)

	// Declare a HTTP server with some sensible timeout settings, which listens on the
	// port provided in the config struct and uses the servemux we created above as the
	// handler.
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s", cfg.Address),
		Handler:      h.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}

}
