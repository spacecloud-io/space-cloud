package eventing

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func TestModule_processStagedEvents(t *testing.T) {
	timeValue, _ := time.Parse(time.RFC3339Nano, "2006-01-02T15:04:05Z07:00")
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		t *time.Time
	}
	tests := []struct {
		name            string
		m               *Module
		args            args
		syncmanMockArgs []mockArgs
		crudMockArgs    []mockArgs
	}{
		{
			name: "config is not enabled",
			m:    &Module{project: "abc", config: &config.Eventing{Enabled: false, DBAlias: "db"}},
			args: args{t: &timeValue},
		},
		{
			name: "error while reading",
			m:    &Module{project: "abc", config: &config.Eventing{Enabled: true, DBAlias: "db"}},
			args: args{t: &timeValue},
			syncmanMockArgs: []mockArgs{
				{
					method:         "GetAssignedTokens",
					paramsReturned: []interface{}{1, 100},
				},
			},
			crudMockArgs: []mockArgs{
				{
					method:         "Read",
					args:           []interface{}{mock.Anything, "db", "event_logs", &model.ReadRequest{Operation: utils.All, Options: &model.ReadOptions{Sort: []string{"ts"}, Limit: &limit}, Find: map[string]interface{}{"status": utils.EventStatusStaged, "token": map[string]interface{}{"$gte": 1, "$lte": 100}}}},
					paramsReturned: []interface{}{[]interface{}{&model.EventDocument{ID: "eventDocID", Timestamp: time.Now().Format(time.RFC3339Nano)}}, errors.New("some error")},
				},
			},
		},
		{
			name: "error while decoding",
			m:    &Module{project: "abc", config: &config.Eventing{Enabled: true, DBAlias: "db"}},
			args: args{t: &timeValue},
			syncmanMockArgs: []mockArgs{
				{
					method:         "GetAssignedTokens",
					paramsReturned: []interface{}{1, 100},
				},
			},
			crudMockArgs: []mockArgs{
				{
					method:         "Read",
					args:           []interface{}{mock.Anything, "db", "event_logs", &model.ReadRequest{Operation: utils.All, Options: &model.ReadOptions{Sort: []string{"ts"}, Limit: &limit}, Find: map[string]interface{}{"status": utils.EventStatusStaged, "token": map[string]interface{}{"$gte": 1, "$lte": 100}}}},
					paramsReturned: []interface{}{[]interface{}{"payload", nil}},
				},
			},
		},
		{
			name: "no error staging events",
			m:    &Module{project: "abc", config: &config.Eventing{Enabled: true, DBAlias: "db"}},
			args: args{t: &timeValue},
			syncmanMockArgs: []mockArgs{
				{
					method:         "GetAssignedTokens",
					paramsReturned: []interface{}{1, 100},
				},
			},
			crudMockArgs: []mockArgs{
				{
					method:         "Read",
					args:           []interface{}{mock.Anything, "db", "event_logs", &model.ReadRequest{Operation: utils.All, Options: &model.ReadOptions{Sort: []string{"ts"}, Limit: &limit}, Find: map[string]interface{}{"status": utils.EventStatusStaged, "token": map[string]interface{}{"$gte": 1, "$lte": 100}}}},
					paramsReturned: []interface{}{[]interface{}{&model.EventDocument{ID: "eventDocID", Timestamp: time.Now().Format(time.RFC3339Nano)}}, nil},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockSyncman := mockSyncmanEventingInterface{}
			mockCrud := mockCrudInterface{}

			for _, m := range tt.syncmanMockArgs {
				mockSyncman.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.crudMockArgs {
				mockCrud.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			tt.m.syncMan = &mockSyncman
			tt.m.crud = &mockCrud

			tt.m.processStagedEvents(tt.args.t)

			mockSyncman.AssertExpectations(t)
			mockCrud.AssertExpectations(t)
		})
	}
}

func TestModule_invokeWebhook(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx        context.Context
		client     model.HTTPEventingInterface
		rule       *config.EventingTrigger
		eventDoc   *model.EventDocument
		cloudEvent *model.CloudEventPayload
	}
	tests := []struct {
		name            string
		m               *Module
		args            args
		authMockArgs    []mockArgs
		crudMockArgs    []mockArgs
		httpMockArgs    []mockArgs
		syncmanMockArgs []mockArgs
		adminMockArgs   []mockArgs
		wantErr         bool
	}{
		{
			name: "error getting sc access token",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype"}},
			args: args{ctx: context.Background(), rule: &config.EventingTrigger{Timeout: 100, URL: "url"}, eventDoc: &model.EventDocument{ID: "id", BatchID: "batchid"}, cloudEvent: &model.CloudEventPayload{Data: "payload"}},
			authMockArgs: []mockArgs{
				{
					method:         "GetSCAccessToken",
					paramsReturned: []interface{}{"", errors.New("some error")},
				},
			},
			wantErr: true,
		},
		{
			name: "error making invocation http request",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype"}},
			args: args{ctx: context.Background(), rule: &config.EventingTrigger{Timeout: 100, URL: "url"}, eventDoc: &model.EventDocument{ID: "id", BatchID: "batchid"}, cloudEvent: &model.CloudEventPayload{Data: "payload"}},
			authMockArgs: []mockArgs{
				{
					method:         "GetSCAccessToken",
					paramsReturned: []interface{}{"scToken", nil},
				},
			},
			crudMockArgs: []mockArgs{
				{
					method:         "InternalCreate",
					args:           []interface{}{mock.Anything, "dbtype", "abc", utils.TableInvocationLogs, &model.CreateRequest{Document: map[string]interface{}{"error_msg": "some error", "event_id": "id", "request_payload": "{\"specversion\":\"\",\"type\":\"\",\"source\":\"\",\"id\":\"\",\"time\":\"\",\"data\":\"payload\"}", "response_body": "", "response_status_code": 0}, Operation: utils.One, IsBatch: true}, false},
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockAuth := mockAuthEventingInterface{}
			mockHTTP := mockHTTPInterface{}
			mockCrud := mockCrudInterface{}
			mockSyncman := mockSyncmanEventingInterface{}

			for _, m := range tt.authMockArgs {
				mockAuth.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.httpMockArgs {
				mockHTTP.On("Do", mock.Anything).Return(m.paramsReturned...)
			}
			for _, m := range tt.crudMockArgs {
				mockCrud.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.syncmanMockArgs {
				mockSyncman.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			tt.m.auth = &mockAuth
			tt.args.client = &mockHTTP
			tt.m.crud = &mockCrud
			tt.m.syncMan = &mockSyncman
			tt.m.updateEventC = make(chan *queueUpdateEvent, 5)

			if err := tt.m.invokeWebhook(context.Background(), "", tt.args.client, tt.args.rule, tt.args.eventDoc, tt.args.cloudEvent); (err != nil) != tt.wantErr {
				t.Errorf("Module.invokeWebhook() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockAuth.AssertExpectations(t)
			mockHTTP.AssertExpectations(t)
			mockCrud.AssertExpectations(t)
			mockSyncman.AssertExpectations(t)
		})
	}
}

func TestModule_processStagedEvent(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		eventDoc *model.EventDocument
	}
	tests := []struct {
		name            string
		m               *Module
		args            args
		crudMockArgs    []mockArgs
		syncmanMockArgs []mockArgs
		authMockArgs    []mockArgs
	}{
		{
			name: "event is already being processed",
			m:    &Module{},
			args: args{eventDoc: &model.EventDocument{ID: "loadedID"}},
		},
		{
			name: "error selecting rule",
			m:    &Module{config: &config.Eventing{Rules: map[string]*config.EventingTrigger{"notSomeRule": {}}, InternalRules: map[string]*config.EventingTrigger{"notSomeRule": {}}}},
			args: args{eventDoc: &model.EventDocument{ID: "eventID", Type: "someType", RuleName: "someRule"}},
		},
		// {
		// 	name: "error invoking webhook",
		// 	m:    &Module{config: &config.Eventing{Rules: map[string]*config.EventingTrigger{"someRule": config.EventingTrigger{}}, InternalRules: map[string]*config.EventingTrigger{"notSomeRule": config.EventingTrigger{}}}},
		// 	args: args{eventDoc: &model.EventDocument{ID: "eventID", Type: "someType", RuleName: "someRule", Payload: "payload"}},
		// 	syncmanMockArgs: []mockArgs{
		// 		mockArgs{
		// 			method:         "GetEventSource",
		// 			paramsReturned: []interface{}{"source"},
		// 		},
		// 	},
		// 	authMockArgs: []mockArgs{
		// 		mockArgs{
		// 			method:         "GetInternalAccessToken",
		// 			paramsReturned: []interface{}{"", errors.New("some error")},
		// 		},
		// 	},
		// 	crudMockArgs: []mockArgs{
		// 		mockArgs{
		// 			method:         "InternalUpdate",
		// 			args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, utils.TableEventingLogs, &model.UpdateRequest{Find: map[string]interface{}{"_id": "eventID"}, Operation: utils.All, Update: map[string]interface{}{"$set": map[string]interface{}{"status": utils.EventStatusFailed, "remark": "Max retires limit reached"}}}},
		// 			paramsReturned: []interface{}{nil},
		// 		},
		// 	},
		// },
		// {
		// 	name: "error invoking webhook and error in internal update",
		// 	m:    &Module{config: &config.Eventing{Rules: map[string]*config.EventingTrigger{"someRule": config.EventingTrigger{Retries: 2}}, InternalRules: map[string]*config.EventingTrigger{"notSomeRule": config.EventingTrigger{}}}},
		// 	args: args{eventDoc: &model.EventDocument{ID: "eventID", Type: "someType", RuleName: "someRule", Payload: "payload"}},
		// 	syncmanMockArgs: []mockArgs{
		// 		mockArgs{
		// 			method:         "GetEventSource",
		// 			paramsReturned: []interface{}{"source"},
		// 		},
		// 	},
		// 	authMockArgs: []mockArgs{
		// 		mockArgs{
		// 			method:         "GetInternalAccessToken",
		// 			paramsReturned: []interface{}{"", errors.New("some error")},
		// 		},
		// 	},
		// 	crudMockArgs: []mockArgs{
		// 		mockArgs{
		// 			method:         "InternalUpdate",
		// 			args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, utils.TableEventingLogs, &model.UpdateRequest{Find: map[string]interface{}{"_id": "eventID"}, Operation: utils.All, Update: map[string]interface{}{"$set": map[string]interface{}{"status": utils.EventStatusFailed, "remark": "Max retires limit reached"}}}},
		// 			paramsReturned: []interface{}{errors.New("some error")},
		// 		},
		// 	},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.m.processingEvents.Store("loadedID", "event")

			mockSyncman := mockSyncmanEventingInterface{}
			mockAuth := mockAuthEventingInterface{}
			mockCrud := mockCrudInterface{}

			for _, m := range tt.syncmanMockArgs {
				mockSyncman.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.authMockArgs {
				mockAuth.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.crudMockArgs {
				mockCrud.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			tt.m.syncMan = &mockSyncman
			tt.m.auth = &mockAuth
			tt.m.crud = &mockCrud
			tt.m.updateEventC = make(chan *queueUpdateEvent, 5)

			tt.m.processStagedEvent(tt.args.eventDoc)

			mockSyncman.AssertExpectations(t)
			mockAuth.AssertExpectations(t)
			mockCrud.AssertExpectations(t)
		})
	}
}
