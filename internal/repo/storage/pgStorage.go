package storage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type pgStorage struct {
	db *pgxpool.Pool
}

func NewPgStorage(db *pgxpool.Pool) *pgStorage {
	return &pgStorage{db: db}
}

func (pg *pgStorage) GetAll() (map[string]float64, map[string]int64) {
	return nil, nil
}

func (pg *pgStorage) Get(t, n string) (string, error) {

	value := ""

	return value, nil
}

func (pg *pgStorage) Set(t, n, v string) error {

	return ErrInvalidPrams
}

// New JSON API
func (pg *pgStorage) GetMetricsByID(id, mtype string) (*Metrics, error) {

	var input Metrics

	return &input, nil
}

func (pg *pgStorage) SetMetrics(ms *Metrics) error {

	return ErrInvalidPrams
}

func (pg *pgStorage) RestoreMetrics(filename string) error {

	return nil
}

func (pg *pgStorage) GetHealth() error {
	// Create a context with a 5-second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := pg.db.Ping(ctx)
	if err != nil {
		return err
	}

	return nil
}
