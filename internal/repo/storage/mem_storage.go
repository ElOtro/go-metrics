package storage

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type memStorage struct {
	mutex    sync.RWMutex
	Gauges   map[string]float64
	Counters map[string]int64
}

func New() *memStorage {
	m := &memStorage{
		mutex:    sync.RWMutex{},
		Gauges:   make(map[string]float64),
		Counters: make(map[string]int64),
	}
	return m
}

func (m *memStorage) GetAll() (map[string]float64, map[string]int64) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.Gauges, m.Counters
}

func (m *memStorage) Get(t, n string) (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	value := ""

	if t == "gauge" {
		v, ok := m.Gauges[n]
		if ok {
			value = fmt.Sprintf("%.3f", v)
		}
	}

	if t == "counter" {
		v, ok := m.Counters[n]
		if ok {
			value = fmt.Sprintf("%d", v)
		}
	}

	if value == "" {
		return "", errors.New("not found")
	}

	return value, nil
}

func (m *memStorage) Set(t, n, v string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if t == "gauge" {
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}

		m.Gauges[n] = value

		return nil
	}

	if t == "counter" {
		value, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}

		_, ok := m.Counters[n]

		if ok {
			m.Counters[n] = m.Counters[n] + value
		} else {
			m.Counters[n] = value
		}

		return nil
	}

	return errors.New("invalid params")
}

// New JSON API
func (m *memStorage) GetMetricsByID(id, mtype string) (*Metrics, error) {

	var input Metrics

	if mtype == "gauge" {
		v, ok := m.Gauges[id]
		if ok {
			input.ID = id
			input.MType = "gauge"
			input.Value = &v
		}
	}

	if mtype == "counter" {
		v, ok := m.Counters[id]
		if ok {
			input.ID = id
			input.MType = "counter"
			input.Delta = &v
		}
	}

	if input.ID == "" {
		return nil, errors.New("not found")
	}

	return &input, nil
}

func (m *memStorage) SetMetrics(ms Metrics) error {
	// fmt.Printf("%+v", ms)
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if ms.MType == "gauge" {
		if ms.Value == nil {
			return errors.New("invalid params")
		}

		m.Gauges[ms.ID] = *ms.Value

		return nil
	}

	if ms.MType == "counter" {
		if ms.Delta == nil {
			return errors.New("invalid params")
		}
		value, ok := m.Counters[ms.ID]

		if ok {
			m.Counters[ms.ID] = value + *ms.Delta
		} else {
			m.Counters[ms.ID] = *ms.Delta
		}

		return nil
	}

	return errors.New("invalid params")
}
