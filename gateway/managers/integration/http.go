package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func invokeHook(ctx context.Context, url, scToken string, params model.RequestParams, ptr interface{}) (int, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return 0, err
	}

	req := &utils.HTTPRequest{URL: url, Params: bytes.NewBuffer(data), Method: http.MethodPost, SCToken: scToken}
	return utils.MakeHTTPRequest(ctx, req, ptr)
}
