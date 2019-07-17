package functions

import (
	"math/rand"
	"sync"
	"time"

	nats "github.com/nats-io/nats.go"

	"github.com/spaceuptech/space-cloud/model"
)

type serviceStub struct {
	clientID    string
	sendPayload SendPayload
}

type servicesStub struct {
	sync.RWMutex
	services     []*serviceStub
	subscription *nats.Subscription
}

func (s *servicesStub) subscribe(nc *nats.Conn, ser *serviceStub, channel chan *nats.Msg, req *model.ServiceRegisterRequest) error {
	s.RLock()
	defer s.RUnlock()

	if s.subscription == nil {
		sub, err := nc.ChanQueueSubscribe(getSubjectName(req.Service), req.Service, channel)
		if err != nil {
			return err
		}
		s.subscription = sub
	}

	if s.services == nil {
		s.services = []*serviceStub{}
	}

	s.services = append(s.services, ser)
	return nil
}

func (s *servicesStub) unsubscribe(services *sync.Map, key interface{}, clientID string) {
	s.Lock()
	defer s.Unlock()

	// Iterate over all clients and delete the client whose id matches
	for i, ser := range s.services {
		if ser.clientID == clientID {
			s.services = remove(s.services, i)
			break
		}
	}

	if len(s.services) == 0 {
		s.subscription.Unsubscribe()
		s.subscription = nil
		services.Delete(key)
	}
}

func (s *servicesStub) getService() *serviceStub {
	s.RLock()
	defer s.RUnlock()

	return s.services[rand.Intn(len(s.services))]
}

func remove(s []*serviceStub, i int) []*serviceStub {
	s[i] = s[len(s)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	return s[:len(s)-1]
}

type pendingRequest struct {
	reply   string
	reqTime time.Time
}
