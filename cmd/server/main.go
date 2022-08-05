package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ElOtro/go-metrics/internal/config"
	"github.com/ElOtro/go-metrics/internal/handlers"
	"github.com/ElOtro/go-metrics/internal/repo"
	"github.com/ElOtro/go-metrics/internal/service"
)

func main() {
	// Declare an instance of the config struct.
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}

	// Print cfg on start
	log.Println(cfg)

	// Initialize a new Storage struct
	rep, err := repo.NewMemStorage()
	if err != nil {
		log.Fatalln(err)
	}

	// Restore metrics from file
	if cfg.Restore {
		err = rep.RestoreMetrics(cfg.StoreFile)
		if err != nil {
			log.Println(err)
		}
	}

	// Initialize a new Handlers struct
	h := handlers.NewHandlers(rep)

	producer, err := service.NewProducer(cfg.StoreInterval, cfg.StoreFile, rep)
	if err != nil {
		log.Fatalln(err)
	}
	// Run producer to write JSON metrics
	go producer.Run()

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
