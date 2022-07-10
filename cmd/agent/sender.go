package main

import (
	"fmt"
	"time"
)

// Run senging metrics for each type (gauge, counter)
func (app *application) postMetricsHandler() {
	cfg := app.config
	client := app.client

	var interval = time.Duration(app.config.reportInterval) * time.Second
	for {
		<-time.After(interval)

		// sending gauge metrics
		for k, v := range app.stats.Gauges {
			url := fmt.Sprintf("http://%s:%d/%s/%s/%s/%.2f", cfg.collectorSrv.address, cfg.collectorSrv.port, "update", "gauge", k, v)

			resp, err := client.Post(url, "text/plain", nil)
			if err != nil {
				fmt.Println("Error")
			}
			resp.Body.Close()
		}

		// sending counter metrics
		for k, v := range app.stats.Counters {
			url := fmt.Sprintf("http://%s:%d/%s/%s/%s/%d", cfg.collectorSrv.address, cfg.collectorSrv.port, "update", "counter", k, v)

			resp, err := client.Post(url, "text/plain", nil)
			if err != nil {
				fmt.Println("Error")
			}
			resp.Body.Close()
		}

	}

}
