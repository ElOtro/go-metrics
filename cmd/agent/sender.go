package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ElOtro/go-metrics/internal/repo/storage"
)

// Run sending metrics for each type (gauge, counter)
func (app *application) postMetrics() {
	cfg := *app.config
	var client = app.client
	var interval = cfg.ReportInterval
	for {
		<-time.After(interval)

		// sending gauge metrics
		sendGauges(client, app.stats.Gauges, cfg.Address)
		// sending counter metrics
		sendCounters(client, app.stats.Counters, cfg.Address)

	}

}

func sendGauges(client http.Client, gauges map[string]float64, address string) {
	for k, v := range gauges {
		url := fmt.Sprintf("http://%s/%s/%s/%s/%.2f", address, "update", storage.Gauge, k, v)

		resp, err := client.Post(url, "text/plain", nil)
		if err != nil {
			log.Println(err)
		} else {
			resp.Body.Close()
		}

	}

}

func sendCounters(client http.Client, counters map[string]int64, address string) {
	for k, v := range counters {
		url := fmt.Sprintf("http://%s/%s/%s/%s/%d", address, "update", storage.Counter, k, v)

		resp, err := client.Post(url, "text/plain", nil)
		if err != nil {
			log.Println(err)
		} else {
			resp.Body.Close()
		}
	}

}
