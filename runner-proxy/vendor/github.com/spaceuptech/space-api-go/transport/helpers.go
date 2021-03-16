package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/spaceuptech/space-api-go/types"
	"github.com/spaceuptech/space-api-go/utils"
)

const contentTypeJSON string = "application/json"

func (t *Transport) makeHTTPRequest(ctx context.Context, token, url string, payload interface{}) (int, types.M, error) {
	// Marshal the payload
	data, err := json.Marshal(payload)
	if err != nil {
		return -1, nil, err
	}

	// Make a http request
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return -1, nil, err
	}

	// Add appropriate headers
	r.Header.Add("Authorization", "Bearer "+token)
	r.Header.Add("Content-Type", contentTypeJSON)

	// Fire the request
	res, err := t.httpClient.Do(r)
	if err != nil {
		return -1, nil, err
	}
	defer utils.CloseTheCloser(res.Body)

	// Unmarshal the response
	result := types.M{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return -1, nil, err
	}

	// Return the final response
	return res.StatusCode, result, nil
}
