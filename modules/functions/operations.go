package functions

import (
	"encoding/json"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/modules/auth"
	"github.com/spaceuptech/space-cloud/utils"
)

// RegisterService registers a new service with the functions module
func (m *Module) RegisterService(client *utils.Client, req *model.ServiceRegisterRequest) error {

	service := new(servicesStub)
	t, _ := m.services.LoadOrStore(req.Service, service)
	service = t.(*servicesStub)

	// Subscribe to nats if not already subscribed
	if service.subscription == nil {
		sub, err := m.nc.ChanQueueSubscribe(getSubjectName(req.Service), req.Service, m.channel)
		if err != nil {
			return err
		}
		service.subscription = sub
		service.clients = []*utils.Client{}
	}

	service.clients = append(service.clients, client)
	m.services.Store(req.Service, service)
	return nil
}

// UnregisterService removes a service from the functions module
func (m *Module) UnregisterService(clientID string) {
	// Delete the service from all groups
	m.services.Range(func(key interface{}, value interface{}) bool {
		service := value.(*servicesStub)

		// Iterate over all clients and delete the client whose id matches
		for i, client := range service.clients {
			if client.ClientID() == clientID {
				remove(service.clients, i)
				break
			}
		}

		// Unsubscribe and delete service group if no clients are present in it
		if len(service.clients) == 0 {
			service.subscription.Unsubscribe()
			m.services.Delete(key)
		}

		return true
	})
}

// HandleServiceResponse handles the service response
func (m *Module) HandleServiceResponse(id string, res *model.FunctionsPayload) {
	// Load the pending request
	t, p := m.pendingRequests.Load(id)
	if !p {
		return
	}
	req := t.(*pendingRequest)

	//  Publish the reply to nats
	data, _ := json.Marshal(res)
	m.nc.Publish(req.reply, data)

	// Remove the pending request from internal map
	m.pendingRequests.Delete(id)
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

	// Create the result to be sent back
	resultBytes, err := json.Marshal(map[string]interface{}{"result": data.Params})
	if err != nil {
		return nil, err
	}
	return resultBytes, nil
}

func remove(s []*utils.Client, i int) []*utils.Client {
	s[i] = s[len(s)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	return s[:len(s)-1]
}
