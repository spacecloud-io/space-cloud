package functions

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/spaceuptech/space-cloud/model"
)

// RegisterService registers a new service with the functions module
func (m *Module) RegisterService(clientID string, req *model.ServiceRegisterRequest, sendPayload SendPayload) {

	service := new(servicesStub)
	t, _ := m.services.LoadOrStore(req.Service, service)
	service = t.(*servicesStub)

	// Subscribe to nats if not already subscribed
	service.subscribe(m.nc, &serviceStub{clientID, sendPayload}, m.channel, req)

	m.services.Store(req.Service, service)
}

// UnregisterService removes a service from the functions module
func (m *Module) UnregisterService(clientID string) {
	// Delete the service from all groups
	m.services.Range(func(key interface{}, value interface{}) bool {
		service := value.(*servicesStub)

		// Remove the client
		service.unsubscribe(&m.services, key, clientID)
		return true
	})
}

// HandleServiceResponse handles the service response
func (m *Module) HandleServiceResponse(res *model.FunctionsPayload) {
	// Load the pending request
	t, p := m.pendingRequests.Load(res.ID)
	if !p {
		return
	}
	req := t.(*pendingRequest)

	//  Publish the reply to nats
	data, _ := json.Marshal(res)
	m.nc.Publish(req.reply, data)

	// Remove the pending request from internal map
	m.pendingRequests.Delete(res.ID)
}

// Call simply calls a function on a service
func (m *Module) Call(service, function string, auth map[string]interface{}, params interface{}, timeout int) (interface{}, error) {
	m.RLock()
	defer m.RUnlock()

	req := &model.FunctionsPayload{Service: service, Func: function, Auth: auth, Params: params}

	// Marshal the object into json
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	// Send request over nats
	subject := getSubjectName(service)
	msg, err := m.nc.Request(subject, data, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}

	res := new(model.FunctionsPayload)
	err = json.Unmarshal(msg.Data, &res)
	if err != nil {
		return nil, err
	}

	// Return an error if response recieved has an error
	if len(res.Error) > 0 {
		return nil, errors.New(res.Error)
	}

	return res.Params, nil
}
