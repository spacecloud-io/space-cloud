package server

import (
	"fmt"
	"math"
	"strings"
	"sync"
)

type aggregator struct {
	lock  sync.Mutex
	count map[string]map[string]*counter
}

func newAggregator() *aggregator {
	return &aggregator{count: map[string]map[string]*counter{}}
}

func (a *aggregator) makeKey(project, service, version string) string {
	return fmt.Sprintf("%s:%s:%s", project, service, version)
}

func (a *aggregator) splitKey(key string) (project, service, version string) {
	array := strings.Split(key, ":")
	return array[0], array[1], array[2]
}

func (a *aggregator) add(project, service, version, nodeID string, value int32) {
	a.lock.Lock()
	defer a.lock.Unlock()

	m, p := a.count[a.makeKey(project, service, version)]
	if !p {
		m = map[string]*counter{}
		a.count[a.makeKey(project, service, version)] = m
	}

	c, p := m[nodeID]
	if !p {
		c = &counter{}
		m[nodeID] = c
	}

	c.value += value
	c.nos++
}

func (a *aggregator) iterate(cb func(project, service, version string, value int32)) {
	a.lock.Lock()
	defer a.lock.Unlock()

	for key, nodes := range a.count {
		var value int32
		for _, c := range nodes {
			value += int32(math.Ceil(float64(c.value) / float64(c.nos)))
		}

		project, service, version := a.splitKey(key)
		cb(project, service, version, value)
	}
}

func (a *aggregator) get(project, service, version string) int32 {
	a.lock.Lock()
	defer a.lock.Unlock()

	nodes, p := a.count[a.makeKey(project, service, version)]
	if !p {
		return 0
	}

	// Aggregate active request of each node
	var value int32
	for _, c := range nodes {
		value += int32(math.Ceil(float64(c.value) / float64(c.nos)))
	}
	return value
}

func (a *aggregator) delete(project, service, version string) {
	a.lock.Lock()
	defer a.lock.Unlock()

	delete(a.count, a.makeKey(project, service, version))
}

// func (a *aggregator) len() int {
// 	a.lock.Lock()
// 	defer a.lock.Unlock()
//
// 	return len(a.count)
// }
