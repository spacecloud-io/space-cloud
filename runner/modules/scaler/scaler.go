package scaler

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

// Scaler is responsible to implement an external-push scaler for keda
type Scaler struct {
	lock sync.RWMutex

	// Map to store is active streams
	isActiveStreams map[string]*isActiveStream

	// Client drivers
	prometheusClient v1.API
}

// New creates an instance of the scaler object
func New(prometheusAddr string) (*Scaler, error) {
	// Create a new prometheus client
	client, err := api.NewClient(api.Config{
		Address: prometheusAddr,
	})
	if err != nil {
		return nil, err
	}

	return &Scaler{
		prometheusClient: v1.NewAPI(client),
		isActiveStreams:  map[string]*isActiveStream{},
	}, nil
}

type isActiveStream struct {
	lock sync.RWMutex

	ch          chan bool
	lastUpdated time.Time
}

func (s *isActiveStream) Notify() {
	// Send a scale up signal only if more than 10 seconds has elapsed since the last update
	if !s.IsActive() {
		s.ch <- true // No need to access channel within a lock

		s.lock.Lock()
		s.lastUpdated = time.Now()
		s.lock.Unlock()
	}
}

func (s *isActiveStream) IsActive() bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.lastUpdated.After(time.Now().Add(-10 * time.Second))
}
