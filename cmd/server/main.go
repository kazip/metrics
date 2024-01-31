package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type MemStorage struct {
	GaugeMutex   sync.RWMutex
	CounterMutex sync.RWMutex
	Gauges       map[string]float64
	Counters     map[string]int64
}

type Config struct {
	ListenUri string
}

var storage *MemStorage

func main() {
	storage = new(MemStorage)
	storage.Gauges = map[string]float64{
		"gauge1": 0.0,
		"gauge2": 0.0,
	}

	storage.Counters = map[string]int64{
		"counter1": 0,
		"counter2": 0,
	}

	config := Config{ListenUri: ":8080"}

	mux := http.NewServeMux()
	mux.HandleFunc("/update/", handleMetric)
	err := http.ListenAndServe(config.ListenUri, mux)
	if err != nil {
		panic(err)
	}
}

func handleGauge(w http.ResponseWriter, r *http.Request, metric string, valueStr string) {

	value, err := strconv.ParseFloat(valueStr, 64)

	if err != nil {
		http.Error(w, "Invalid parameter value", http.StatusBadRequest)
		return
	}

	storage.GaugeMutex.Lock()

	if _, ok := storage.Gauges[metric]; !ok {
		http.Error(w, "Invalid gaguge name", http.StatusBadRequest)
		return
	}

	storage.Gauges[metric] = value

	storage.GaugeMutex.Unlock()

	w.WriteHeader(http.StatusOK)
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

	fmt.Println(r.URL.Path)

	metricType := ""

	path, ok := strings.CutPrefix(r.URL.Path, "/update/counter/")
	if ok {
		metricType = "counter"
	} else {
		path, ok = strings.CutPrefix(r.URL.Path, "/update/gauge/")
		if ok {
			metricType = "gauge"
		}
	}

	if metricType == "" {
		http.Error(w, "Bad Metric Type", http.StatusBadRequest)
		return
	}

	fmt.Println(path)
	fragments := strings.Split(path, "/")
	if len(fragments) != 2 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	fmt.Println(fragments)

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
		http.Error(w, "Invalid counter name", http.StatusBadRequest)
		return
	}

	storage.Counters[metric] += value

	storage.CounterMutex.Unlock()

	w.WriteHeader(http.StatusOK)

	storage.CounterMutex.RLock()
	io.WriteString(w, fmt.Sprintf("%v", storage.Counters[metric]))
	storage.CounterMutex.RUnlock()
}
