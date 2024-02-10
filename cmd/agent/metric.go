package main

import "runtime"

type GaugeMetric struct {
	name  string
	value float64
}

type CounterMetric struct {
	name  string
	value int64
}

type Metrics struct {
	gauges   []GaugeMetric
	counters []CounterMetric
}

func (m *Metrics) fromRtm(rtm runtime.MemStats, pollCount int64, randomValue float64) {

	m.counters = []CounterMetric{
		{
			name:  "PollCount",
			value: pollCount,
		},
	}

	m.gauges = []GaugeMetric{
		{
			name:  "Alloc",
			value: float64(rtm.Alloc),
		},
		{
			name:  "BuckHashSys",
			value: float64(rtm.BuckHashSys),
		},
		{
			name:  "Frees",
			value: float64(rtm.Frees),
		},
		{
			name:  "GCCPUFraction",
			value: float64(rtm.GCCPUFraction),
		},
		{
			name:  "GCSys",
			value: float64(rtm.GCSys),
		},
		{
			name:  "HeapAlloc",
			value: float64(rtm.HeapAlloc),
		},
		{
			name:  "HeapIdle",
			value: float64(rtm.HeapIdle),
		},
		{
			name:  "HeapInuse",
			value: float64(rtm.HeapInuse),
		},
		{
			name:  "HeapObjects",
			value: float64(rtm.HeapObjects),
		},
		{
			name:  "HeapReleased",
			value: float64(rtm.HeapReleased),
		},
		{
			name:  "HeapSys",
			value: float64(rtm.HeapSys),
		},
		{
			name:  "LastGC",
			value: float64(rtm.LastGC),
		},
		{
			name:  "Lookups",
			value: float64(rtm.Lookups),
		},
		{
			name:  "MCacheInuse",
			value: float64(rtm.MCacheInuse),
		},
		{
			name:  "MCacheSys",
			value: float64(rtm.MCacheSys),
		},
		{
			name:  "MSpanInuse",
			value: float64(rtm.MSpanInuse),
		},
		{
			name:  "MSpanSys",
			value: float64(rtm.MSpanSys),
		},
		{
			name:  "Mallocs",
			value: float64(rtm.Mallocs),
		},
		{
			name:  "NextGC",
			value: float64(rtm.NextGC),
		},
		{
			name:  "NumForcedGC",
			value: float64(rtm.NumForcedGC),
		},
		{
			name:  "NumGC",
			value: float64(rtm.NumGC),
		},
		{
			name:  "OtherSys",
			value: float64(rtm.OtherSys),
		},
		{
			name:  "PauseTotalNs",
			value: float64(rtm.PauseTotalNs),
		},
		{
			name:  "StackInuse",
			value: float64(rtm.StackInuse),
		},
		{
			name:  "StackSys",
			value: float64(rtm.StackSys),
		},
		{
			name:  "Sys",
			value: float64(rtm.Sys),
		},
		{
			name:  "TotalAlloc",
			value: float64(rtm.TotalAlloc),
		},
		{
			name:  "RandomValue",
			value: randomValue,
		},
	}
}
