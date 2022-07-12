package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) CreateMetricHandler(w http.ResponseWriter, r *http.Request) {
	t := chi.URLParam(r, "type")
	n := chi.URLParam(r, "name")
	v := chi.URLParam(r, "value")

	if t == "" && n == "" && v == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := app.rep.Set(t, n, v)
	if err != nil {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) GetMetricHandler(w http.ResponseWriter, r *http.Request) {
	s := app.rep.Get()
	w.Write([]byte(s))
}
