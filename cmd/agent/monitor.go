package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

type MemStatsReader interface {
	Read(*runtime.MemStats)
}

func NewMonitor(duration int, metrics *Metrics, done <-chan bool, reader MemStatsReader) {
	var interval = time.Duration(duration) * time.Second
	var pollCount int64 = 0
	var rtm runtime.MemStats

loop:
	for {
		select {
		case <-time.After(interval):
		case <-done:
			break loop
		}

		pollCount++
		fmt.Println("Collecting metrics")
		mutex.Lock()
		reader.Read(&rtm)
		metrics.fromRtm(rtm, pollCount, rand.Float64())
		mutex.Unlock()
	}
}
