package eventing

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func TestModule_processCreateDocs(t *testing.T) {

	payload := model.DatabaseEventMessage{DBType: "db", Col: "col", Doc: map[string]interface{}{"key": "value"}, Find: map[string]interface{}{}}

	data, _ := json.Marshal(payload)

	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		token      int
		eventDocID string
		batchID    string
		dbAlias    string
		col        string
		rows       []interface{}
	}
	tests := []struct {
		name           string
		m              *Module
		args           args
		schemaMockArgs []mockArgs
		want           []*model.EventDocument
	}{
		{
			name: "length of rules is 0",
			m:    &Module{config: &config.Eventing{Rules: map[string]config.EventingRule{"rule": {Type: "someType", Options: map[string]string{"col": "col", "db": "db"}}}}},
			args: args{token: 50, eventDocID: "eventid", batchID: "batchid", dbAlias: "db", col: "col", rows: nil},
			want: nil,
		},
		{
			name: "rows are nil",
			m:    &Module{config: &config.Eventing{Rules: map[string]config.EventingRule{"rule": {Type: utils.EventDBCreate, Options: map[string]string{"col": "col", "db": "db"}}}}},
			args: args{token: 50, eventDocID: "eventid", batchID: "batchid", dbAlias: "db", col: "col", rows: nil},
			want: []*model.EventDocument{},
		},
		{
			name: "eventing is not possible",
			m:    &Module{config: &config.Eventing{Rules: map[string]config.EventingRule{"rule": {Type: utils.EventDBCreate, Options: map[string]string{"col": "col", "db": "db"}}}}},
			args: args{token: 50, eventDocID: "eventid", batchID: "batchid", dbAlias: "db", col: "col", rows: []interface{}{map[string]interface{}{"key": "value"}}},
			schemaMockArgs: []mockArgs{
				{
					method:         "CheckIfEventingIsPossible",
					args:           []interface{}{"db", "col", map[string]interface{}{"key": "value"}, false},
					paramsReturned: []interface{}{nil, false},
				},
			},
			want: nil,
		},
		{
			name: "doc are created",
			m:    &Module{config: &config.Eventing{Rules: map[string]config.EventingRule{"rule": {Type: utils.EventDBCreate, Options: map[string]string{"col": "col", "db": "db"}}}}},
			args: args{token: 50, batchID: "batchid", dbAlias: "db", col: "col", rows: []interface{}{map[string]interface{}{"key": "value"}}},
			schemaMockArgs: []mockArgs{
				{
					method:         "CheckIfEventingIsPossible",
					args:           []interface{}{"db", "col", map[string]interface{}{"key": "value"}, false},
					paramsReturned: []interface{}{map[string]interface{}{}, true},
				},
			},
			want: []*model.EventDocument{{BatchID: "batchid", Type: utils.EventDBCreate, RuleName: "rule", Token: 50, Timestamp: time.Now().Format(time.RFC3339), Payload: string(data), Status: utils.EventStatusIntent}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockSchema := mockSchemaEventingInterface{}

			for _, m := range tt.schemaMockArgs {
				mockSchema.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			tt.m.schema = &mockSchema

			got := tt.m.processCreateDocs(tt.args.token, tt.args.batchID, tt.args.dbAlias, tt.args.col, tt.args.rows)
			if got != nil && !reflect.DeepEqual(got, []*model.EventDocument{}) {
				if !reflect.DeepEqual(got[0].BatchID, tt.want[0].BatchID) {
					t.Errorf("Module.processCreateDocs() = %v, want %v", got[0].BatchID, tt.want[0].BatchID)
				}
				if !reflect.DeepEqual(got[0].Type, tt.want[0].Type) {
					t.Errorf("Module.processCreateDocs() = %v, want %v", got[0].Type, tt.want[0].Type)
				}
				if !reflect.DeepEqual(got[0].RuleName, tt.want[0].RuleName) {
					t.Errorf("Module.processCreateDocs() = %v, want %v", got[0].RuleName, tt.want[0].RuleName)
				}
				if !reflect.DeepEqual(got[0].Token, tt.want[0].Token) {
					t.Errorf("Module.processCreateDocs() = %v, want %v", got[0].Token, tt.want[0].Token)
				}
				if !reflect.DeepEqual(got[0].Timestamp, tt.want[0].Timestamp) {
					t.Errorf("Module.processCreateDocs() = %v, want %v", got[0].Timestamp, tt.want[0].Timestamp)
				}
				if !reflect.DeepEqual(got[0].Payload, tt.want[0].Payload) {
					t.Errorf("Module.processCreateDocs() = %v, want %v", got[0].Payload, tt.want[0].Payload)
				}
				if !reflect.DeepEqual(got[0].Status, tt.want[0].Status) {
					t.Errorf("Module.processCreateDocs() = %v, want %v", got[0].Status, tt.want[0].Status)
				}
			}

			mockSchema.AssertExpectations(t)
		})
	}
}

func TestModule_processUpdateDeleteHook(t *testing.T) {

	payload := model.DatabaseEventMessage{DBType: "db", Col: "col", Find: map[string]interface{}{}}
	data, _ := json.Marshal(payload)

	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		token     int
		eventType string
		batchID   string
		dbAlias   string
		col       string
		find      map[string]interface{}
	}
	tests := []struct {
		name           string
		m              *Module
		args           args
		schemaMockArgs []mockArgs
		want           []*model.EventDocument
		want1          bool
	}{
		{
			name:  "no rules matched",
			m:     &Module{config: &config.Eventing{Rules: map[string]config.EventingRule{"rule": {Type: "nottype", Options: map[string]string{"col": "col", "db": "db"}}}}},
			args:  args{token: 50, eventType: "type", batchID: "batchid", dbAlias: "db", col: "col", find: map[string]interface{}{"key": "value"}},
			want:  nil,
			want1: false,
		},
		{
			name: "eventing is not possible",
			m:    &Module{config: &config.Eventing{Rules: map[string]config.EventingRule{"rule": {Type: "type", Options: map[string]string{"col": "col", "db": "db"}}}}},
			args: args{token: 50, eventType: "type", batchID: "batchid", dbAlias: "db", col: "col", find: map[string]interface{}{"key": "value"}},
			schemaMockArgs: []mockArgs{
				{
					method:         "CheckIfEventingIsPossible",
					args:           []interface{}{"db", "col", map[string]interface{}{"key": "value"}, true},
					paramsReturned: []interface{}{nil, false},
				},
			},
			want:  nil,
			want1: false,
		},
		{
			name: "docs are updated",
			m:    &Module{config: &config.Eventing{Rules: map[string]config.EventingRule{"rule": {Type: "type", Options: map[string]string{"col": "col", "db": "db"}}}}},
			args: args{token: 50, eventType: "type", batchID: "batchid", dbAlias: "db", col: "col", find: map[string]interface{}{"key": "value"}},
			schemaMockArgs: []mockArgs{
				{
					method:         "CheckIfEventingIsPossible",
					args:           []interface{}{"db", "col", map[string]interface{}{"key": "value"}, true},
					paramsReturned: []interface{}{map[string]interface{}{}, true},
				},
			},
			want:  []*model.EventDocument{{BatchID: "batchid", Type: "type", RuleName: "rule", Token: 50, Timestamp: time.Now().Format(time.RFC3339), Payload: string(data), Status: utils.EventStatusIntent}},
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockSchema := mockSchemaEventingInterface{}

			for _, m := range tt.schemaMockArgs {
				mockSchema.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			tt.m.schema = &mockSchema

			got, got1 := tt.m.processUpdateDeleteHook(tt.args.token, tt.args.eventType, tt.args.batchID, tt.args.dbAlias, tt.args.col, tt.args.find)
			if got != nil && !reflect.DeepEqual(got, []*model.EventDocument{}) {
				if !reflect.DeepEqual(got[0].BatchID, tt.want[0].BatchID) {
					t.Errorf("Module.processCreateDocs() = %v, want %v", got[0].BatchID, tt.want[0].BatchID)
				}
				if !reflect.DeepEqual(got[0].Type, tt.want[0].Type) {
					t.Errorf("Module.processCreateDocs() = %v, want %v", got[0].Type, tt.want[0].Type)
				}
				if !reflect.DeepEqual(got[0].RuleName, tt.want[0].RuleName) {
					t.Errorf("Module.processCreateDocs() = %v, want %v", got[0].RuleName, tt.want[0].RuleName)
				}
				if !reflect.DeepEqual(got[0].Token, tt.want[0].Token) {
					t.Errorf("Module.processCreateDocs() = %v, want %v", got[0].Token, tt.want[0].Token)
				}
				if !reflect.DeepEqual(got[0].Timestamp, tt.want[0].Timestamp) {
					t.Errorf("Module.processCreateDocs() = %v, want %v", got[0].Timestamp, tt.want[0].Timestamp)
				}
				if !reflect.DeepEqual(got[0].Payload, tt.want[0].Payload) {
					t.Errorf("Module.processCreateDocs() = %v, want %v", got[0].Payload, tt.want[0].Payload)
				}
				if !reflect.DeepEqual(got[0].Status, tt.want[0].Status) {
					t.Errorf("Module.processCreateDocs() = %v, want %v", got[0].Status, tt.want[0].Status)
				}
			}
			if got1 != tt.want1 {
				t.Errorf("Module.processUpdateDeleteHook() got1 = %v, want %v", got1, tt.want1)
			}

			mockSchema.AssertExpectations(t)
		})
	}
}

func TestModule_HookStage(t *testing.T) {
	var res interface{}
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx    context.Context
		intent *model.EventIntent
		err    error
	}
	tests := []struct {
		name          string
		m             *Module
		args          args
		crudMockArgs  []mockArgs
		syncMockArgs  []mockArgs
		adminMockArgs []mockArgs
		authMockArgs  []mockArgs
	}{
		{
			name: "intent is invalid",
			m:    &Module{},
			args: args{ctx: context.Background(), intent: &model.EventIntent{Invalid: true}},
		},
		{
			name: "error is not nil",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype"}},
			args: args{ctx: context.Background(), intent: &model.EventIntent{BatchID: "batchid", Token: 50, Invalid: false, Docs: []*model.EventDocument{{Type: "notUpdate"}}}, err: errors.New("some error")},
			crudMockArgs: []mockArgs{
				{
					method:         "InternalUpdate",
					args:           []interface{}{mock.Anything, "dbtype", "abc", utils.TableEventingLogs, &model.UpdateRequest{Find: map[string]interface{}{"batchid": "batchid"}, Operation: utils.All, Update: map[string]interface{}{"$set": map[string]interface{}{"status": "cancel", "remark": "some error"}}}},
					paramsReturned: []interface{}{nil},
				},
			},
		},
		{
			name: "error is not nil and unable to update internal",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype"}},
			args: args{ctx: context.Background(), intent: &model.EventIntent{BatchID: "batchid", Token: 50, Invalid: false, Docs: []*model.EventDocument{{Type: "notUpdate"}}}, err: errors.New("some error")},
			crudMockArgs: []mockArgs{
				{
					method:         "InternalUpdate",
					args:           []interface{}{mock.Anything, "dbtype", "abc", utils.TableEventingLogs, &model.UpdateRequest{Find: map[string]interface{}{"batchid": "batchid"}, Operation: utils.All, Update: map[string]interface{}{"$set": map[string]interface{}{"status": "cancel", "remark": "some error"}}}},
					paramsReturned: []interface{}{errors.New("some error")},
				},
			},
		},
		{
			name: "error is nil and unable to update internal",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype"}},
			args: args{ctx: context.Background(), intent: &model.EventIntent{BatchID: "batchid", Token: 50, Invalid: false, Docs: []*model.EventDocument{{Type: "notUpdate"}}}},
			crudMockArgs: []mockArgs{
				{
					method:         "InternalUpdate",
					args:           []interface{}{mock.Anything, "dbtype", "abc", utils.TableEventingLogs, &model.UpdateRequest{Find: map[string]interface{}{"batchid": "batchid"}, Operation: utils.All, Update: map[string]interface{}{"$set": map[string]interface{}{"status": "staged"}}}},
					paramsReturned: []interface{}{errors.New("some error")},
				},
			},
		},
		{
			name: "event is staged not of type DB_UPDATE",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype"}},
			args: args{ctx: context.Background(), intent: &model.EventIntent{BatchID: "batchid", Token: 50, Invalid: false, Docs: []*model.EventDocument{{Type: "notUpdate"}}}},
			crudMockArgs: []mockArgs{
				{
					method:         "InternalUpdate",
					args:           []interface{}{mock.Anything, "dbtype", "abc", utils.TableEventingLogs, &model.UpdateRequest{Find: map[string]interface{}{"batchid": "batchid"}, Operation: utils.All, Update: map[string]interface{}{"$set": map[string]interface{}{"status": "staged"}}}},
					paramsReturned: []interface{}{nil},
				},
			},
			syncMockArgs: []mockArgs{
				{
					method:         "GetAssignedSpaceCloudURL",
					args:           []interface{}{mock.Anything, "abc", 50},
					paramsReturned: []interface{}{"url", nil},
				},
				{
					method:         "MakeHTTPRequest",
					args:           []interface{}{mock.Anything, "POST", "url", mock.Anything, mock.Anything, []*model.EventDocument{{Status: "staged", Type: "notUpdate"}}, &res},
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
		},
		{
			name: "type DB_UPDATE but error unmarshalling",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype"}},
			args: args{ctx: context.Background(), intent: &model.EventIntent{BatchID: "batchid", Token: 50, Invalid: false, Docs: []*model.EventDocument{{Type: utils.EventDBUpdate, Payload: "payload"}}}},
			crudMockArgs: []mockArgs{
				{
					method:         "InternalUpdate",
					args:           []interface{}{mock.Anything, "dbtype", "abc", utils.TableEventingLogs, &model.UpdateRequest{Find: map[string]interface{}{"batchid": "batchid"}, Operation: utils.All, Update: map[string]interface{}{"$set": map[string]interface{}{"status": "staged"}}}},
					paramsReturned: []interface{}{nil},
				},
			},
			syncMockArgs: []mockArgs{
				{
					method:         "GetAssignedSpaceCloudURL",
					args:           []interface{}{mock.Anything, "abc", 50},
					paramsReturned: []interface{}{"url", nil},
				},
				{
					method:         "MakeHTTPRequest",
					args:           []interface{}{mock.Anything, "POST", "url", mock.Anything, mock.Anything, []*model.EventDocument{{Status: "staged", Type: "DB_UPDATE", Payload: "payload"}}, &res},
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
		},
		{
			name: "error reading",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype"}},
			args: args{ctx: context.Background(), intent: &model.EventIntent{BatchID: "batchid", Token: 50, Invalid: false, Docs: []*model.EventDocument{{Type: utils.EventDBUpdate, Payload: `{"db": "db", "col": "col", "find": {"key1": "value1"}}`}}}},
			crudMockArgs: []mockArgs{
				{
					method:         "Read",
					args:           []interface{}{mock.Anything, "db", "col", &model.ReadRequest{Find: map[string]interface{}{"key1": "value1"}, Operation: utils.One}},
					paramsReturned: []interface{}{map[string]interface{}{}, errors.New("some error")},
				},
				{
					method:         "InternalUpdate",
					args:           []interface{}{mock.Anything, "dbtype", "abc", utils.TableEventingLogs, mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			syncMockArgs: []mockArgs{
				{
					method:         "GetAssignedSpaceCloudURL",
					args:           []interface{}{mock.Anything, "abc", 50},
					paramsReturned: []interface{}{"url", nil},
				},
				{
					method:         "MakeHTTPRequest",
					args:           []interface{}{mock.Anything, "POST", "url", mock.Anything, mock.Anything, []*model.EventDocument{{Status: "staged", Type: "DB_UPDATE", Timestamp: "", Payload: "{\"db\": \"db\", \"col\": \"col\", \"find\": {\"key1\": \"value1\"}}"}}, &res},
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
		},
		{
			name: "type DB_UPDATE event is staged",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype"}},
			args: args{ctx: context.Background(), intent: &model.EventIntent{BatchID: "batchid", Token: 50, Invalid: false, Docs: []*model.EventDocument{{Type: utils.EventDBUpdate, Payload: `{"db": "db", "col": "col", "find": {"key1": "value1"}}`}}}},
			crudMockArgs: []mockArgs{
				{
					method:         "Read",
					args:           []interface{}{mock.Anything, "db", "col", &model.ReadRequest{Find: map[string]interface{}{"key1": "value1"}, Operation: utils.One}},
					paramsReturned: []interface{}{map[string]interface{}{}, nil},
				},
				{
					method:         "InternalUpdate",
					args:           []interface{}{mock.Anything, "dbtype", "abc", utils.TableEventingLogs, mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			syncMockArgs: []mockArgs{
				{
					method:         "GetAssignedSpaceCloudURL",
					args:           []interface{}{mock.Anything, "abc", 50},
					paramsReturned: []interface{}{"url", nil},
				},
				{
					method:         "MakeHTTPRequest",
					args:           []interface{}{mock.Anything, "POST", "url", mock.Anything, mock.Anything, []*model.EventDocument{{Status: "staged", Type: "DB_UPDATE", Timestamp: time.Now().Format(time.RFC3339), Payload: "{\"db\":\"db\",\"col\":\"col\",\"doc\":{},\"find\":{\"key1\":\"value1\"}}"}}, &res},
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockCrud := mockCrudInterface{}
			mockAuth := mockAuthEventingInterface{}
			mockAdmin := mockAdminEventingInterface{}
			mockSyncman := mockSyncmanEventingInterface{}

			for _, m := range tt.crudMockArgs {
				mockCrud.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.syncMockArgs {
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

			tt.m.HookStage(tt.args.ctx, tt.args.intent, tt.args.err)

			mockCrud.AssertExpectations(t)
			mockSyncman.AssertExpectations(t)
			mockAdmin.AssertExpectations(t)
			mockAuth.AssertExpectations(t)
		})
	}
}

func TestModule_hookDBUpdateDeleteIntent(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx       context.Context
		eventType string
		dbAlias   string
		col       string
		find      map[string]interface{}
	}
	tests := []struct {
		name            string
		m               *Module
		args            args
		syncmanMockArgs []mockArgs
		schemaMockArgs  []mockArgs
		crudMockArgs    []mockArgs
		want            *model.EventIntent
		wantErr         bool
	}{
		{
			name: "update delete hook not ok",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype", Rules: map[string]config.EventingRule{"rule": {Type: "notEvType", Options: map[string]string{"col": "col", "db": "db"}}}}},
			args: args{ctx: context.Background(), eventType: "evType", dbAlias: "db", col: "col", find: map[string]interface{}{}},
			syncmanMockArgs: []mockArgs{
				{
					method:         "GetNodeID",
					paramsReturned: []interface{}{"nodeid"},
				},
			},
			want: &model.EventIntent{Invalid: true},
		},
		{
			name: "error in internal create",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype", Rules: map[string]config.EventingRule{"rule": {Type: "evType", Options: map[string]string{"col": "col", "db": "db"}}}}},
			args: args{ctx: context.Background(), eventType: "evType", dbAlias: "db", col: "col", find: map[string]interface{}{}},
			syncmanMockArgs: []mockArgs{
				{
					method:         "GetNodeID",
					paramsReturned: []interface{}{"nodeid"},
				},
			},
			schemaMockArgs: []mockArgs{
				{
					method:         "CheckIfEventingIsPossible",
					args:           []interface{}{"db", "col", map[string]interface{}{}, true},
					paramsReturned: []interface{}{map[string]interface{}{}, true},
				},
			},
			crudMockArgs: []mockArgs{
				{
					method:         "InternalCreate",
					args:           []interface{}{mock.Anything, "dbtype", "abc", utils.TableEventingLogs, mock.Anything, false},
					paramsReturned: []interface{}{errors.New("some error")},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "no errors",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype", Rules: map[string]config.EventingRule{"rule": {Type: "evType", Options: map[string]string{"col": "col", "db": "db"}}}}},
			args: args{ctx: context.Background(), eventType: "evType", dbAlias: "db", col: "col", find: map[string]interface{}{}},
			syncmanMockArgs: []mockArgs{
				{
					method:         "GetNodeID",
					paramsReturned: []interface{}{"nodeid"},
				},
			},
			schemaMockArgs: []mockArgs{
				{
					method:         "CheckIfEventingIsPossible",
					args:           []interface{}{"db", "col", map[string]interface{}{}, true},
					paramsReturned: []interface{}{map[string]interface{}{}, true},
				},
			},
			crudMockArgs: []mockArgs{
				{
					method:         "InternalCreate",
					args:           []interface{}{mock.Anything, "dbtype", "abc", utils.TableEventingLogs, mock.Anything, false},
					paramsReturned: []interface{}{nil},
				},
			},
			want: &model.EventIntent{Docs: []*model.EventDocument{{Type: "evType", RuleName: "rule", Timestamp: time.Now().Format(time.RFC3339), Payload: `{"db":"db","col":"col","doc":null,"find":{}}`, Status: "intent"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockSyncman := mockSyncmanEventingInterface{}
			mockSchema := mockSchemaEventingInterface{}
			mockCrud := mockCrudInterface{}

			for _, m := range tt.syncmanMockArgs {
				mockSyncman.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.schemaMockArgs {
				mockSchema.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.crudMockArgs {
				mockCrud.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			tt.m.syncMan = &mockSyncman
			tt.m.schema = &mockSchema
			tt.m.crud = &mockCrud

			got, err := tt.m.hookDBUpdateDeleteIntent(tt.args.ctx, tt.args.eventType, tt.args.dbAlias, tt.args.col, tt.args.find)
			if (err != nil) != tt.wantErr {
				t.Errorf("Module.hookDBUpdateDeleteIntent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != nil && tt.want != nil {
				invalid := got.Invalid

				if !reflect.DeepEqual(invalid, tt.want.Invalid) {
					t.Errorf("Module.hookDBUpdateDeleteIntent() = %v, want %v", invalid, tt.want.Invalid)
				}

				docs := got.Docs
				if docs != nil {
					if !reflect.DeepEqual(docs[0].Type, tt.want.Docs[0].Type) {
						t.Errorf("Module.hookDBUpdateDeleteIntent() = %v, want %v", docs[0].Type, tt.want.Docs[0].Type)
					}
					if !reflect.DeepEqual(docs[0].RuleName, tt.want.Docs[0].RuleName) {
						t.Errorf("Module.hookDBUpdateDeleteIntent() = %v, want %v", docs[0].RuleName, tt.want.Docs[0].RuleName)
					}
					if !reflect.DeepEqual(docs[0].Timestamp, tt.want.Docs[0].Timestamp) {
						t.Errorf("Module.hookDBUpdateDeleteIntent() = %v, want %v", docs[0].Timestamp, tt.want.Docs[0].Timestamp)
					}
					if !reflect.DeepEqual(docs[0].Payload, tt.want.Docs[0].Payload) {
						t.Errorf("Module.hookDBUpdateDeleteIntent() = %v, want %v", docs[0].Payload, tt.want.Docs[0].Payload)
					}
					if !reflect.DeepEqual(docs[0].Status, tt.want.Docs[0].Status) {
						t.Errorf("Module.hookDBUpdateDeleteIntent() = %v, want %v", docs[0].Status, tt.want.Docs[0].Status)
					}
				}
			}

			mockSyncman.AssertExpectations(t)
			mockSchema.AssertExpectations(t)
			mockCrud.AssertExpectations(t)
		})
	}
}

func TestModule_HookDBDeleteIntent(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx     context.Context
		dbAlias string
		col     string
		req     *model.DeleteRequest
	}
	tests := []struct {
		name            string
		m               *Module
		args            args
		syncmanMockArgs []mockArgs
		want            *model.EventIntent
		wantErr         bool
	}{
		{
			name: " eventing not enabled",
			m:    &Module{config: &config.Eventing{Enabled: false, Rules: map[string]config.EventingRule{"rule": {Type: "not delete"}}}},
			args: args{context.Background(), "db", "col", nil},
			want: &model.EventIntent{Invalid: true},
		},
		{
			name: "eventing is enabled",
			m:    &Module{config: &config.Eventing{Enabled: true}},
			args: args{ctx: context.Background(), dbAlias: "db", col: "col", req: &model.DeleteRequest{Find: map[string]interface{}{}}},
			syncmanMockArgs: []mockArgs{
				{
					method:         "GetNodeID",
					paramsReturned: []interface{}{"nodeid"},
				},
			},
			want: &model.EventIntent{Invalid: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockSyncman := mockSyncmanEventingInterface{}

			for _, m := range tt.syncmanMockArgs {
				mockSyncman.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			tt.m.syncMan = &mockSyncman

			got, err := tt.m.HookDBDeleteIntent(tt.args.ctx, tt.args.dbAlias, tt.args.col, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Module.HookDBDeleteIntent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Module.HookDBDeleteIntent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModule_HookDBUpdateIntent(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx     context.Context
		dbAlias string
		col     string
		req     *model.UpdateRequest
	}
	tests := []struct {
		name            string
		m               *Module
		args            args
		syncmanMockArgs []mockArgs
		want            *model.EventIntent
		wantErr         bool
	}{
		{
			name: " eventing not enabled",
			m:    &Module{config: &config.Eventing{Enabled: false, Rules: map[string]config.EventingRule{"rule": {Type: "not update"}}}},
			args: args{context.Background(), "db", "col", nil},
			want: &model.EventIntent{Invalid: true},
		},
		{
			name: "eventing is enabled",
			m:    &Module{config: &config.Eventing{Enabled: true}},
			args: args{ctx: context.Background(), dbAlias: "db", col: "col", req: &model.UpdateRequest{Find: map[string]interface{}{}}},
			syncmanMockArgs: []mockArgs{
				{
					method:         "GetNodeID",
					paramsReturned: []interface{}{"nodeid"},
				},
			},
			want: &model.EventIntent{Invalid: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockSyncman := mockSyncmanEventingInterface{}

			for _, m := range tt.syncmanMockArgs {
				mockSyncman.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			tt.m.syncMan = &mockSyncman

			got, err := tt.m.HookDBUpdateIntent(tt.args.ctx, tt.args.dbAlias, tt.args.col, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Module.HookDBUpdateIntent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Module.HookDBUpdateIntent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModule_HookDBCreateIntent(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx     context.Context
		dbAlias string
		col     string
		req     *model.CreateRequest
	}
	tests := []struct {
		name            string
		m               *Module
		args            args
		syncmanMockArgs []mockArgs
		schemaMockArgs  []mockArgs
		crudMockArgs    []mockArgs
		want            *model.EventIntent
		wantErr         bool
	}{
		{
			name: "eventing not enabled",
			m:    &Module{config: &config.Eventing{Enabled: false}},
			args: args{ctx: context.Background(), dbAlias: "db", col: "col", req: nil},
			want: &model.EventIntent{Invalid: true},
		},
		{
			name: "len of eventDocs is 0",
			m:    &Module{config: &config.Eventing{Enabled: true}},
			args: args{ctx: context.Background(), dbAlias: "db", col: "col", req: &model.CreateRequest{Operation: "default", Document: map[string]interface{}{}, IsBatch: true}},
			syncmanMockArgs: []mockArgs{
				{
					method:         "GetNodeID",
					paramsReturned: []interface{}{"nodeid"},
				},
			},
			want: &model.EventIntent{Invalid: true},
		},
		{
			name: "error creating internal",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype", Enabled: true, Rules: map[string]config.EventingRule{"rule": {Type: utils.EventDBCreate, Options: map[string]string{"db": "db", "col": "col"}}}}},
			args: args{ctx: context.Background(), dbAlias: "db", col: "col", req: &model.CreateRequest{Operation: "one", Document: map[string]interface{}{"key": "value"}, IsBatch: true}},
			syncmanMockArgs: []mockArgs{
				{
					method:         "GetNodeID",
					paramsReturned: []interface{}{"nodeid"},
				},
			},
			schemaMockArgs: []mockArgs{
				{
					method:         "CheckIfEventingIsPossible",
					args:           []interface{}{"db", "col", map[string]interface{}{"key": "value"}, false},
					paramsReturned: []interface{}{map[string]interface{}{}, true},
				},
			},
			crudMockArgs: []mockArgs{
				{
					method:         "InternalCreate",
					args:           []interface{}{mock.Anything, "dbtype", "abc", utils.TableEventingLogs, mock.Anything, false},
					paramsReturned: []interface{}{errors.New("some error")},
				},
			},
			wantErr: true,
		},
		{
			name: "no errors",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype", Enabled: true, Rules: map[string]config.EventingRule{"rule": {Type: utils.EventDBCreate, Options: map[string]string{"db": "db", "col": "col"}}}}},
			args: args{ctx: context.Background(), dbAlias: "db", col: "col", req: &model.CreateRequest{Operation: "one", Document: map[string]interface{}{"key": "value"}, IsBatch: true}},
			syncmanMockArgs: []mockArgs{
				{
					method:         "GetNodeID",
					paramsReturned: []interface{}{"nodeid"},
				},
			},
			schemaMockArgs: []mockArgs{
				{
					method:         "CheckIfEventingIsPossible",
					args:           []interface{}{"db", "col", map[string]interface{}{"key": "value"}, false},
					paramsReturned: []interface{}{map[string]interface{}{}, true},
				},
			},
			crudMockArgs: []mockArgs{
				{
					method:         "InternalCreate",
					args:           []interface{}{mock.Anything, "dbtype", "abc", utils.TableEventingLogs, mock.Anything, false},
					paramsReturned: []interface{}{nil},
				},
			},
			want: &model.EventIntent{Docs: []*model.EventDocument{{Type: "DB_INSERT", RuleName: "rule", Timestamp: time.Now().Format(time.RFC3339), Payload: `{"db":"db","col":"col","doc":{"key":"value"},"find":{}}`, Status: "intent"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockSyncman := mockSyncmanEventingInterface{}
			mockSchema := mockSchemaEventingInterface{}
			mockCrud := mockCrudInterface{}

			for _, m := range tt.syncmanMockArgs {
				mockSyncman.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.schemaMockArgs {
				mockSchema.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.crudMockArgs {
				mockCrud.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			tt.m.syncMan = &mockSyncman
			tt.m.schema = &mockSchema
			tt.m.crud = &mockCrud

			got, err := tt.m.HookDBCreateIntent(tt.args.ctx, tt.args.dbAlias, tt.args.col, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Module.HookDBCreateIntent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && tt.want != nil {
				invalid := got.Invalid

				if !reflect.DeepEqual(invalid, tt.want.Invalid) {
					t.Errorf("Module.HookDBCreateIntent() = %v, want %v", invalid, tt.want.Invalid)
				}

				docs := got.Docs
				if docs != nil {
					if !reflect.DeepEqual(docs[0].Type, tt.want.Docs[0].Type) {
						t.Errorf("Module.HookDBCreateIntent() = %v, want %v", docs[0].Type, tt.want.Docs[0].Type)
					}
					if !reflect.DeepEqual(docs[0].RuleName, tt.want.Docs[0].RuleName) {
						t.Errorf("Module.HookDBCreateIntent() = %v, want %v", docs[0].RuleName, tt.want.Docs[0].RuleName)
					}
					if !reflect.DeepEqual(docs[0].Timestamp, tt.want.Docs[0].Timestamp) {
						t.Errorf("Module.HookDBCreateIntent() = %v, want %v", docs[0].Timestamp, tt.want.Docs[0].Timestamp)
					}
					if !reflect.DeepEqual(docs[0].Payload, tt.want.Docs[0].Payload) {
						t.Errorf("Module.HookDBCreateIntent() = %v, want %v", docs[0].Payload, tt.want.Docs[0].Payload)
					}
					if !reflect.DeepEqual(docs[0].Status, tt.want.Docs[0].Status) {
						t.Errorf("Module.HookDBCreateIntent() = %v, want %v", docs[0].Status, tt.want.Docs[0].Status)
					}
				}
			}

			mockSyncman.AssertExpectations(t)
			mockSchema.AssertExpectations(t)
			mockCrud.AssertExpectations(t)
		})
	}
}

func TestModule_HookDBBatchIntent(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx     context.Context
		dbAlias string
		req     *model.BatchRequest
	}
	tests := []struct {
		name            string
		m               *Module
		args            args
		syncmanMockArgs []mockArgs
		schemaMockArgs  []mockArgs
		crudMockArgs    []mockArgs
		want            *model.EventIntent
		wantErr         bool
	}{
		{
			name: "eventing is not enabled",
			m:    &Module{config: &config.Eventing{Enabled: false}},
			args: args{ctx: context.Background(), dbAlias: "db", req: &model.BatchRequest{Requests: []*model.AllRequest{{Col: "col"}}}},
			want: &model.EventIntent{Invalid: true},
		},
		{
			name: "r.type default case",
			m:    &Module{config: &config.Eventing{Enabled: true}},
			args: args{ctx: context.Background(), dbAlias: "db", req: &model.BatchRequest{Requests: []*model.AllRequest{{Type: "default", Col: "col"}}}},
			syncmanMockArgs: []mockArgs{
				{
					method:         "GetNodeID",
					paramsReturned: []interface{}{"nodeID"},
				},
			},
			wantErr: true,
		},
		{
			name: "error creating internal",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype", Enabled: true, Rules: map[string]config.EventingRule{"rule": {Type: utils.EventDBCreate, Options: map[string]string{"db": "db", "col": "col"}}}}},
			args: args{ctx: context.Background(), dbAlias: "db", req: &model.BatchRequest{Requests: []*model.AllRequest{{Type: "create", Col: "col", Document: map[string]interface{}{"key": "value"}, Operation: "one"}}}},
			syncmanMockArgs: []mockArgs{
				{
					method:         "GetNodeID",
					paramsReturned: []interface{}{"nodeID"},
				},
			},
			schemaMockArgs: []mockArgs{
				{
					method:         "CheckIfEventingIsPossible",
					args:           []interface{}{"db", "col", map[string]interface{}{"key": "value"}, false},
					paramsReturned: []interface{}{map[string]interface{}{}, true},
				},
			},
			crudMockArgs: []mockArgs{
				{
					method:         "InternalCreate",
					args:           []interface{}{mock.Anything, "dbtype", "abc", utils.TableEventingLogs, mock.Anything, false},
					paramsReturned: []interface{}{errors.New("some error")},
				},
			},
			wantErr: true,
		},
		{
			name: "r.type create case",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype", Enabled: true, Rules: map[string]config.EventingRule{"rule": {Type: utils.EventDBCreate, Options: map[string]string{"db": "db", "col": "col"}}}}},
			args: args{ctx: context.Background(), dbAlias: "db", req: &model.BatchRequest{Requests: []*model.AllRequest{{Type: "create", Col: "col", Document: map[string]interface{}{"key": "value"}, Operation: "one"}}}},
			syncmanMockArgs: []mockArgs{
				{
					method:         "GetNodeID",
					paramsReturned: []interface{}{"nodeID"},
				},
			},
			schemaMockArgs: []mockArgs{
				{
					method:         "CheckIfEventingIsPossible",
					args:           []interface{}{"db", "col", map[string]interface{}{"key": "value"}, false},
					paramsReturned: []interface{}{map[string]interface{}{}, true},
				},
			},
			crudMockArgs: []mockArgs{
				{
					method:         "InternalCreate",
					args:           []interface{}{mock.Anything, "dbtype", "abc", utils.TableEventingLogs, mock.Anything, false},
					paramsReturned: []interface{}{nil},
				},
			},
			want: &model.EventIntent{Docs: []*model.EventDocument{{Type: "DB_INSERT", RuleName: "rule", Timestamp: time.Now().Format(time.RFC3339), Payload: `{"db":"db","col":"col","doc":{"key":"value"},"find":{}}`, Status: "intent"}}},
		},
		{
			name: "r.type update case where not ok",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype", Enabled: true, Rules: map[string]config.EventingRule{"rule": {Type: "not update", Options: map[string]string{"db": "db", "col": "col"}}}}},
			args: args{ctx: context.Background(), dbAlias: "db", req: &model.BatchRequest{Requests: []*model.AllRequest{{Type: "update", Col: "col", Find: map[string]interface{}{"key": "value"}}}}},
			syncmanMockArgs: []mockArgs{
				{
					method:         "GetNodeID",
					paramsReturned: []interface{}{"nodeID"},
				},
			},
			want: &model.EventIntent{Invalid: true},
		},
		{
			name: "r.type update case",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype", Enabled: true, Rules: map[string]config.EventingRule{"rule": {Type: utils.EventDBUpdate, Options: map[string]string{"db": "db", "col": "col"}}}}},
			args: args{ctx: context.Background(), dbAlias: "db", req: &model.BatchRequest{Requests: []*model.AllRequest{{Type: "update", Col: "col", Find: map[string]interface{}{"key": "value"}}}}},
			syncmanMockArgs: []mockArgs{
				{
					method:         "GetNodeID",
					paramsReturned: []interface{}{"nodeID"},
				},
			},
			schemaMockArgs: []mockArgs{
				{
					method:         "CheckIfEventingIsPossible",
					args:           []interface{}{"db", "col", map[string]interface{}{"key": "value"}, true},
					paramsReturned: []interface{}{map[string]interface{}{}, true},
				},
			},
			crudMockArgs: []mockArgs{
				{
					method:         "InternalCreate",
					args:           []interface{}{mock.Anything, "dbtype", "abc", utils.TableEventingLogs, mock.Anything, false},
					paramsReturned: []interface{}{nil},
				},
			},
			want: &model.EventIntent{Docs: []*model.EventDocument{{Type: "DB_UPDATE", RuleName: "rule", Timestamp: time.Now().Format(time.RFC3339), Payload: `{"db":"db","col":"col","doc":null,"find":{}}`, Status: "intent"}}},
		},
		{
			name: "r.type delete case",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype", Enabled: true, Rules: map[string]config.EventingRule{"rule": {Type: utils.EventDBDelete, Options: map[string]string{"db": "db", "col": "col"}}}}},
			args: args{ctx: context.Background(), dbAlias: "db", req: &model.BatchRequest{Requests: []*model.AllRequest{{Type: "delete", Col: "col", Find: map[string]interface{}{"key": "value"}}}}},
			syncmanMockArgs: []mockArgs{
				{
					method:         "GetNodeID",
					paramsReturned: []interface{}{"nodeID"},
				},
			},
			schemaMockArgs: []mockArgs{
				{
					method:         "CheckIfEventingIsPossible",
					args:           []interface{}{"db", "col", map[string]interface{}{"key": "value"}, true},
					paramsReturned: []interface{}{map[string]interface{}{}, true},
				},
			},
			crudMockArgs: []mockArgs{
				{
					method:         "InternalCreate",
					args:           []interface{}{mock.Anything, "dbtype", "abc", utils.TableEventingLogs, mock.Anything, false},
					paramsReturned: []interface{}{nil},
				},
			},
			want: &model.EventIntent{Docs: []*model.EventDocument{{Type: "DB_DELETE", RuleName: "rule", Timestamp: time.Now().Format(time.RFC3339), Payload: `{"db":"db","col":"col","doc":null,"find":{}}`, Status: "intent"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockSyncman := mockSyncmanEventingInterface{}
			mockSchema := mockSchemaEventingInterface{}
			mockCrud := mockCrudInterface{}

			for _, m := range tt.syncmanMockArgs {
				mockSyncman.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			for _, m := range tt.schemaMockArgs {
				mockSchema.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			for _, m := range tt.crudMockArgs {
				mockCrud.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			tt.m.syncMan = &mockSyncman
			tt.m.schema = &mockSchema
			tt.m.crud = &mockCrud

			got, err := tt.m.HookDBBatchIntent(tt.args.ctx, tt.args.dbAlias, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Module.HookDBBatchIntent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && tt.want != nil {
				invalid := got.Invalid

				if !reflect.DeepEqual(invalid, tt.want.Invalid) {
					t.Errorf("Module.HookDBBatchIntent() = %v, want %v", invalid, tt.want.Invalid)
				}

				docs := got.Docs
				if docs != nil {
					if !reflect.DeepEqual(docs[0].Type, tt.want.Docs[0].Type) {
						t.Errorf("Module.HookDBBatchIntent() = %v, want %v", docs[0].Type, tt.want.Docs[0].Type)
					}
					if !reflect.DeepEqual(docs[0].RuleName, tt.want.Docs[0].RuleName) {
						t.Errorf("Module.HookDBBatchIntent() = %v, want %v", docs[0].RuleName, tt.want.Docs[0].RuleName)
					}
					if !reflect.DeepEqual(docs[0].Timestamp, tt.want.Docs[0].Timestamp) {
						t.Errorf("Module.HookDBBatchIntent() = %v, want %v", docs[0].Timestamp, tt.want.Docs[0].Timestamp)
					}
					if !reflect.DeepEqual(docs[0].Payload, tt.want.Docs[0].Payload) {
						t.Errorf("Module.HookDBBatchIntent() = %v, want %v", docs[0].Payload, tt.want.Docs[0].Payload)
					}
					if !reflect.DeepEqual(docs[0].Status, tt.want.Docs[0].Status) {
						t.Errorf("Module.HookDBBatchIntent() = %v, want %v", docs[0].Status, tt.want.Docs[0].Status)
					}
				}
			}
			mockSyncman.AssertExpectations(t)
			mockSchema.AssertExpectations(t)
			mockCrud.AssertExpectations(t)
		})
	}
}
