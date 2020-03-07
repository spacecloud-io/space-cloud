package eventing

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/net/context"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func interfaceToByteArray(params interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(params)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MakeInvocationHTTPRequest fires an http request and returns a response
func (s *Module) MakeInvocationHTTPRequest(ctx context.Context, method, url, token, scToken string, params, vPtr interface{}) error {
	// Marshal json into byte array
	data, _ := json.Marshal(params)

	// Make a request object
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(data))
	payload, er := interfaceToByteArray(params)
	if er != nil {
		invocationDoc := &model.InvocationDocument{
			InvocationTime:     time.Now().String(),
			ResponseStatusCode: 0,
			ErrorMessage:       er.Error(),
		}
		createRequest := &model.CreateRequest{Document: invocationDoc, Operation: utils.All, IsBatch: true}
		if err := s.crud.InternalCreate(ctx, s.config.DBType, s.project, invocationLogs, createRequest, false); err != nil {
			return errors.New("eventing module couldn't log the request - " + err.Error())
		}
		return er
	}
	if err != nil {
		invocationDoc := &model.InvocationDocument{
			InvocationTime:     time.Now().String(),
			RequestPayload:     string(payload),
			ResponseStatusCode: 0,
			ErrorMessage:       err.Error(),
		}
		createRequest := &model.CreateRequest{Document: invocationDoc, Operation: utils.All, IsBatch: true}
		if err := s.crud.InternalCreate(ctx, s.config.DBType, s.project, invocationLogs, createRequest, false); err != nil {
			return errors.New("eventing module couldn't log the request - " + err.Error())
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

	// Create a http client and fire the request
	client := &http.Client{}

	// if s.storeType && s.isConsulConnectEnabled && strings.Contains(url, "https") && strings.Contains(url, ".consul") {
	// 	 client = s.consulService.HTTPClient()
	// }

	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	responseBody, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		invocationDoc := &model.InvocationDocument{
			InvocationTime:     time.Now().String(),
			RequestPayload:     string(payload),
			ResponseStatusCode: 0,
			ErrorMessage:       e.Error(),
		}
		createRequest := &model.CreateRequest{Document: invocationDoc, Operation: utils.All, IsBatch: true}
		if err := s.crud.InternalCreate(ctx, s.config.DBType, s.project, invocationLogs, createRequest, false); err != nil {
			return errors.New("eventing module couldn't log the request - " + err.Error())
		}
		return e
	}
	if err != nil {
		invocationDoc := &model.InvocationDocument{
			InvocationTime:     time.Now().String(),
			RequestPayload:     string(payload),
			ResponseStatusCode: resp.StatusCode,
			ResponseBody:       string(responseBody),
			ErrorMessage:       err.Error(),
		}
		createRequest := &model.CreateRequest{Document: invocationDoc, Operation: utils.All, IsBatch: true}
		if err := s.crud.InternalCreate(ctx, s.config.DBType, s.project, invocationLogs, createRequest, false); err != nil {
			return errors.New("eventing module couldn't log the request - " + err.Error())
		}
		return err
	}
	defer utils.CloseTheCloser(resp.Body)

	if err := json.NewDecoder(resp.Body).Decode(vPtr); err != nil {
		invocationDoc := &model.InvocationDocument{
			InvocationTime:     time.Now().String(),
			RequestPayload:     string(payload),
			ResponseStatusCode: resp.StatusCode,
			ResponseBody:       string(responseBody),
			ErrorMessage:       err.Error(),
		}
		createRequest := &model.CreateRequest{Document: invocationDoc, Operation: utils.All, IsBatch: true}
		if err := s.crud.InternalCreate(ctx, s.config.DBType, s.project, invocationLogs, createRequest, false); err != nil {
			return errors.New("eventing module couldn't log the request - " + err.Error())
		}
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		invocationDoc := &model.InvocationDocument{
			InvocationTime:     time.Now().String(),
			RequestPayload:     string(payload),
			ResponseStatusCode: resp.StatusCode,
			ResponseBody:       string(responseBody),
			ErrorMessage:       err.Error(),
		}
		createRequest := &model.CreateRequest{Document: invocationDoc, Operation: utils.All, IsBatch: true}
		if err := s.crud.InternalCreate(ctx, s.config.DBType, s.project, invocationLogs, createRequest, false); err != nil {
			return errors.New("eventing module couldn't log the request - " + err.Error())
		}
		return errors.New("service responded with status code " + strconv.Itoa(resp.StatusCode))
	}

	return nil
}
