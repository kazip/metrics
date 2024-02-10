package main

import "runtime"

type MemStats struct {

}

func (m *MemStats) Read(rtm *runtime.MemStats) {
	runtime.ReadMemStats(rtm)
}