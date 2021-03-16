package transport

import (
	"context"
	"fmt"
	"github.com/spaceuptech/space-api-go/types"
)

func (t *Transport) TriggerEvent(ctx context.Context, meta *types.Meta, req interface{}) (*types.Response, error) {
	url := t.generateEventingURL(meta)
	// Fire the http request
	status, result, err := t.makeHTTPRequest(ctx, meta.Token, url, req)
	if err != nil {
		return nil, err
	}

	if status >= 200 && status < 300 {
		return &types.Response{Status: status, Data: result}, nil
	}

	return &types.Response{Status: status, Error: result["error"].(string)}, nil

}

func (t *Transport) generateEventingURL(meta *types.Meta) string {
	scheme := "http"
	if t.sslEnabled {
		scheme = "https"
	}

	return fmt.Sprintf("%s://%s/v1/api/%s/eventing/queue", scheme, t.addr, meta.Project)
}
