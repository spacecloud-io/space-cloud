package server

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
	l sync.Mutex
}

func tick() (m *TTLMap) {
	m = &TTLMap{m: make(map[string]*Item)}
	go func() {
		for now := range time.Tick(time.Second) {
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

func (m *TTLMap) put(k string) {
	m.l.Lock()
	it, ok := m.m[k]
	if !ok {
		it = &Item{lastAccess: time.Now().Unix()}
		m.m[k] = it
	} else {
		m.m[k] = &Item{lastAccess: time.Now().Unix()}
	}
	m.l.Unlock()
}
