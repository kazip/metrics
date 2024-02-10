package main

import (
	"errors"
	"fmt"
	"sync"
)

type repository interface {
	SetGauge(name string, value float64) error
	SetCounter(name string, value int64) error
	GetGauge(name string) (float64, error)
	GetCounter(name string) (int64, error)
	GetCounters() map[string]int64
	GetGauges() map[string]float64
}

type MemStorage struct {
	GaugeMutex   sync.RWMutex
	CounterMutex sync.RWMutex
	Gauges       map[string]float64
	Counters     map[string]int64
}

func (m *MemStorage) GetGauges() map[string]float64 {
	m.GaugeMutex.RLock()
	gauges := m.Gauges
	m.GaugeMutex.RUnlock()

	return gauges
}

func (m *MemStorage) GetCounters() map[string]int64 {
	m.CounterMutex.RLock()
	counters := m.Counters
	m.CounterMutex.RUnlock()

	return counters
}

func (m *MemStorage) GetGauge(name string) (float64, error) {
	m.GaugeMutex.RLock()
	value, ok := m.Gauges[name]
	m.GaugeMutex.RUnlock()

	if !ok {
		return 0, fmt.Errorf("unknown gauge metric")
	}
	return value, nil
}

func (m *MemStorage) GetCounter(name string) (int64, error) {
	m.CounterMutex.RLock()
	value, ok := m.Counters[name]
	m.CounterMutex.RUnlock()

	if !ok {
		return 0, fmt.Errorf("unknown counter metric")
	}
	return value, nil
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
