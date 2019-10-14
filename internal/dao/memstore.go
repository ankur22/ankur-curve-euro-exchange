package dao

import (
	"time"
)

// DatabaseDAO - Interace to store and retrieve exchange rates
type DatabaseDAO interface {
	Store(from string, to string, oneUnit float32, shouldExchange bool, now time.Time)
	Get(from string, to string) (float32, bool, time.Time)
}

type memstore struct {
	oneUnit        map[string]float32
	shouldExchange map[string]bool
	dt             map[string]time.Time
}

// CreateNewMemstore - Cache in memory exchange data from the internet
func CreateNewMemstore() *memstore {
	m := memstore{oneUnit: make(map[string]float32), shouldExchange: make(map[string]bool), dt: make(map[string]time.Time)}
	return &m
}

// StoreOneUnit - Store exchange from one unit of currency (e.g. EUR)
// 				  to another currency (e.g. GBP) and whether it's a
// 				  good time to buy
func (m *memstore) Store(from string, to string, oneUnit float32, shouldExchange bool, now time.Time) {
	key := from + to
	m.oneUnit[key] = oneUnit
	m.shouldExchange[key] = shouldExchange
	m.dt[key] = now
}

// GetOneUnit - Get stored exchange from one unit of currency (e.g. EUR)
// 			    to another currency (e.g. GBP) and whether it's a good
//				time to buy (time when stored)
func (m *memstore) Get(from string, to string) (float32, bool, time.Time) {
	key := from + to
	return m.oneUnit[key], m.shouldExchange[key], m.dt[key]
}
