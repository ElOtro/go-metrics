package repo

import (
	"errors"

	memStorage "github.com/ElOtro/go-metrics/internal/repo/storage"
)

type Options struct {
	Environment string
}

var ErrEmptyOptions = errors.New("empty options")
var ErrInvalidOptions = errors.New("invalid options")

type Getter interface {
	GetAll() (map[string]float64, map[string]int64)
	Get(t, n string) (string, error)
	Set(t, n, v string) error
	GetMetricsByID(id, mtype string) (*memStorage.Metrics, error)
	SetMetrics(*memStorage.Metrics) error
	RestoreMetrics(filename string) error
}

func NewMemStorage() (Getter, error) {
	return memStorage.New(), nil
}
