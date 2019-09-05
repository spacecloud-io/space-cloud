package pubsub

import (
	"sync"
	"net/http"
	"errors"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

// pubsubBroker abstracts the implementation of the pubsub clients
type pubsubBroker interface {
	Publish(subject string, msg *model.PubsubMsg) error
	QueueSubscribe(subject, queue string, ch chan *model.PubsubMsg) (model.PubsubUnsubscribe, error)
}

// pubsubSubscription stores the details of the subscription
type pubsubSubscription struct {
	unsubscribeFunc model.PubsubUnsubscribe
}

// Publish publishes a byte array to a particular subject, if its permitted
func (m *Module) Publish(project, token, subject string, data interface{}) (int, error) {
	// Exit if pubsub is not enabled
	if !m.IsEnabled() {
		return http.StatusNotFound, errors.New("This feature isn't enabled")
	}
	
	subject = utils.SingleLeadingTrailing(subject, "/")

	// Check if the user is authorised to make this request
	err := m.auth.IsPublishAuthorised(project, token, subject, map[string]interface{}{})
	if err != nil {
		return http.StatusForbidden, errors.New("You are not authorized to make this request")
	}

	// m.RLock()
	// defer m.RUnlock()

	err = m.connection.Publish(subject, &model.PubsubMsg{subject, data})
	if err != nil {
		return http.StatusInternalServerError, err
	} else {
		return http.StatusOK, nil
	}
}

// Subscribe subscribes to a particular subject (and its children), if its permitted
func (m *Module) Subscribe(project, token, clientID, subject string, cb model.PubsubCallback) (int, error) {
	return m.QueueSubscribe(project, token, clientID, subject, clientID, cb)
}

// QueueSubscribe subscribes to a particular subject (and its children) using a queue, if its permitted
func (m *Module) QueueSubscribe(project, token, clientID, subject, queue string, cb model.PubsubCallback) (int, error) {
	// Exit if pubsub is not enabled
	if !m.IsEnabled() {
		return http.StatusNotFound, errors.New("This feature isn't enabled")
	}
	
	subject = utils.SingleLeadingTrailing(subject, "/")

	// Check if the user is authorised to make this request
	err := m.auth.IsSubscribeAuthorised(project, token, subject, queue, map[string]interface{}{})
	if err != nil {
		return http.StatusForbidden, errors.New("You are not authorized to make this request")
	}

	receiveChannel := make(chan *model.PubsubMsg)
	subs, err := m.connection.QueueSubscribe(subject, queue, receiveChannel)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	subscription := pubsubSubscription{subs}
	err = m.storeSubs(subject, clientID, &subscription)
	if err != nil {
		subscription.unsubscribeFunc()
		return http.StatusInternalServerError, err
	}
	go func() {
		for msg := range receiveChannel {
        	cb(msg)
		}
	}()
	return http.StatusOK, nil
}

// Unsubscribe unsubscribes a client from a particular subject
func (m *Module) Unsubscribe(clientID, subject string) (int, error) {
	// Exit if pubsub is not enabled
	if !m.IsEnabled() {
		return http.StatusNotFound, errors.New("This feature isn't enabled")
	}
	
	subject = utils.SingleLeadingTrailing(subject, "/")

	m.Lock()
	defer m.Unlock()

	subjM, ok := m.storage.Load(subject)
	if ok {
		subjectMap := subjM.(sync.Map)
		subs, ok := subjectMap.Load(clientID)
		if ok {
			subscription := subs.(*pubsubSubscription)
			subjectMap.Delete(clientID)
			if mapLen(&subjectMap) == 0 {
				m.storage.Delete(subject)
			}
			(*subscription).unsubscribeFunc()
		}
	}
	return http.StatusOK, nil
}

// UnsubscribeAll unsubscribes a client from all subjects
func (m *Module) UnsubscribeAll(clientID string) (int, error) {
	// Exit if pubsub is not enabled
	if !m.IsEnabled() {
		return http.StatusNotFound, errors.New("This feature isn't enabled")
	}

	// m.RLock()
	// defer m.RUnlock()
	
	var err error = nil
	status := 200

	m.storage.Range(func(subject, v interface{}) bool {
		s, e := m.Unsubscribe(clientID, subject.(string))
		if e != nil {
			status = s
			err = e
			return true
		}
		return true
	})
	return status, err
}

// 1 client cannot subs to same suject, using 2 diff queues
// Map subject -> clientID -> pubsubSubscription
// storeSubs stores a particular subscription's details
func (m *Module) storeSubs(subject, clientID string, subs *pubsubSubscription) error {
	m.Lock()
	defer m.Unlock()

	var temp sync.Map
	temp.Store(clientID, subs)
	clientMap, alreadyPresent := m.storage.LoadOrStore(subject, temp)
	if alreadyPresent {
		cli := clientMap.(sync.Map)
		_, alreadyPresent := cli.LoadOrStore(clientID, subs)
		if alreadyPresent {
			return errors.New("Already subscribed to this channel")
		}
	}
	return nil
}

// mapLen calculates the number of keys in a sync.Map
func mapLen(m *sync.Map) int {
	counter := 0
	m.Range(func(k, v interface{}) bool {
		counter++
		return true
	})
	return counter
}
