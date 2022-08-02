package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

// Run sending metrics for each type (gauge, counter)
func (app *application) postMetrics() {
	cfg := app.config
	var client = app.client
	var interval = time.Duration(app.config.reportInterval) * time.Second
	for {
		<-time.After(interval)

		// sending gauge metrics
		sendGauges(client, app.stats.Gauges, cfg.collectorSrv.address, cfg.collectorSrv.port)
		sendJSONGauges(client, app.stats.Gauges, cfg.collectorSrv.address, cfg.collectorSrv.port)
		// sending counter metrics
		sendCounters(client, app.stats.Counters, cfg.collectorSrv.address, cfg.collectorSrv.port)
		sendJSONCounters(client, app.stats.Counters, cfg.collectorSrv.address, cfg.collectorSrv.port)

	}

}

func sendGauges(client http.Client, gauges map[string]float64, address string, port int) {
	for k, v := range gauges {
		url := fmt.Sprintf("http://%s:%d/%s/%s/%s/%.2f", address, port, "update", "gauge", k, v)

		resp, err := client.Post(url, "text/plain", nil)
		if err != nil {
			panic(err)
		}
		resp.Body.Close()
	}

}

func sendCounters(client http.Client, counters map[string]int64, address string, port int) {
	for k, v := range counters {
		url := fmt.Sprintf("http://%s:%d/%s/%s/%s/%d", address, port, "update", "counter", k, v)

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
		if err != nil {
			fmt.Println(err.Error())
		}
		request.Header.Set("Content-Type", "text/plain")

		resp, err := client.Do(request)
		if err != nil {
			panic(err)
		}
		resp.Body.Close()
	}

}

func sendJSONGauges(client http.Client, gauges map[string]float64, address string, port int) {
	for k, v := range gauges {
		url := fmt.Sprintf("http://%s:%d/%s/", address, port, "update")

		m := Metrics{
			ID:    k,
			MType: "gauge",
			Value: &v,
		}

		// преобразуем m в JSON-формат
		js, err := json.Marshal(m)
		if err != nil {
			panic(err)
		}

		request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(js))
		if err != nil {
			fmt.Println(err.Error())
		}
		request.Header.Set("Content-Type", "application/json; charset=UTF-8")

		resp, err := client.Do(request)
		if err != nil {
			panic(err)
		}
		resp.Body.Close()
	}

}

func sendJSONCounters(client http.Client, counters map[string]int64, address string, port int) {
	for k, v := range counters {
		url := fmt.Sprintf("http://%s:%d/%s/", address, port, "update")

		m := Metrics{
			ID:    k,
			MType: "counter",
			Delta: &v,
		}

		// преобразуем m в JSON-формат
		js, err := json.Marshal(m)
		if err != nil {
			panic(err)
		}

		// fmt.Println(string(js))

		request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(js))
		if err != nil {
			fmt.Println(err.Error())
		}
		request.Header.Set("Content-Type", "application/json; charset=UTF-8")

		resp, err := client.Do(request)
		if err != nil {
			panic(err)
		}
		resp.Body.Close()
	}

}
