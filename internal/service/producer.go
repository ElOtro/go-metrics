package service

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/ElOtro/go-metrics/internal/repo"
	"github.com/ElOtro/go-metrics/internal/repo/storage"
)

type OutputMetrics struct {
	ID    string  `json:"id"`              // имя метрики
	MType string  `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type producer struct {
	storeInterval time.Duration
	filename      string
	repo          repo.Getter
}

func NewProducer(duration time.Duration, filename string, repo repo.Getter) (*producer, error) {

	p := &producer{
		storeInterval: duration,
		filename:      filename,
		repo:          repo,
	}

	return p, nil
}

func (p *producer) Run() {
	for {
		<-time.After(p.storeInterval)

		if err := p.WriteMetric(); err != nil {
			log.Println(err)
		}

	}

}

func (p *producer) WriteMetric() error {
	var metrics []OutputMetrics

	file, err := os.OpenFile(p.filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}

	gauges, counters := p.repo.GetAll()
	for k, v := range gauges {
		var metric = OutputMetrics{}

		metric.ID = k
		metric.MType = storage.Gauge
		metric.Value = v

		metrics = append(metrics, metric)
	}

	for k, v := range counters {
		var metric = OutputMetrics{}

		metric.ID = k
		metric.MType = storage.Counter
		metric.Delta = v

		metrics = append(metrics, metric)
	}

	js, err := json.Marshal(&metrics)
	if err != nil {
		return err
	}

	if _, err := file.Write(js); err != nil {
		file.Close()
		log.Println(err)
	}

	return nil
}
