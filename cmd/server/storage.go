package main

import (
	"errors"
	"sync"
)

type repository interface {
	SetGauge(name string, value float64) error
	SetCounter(name string, value int64) error
}

type MemStorage struct {
	GaugeMutex   sync.RWMutex
	CounterMutex sync.RWMutex
	Gauges       map[string]float64
	Counters     map[string]int64
}

func (m *MemStorage) SetGauge(name string, value float64) error {

	if name == "" {
		return errors.New("empty name")
	}

	m.GaugeMutex.Lock()
	m.Gauges[name] = value
	m.GaugeMutex.Unlock()
	return nil
}

func (m *MemStorage) SetCounter(name string, value int64) error {

	if name == "" {
		return errors.New("empty name")
	}

	m.CounterMutex.Lock()

	if _, ok := m.Counters[name]; !ok {
		m.Counters[name] = 0
	}

	m.Counters[name] += value
	m.CounterMutex.Unlock()
	return nil
}

func NewMemStorage() *MemStorage {
	storage := MemStorage{}
	storage.Gauges = make(map[string]float64)
	storage.Counters = make(map[string]int64)
	return &storage
}
