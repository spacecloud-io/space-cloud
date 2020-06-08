package routing

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func authorizeRequest(ctx context.Context, modules modulesInterface, route *config.Route, req *http.Request) error {
	// Return if the rule is allow
	if route.Rule == nil || route.Rule.Rule == "allow" {
		return nil
	}

	// Extract the params only if content-type is `application/json`
	var params interface{}
	var data []byte
	var err error
	if req.Header.Get("Content-Type") == "application/json" {
		data, err = ioutil.ReadAll(req.Body)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(data, &params); err != nil {
			return utils.LogError("Unable to unmarshal body to JSON", module, handle, err)
		}
	}

	// Extract the token
	token := utils.GetTokenFromHeader(req)

	// Finally we authorize the request
	a := modules.Auth()
	_, err = a.AuthorizeRequest(ctx, route.Rule, route.Project, token, params)

	// Don't forget to reset the body
	if data != nil {
		req.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	}

	return err
}
