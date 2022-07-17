package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handlers) GetAllMetricsHandler(w http.ResponseWriter, r *http.Request) {
	s := h.repo.GetAll()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(s))
}

func (h *Handlers) CreateMetricHandler(w http.ResponseWriter, r *http.Request) {
	t := chi.URLParam(r, "type")
	n := chi.URLParam(r, "name")
	v := chi.URLParam(r, "value")

	if t == "" && n == "" && v == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if t != "gauge" && t != "counter" {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	err := h.repo.Set(t, n, v)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func (h *Handlers) GetMetricHandler(w http.ResponseWriter, r *http.Request) {
	t := chi.URLParam(r, "type")
	n := chi.URLParam(r, "name")

	if t == "" && n == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s, err := h.repo.Get(t, n)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(s))
}
