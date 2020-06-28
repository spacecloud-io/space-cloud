package eventing

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"golang.org/x/net/context"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (s *Module) logInvocation(ctx context.Context, eventID string, payload []byte, responseStatusCode int, responseBody, errorMsg string) error {
	invocationDoc := map[string]interface{}{
		"event_id":             eventID,
		"request_payload":      string(payload),
		"response_status_code": responseStatusCode,
		"response_body":        responseBody,
		"error_msg":            errorMsg,
	}
	createRequest := &model.CreateRequest{Document: invocationDoc, Operation: utils.One, IsBatch: true}
	if err := s.crud.InternalCreate(ctx, s.config.DBAlias, s.project, utils.TableInvocationLogs, createRequest, false); err != nil {
		return errors.New("eventing module couldn't log the request - " + err.Error())
	}
	return nil
}

// MakeInvocationHTTPRequest fires an http request and returns a response
func (s *Module) MakeInvocationHTTPRequest(ctx context.Context, client model.HTTPEventingInterface, method, url, eventID, token, scToken string, payload, vPtr interface{}) error {
	// Marshal json into byte array
	data, _ := json.Marshal(payload)

	// Make a request object
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(data))
	if err != nil {
		if err := s.logInvocation(ctx, eventID, data, 0, "", err.Error()); err != nil {
			logrus.Errorf("eventing module couldn't log the invocation - %s", err.Error())
			return err
		}
		return err
	}

	// Add the headers
	if token != "" {
		// Add the token only if its provided
		req.Header.Add("Authorization", "Bearer "+token)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-sc-token", "Bearer "+scToken)

	// if s.storeType && s.isConsulConnectEnabled && strings.Contains(url, "https") && strings.Contains(url, ".consul") {
	// 	 client = s.consulService.HTTPClient()
	// }

	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		if err := s.logInvocation(ctx, eventID, data, 0, "", err.Error()); err != nil {
			logrus.Errorf("eventing module couldn't log the invocation - %s", err.Error())
			return err
		}
		return err
	}
	defer utils.CloseTheCloser(resp.Body)
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if err := s.logInvocation(ctx, eventID, data, 0, "", err.Error()); err != nil {
			logrus.Errorf("eventing module couldn't log the invocation - %s", err.Error())
			return err
		}
		return err
	}

	if err := json.Unmarshal(responseBody, vPtr); err != nil {
		if err := s.logInvocation(ctx, eventID, data, resp.StatusCode, string(responseBody), err.Error()); err != nil {
			logrus.Errorf("eventing module couldn't log the invocation - %s", err.Error())
			return err
		}
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if err := s.logInvocation(ctx, eventID, data, resp.StatusCode, string(responseBody), errors.New("invalid status code received").Error()); err != nil {
			logrus.Errorf("eventing module couldn't log the invocation - %s", err.Error())
			return err
		}
		return fmt.Errorf("service responded with status code - %s", strconv.Itoa(resp.StatusCode))
	}

	if err := s.logInvocation(ctx, eventID, data, resp.StatusCode, string(responseBody), ""); err != nil {
		logrus.Errorf("eventing module couldn't log the invocation - %s", err.Error())
		return err
	}

	return nil
}
