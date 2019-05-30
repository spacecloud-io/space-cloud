package functions

import (
	"math/rand"
	"sync"
	"time"

	nats "github.com/nats-io/go-nats"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils/client"
)

type servicesStub struct {
	sync.RWMutex
	clients      []client.Client
	subscription *nats.Subscription
}

func (s *servicesStub) subscribe(nc *nats.Conn, c client.Client, channel chan *nats.Msg, req *model.ServiceRegisterRequest) error {
	s.RLock()
	defer s.RUnlock()

	if s.subscription == nil {
		sub, err := nc.ChanQueueSubscribe(getSubjectName(req.Service), req.Service, channel)
		if err != nil {
			return err
		}
		s.subscription = sub
		s.clients = []client.Client{}
	}

	if s.clients == nil {
		s.clients = []client.Client{}
	}
	s.clients = append(s.clients, c)

	return nil
}

func (s *servicesStub) unsubscribe(services *sync.Map, key interface{}, clientID string) {
	s.Lock()
	defer s.Unlock()

	// Iterate over all clients and delete the client whose id matches
	for i, client := range s.clients {
		if client.ClientID() == clientID {
			s.clients = remove(s.clients, i)
			break
		}
	}

	if len(s.clients) == 0 {
		s.subscription.Unsubscribe()
		s.subscription = nil
		services.Delete(key)
	}
}

func (s *servicesStub) getClient() client.Client {
	s.RLock()
	defer s.RUnlock()

	return s.clients[rand.Intn(len(s.clients))]
}

func remove(s []client.Client, i int) []client.Client {
	s[i] = s[len(s)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	return s[:len(s)-1]
}

type pendingRequest struct {
	reply   string
	reqTime time.Time
}
