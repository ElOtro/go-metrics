package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"sync"
)

type memStorage struct {
	mutex    sync.RWMutex
	Gauges   map[string]float64
	Counters map[string]int64
}

func NewMemStorage() *memStorage {
	m := &memStorage{
		mutex:    sync.RWMutex{},
		Gauges:   make(map[string]float64),
		Counters: make(map[string]int64),
	}
	return m
}

func (m *memStorage) List() ([]*Metrics, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	metrics := []*Metrics{}

	for k, v := range m.Gauges {
		value := v
		metrics = append(metrics, &Metrics{ID: k, MType: Gauge, Value: &value})
	}

	for k, v := range m.Counters {
		value := v
		metrics = append(metrics, &Metrics{ID: k, MType: Counter, Delta: &value})
	}

	return metrics, nil
}

func (m *memStorage) Get(t, n string) (*Metrics, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	metric := Metrics{}

	if t == Gauge {
		v, ok := m.Gauges[n]
		if ok {
			metric.ID = n
			metric.MType = Gauge
			metric.Value = &v
			// value = fmt.Sprintf("%.3f", v)
		}
	}

	if t == Counter {
		v, ok := m.Counters[n]
		if ok {
			metric.ID = n
			metric.MType = Counter
			metric.Delta = &v
			// value = fmt.Sprintf("%d", v)
		}
	}

	if metric.ID == "" {
		return nil, ErrNotFound
	}

	return &metric, nil
}

func (m *memStorage) Set(t, n, v string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	fmt.Println(t, n, v)

	if t == Gauge {
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}

		m.Gauges[n] = value

		return nil
	}

	if t == Counter {
		value, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}

		v, ok := m.Counters[n]

		if ok {
			m.Counters[n] = v + value
		} else {
			m.Counters[n] = value
		}

		return nil
	}

	return ErrInvalidPrams
}

// New JSON API
func (m *memStorage) GetMetricsByID(id, mtype string) (*Metrics, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var input Metrics

	if mtype == Gauge {
		v, ok := m.Gauges[id]
		if ok {
			input.ID = id
			input.MType = Gauge
			input.Value = &v
		}
	}

	if mtype == Counter {
		v, ok := m.Counters[id]
		if ok {
			input.ID = id
			input.MType = Counter
			input.Delta = &v
		}
	}

	if input.ID == "" {
		return nil, ErrNotFound
	}

	return &input, nil
}

func (m *memStorage) SetMetrics(ms *Metrics) error {
	// fmt.Printf("%+v", ms)
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if ms.MType == Gauge {
		if ms.Value == nil {
			return ErrInvalidPrams
		}

		m.Gauges[ms.ID] = *ms.Value

		return nil
	}

	if ms.MType == Counter {
		if ms.Delta == nil {
			return ErrInvalidPrams
		}
		value, ok := m.Counters[ms.ID]

		if ok {
			m.Counters[ms.ID] = value + *ms.Delta
		} else {
			m.Counters[ms.ID] = *ms.Delta
		}

		return nil
	}

	return ErrInvalidPrams
}

// update/insert bulk Metrics
func (m *memStorage) SetBatchMetrics([]*Metrics) error {
	return nil
}

func (m *memStorage) RestoreMetrics(filename string) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	metrics := []Metrics{}
	err = json.Unmarshal([]byte(file), &metrics)
	if err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, metric := range metrics {
		if metric.MType == Gauge {
			if metric.Value != nil {
				m.Gauges[metric.ID] = *metric.Value
			}
		}

		if metric.MType == Counter {
			if metric.Delta != nil {
				m.Counters[metric.ID] = *metric.Delta
			}
		}
	}

	return nil
}

func (m *memStorage) GetHealth() error {
	// nothing to return
	return nil
}
