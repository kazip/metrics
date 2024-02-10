package main

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetrics_fromRtm(t *testing.T) {
	type args struct {
		rtm         runtime.MemStats
		pollCount   int64
		randomValue float64
	}

	type wantGaugeValue struct {
		key string
		value float64
	}

	type wantCounterValue struct {
		key string
		value int64
	}

	type want struct {
		gaugeValues []wantGaugeValue
		counterValues []wantCounterValue
	}

	tests := []struct {
		name   string		
		args   args
		want want
	}{
		{
			name: "test 1",
			args: args{
				rtm: runtime.MemStats{
					Alloc: 1,
				},
				pollCount: 1,
				randomValue: 0.1,
			},
			want: want{
				gaugeValues: []wantGaugeValue {
					{
						key: "Alloc",
						value: 1,
					},
					{
						key: "RandomValue",
						value: 0.1,
					},
				},
				counterValues: []wantCounterValue {
					{
						key: "PollCount",
						value: 1,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metrics{}
			m.fromRtm(tt.args.rtm, tt.args.pollCount, tt.args.randomValue)			

			gaugesMap := map[string]float64{}
			countersMap := map[string]int64{}
			
			for _,counter := range m.counters {
				countersMap[counter.name] = counter.value
			}

			for _,gauge := range m.gauges {
				gaugesMap[gauge.name] = gauge.value				
			}

			gauges := []string {
				"Alloc",					
				"BuckHashSys",
				"Frees",
				"GCCPUFraction",
				"GCSys",
				"HeapAlloc",
				"HeapIdle",
				"HeapInuse",
				"HeapObjects",
				"HeapReleased",
				"HeapSys",
				"LastGC",
				"Lookups",
				"MCacheInuse",
				"MCacheSys",
				"MSpanInuse",
				"MSpanSys",
				"Mallocs",
				"NextGC",
				"NumForcedGC",
				"NumGC",
				"OtherSys",
				"PauseTotalNs",
				"StackInuse",
				"StackSys",
				"Sys",
				"TotalAlloc",
				"RandomValue",
			}

			counters := []string {
				"PollCount",
			}
			
			for _, gauge := range gauges {
				assert.Contains(t, gaugesMap, gauge)	
			}

			for _, counter := range counters {
				assert.Contains(t, countersMap, counter)
			}

			for _, gauge := range tt.want.gaugeValues {
				if assert.Contains(t, gaugesMap, gauge.key) {
					assert.Equal(t, gauge.value, gaugesMap[gauge.key])
				}				
			}

			for _, counter := range tt.want.counterValues {
				if assert.Contains(t, countersMap, counter.key) {
					assert.Equal(t, counter.value, countersMap[counter.key])
				}				
			}
			
		})
	}
}
