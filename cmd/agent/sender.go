package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ElOtro/go-metrics/internal/repo/storage"
	"github.com/ElOtro/go-metrics/internal/service"
)

// Run sending metrics for each type (gauge, counter)
func (app *application) postMetrics() {
	cfg := *app.config
	var client = app.client
	var interval = cfg.ReportInterval
	var hm = service.NewHashMetric(cfg.Key)

	for {
		<-time.After(interval)

		// sending gauge metrics
		err := postGauges(client, app.stats.Gauges, cfg.Address, hm)
		if err != nil {
			log.Println(err)
		}
		// sending counter metrics
		err = postCounters(client, app.stats.Counters, cfg.Address, hm)
		if err != nil {
			log.Println(err)
		}

	}

}

func postGauges(client http.Client, gauges map[string]float64, address string, hm *service.HashMetric) error {
	for k, v := range gauges {
		url := fmt.Sprintf("http://%s/%s/%s/%s/%.2f", address, "update", storage.Gauge, k, v)

		req, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			return err
		}

		if hm.UseHash {
			m := storage.Metrics{
				ID:    k,
				MType: storage.Gauge,
				Value: &v,
			}

			req.Header.Add("Hash", hm.Hash(&m))
		}

		req.Header.Add("Content-Type", "text/plain")
		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		resp.Body.Close()

	}

	return nil

}

func postCounters(client http.Client, counters map[string]int64, address string, hm *service.HashMetric) error {
	for k, v := range counters {
		url := fmt.Sprintf("http://%s/%s/%s/%s/%d", address, "update", storage.Counter, k, v)

		req, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			return err
		}

		if hm.UseHash {
			m := storage.Metrics{
				ID:    k,
				MType: storage.Counter,
				Delta: &v,
			}

			req.Header.Add("Hash", hm.Hash(&m))
		}

		req.Header.Add("Content-Type", "text/plain")
		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		resp.Body.Close()
	}

	return nil

}
