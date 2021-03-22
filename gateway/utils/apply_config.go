package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/spaceuptech/helpers"
	"golang.org/x/net/context"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

// ApplySpec takes a spec object and applies it
func ApplySpec(ctx context.Context, token, hostAddr string, specObj *model.SpecObject) error {
	requestBody, err := json.Marshal(specObj.Spec)
	if err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "error while applying service unable to marshal spec", err, nil)
	}
	url, err := adjustPath(fmt.Sprintf("%s%s", hostAddr, specObj.API), specObj.Meta)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "error while applying service unable to send http request", err, nil)
	}

	v := map[string]interface{}{}
	_ = json.NewDecoder(resp.Body).Decode(&v)
	CloseTheCloser(req.Body)

	if resp.StatusCode == http.StatusAccepted {
		// Make checker send this status
		helpers.Logger.LogInfo(helpers.GetRequestID(ctx), fmt.Sprintf("Successfully queued %s", specObj.Type), nil)
	} else if resp.StatusCode == http.StatusOK {
		helpers.Logger.LogInfo(helpers.GetRequestID(ctx), fmt.Sprintf("Successfully applied %s", specObj.Type), nil)
	} else {
		_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("error while applying service got http status code %s", resp.Status), fmt.Errorf("%s", v["error"]), nil)
		return fmt.Errorf("%v", v["error"])
	}
	return nil
}

func adjustPath(path string, meta map[string]string) (string, error) {
	newPath := path
	for {
		pre := strings.IndexRune(newPath, '{')
		if pre < 0 {
			return newPath, nil
		}
		post := strings.IndexRune(newPath, '}')

		key := strings.TrimSuffix(strings.TrimPrefix(newPath[pre:post], "{"), "}")
		value, p := meta[key]
		if !p {
			return "", fmt.Errorf("provided key (%s) does not exist in metadata", key)
		}

		newPath = newPath[:pre] + value + newPath[post+1:]
	}
}
