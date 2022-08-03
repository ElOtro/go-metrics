package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ElOtro/go-metrics/internal/config"
	"github.com/ElOtro/go-metrics/internal/handlers"
	"github.com/ElOtro/go-metrics/internal/repo"
	"github.com/caarlos0/env/v6"
)

func main() {
	// Declare an instance of the config struct.
	var cfg config.Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(cfg)

	rep, err := repo.NewMemStorage(&repo.Options{Environment: cfg.Environment})
	if err != nil {
		log.Fatalln(err)
	}

	// Initialize a new Handlers struct
	h := handlers.NewHandlers(rep)

	// Declare a HTTP server with some sensible timeout settings, which listens on the
	// port provided in the config struct and uses the servemux we created above as the
	// handler.
	srv := &http.Server{
		Addr:         cfg.Address,
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
