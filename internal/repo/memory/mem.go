package memory

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
)

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

func (m *memStorage) GetAll() string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	g, _ := json.Marshal(m.Gauges)
	c, _ := json.Marshal(m.Counters)
	return fmt.Sprintf("gauges: %s\r\ncounters: %s\r\n", string(g), string(c))
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
	m.mutex.RLock()
	defer m.mutex.RUnlock()

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
