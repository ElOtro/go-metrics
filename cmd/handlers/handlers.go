package handlers

import (
	"github.com/ElOtro/go-metrics/internal/repo"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Repo interface {
	repo.Getter
}

// Create a Handlers struct which wraps all models.
type Handlers struct {
	repo Repo
}

// For ease of use, we also add a NewHandlers() method which
// returns a Handlers struct
func NewHandlers(repo repo.Getter) *Handlers {
	return &Handlers{repo: repo}
}

func (a *Handlers) Routes() *chi.Mux {
	r := chi.NewRouter()
	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// r.Use(app.getQueryParams)

	// RESTy routes for "articles" resource
	r.Get("/", a.GetAllMetricsHandler)

	r.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", a.GetMetricHandler)
	})

	r.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", a.CreateMetricHandler)
	})

	return r

}
