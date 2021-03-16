package transport

import (
	"context"
	"fmt"

	"github.com/spaceuptech/space-api-go/types"
)

// Batch triggers the gRPC batch function on space cloud
func (t *Transport) DoDBRequest(ctx context.Context, meta *types.Meta, req interface{}) (*types.Response, error) {
	url := t.generateDatabaseURL(meta)

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

func (t *Transport) generateDatabaseURL(meta *types.Meta) string {
	scheme := "http"
	if t.sslEnabled {
		scheme = "https"
	}

	if meta.Operation == types.Batch {
		return fmt.Sprintf("%s://%s/v1/api/%s/crud/%s/batch", scheme, t.addr, meta.Project, meta.DbType)
	}

	return fmt.Sprintf("%s://%s/v1/api/%s/crud/%s/%s/%s", scheme, t.addr, meta.Project, meta.DbType, meta.Col, meta.Operation)
}
