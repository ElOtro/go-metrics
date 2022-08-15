package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ElOtro/go-metrics/internal/config"
	"github.com/ElOtro/go-metrics/internal/handlers"
	"github.com/ElOtro/go-metrics/internal/repo"
	"github.com/ElOtro/go-metrics/internal/service"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	// Declare an instance of the config struct.
	cfg, err := config.NewServerConfig()
	if err != nil {
		log.Fatalln(err)
	}

	// Print cfg on start
	log.Printf("%+v", cfg)

	// set default repository in memory
	repoOptions := &repo.Options{Memory: true}
	// Call the openDB() helper function (see below) to create the connection pool,
	// passing in the Dsn. If this returns an error, we log it and exit the
	// application immediately.
	if cfg.Dsn != "" {
		db, err := openDB(cfg.Dsn)
		if err != nil {
			log.Fatal(err)
		}
		// Add pgxpool.Pool to options
		repoOptions.DB = db
		repoOptions.Memory = false
		// Defer a call to db.Close() so that the connection pool is closed before the
		// main() function exits.
		defer db.Close()
	}

	// Initialize a new Storage struct
	rep, err := repo.NewRepo(repoOptions)
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

	// HashMetric service
	hm := service.NewHashMetric(cfg.Key)

	// Initialize a new Handlers struct
	h := handlers.NewHandlers(rep, *hm)

	producer, err := service.NewProducer(cfg.StoreInterval, cfg.StoreFile, rep)
	if err != nil {
		log.Fatalln(err)
	}
	// Run producer to write JSON metrics
	go producer.Run()

	// Declare a HTTP server with some sensible timeout settings, which listens on the
	// port provided in the config struct and uses the servemux we created above as the
	// handler.
	fmt.Println(cfg.Address)
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

// The openDB() function returns a sql.DB connection pool.
func openDB(dsn string) (*pgxpool.Pool, error) {

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	dbpool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	// Create a context with a 5-second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use Ping() to establish a new connection to the database, passing in the
	// context we created above as a parameter. If the connection couldn't be
	// established successfully within the 5 second deadline, then this will return an
	// error.
	err = dbpool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	// Return the pgxpool.Pool connection pool.
	return dbpool, nil
}
