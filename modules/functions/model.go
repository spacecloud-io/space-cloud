package functions

import (
	"time"

	nats "github.com/nats-io/go-nats"

	"github.com/spaceuptech/space-cloud/utils/client"
)

type servicesStub struct {
	clients      []client.Client
	subscription *nats.Subscription
}

type pendingRequest struct {
	reply   string
	reqTime time.Time
}
