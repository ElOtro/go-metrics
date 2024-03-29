package storage

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type pgStorage struct {
	db *pgxpool.Pool
}

func NewPgStorage(db *pgxpool.Pool) *pgStorage {
	return &pgStorage{db: db}
}

func (pg *pgStorage) List() ([]*Metrics, error) {
	// Construct the SQL query to retrieve all records.
	query := "SELECT name, type, delta, value FROM metrics"

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := pg.db.Query(ctx, query)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer rows.Close()

	metrics := []*Metrics{}

	// Use rows.Next to iterate through the rows in the resultset.
	for rows.Next() {
		// Initialize an empty OutputMetric struct to hold the data for an individual metric.
		metric := Metrics{}

		// Scan the values from the row into the OutputMetric struct.
		err := rows.Scan(
			&metric.ID,
			&metric.MType,
			&metric.Delta,
			&metric.Value,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		// Add the Unit struct to the slice.
		metrics = append(metrics, &metric)
	}

	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return metrics, nil
}

func (pg *pgStorage) Get(t, n string) (*Metrics, error) {
	// Define the SQL query for retrieving data.
	query := "SELECT name, type, delta, value FROM metrics WHERE type = $1 AND name = $2"
	// Declare a Metrics struct to hold the data returned by the query.
	var metric Metrics

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	// Importantly, use defer to make sure that we cancel the context before the Get()
	// method returns.
	defer cancel()

	// Execute the query using the QueryRow() method, passing in the provided id value
	err := pg.db.QueryRow(ctx, query, t, n).Scan(
		&metric.ID,
		&metric.MType,
		&metric.Delta,
		&metric.Value,
	)

	// Handle any errors. If there was no matching found, Scan() will return
	// a sql.ErrNoRows error. We check for this and return our custom ErrRecordNotFound
	// error instead.
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &metric, nil
}

func (pg *pgStorage) Set(t, n, v string) error {
	var delta *int64
	var value *float64

	if t == Counter {
		val, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
		delta = &val
	}

	if t == Gauge {
		val, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}
		value = &val
	}

	metric := Metrics{ID: n, MType: t, Delta: delta, Value: value}

	// Define the SQL query for inserting a new record
	query := setMetricQuery()
	// Create an args slice containing the values for the placeholder parameters.
	args := []interface{}{metric.ID, metric.MType, metric.Delta, metric.Value}
	// Use the QueryRow() method to execute the query, passing in the args slice as a
	// variadic parameter and scanning the new version value into the metric struct.
	err := pg.db.QueryRow(context.Background(), query, args...).Scan(&metric.Delta, &metric.Value)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// New JSON API
func (pg *pgStorage) GetMetricsByID(id, mtype string) (*Metrics, error) {

	// Define the SQL query for retrieving data.
	query := "SELECT name, type, delta, value FROM metrics WHERE type = $1 AND name = $2"
	// Declare a Metric struct to hold the data returned by the query.
	metric := Metrics{}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	// Importantly, use defer to make sure that we cancel the context before the Get()
	// method returns.
	defer cancel()

	// Execute the query using the QueryRow() method, passing in the provided id value
	err := pg.db.QueryRow(ctx, query, mtype, id).Scan(
		&metric.ID,
		&metric.MType,
		&metric.Delta,
		&metric.Value,
	)

	// Handle any errors. If there was no matching found, Scan() will return
	// a pgx.ErrNoRows error. We check for this and return our custom ErrNotFound
	// error instead.
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &metric, nil
}

func (pg *pgStorage) SetMetrics(m *Metrics) error {
	// Define the SQL query for inserting a new record
	query := setMetricQuery()
	// Create an args slice containing the values for the placeholder parameters.
	args := []interface{}{m.ID, m.MType, m.Delta, m.Value}
	// Use the QueryRow() method to execute the query, passing in the args slice as a
	// variadic parameter and scanning the new version value into the metric struct.
	err := pg.db.QueryRow(context.Background(), query, args...).Scan(&m.Delta, &m.Value)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// SetBatchMetrics updates/inserts bulk of Metrics
func (pg *pgStorage) SetBatchMetrics(metrics []*Metrics) error {
	// Define the SQL query for inserting a new record
	query := `INSERT INTO metrics (name, type, delta, value) 
	          VALUES ($1, $2, $3, $4)
			  ON CONFLICT (name) DO UPDATE SET type=$2, delta=$3 + COALESCE(metrics.delta, 0), value=$4`

	batch := &pgx.Batch{}
	for _, v := range metrics {
		batch.Queue(query, v.ID, v.MType, v.Delta, v.Value)
	}

	br := pg.db.SendBatch(context.Background(), batch)

	_, err := br.Exec()
	if err != nil {
		log.Println(err)
		return err
	}

	err = br.Close()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// RestoreMetrics restores metrics from an json file
func (pg *pgStorage) RestoreMetrics(filename string) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	metrics := []Metrics{}
	err = json.Unmarshal([]byte(file), &metrics)
	if err != nil {
		return err
	}

	for _, metric := range metrics {
		// Define the SQL query for inserting a new record
		query := `INSERT INTO metrics (name, type, delta, value) 
		          VALUES ($1, $2, $3, $4) 
				  ON CONFLICT (name) DO UPDATE SET type=$2, delta=$3, value=$4
				  RETURNING delta, value`
		// Create an args slice containing the values for the placeholder parameters.
		args := []interface{}{metric.ID, metric.MType, metric.Delta, metric.Value}
		// Use the QueryRow() method to execute the SQL query on our connection pool
		err := pg.db.QueryRow(context.Background(), query, args...).Scan(&metric.Delta, &metric.Value)
		if err != nil {
			log.Println(err)
			return err
		}
	}

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

func setMetricQuery() string {
	// Define the SQL query for inserting a new record
	query := `INSERT INTO metrics (name, type, delta, value) 
	          VALUES ($1, $2, $3, $4)
			  ON CONFLICT (name) DO UPDATE SET type=$2, delta=$3 + COALESCE(metrics.delta, 0), value=$4
			  RETURNING delta, value`

	return query
}
