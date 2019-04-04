package faas

import(
	"encoding/json"

	"github.com/spaceuptech/space-cloud/modules/auth"
	
)

func (m *Module)Operation(auth *auth.Module, token,engine,function string, timeout int) ([]byte, error) {
	var params interface{}
	authObj,_ := auth.GetAuthObj(token)
	dataBytes,err := m.Request(engine, function, int(timeout), map[string]interface{}{"auth": authObj, "params": params})
	if err != nil {
		return nil,err	
	}
	data := map[string]interface{}{}
	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		return nil,err
	}

	// Create the result to be sent back
	resultBytes, err := json.Marshal(map[string]interface{}{"result": data})
	if err != nil {
		return nil,err
	}
	return resultBytes,nil
}

