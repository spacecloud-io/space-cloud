package eventing

import (
	"errors"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
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
			s:    &Module{project: "abc", config: &config.Eventing{DBType: "dbtype"}},
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
			s:    &Module{project: "abc", config: &config.Eventing{DBType: "dbtype"}},
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

type mockCrudInterface struct {
	mock.Mock
}

func (m *mockCrudInterface) InternalCreate(ctx context.Context, dbAlias, project, col string, req *model.CreateRequest, isIgnoreMetrics bool) error {
	c := m.Called(ctx, dbAlias, project, col, req, isIgnoreMetrics)
	if err := c.Error(0); err != nil {
		return err
	}
	return nil
}

func (m *mockCrudInterface) Read(ctx context.Context, dbAlias, project, col string, req *model.ReadRequest) (interface{}, error) {
	return nil, nil
}

func (m *mockCrudInterface) InternalUpdate(ctx context.Context, dbAlias, project, col string, req *model.UpdateRequest) error {
	return nil
}

// TODO: MakeInvocationHTTPRequest
