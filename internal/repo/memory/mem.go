package memory

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type Item struct {
	Type  string
	Value interface{}
}

type memstorage struct {
	Metrics map[string]Item
}

func New() *memstorage {
	return &memstorage{make(map[string]Item)}
}

func (m memstorage) GetAll() string {
	b, _ := json.Marshal(m.Metrics)
	return fmt.Sprintf("memory: %s", string(b))
}

func (m memstorage) Get(t, n string) (string, error) {
	metrics := m.Metrics

	for k, v := range metrics {
		if k == n && v.Type == t {
			return fmt.Sprintf("%v", v.Value), nil
		} else {
			return "", errors.New("not found")
		}
	}

	return "", nil
}

func (m memstorage) Set(t, n, v string) error {
	metrics := m.Metrics

	if t == "gauge" {
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}

		metrics[n] = Item{t, value}

		return nil
	}

	if t == "counter" {
		value, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}

		_, ok := metrics[n]

		if ok {
			v := metrics[n].Value.(int64)
			metrics[n] = Item{t, v + value}
		} else {
			metrics[n] = Item{t, value}
		}

		v := metrics[n].Value.(int64)
		metrics[n] = Item{t, v + value}

		return nil
	}

	return errors.New("invalid params")
}
