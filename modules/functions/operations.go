package functions

import (
	"encoding/json"

	"github.com/spaceuptech/space-cloud/modules/auth"
)

// Operation handles the function call operation
func (m *Module) Operation(auth *auth.Module, token, service, function string, timeout int) ([]byte, error) {
	var params interface{}
	authObj, _ := auth.GetAuthObj(token)
	dataBytes, err := m.Request(service, function, int(timeout), map[string]interface{}{"auth": authObj, "params": params})
	if err != nil {
		return nil, err
	}
	data := map[string]interface{}{}
	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		return nil, err
	}

	// Create the result to be sent back
	resultBytes, err := json.Marshal(map[string]interface{}{"result": data})
	if err != nil {
		return nil, err
	}
	return resultBytes, nil
}
