package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ElOtro/go-metrics/internal/repo/storage"
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

//  New API
func (h *Handlers) GetMetricsJSONHandler(w http.ResponseWriter, r *http.Request) {

	var input storage.Metrics

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&input)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	m, err := h.repo.GetMetricsByID(input.ID, input.MType)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// преобразуем m в JSON-формат
	js, err := json.Marshal(m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)

}

func (h *Handlers) CreateMetricsJSONHandler(w http.ResponseWriter, r *http.Request) {

	var input storage.Metrics

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&input)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.repo.SetMetrics(input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}
