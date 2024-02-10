package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func Router(storage repository) chi.Router {
	r := chi.NewRouter()
	r.Route("/update", func(r chi.Router) {
		r.Post("/counter/{metric}/{value}", handleCounterFunc(storage))
		r.Post("/gauge/{metric}/{value}", handleGaugeFunc(storage))
		r.Post("/counter/{metric}", handleInvalidRequest)
		r.Post("/gauge/{metric}", handleInvalidRequest)
		r.Post("/counter/", handleUnknownMetric)
		r.Post("/gauge/", handleUnknownMetric)
		r.Post("/{type}/*", handleUnknownType)
	})

	r.Get("/value/{metric_type}/{metric}", handleMetricValueFunc(storage))
	r.Get("/", handleHTMLMetricFunc(storage))

	return r
}

func handleUnknownType(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Bad request", http.StatusBadRequest)
}

func handleUnknownMetric(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Invalid metrics", http.StatusNotFound)
}

func handleInvalidRequest(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Invalid request", http.StatusBadRequest)
}

func handleHTMLMetricFunc(storage repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		gauges := storage.GetGauges()
		counters := storage.GetCounters()

		html := "<html><head><title>Metric list</title><head><body>"
		html += "<h1>All metrics</h1>"
		html += "<h2>Gauges</h2>"
		html += "<table><tr><th>Metric</th><th>Value</th></tr>"
		for k, v := range gauges {
			html += fmt.Sprintf("<tr><td>%s</td><td>%f</td></tr>", k, v)
		}

		html += "</table>"

		html += "<h2>Counters</h2>"
		html += "<table><tr><th>Metric</th><th>Value</th></tr>"
		for k, v := range counters {
			html += fmt.Sprintf("<tr><td>%s</td><td>%d</td></tr>", k, v)
		}
		html += "</table>"
		html += "</body></html>"

		w.Write([]byte(html))
		w.WriteHeader(http.StatusOK)

	}
}

func handleMetricValueFunc(storage repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		metricType := chi.URLParam(r, "metric_type")

		if metricType != "counter" && metricType != "gauge" {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		metric := chi.URLParam(r, "metric")
		if metricType == "counter" {
			value, err := storage.GetCounter(metric)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(err.Error()))
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(fmt.Sprintf("%d", value)))
			}
			return
		}

		if metricType == "gauge" {
			value, err := storage.GetGauge(metric)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(err.Error()))
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(fmt.Sprintf("%f", value)))
			}

			return
		}

	}
}

func handleCounterFunc(storage repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		metric := chi.URLParam(r, "metric")
		valueStr := chi.URLParam(r, "value")

		value, err := strconv.ParseInt(valueStr, 10, 64)

		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid parameter value %s", r.URL.Path), http.StatusBadRequest)
			return
		}

		storage.SetCounter(metric, value)
		fmt.Printf("Set counter: %s = %d\n", metric, value)
		w.WriteHeader(http.StatusOK)
	}
}

func handleGaugeFunc(storage repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		metric := chi.URLParam(r, "metric")
		valueStr := chi.URLParam(r, "value")

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
