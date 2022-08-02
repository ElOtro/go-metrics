package data

import (
	"math/rand"
	"runtime"
	"sync"
	"time"
)

// Keeps all types metrics
type Stats struct {
	mutex    sync.RWMutex
	Gauges   map[string]float64
	Counters map[string]int64
}

// Initialize new Stats
func NewStats(duration int) *Stats {

	s := &Stats{
		mutex:    sync.RWMutex{},
		Gauges:   make(map[string]float64),
		Counters: make(map[string]int64),
	}

	if duration > 0 {
		go runStats(s, duration)
	}

	return s
}

// Run to collect metrics with certain interval
func runStats(s *Stats, duration int) {
	var rtm runtime.MemStats
	var interval = time.Duration(duration) * time.Second

	for {
		<-time.After(interval)

		runtime.ReadMemStats(&rtm)
		updateStats(s, &rtm)
	}

}

// Update metrics from Stats struct
func updateStats(s *Stats, rtm *runtime.MemStats) {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.Gauges["Alloc"] = float64(rtm.Alloc)
	s.Gauges["BuckHashSys"] = float64(rtm.BuckHashSys)
	s.Gauges["Frees"] = float64(rtm.Frees)
	s.Gauges["GCCPUFraction"] = rtm.GCCPUFraction
	s.Gauges["GCSys"] = float64(rtm.GCSys)

	s.Gauges["HeapAlloc"] = float64(rtm.HeapAlloc)
	s.Gauges["HeapIdle"] = float64(rtm.HeapIdle)
	s.Gauges["HeapInuse"] = float64(rtm.HeapInuse)
	s.Gauges["HeapObjects"] = float64(rtm.HeapObjects)
	s.Gauges["HeapObjects"] = float64(rtm.HeapObjects)
	s.Gauges["HeapReleased"] = float64(rtm.HeapReleased)
	s.Gauges["HeapSys"] = float64(rtm.HeapSys)

	s.Gauges["LastGC"] = float64(rtm.LastGC)
	s.Gauges["Lookups"] = float64(rtm.Lookups)
	s.Gauges["MCacheInuse"] = float64(rtm.MCacheInuse)
	s.Gauges["MCacheSys"] = float64(rtm.MCacheSys)
	s.Gauges["MSpanInuse"] = float64(rtm.MSpanInuse)
	s.Gauges["MSpanSys"] = float64(rtm.MSpanSys)
	s.Gauges["Mallocs"] = float64(rtm.Mallocs)
	s.Gauges["NextGC"] = float64(rtm.NextGC)
	s.Gauges["NumForcedGC"] = float64(rtm.NumForcedGC)
	s.Gauges["NumGC"] = float64(rtm.NumGC)

	s.Gauges["OtherSys"] = float64(rtm.OtherSys)
	s.Gauges["PauseTotalNs"] = float64(rtm.PauseTotalNs)
	s.Gauges["StackInuse"] = float64(rtm.StackInuse)
	s.Gauges["StackSys"] = float64(rtm.StackSys)
	s.Gauges["Sys"] = float64(rtm.Sys)
	s.Gauges["TotalAlloc"] = float64(rtm.TotalAlloc)

	rand.Seed(time.Now().UnixNano())
	s.Gauges["RandomValue"] = rand.Float64()

	s.Counters["PollCount"] = s.Counters["PollCount"] + 1
}
