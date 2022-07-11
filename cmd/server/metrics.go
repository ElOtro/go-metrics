package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) createMetricHandler(w http.ResponseWriter, r *http.Request) {
	t := chi.URLParam(r, "type")
	n := chi.URLParam(r, "name")
	v := chi.URLParam(r, "value")

	if t == "" && n == "" && v == "" {
		fmt.Println("is zero value")
	}

	err := app.rep.Set(t, n, v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) getMetricHandler(w http.ResponseWriter, r *http.Request) {
	s := app.rep.Get()
	w.Write([]byte(s))
}
