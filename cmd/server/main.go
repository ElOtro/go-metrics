package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ElOtro/go-metrics/internal/repo"
)

// Define a config struct to hold all the configuration settings for our application.
type config struct {
	address    string
	port       int
	enviroment string
}

// Define an application struct to hold the dependencies
type application struct {
	config config
	rep    repo.Getter
}

func main() {
	// Declare an instance of the config struct.
	var cfg config

	// Read the value of the port and env command-line flags into the config struct.
	flag.StringVar(&cfg.address, "address", "127.0.0.1", "API server address")
	flag.IntVar(&cfg.port, "port", 8080, "API server port")
	flag.StringVar(&cfg.enviroment, "enviroment", "debug", "API server mode")

	flag.Parse()

	rep, err := repo.NewGetter(&repo.Options{Environment: cfg.enviroment})
	if err != nil {
		//  в мейн паниковать можно
		log.Fatalln(err)
	}

	// Declare an instance of the application struct, containing the config struct and
	// the logger.
	app := &application{
		config: cfg,
		rep:    rep,
	}

	// Declare a HTTP server with some sensible timeout settings, which listens on the
	// port provided in the config struct and uses the servemux we created above as the
	// handler.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}

}
