package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func splitFragments(path string) []string {
	path, _ = strings.CutPrefix(path, "/")
	return strings.Split(path, "/")
}

func extractMetricAndValue(path string) (string, string, error) {

	fragments := splitFragments(path)

	if len(fragments) != 2 {
		return "", "", HandleError{"Invalid request", http.StatusBadRequest}
	}

	metric := fragments[0]

	if metric == "" {
		return "", "", HandleError{"Please specify metric name", http.StatusNotFound}
	}

	valueStr := fragments[1]

	return metric, valueStr, nil
}

func handleCounterFunc(storage repository) func (w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
			return
		}
	
		path, _ := strings.CutPrefix(r.URL.Path, "/update/counter/")
		if path == "" {
			http.Error(w, "Invalid metrics", http.StatusNotFound)
			return
		}
	
		metric, valueStr, err := extractMetricAndValue(path)
		if err != nil {
			if err1, ok := err.(HandleError); ok {
				http.Error(w, err1.Error(), err1.Status())
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	
		value, err := strconv.ParseInt(valueStr, 10, 64)
	
		if err != nil {
			http.Error(w, "Invalid parameter value", http.StatusBadRequest)
			return
		}
			
		storage.SetCounter(metric, value)
		fmt.Printf("Set counter: %s = %d\n", metric, value)
		w.WriteHeader(http.StatusOK)
	}
}

func handleGaugeFunc(storage repository) func (w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
			return
		}
		
		path, _ := strings.CutPrefix(r.URL.Path, "/update/gauge/")
		if path == "" {
			http.Error(w, "Invalid metrics", http.StatusNotFound)
			return
		}
	
		metric, valueStr, err := extractMetricAndValue(path)
		if err != nil {
			if err1, ok := err.(HandleError); ok {
				http.Error(w, err1.Error(), err1.Status())
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	
		value, err := strconv.ParseFloat(valueStr, 64)
	
		if err != nil {
			http.Error(w, "Invalid parameter value", http.StatusBadRequest)
			return
		}
		
		storage.SetGauge(metric, value)	
		fmt.Printf("Set gauge: %s = %f\n", metric, value)
		w.WriteHeader(http.StatusOK)
	}
}
