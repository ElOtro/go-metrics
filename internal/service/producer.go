package service

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ElOtro/go-metrics/internal/repo"
)

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

		fmt.Println("-------------------------")
		fmt.Println("File is saved")

		if err := p.WriteMetric(); err != nil {
			log.Fatal(err)
		}

	}

}

func (p *producer) WriteMetric() error {
	if err := os.Remove(p.filename); err != nil {
		log.Fatal(err)
		return err
	}

	file, err := os.OpenFile(p.filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	metrics := p.repo.GetMetrics()
	for _, metric := range metrics {
		if err := encoder.Encode(&metric); err != nil {
			log.Fatal(err)
			return err
		}
	}

	if err := file.Close(); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

// func (p *producer) Close() error {
// 	return p.file.Close()
// }
