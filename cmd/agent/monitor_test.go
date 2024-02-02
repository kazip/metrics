package main

import (
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockMemStatsReader struct {
	Alloc int64
}

func (m *MockMemStatsReader) Read(rtm *runtime.MemStats) {
	rtm.Alloc = uint64(m.Alloc)
}

func TestNewMonitor(t *testing.T) {
	type args struct {
		duration int
		metrics  *Metrics
		done     chan bool
	}

	type want struct {
		gaugeValues map[string]float64
		counterValues map[string]int64
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "test 1",
			args: args {
				done: make(chan bool, 1),
				duration: 2,
				metrics: &Metrics{
					gauges: []GaugeMetric{

					},
					counters: []CounterMetric {
						{
							name: "PollCount",	
							value: 0,
						},
					},
				},
			},

			want: want {
				counterValues: map[string]int64{
					"PollCount": 1,
				},
				gaugeValues: map[string]float64 {
					
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {


			
			go NewMonitor(tt.args.duration, tt.args.metrics, tt.args.done, &MockMemStatsReader{
				Alloc: 12345,
			})
			time.Sleep(time.Duration(3) * time.Second)
			tt.args.done <- true

			countersMap := map[string]int64 {}
			for _, counter := range tt.args.metrics.counters {
				countersMap[counter.name] = counter.value
			}

			gaugesMap := map[string]float64 {}
			for _, gauge := range tt.args.metrics.gauges {
				gaugesMap[gauge.name] = gauge.value
			}

			assert.Equal(t, tt.want.counterValues["PollCount"], countersMap["PollCount"])
			assert.Equal(t, float64(12345), gaugesMap["Alloc"])

		})
	}
}
