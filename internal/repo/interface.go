package repo

import (
	"github.com/ElOtro/go-metrics/internal/repo/storage"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Options struct {
	Memory bool
	DB     *pgxpool.Pool
}

type Getter interface {
	GetAll() (map[string]float64, map[string]int64)
	Get(t, n string) (string, error)
	Set(t, n, v string) error
	GetMetricsByID(id, mtype string) (*storage.Metrics, error)
	SetMetrics(*storage.Metrics) error
	RestoreMetrics(filename string) error
	GetHealth() error
}

func NewRepo(options *Options) (Getter, error) {
	if !options.Memory {
		return storage.NewPgStorage(options.DB), nil
	}

	return storage.NewMemStorage(), nil

}
