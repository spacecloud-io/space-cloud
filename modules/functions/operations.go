package functions

import (
	"encoding/json"
	"errors"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/modules/auth"
	"github.com/spaceuptech/space-cloud/utils"
	"github.com/spaceuptech/space-cloud/utils/client"
)

// RegisterService registers a new service with the functions module
func (m *Module) RegisterService(reqID string, c client.Client, req *model.ServiceRegisterRequest) {

	service := new(servicesStub)
	t, _ := m.services.LoadOrStore(req.Service, service)
	service = t.(*servicesStub)

	// Subscribe to nats if not already subscribed
	service.subscribe(m.nc, c, m.channel, req)

	m.services.Store(req.Service, service)

	c.Write(&model.Message{
		ID:   reqID,
		Type: utils.TypeServiceRegister,
		Data: map[string]interface{}{"ack": true},
	})
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

// Operation handles the function call operation
func (m *Module) Operation(auth *auth.Module, token, service, function string, params interface{}, timeout int) ([]byte, error) {
	authObj, _ := auth.GetAuthObj(token)
	dataBytes, err := m.Request(service, int(timeout), &model.FunctionsPayload{Service: service, Func: function, Auth: authObj, Params: params})
	if err != nil {
		return nil, err
	}
	data := new(model.FunctionsPayload)
	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		return nil, err
	}

	// Return an error if response recieved has an error
	if len(data.Error) > 0 {
		return nil, errors.New(data.Error)
	}

	// Create the result to be sent back
	resultBytes, err := json.Marshal(map[string]interface{}{"result": data.Params})
	if err != nil {
		return nil, err
	}
	return resultBytes, nil
}
