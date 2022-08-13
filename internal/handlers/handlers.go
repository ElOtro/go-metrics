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

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
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

	hash := r.Header.Get("Hash")

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

	if t != storage.Gauge && t != storage.Counter {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	// Check if a hash from a header request is not valid
	if h.hm.UseHash && !h.hm.ValidAgentHash(hash, t, n, v) {
		w.WriteHeader(http.StatusBadRequest)
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
// GetMetricsJSONHandler respond to POST /value/
func (h *Handlers) GetMetricsJSONHandler(w http.ResponseWriter, r *http.Request) {

	var input storage.Metrics

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&input)

	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	m, err := h.repo.GetMetricsByID(input.ID, input.MType)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if h.hm.UseHash {
		hash := h.hm.Hash(m)
		m.Hash = hash
		w.Header().Set("Hash", hash)
	}

	// преобразуем m в JSON-формат
	js, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(js)
	if err != nil {
		log.Println(err)
	}

}

func (h *Handlers) CreateMetricsJSONHandler(w http.ResponseWriter, r *http.Request) {
	// hash := r.Header.Get("Hash")
	var input *storage.Metrics

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&input)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// if h.hm.UseHash && !h.hm.Valid(hash, input) {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	err = h.repo.SetMetrics(input)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if h.hm.UseHash {
		w.Header().Set("Hash", h.hm.Hash(input))
	}
	w.WriteHeader(http.StatusOK)

}
