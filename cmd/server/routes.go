package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() *chi.Mux {
	r := chi.NewRouter()
	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// r.Use(app.getQueryParams)

	// RESTy routes for "articles" resource
	r.Get("/", app.GetAllMetricsHandler)

	r.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", app.CreateMetricHandler)
	})

	r.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", app.GetMetricHandler)
	})

	// Return the router instance.
	return r
}
