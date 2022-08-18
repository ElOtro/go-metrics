package service

import (
	"encoding/json"
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

		if err := p.WriteMetric(); err != nil {
			log.Println(err)
		}

	}

}

func (p *producer) WriteMetric() error {
	file, err := os.OpenFile(p.filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}

	metrics, err := p.repo.List()
	if err != nil {
		return err
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
