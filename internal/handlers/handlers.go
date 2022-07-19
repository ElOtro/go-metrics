package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handlers) GetAllMetricsHandler(w http.ResponseWriter, r *http.Request) {
	gauges, counters := h.repo.GetAll()

	g, _ := json.Marshal(gauges)
	c, _ := json.Marshal(counters)

	s := fmt.Sprintf("gauges: %s\r\ncounters: %s\r\n", string(g), string(c))

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(s))
	if err != nil {
		log.Fatalln(err)
	}
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
	_, err = w.Write([]byte(s))
	if err != nil {
		log.Fatalln(err)
	}

}

func (h *Handlers) CreateMetricHandler(w http.ResponseWriter, r *http.Request) {
	t := chi.URLParam(r, "type")
	n := chi.URLParam(r, "name")
	v := chi.URLParam(r, "value")

	if t == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if n == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if v == "" {
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
