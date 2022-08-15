package repo

import (
	"errors"

	"github.com/ElOtro/go-metrics/internal/repo/storage"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Options struct {
	Memory bool
	DB     *pgxpool.Pool
}

var ErrEmptyOptions = errors.New("empty options")
var ErrInvalidOptions = errors.New("invalid options")

type Getter interface {
	GetAll() (map[string]float64, map[string]int64)
	Get(t, n string) (string, error)
	Set(t, n, v string) error
	GetMetricsByID(id, mtype string) (*storage.Metrics, error)
	SetMetrics(*storage.Metrics) error
	RestoreMetrics(filename string) error
	GetHealth() error
}

func NewRepo(opts *Options) (Getter, error) {
	return storage.NewMemStorage(), nil
	// switch opts.Memory {
	// case true:
	// 	return storage.NewMemStorage(), nil
	// case false:
	// 	return storage.NewPgStorage(opts.DB), nil
	// default:
	// 	return nil, ErrInvalidOptions
	// }

}
