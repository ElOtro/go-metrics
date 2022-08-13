package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"strconv"

	"github.com/ElOtro/go-metrics/internal/repo/storage"
)

// Define a new HashMetric type which contains a key.
type HashMetric struct {
	key     []byte
	UseHash bool
}

// New is a helper which creates a new HashMetric instance with a key.
func NewHashMetric(key string) *HashMetric {
	var useHash bool
	if key != "" {
		useHash = true
	}
	return &HashMetric{key: []byte(key), UseHash: useHash}
}

// ValidAgentHash from an agent header.
func (hm *HashMetric) ValidAgentHash(hash, mtype, id, value string) bool {
	expectedHash := ""

	if mtype == storage.Counter {
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return false
		}
		m := storage.Metrics{
			ID:    id,
			MType: storage.Counter,
			Delta: &v,
		}

		expectedHash = hm.Hash(&m)

	}

	if mtype == storage.Gauge {
		value, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return false
		}
		m := storage.Metrics{
			ID:    id,
			MType: storage.Gauge,
			Value: &value,
		}

		expectedHash = hm.Hash(&m)

	}

	return hash == expectedHash

}

func (hm *HashMetric) Valid(hash string, m *storage.Metrics) bool {
	expectedHash := hm.Hash(m)

	return hash == expectedHash
}

func (hm *HashMetric) Hash(m *storage.Metrics) string {
	var metric string
	if m.MType == storage.Counter {
		metric = fmt.Sprintf("%s:%s:%d", m.ID, m.MType, *m.Delta)
	}

	if m.MType == storage.Gauge {
		metric = fmt.Sprintf("%s:%s:%f", m.ID, m.MType, *m.Value)
	}

	h := hmac.New(sha256.New, hm.key)
	h.Write([]byte(metric))

	return fmt.Sprintf("%x", h.Sum(nil))
}
