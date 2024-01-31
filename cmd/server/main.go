package main

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type MemStorage struct {
	GaugeMutex   *sync.RWMutex
	CounterMutex *sync.RWMutex
	Gauges       map[string]float64
	Counters     map[string]int64
}

func NewMemStorage() *MemStorage {
	storage := new(MemStorage)
	storage.Gauges = make(map[string]float64)
	storage.Counters = make(map[string]int64)
	storage.GaugeMutex = new(sync.RWMutex)
	storage.CounterMutex = new(sync.RWMutex)
	return storage
}

type Config struct {
	ListenUri string
}

var storage *MemStorage

func main() {
	storage = NewMemStorage()

	config := Config{ListenUri: ":8080"}

	mux := http.NewServeMux()
	mux.HandleFunc("/update/", handleMetric)
	err := http.ListenAndServe(config.ListenUri, mux)
	if err != nil {
		panic(err)
	}
}

func handleMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != `text/plain` {
		http.Error(w, "Wrong content-type", http.StatusBadRequest)
		return
	}

	metricType := ""

	path, ok := strings.CutPrefix(r.URL.Path, "/update/counter")
	if ok {
		metricType = "counter"
	} else if path, ok = strings.CutPrefix(r.URL.Path, "/update/gauge"); ok {
		metricType = "gauge"
	}

	path, _ = strings.CutPrefix(path, "/")

	if metricType == "" {
		http.Error(w, "Bad Metric Type", http.StatusBadRequest)
		return
	}

	fragments := strings.Split(path, "/")

	if len(fragments) == 0 || len(fragments) == 1 && fragments[0] == "" {
		http.Error(w, "Metric not found", http.StatusNotFound)
		return
	}

	if len(fragments) != 2 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	metric := fragments[0]

	if metric == "" {
		http.Error(w, "Please specify metric name", http.StatusNotFound)
		return
	}

	valueStr := fragments[1]

	switch metricType {
	case "counter":
		handleCounter(w, r, metric, valueStr)
	case "gauge":
		handleGauge(w, r, metric, valueStr)
	default:
		http.Error(w, "Bad Metric Type", http.StatusBadRequest)
	}
}

func handleCounter(w http.ResponseWriter, r *http.Request, metric string, valueStr string) {

	value, err := strconv.ParseInt(valueStr, 10, 64)

	if err != nil {
		http.Error(w, "Invalid parameter value", http.StatusBadRequest)
		return
	}

	storage.CounterMutex.Lock()

	if _, ok := storage.Counters[metric]; !ok {
		storage.Counters[metric] = 0
	}

	storage.Counters[metric] += value
	storage.CounterMutex.Unlock()
	w.WriteHeader(http.StatusOK)
}

func handleGauge(w http.ResponseWriter, r *http.Request, metric string, valueStr string) {

	value, err := strconv.ParseFloat(valueStr, 64)

	if err != nil {
		http.Error(w, "Invalid parameter value", http.StatusBadRequest)
		return
	}

	storage.GaugeMutex.Lock()
	storage.Gauges[metric] = value
	storage.GaugeMutex.Unlock()
	w.WriteHeader(http.StatusOK)
}
