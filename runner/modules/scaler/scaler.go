package scaler

import (
	"os"
	"sync"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"

	"github.com/spaceuptech/space-cloud/runner/modules/pubsub"
)

// Scaler is responsible to implement an external-push scaler for keda
type Scaler struct {
	lock sync.RWMutex

	// Map to store is active streams
	isActiveStreams map[string]*isActiveStream

	// Client drivers
	prometheusClient v1.API
	pubsubClient     *pubsub.Module
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

	// Create a new pubsub client
	pubsubClient, err := pubsub.New("runner", os.Getenv("REDIS_CONN"))
	if err != nil {
		return nil, err
	}

	return &Scaler{
		prometheusClient: v1.NewAPI(client),
		pubsubClient:     pubsubClient,
		isActiveStreams:  map[string]*isActiveStream{},
	}, nil
}
