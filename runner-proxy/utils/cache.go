package utils

import (
	"sync"
	"time"
)

// Item keeps track of last time service was called
type Item struct {
	lastAccess int64
}

// TTLMap keeps track of which service was called at what time
type TTLMap struct {
	m map[string]*Item
	l sync.RWMutex
}

// Tick creates a TTLMap object and constantly check for outdated elements inside the object.
func Tick() (m *TTLMap) {
	m = &TTLMap{m: make(map[string]*Item)}
	go func() {
		for now := range time.Tick(10 * time.Second) {
			m.l.Lock()
			for k, v := range m.m {
				if now.Unix()-v.lastAccess > int64(120) {
					delete(m.m, k)
				}
			}
			m.l.Unlock()
		}
	}()
	return m
}

func (m *TTLMap) len() int {
	return len(m.m)
}

// Put insert new value in TTLMap
func (m *TTLMap) Put(k string) {
	m.l.RLock()
	defer m.l.RUnlock()

	it, ok := m.m[k]
	if !ok {
		it = &Item{lastAccess: time.Now().Unix()}
		m.m[k] = it
	} else {
		m.m[k] = &Item{lastAccess: time.Now().Unix()}
	}
}

// GetDeployment check if the element exist in the TTLMap
func (m *TTLMap) GetDeployment(k string) bool {
	m.l.Lock()
	defer m.l.Unlock()

	it, ok := m.m[k]
	if ok {
		if time.Now().Unix()-it.lastAccess < int64(120) {
			return true
		}
	}
	return false
}
