package eventing

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func TestModule_logInvocation(t *testing.T) {
	type args struct {
		ctx                context.Context
		eventID            string
		payload            []byte
		responseStatusCode int
		responseBody       string
		errorMsg           string
	}
	type mockArgs struct {
		method        string
		args          []interface{}
		paramReturned []interface{}
	}
	tests := []struct {
		name         string
		s            *Module
		args         args
		crudMockArgs []mockArgs
		wantErr      bool
	}{
		{
			name: "invocation is logged",
			s:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype"}},
			crudMockArgs: []mockArgs{
				{
					method:        "InternalCreate",
					args:          []interface{}{context.Background(), "dbtype", "abc", "invocation_logs", &model.CreateRequest{Document: map[string]interface{}{"error_msg": "error", "event_id": "eventID", "request_payload": "", "response_body": "body", "response_status_code": 200}, Operation: "one", IsBatch: true}, false},
					paramReturned: []interface{}{nil},
				},
			},
			args:    args{ctx: context.Background(), eventID: "eventID", payload: []byte{}, responseStatusCode: 200, responseBody: "body", errorMsg: "error"},
			wantErr: false,
		},
		{
			name: "invocation is not logged",
			s:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype"}},
			crudMockArgs: []mockArgs{
				{
					method:        "InternalCreate",
					args:          []interface{}{context.Background(), "dbtype", "abc", "invocation_logs", &model.CreateRequest{Document: map[string]interface{}{"error_msg": "error", "event_id": "eventID", "request_payload": "", "response_body": "body", "response_status_code": 200}, Operation: "one", IsBatch: true}, false},
					paramReturned: []interface{}{errors.New("eventing module couldn't log the request - ")},
				},
			},
			args:    args{ctx: context.Background(), eventID: "eventID", payload: []byte{}, responseStatusCode: 200, responseBody: "body", errorMsg: "error"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockCrud := mockCrudInterface{}

			for _, m := range tt.crudMockArgs {
				mockCrud.On(m.method, m.args...).Return(m.paramReturned...)
			}

			tt.s.crud = &mockCrud

			if err := tt.s.logInvocation(tt.args.ctx, tt.args.eventID, tt.args.payload, tt.args.responseStatusCode, tt.args.responseBody, tt.args.errorMsg); (err != nil) != tt.wantErr {
				t.Errorf("Module.logInvocation() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockCrud.AssertExpectations(t)
		})
	}
}

func TestModule_MakeInvocationHTTPRequest(t *testing.T) {
	var eventResponse model.EventResponse
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx     context.Context
		client  model.HTTPEventingInterface
		method  string
		url     string
		eventID string
		token   string
		scToken string
		payload interface{}
		vPtr    interface{}
	}
	tests := []struct {
		name         string
		s            *Module
		args         args
		crudMockArgs []mockArgs
		httpMockArgs []mockArgs
		wantErr      bool
	}{
		{
			name: "error making new request with context and invocation is logged",
			s:    &Module{config: &config.Eventing{DBAlias: mock.Anything}, project: mock.Anything},
			args: args{method: "some-method", url: "url", eventID: "id", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImludGVybmFsLXNjLXVzZXIifQ.k3OcidcCnshBOGtzpprfV5Fhl2xWb6sjzPZH3omDDpw", scToken: "scToken", payload: "payload", vPtr: eventResponse},
			crudMockArgs: []mockArgs{
				{
					method:         "InternalCreate",
					args:           []interface{}{nil, mock.Anything, mock.Anything, utils.TableInvocationLogs, &model.CreateRequest{Document: map[string]interface{}{"event_id": "id", "request_payload": "\"payload\"", "response_status_code": 0, "response_body": "", "error_msg": "net/http: nil Context"}, Operation: utils.One, IsBatch: true}, false},
					paramsReturned: []interface{}{nil},
				},
			},
			wantErr: true,
		},
		{
			name: "error making new request with context and invocation is not logged",
			s:    &Module{config: &config.Eventing{DBAlias: mock.Anything}, project: mock.Anything},
			args: args{method: "some-method", url: "url", eventID: "id", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImludGVybmFsLXNjLXVzZXIifQ.k3OcidcCnshBOGtzpprfV5Fhl2xWb6sjzPZH3omDDpw", scToken: "scToken", payload: "payload", vPtr: eventResponse},
			crudMockArgs: []mockArgs{
				{
					method:         "InternalCreate",
					args:           []interface{}{nil, mock.Anything, mock.Anything, utils.TableInvocationLogs, &model.CreateRequest{Document: map[string]interface{}{"event_id": "id", "request_payload": "\"payload\"", "response_status_code": 0, "response_body": "", "error_msg": "net/http: nil Context"}, Operation: utils.One, IsBatch: true}, false},
					paramsReturned: []interface{}{errors.New("some error")},
				},
			},
			wantErr: true,
		},
		{
			name: "error doing the request and invocation is logged",
			s:    &Module{config: &config.Eventing{DBAlias: mock.Anything}, project: mock.Anything},
			args: args{ctx: context.Background(), method: "method", url: "url", eventID: "id", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImludGVybmFsLXNjLXVzZXIifQ.k3OcidcCnshBOGtzpprfV5Fhl2xWb6sjzPZH3omDDpw", scToken: "scToken", payload: "payload", vPtr: eventResponse},
			crudMockArgs: []mockArgs{
				{
					method:         "InternalCreate",
					args:           []interface{}{context.Background(), mock.Anything, mock.Anything, utils.TableInvocationLogs, &model.CreateRequest{Document: map[string]interface{}{"event_id": "id", "request_payload": "\"payload\"", "response_status_code": 0, "response_body": "", "error_msg": "some error"}, Operation: utils.One, IsBatch: true}, false},
					paramsReturned: []interface{}{nil},
				},
			},
			httpMockArgs: []mockArgs{
				{
					paramsReturned: []interface{}{nil, errors.New("some error")},
				},
			},
			wantErr: true,
		},
		{
			name: "error doing the request and invocation is not logged",
			s:    &Module{config: &config.Eventing{DBAlias: mock.Anything}, project: mock.Anything},
			args: args{ctx: context.Background(), method: "method", url: "url", eventID: "id", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImludGVybmFsLXNjLXVzZXIifQ.k3OcidcCnshBOGtzpprfV5Fhl2xWb6sjzPZH3omDDpw", scToken: "scToken", payload: "payload", vPtr: eventResponse},
			crudMockArgs: []mockArgs{
				{
					method:         "InternalCreate",
					args:           []interface{}{context.Background(), mock.Anything, mock.Anything, utils.TableInvocationLogs, &model.CreateRequest{Document: map[string]interface{}{"event_id": "id", "request_payload": "\"payload\"", "response_status_code": 0, "response_body": "", "error_msg": "some error"}, Operation: utils.One, IsBatch: true}, false},
					paramsReturned: []interface{}{errors.New("some error")},
				},
			},
			httpMockArgs: []mockArgs{
				{
					paramsReturned: []interface{}{nil, errors.New("some error")},
				},
			},
			wantErr: true,
		},
		{
			name: "error unmarshalling and invocation is logged",
			s:    &Module{config: &config.Eventing{DBAlias: mock.Anything}, project: mock.Anything},
			args: args{ctx: context.Background(), method: "method", url: "url", eventID: "id", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImludGVybmFsLXNjLXVzZXIifQ.k3OcidcCnshBOGtzpprfV5Fhl2xWb6sjzPZH3omDDpw", scToken: "scToken", payload: "payload", vPtr: eventResponse},
			crudMockArgs: []mockArgs{
				{
					method:         "InternalCreate",
					args:           []interface{}{context.Background(), mock.Anything, mock.Anything, utils.TableInvocationLogs, &model.CreateRequest{Document: map[string]interface{}{"event_id": "id", "request_payload": "\"payload\"", "response_status_code": 200, "response_body": "{\"event\": {\"type\": \"someType\"}, \"response\": \"response\"}", "error_msg": "json: Unmarshal(non-pointer model.EventResponse)"}, Operation: utils.One, IsBatch: true}, false},
					paramsReturned: []interface{}{nil},
				},
			},
			httpMockArgs: []mockArgs{
				{
					paramsReturned: []interface{}{&http.Response{Body: http.NoBody}, nil},
				},
			},
			wantErr: true,
		},
		{
			name: "error unmarshalling and invocation is not logged",
			s:    &Module{config: &config.Eventing{DBAlias: mock.Anything}, project: mock.Anything},
			args: args{ctx: context.Background(), method: "method", url: "url", eventID: "id", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImludGVybmFsLXNjLXVzZXIifQ.k3OcidcCnshBOGtzpprfV5Fhl2xWb6sjzPZH3omDDpw", scToken: "scToken", payload: "payload", vPtr: eventResponse},
			crudMockArgs: []mockArgs{
				{
					method:         "InternalCreate",
					args:           []interface{}{context.Background(), mock.Anything, mock.Anything, utils.TableInvocationLogs, &model.CreateRequest{Document: map[string]interface{}{"event_id": "id", "request_payload": "\"payload\"", "response_status_code": 200, "response_body": "{\"event\": {\"type\": \"someType\"}, \"response\": \"response\"}", "error_msg": "json: Unmarshal(non-pointer model.EventResponse)"}, Operation: utils.One, IsBatch: true}, false},
					paramsReturned: []interface{}{errors.New("some error")},
				},
			},
			httpMockArgs: []mockArgs{
				{
					paramsReturned: []interface{}{&http.Response{Body: http.NoBody}, nil},
				},
			},
			wantErr: true,
		},
		{
			name: "no error and invocation is logged",
			s:    &Module{config: &config.Eventing{DBAlias: mock.Anything}, project: mock.Anything},
			args: args{ctx: context.Background(), method: "method", url: "url", eventID: "id", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImludGVybmFsLXNjLXVzZXIifQ.k3OcidcCnshBOGtzpprfV5Fhl2xWb6sjzPZH3omDDpw", scToken: "scToken", payload: "payload", vPtr: &eventResponse},
			crudMockArgs: []mockArgs{
				{
					method:         "InternalCreate",
					args:           []interface{}{context.Background(), mock.Anything, mock.Anything, utils.TableInvocationLogs, &model.CreateRequest{Document: map[string]interface{}{"event_id": "id", "request_payload": "\"payload\"", "response_status_code": 200, "response_body": "{\"event\": {\"type\": \"someType\"}, \"response\": \"response\"}", "error_msg": ""}, Operation: utils.One, IsBatch: true}, false},
					paramsReturned: []interface{}{nil},
				},
			},
			httpMockArgs: []mockArgs{
				{
					paramsReturned: []interface{}{&http.Response{
						Proto:      "HTTP/1.1",
						ProtoMajor: 1,
						ProtoMinor: 1,
					}, nil},
				},
			},
		},
		{
			name: "no error and invocation is not logged",
			s:    &Module{config: &config.Eventing{DBAlias: mock.Anything}, project: mock.Anything},
			args: args{ctx: context.Background(), method: "method", url: "url", eventID: "id", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImludGVybmFsLXNjLXVzZXIifQ.k3OcidcCnshBOGtzpprfV5Fhl2xWb6sjzPZH3omDDpw", scToken: "scToken", payload: "payload", vPtr: &eventResponse},
			crudMockArgs: []mockArgs{
				{
					method:         "InternalCreate",
					args:           []interface{}{context.Background(), mock.Anything, mock.Anything, utils.TableInvocationLogs, &model.CreateRequest{Document: map[string]interface{}{"event_id": "id", "request_payload": "\"payload\"", "response_status_code": 200, "response_body": "{\"event\": {\"type\": \"someType\"}, \"response\": \"response\"}", "error_msg": ""}, Operation: utils.One, IsBatch: true}, false},
					paramsReturned: []interface{}{errors.New("some error")},
				},
			},
			httpMockArgs: []mockArgs{
				{
					paramsReturned: []interface{}{&http.Response{
						Proto:      "HTTP/1.1",
						ProtoMajor: 1,
						ProtoMinor: 1,
					}, nil},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockCrud := mockCrudInterface{}
			mockHTTP := mockHTTPInterface{}

			for _, m := range tt.crudMockArgs {
				mockCrud.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			for _, m := range tt.httpMockArgs {
				mockHTTP.On("Do", mock.Anything).Return(m.paramsReturned...)
			}

			tt.args.client = &mockHTTP
			tt.s.crud = &mockCrud

			if err := tt.s.MakeInvocationHTTPRequest(tt.args.ctx, tt.args.client, tt.args.method, tt.args.url, tt.args.eventID, tt.args.token, tt.args.scToken, tt.args.payload, tt.args.vPtr); (err != nil) != tt.wantErr {
				t.Errorf("Module.MakeInvocationHTTPRequest() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockCrud.AssertExpectations(t)
			mockHTTP.AssertExpectations(t)
		})
	}
}

// TODO: Write test cases for error in ReadAll and for statusCode not <200 or >300
