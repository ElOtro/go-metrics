package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Run sending metrics for each type (gauge, counter)
func (app *application) postMetrics() {
	cfg := app.config
	var client = app.client
	var interval = time.Duration(app.config.ReportInterval) * time.Second
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
		url := fmt.Sprintf("http://%s/%s/%s/%s/%.2f", address, "update", "gauge", k, v)

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
		url := fmt.Sprintf("http://%s/%s/%s/%s/%d", address, "update", "counter", k, v)

		resp, err := client.Post(url, "text/plain", nil)
		if err != nil {
			log.Println(err)
		} else {
			resp.Body.Close()
		}
	}

}
