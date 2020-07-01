package eventing

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func TestModule_SetRealtimeTriggers(t *testing.T) {
	type args struct {
		eventingRules []config.EventingRule
	}
	tests := []struct {
		name string
		m    *Module
		args args
		want map[string]config.EventingRule
	}{
		{
			name: "no rules with prefix 'realtime'",
			m:    &Module{config: &config.Eventing{InternalRules: map[string]config.EventingRule{"notrealtime": {Type: "type", Options: map[string]string{"db": "db", "col": "col"}}}}},
			args: args{eventingRules: []config.EventingRule{{Type: "type", Options: map[string]string{"db": "db", "col": "col"}}}},
			want: map[string]config.EventingRule{"notrealtime": {Type: "type", Options: map[string]string{"db": "db", "col": "col"}}, "realtime-db-col-type": {Type: "type", Options: map[string]string{"db": "db", "col": "col"}}},
		},
		{
			name: "rules with prefix 'realtime'",
			m:    &Module{config: &config.Eventing{InternalRules: map[string]config.EventingRule{"realtime-abc": {Type: "type", Options: map[string]string{"db": "db", "col": "col"}}}}},
			args: args{eventingRules: []config.EventingRule{{Type: "type", Options: map[string]string{"db": "db", "col": "col"}}}},
			want: map[string]config.EventingRule{"realtime-db-col-type": {Type: "type", Options: map[string]string{"db": "db", "col": "col"}}},
		},
		{
			name: "add eventing rules when no internal rules exist",
			m:    &Module{config: &config.Eventing{InternalRules: make(map[string]config.EventingRule)}},
			args: args{eventingRules: []config.EventingRule{{Type: "type", Options: map[string]string{"db": "db", "col": "col"}}, {Type: "type1", Options: map[string]string{"db": "db1", "col": "col1"}}}},
			want: map[string]config.EventingRule{"realtime-db-col-type": {Type: "type", Options: map[string]string{"db": "db", "col": "col"}}, "realtime-db1-col1-type1": {Type: "type1", Options: map[string]string{"db": "db1", "col": "col1"}}},
		},
		{
			name: "add eventing rules when no realtime internal rules exist",
			m:    &Module{config: &config.Eventing{InternalRules: map[string]config.EventingRule{"notrealtime": {Type: "type", Options: map[string]string{"db": "db", "col": "col"}}}}},
			args: args{eventingRules: []config.EventingRule{{Type: "type", Options: map[string]string{"db": "db", "col": "col"}}, {Type: "type1", Options: map[string]string{"db": "db1", "col": "col1"}}}},
			want: map[string]config.EventingRule{"notrealtime": {Type: "type", Options: map[string]string{"db": "db", "col": "col"}}, "realtime-db-col-type": {Type: "type", Options: map[string]string{"db": "db", "col": "col"}}, "realtime-db1-col1-type1": {Type: "type1", Options: map[string]string{"db": "db1", "col": "col1"}}},
		},
		{
			name: "add eventing rules when realtime internal rules exist",
			m:    &Module{config: &config.Eventing{InternalRules: map[string]config.EventingRule{"realtime-abc": {Type: "type", Options: map[string]string{"db": "db", "col": "col"}}, "realtime-def": {Type: "type", Options: map[string]string{"db": "db", "col": "col"}}}}},
			args: args{eventingRules: []config.EventingRule{{Type: "type", Options: map[string]string{"db": "db", "col": "col"}}, {Type: "type1", Options: map[string]string{"db": "db1", "col": "col1"}}}},
			want: map[string]config.EventingRule{"realtime-db-col-type": {Type: "type", Options: map[string]string{"db": "db", "col": "col"}}, "realtime-db1-col1-type1": {Type: "type1", Options: map[string]string{"db": "db1", "col": "col1"}}},
		},
		{
			name: "add eventing rules when realtime and non-realtime internal rules exist",
			m:    &Module{config: &config.Eventing{InternalRules: map[string]config.EventingRule{"realtime-abc": {Type: "type", Options: map[string]string{"db": "db", "col": "col"}}, "nonrealtime-def": {Type: "type", Options: map[string]string{"db": "db", "col": "col"}}}}},
			args: args{eventingRules: []config.EventingRule{{Type: "type", Options: map[string]string{"db": "db", "col": "col"}}, {Type: "type1", Options: map[string]string{"db": "db1", "col": "col1"}}}},
			want: map[string]config.EventingRule{"nonrealtime-def": {Type: "type", Options: map[string]string{"db": "db", "col": "col"}}, "realtime-db-col-type": {Type: "type", Options: map[string]string{"db": "db", "col": "col"}}, "realtime-db1-col1-type1": {Type: "type1", Options: map[string]string{"db": "db1", "col": "col1"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.SetRealtimeTriggers(tt.args.eventingRules)

			if !reflect.DeepEqual(tt.m.config.InternalRules, tt.want) {
				t.Errorf("Error: got %v; wanted %v", tt.m.config.InternalRules, tt.want)
				return
			}
		})
	}
}

func TestModule_IsEnabled(t *testing.T) {
	tests := []struct {
		name string
		m    *Module
		want bool
	}{
		{
			name: "config is nil",
			m:    &Module{},
			want: false,
		},
		{
			name: "config is not nil but enabled is nil",
			m:    &Module{config: &config.Eventing{}},
			want: false,
		},
		{
			name: "enabled is true",
			m:    &Module{config: &config.Eventing{Enabled: true}},
			want: true,
		},
		{
			name: "enabled is false",
			m:    &Module{config: &config.Eventing{Enabled: false}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.IsEnabled(); got != tt.want {
				t.Errorf("Module.IsEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModule_SendEventResponse(t *testing.T) {
	resultChan := make(chan interface{}, 1)
	type args struct {
		batchID string
		payload interface{}
	}
	tests := []struct {
		name string
		m    *Module
		args args
		rcv  interface{}
	}{
		{
			name: "key not present in eventChanMap",
			m:    &Module{},
			args: args{batchID: "notbatchid", payload: "payload"},
		},
		{
			name: "event response is sent",
			m:    &Module{},
			args: args{batchID: "batchid", payload: "payload"},
			rcv:  "payload",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.m.eventChanMap.Store("batchid", eventResponse{time: time.Now(), response: resultChan})

			tt.m.SendEventResponse(tt.args.batchID, tt.args.payload)

			if len(resultChan) > 1 {
				close(resultChan)
				for response := range resultChan {
					if !reflect.DeepEqual(response, tt.rcv) {
						t.Fatalf("SendEventResponse() = got - %v; wanted - %v", response, tt.rcv)
					}
				}
			}
		})
	}
}

func TestModule_QueueEvent(t *testing.T) {
	var res interface{}
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx     context.Context
		project string
		token   string
		req     *model.QueueEventRequest
	}
	tests := []struct {
		name            string
		m               *Module
		args            args
		crudMockArgs    []mockArgs
		syncmanMockArgs []mockArgs
		adminMockArgs   []mockArgs
		authMockArgs    []mockArgs
		want            interface{}
		wantErr         bool
	}{
		{
			name: "error validating",
			m:    &Module{project: mock.Anything, config: &config.Eventing{DBAlias: mock.Anything, Rules: map[string]config.EventingRule{"rule": {Type: "someType", Options: make(map[string]string)}}}},
			args: args{ctx: context.Background(), project: "project", token: "token", req: &model.QueueEventRequest{Type: "someType", Delay: int64(0), Timestamp: time.Now().Format(time.RFC3339), Payload: "payload", Options: make(map[string]string), IsSynchronous: false}},
			authMockArgs: []mockArgs{
				{
					method:         "IsEventingOpAuthorised",
					args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{errors.New("some error")},
				},
			},
			wantErr: true,
		},
		{
			name: "error batching requests",
			m:    &Module{project: mock.Anything, config: &config.Eventing{DBAlias: mock.Anything, Rules: map[string]config.EventingRule{"rule": {Type: "DB_INSERT", Options: make(map[string]string)}}}},
			args: args{ctx: context.Background(), project: "project", token: "token", req: &model.QueueEventRequest{Type: "DB_INSERT", Delay: int64(0), Timestamp: time.Now().Format(time.RFC3339), Payload: "payload", Options: make(map[string]string), IsSynchronous: false}},
			syncmanMockArgs: []mockArgs{
				{
					method:         "GetNodeID",
					args:           []interface{}{},
					paramsReturned: []interface{}{"nodeID"},
				},
			},
			crudMockArgs: []mockArgs{
				{
					method:         "InternalCreate",
					args:           []interface{}{context.Background(), mock.Anything, mock.Anything, utils.TableEventingLogs, mock.Anything, false},
					paramsReturned: []interface{}{errors.New("some error")},
				},
			},
			wantErr: true,
		},
		{
			name: "event is queued",
			m:    &Module{metricHook: func(project, eventingType string) {}, project: mock.Anything, config: &config.Eventing{DBAlias: mock.Anything, Rules: map[string]config.EventingRule{"rule": {Type: "DB_INSERT", Options: make(map[string]string)}}}},
			args: args{ctx: context.Background(), project: "project", token: "token", req: &model.QueueEventRequest{Type: "DB_INSERT", Delay: int64(0), Timestamp: time.Now().Format(time.RFC3339), Payload: "payload", Options: make(map[string]string), IsSynchronous: false}},
			crudMockArgs: []mockArgs{
				{
					method:         "InternalCreate",
					args:           []interface{}{context.Background(), mock.Anything, mock.Anything, utils.TableEventingLogs, mock.Anything, false},
					paramsReturned: []interface{}{nil},
				},
			},
			syncmanMockArgs: []mockArgs{
				{
					method:         "GetNodeID",
					args:           []interface{}{},
					paramsReturned: []interface{}{"nodeID"},
				},
				{
					method:         "GetAssignedSpaceCloudURL",
					args:           []interface{}{mock.Anything, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{mock.Anything, nil},
				},
				{
					method:         "MakeHTTPRequest",
					args:           []interface{}{mock.Anything, "POST", mock.Anything, mock.Anything, mock.Anything, mock.Anything, &res},
					paramsReturned: []interface{}{nil},
				},
			},
			adminMockArgs: []mockArgs{
				{
					method:         "GetInternalAccessToken",
					args:           []interface{}{},
					paramsReturned: []interface{}{mock.Anything, nil},
				},
			},
			authMockArgs: []mockArgs{
				{
					method:         "GetSCAccessToken",
					args:           []interface{}{},
					paramsReturned: []interface{}{mock.Anything, nil},
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockCrud := mockCrudInterface{}
			mockSyncman := mockSyncmanEventingInterface{}
			mockAdmin := mockAdminEventingInterface{}
			mockAuth := mockAuthEventingInterface{}

			for _, m := range tt.crudMockArgs {
				mockCrud.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.syncmanMockArgs {
				mockSyncman.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.adminMockArgs {
				mockAdmin.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.authMockArgs {
				mockAuth.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			tt.m.crud = &mockCrud
			tt.m.syncMan = &mockSyncman
			tt.m.adminMan = &mockAdmin
			tt.m.auth = &mockAuth

			got, err := tt.m.QueueEvent(tt.args.ctx, tt.args.project, tt.args.token, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Module.QueueEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Module.QueueEvent() = %v, want %v", got, tt.want)
			}

			mockCrud.AssertExpectations(t)
			mockSyncman.AssertExpectations(t)
			mockAdmin.AssertExpectations(t)
			mockAuth.AssertExpectations(t)
		})
	}
}

// TODO: write test case in QueueEvent where request is synchronous
