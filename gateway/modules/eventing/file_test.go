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

func TestModule_CreateFileIntentHook(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx context.Context
		req *model.CreateFileRequest
	}
	tests := []struct {
		name            string
		m               *Module
		args            args
		syncmanMockArgs []mockArgs
		crudMockArgs    []mockArgs
		want            *model.EventIntent
		wantErr         bool
	}{
		{
			name: "eventing is not enabled",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype", Enabled: false}},
			args: args{ctx: context.Background(), req: &model.CreateFileRequest{Meta: map[string]interface{}{}, Path: "path"}},
			want: &model.EventIntent{Invalid: true},
		},
		{
			name: "no rules match",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype", Enabled: true, Rules: map[string]*config.EventingTrigger{"rule": {Type: "not file create"}}}},
			args: args{ctx: context.Background(), req: &model.CreateFileRequest{Meta: map[string]interface{}{}, Path: "path"}},

			want: &model.EventIntent{Invalid: true},
		},
		{
			name: "error creating internal",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype", Enabled: true, Rules: map[string]*config.EventingTrigger{"rule": {Type: utils.EventFileCreate}}}},
			args: args{ctx: context.Background(), req: &model.CreateFileRequest{Meta: map[string]interface{}{}, Path: "path"}},
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
			name: "file intent request handled",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype", Enabled: true, Rules: map[string]*config.EventingTrigger{"rule": {Type: utils.EventFileCreate}}}},
			args: args{ctx: context.Background(), req: &model.CreateFileRequest{Meta: map[string]interface{}{}, Path: "path"}},
			crudMockArgs: []mockArgs{
				{
					method:         "InternalCreate",
					args:           []interface{}{mock.Anything, "dbtype", "abc", utils.TableEventingLogs, mock.Anything, false},
					paramsReturned: []interface{}{nil},
				},
			},
			want: &model.EventIntent{Docs: []*model.EventDocument{{Type: utils.EventFileCreate, RuleName: "rule", Timestamp: time.Now().Format(time.RFC3339Nano), Payload: `{"meta":{},"path":"path"}`, Status: "intent"}}},
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

			got, err := tt.m.CreateFileIntentHook(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Module.CreateFileIntentHook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && tt.want != nil {
				invalid := got.Invalid

				if !reflect.DeepEqual(invalid, tt.want.Invalid) {
					t.Errorf("Module.CreateFileIntentHook() = %v, want %v", invalid, tt.want.Invalid)
				}

				docs := got.Docs
				if docs != nil {
					if !reflect.DeepEqual(docs[0].Type, tt.want.Docs[0].Type) {
						t.Errorf("Module.CreateFileIntentHook() = %v, want %v", docs[0].Type, tt.want.Docs[0].Type)
					}
					if !reflect.DeepEqual(docs[0].RuleName, tt.want.Docs[0].RuleName) {
						t.Errorf("Module.CreateFileIntentHook() = %v, want %v", docs[0].RuleName, tt.want.Docs[0].RuleName)
					}
					t1, _ := time.Parse(time.RFC3339, docs[0].Timestamp)
					t1 = t1.Truncate(time.Second)
					t2, _ := time.Parse(time.RFC3339, tt.want.Docs[0].Timestamp)
					t2 = t2.Truncate(time.Second)
					if !t1.Equal(t2) {
						t.Errorf("Module.DeleteFileIntentHook() = %v, want %v", docs[0].Timestamp, tt.want.Docs[0].Timestamp)
					}
					if !reflect.DeepEqual(docs[0].Payload, tt.want.Docs[0].Payload) {
						t.Errorf("Module.CreateFileIntentHook() = %v, want %v", docs[0].Payload, tt.want.Docs[0].Payload)
					}
					if !reflect.DeepEqual(docs[0].Status, tt.want.Docs[0].Status) {
						t.Errorf("Module.CreateFileIntentHook() = %v, want %v", docs[0].Status, tt.want.Docs[0].Status)
					}
				}
			}

			mockSyncman.AssertExpectations(t)
			mockCrud.AssertExpectations(t)
		})
	}
}

func TestModule_DeleteFileIntentHook(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx  context.Context
		path string
		meta map[string]interface{}
	}
	tests := []struct {
		name            string
		m               *Module
		args            args
		syncmanMockArgs []mockArgs
		crudMockArgs    []mockArgs
		want            *model.EventIntent
		wantErr         bool
	}{
		{
			name: "eventing is not enabled",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype", Enabled: false}},
			args: args{ctx: context.Background(), meta: map[string]interface{}{}, path: "path"},
			want: &model.EventIntent{Invalid: true},
		},
		{
			name: "no rules match",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype", Enabled: true, Rules: map[string]*config.EventingTrigger{"rule": {Type: "not file delete"}}}},
			args: args{ctx: context.Background(), meta: map[string]interface{}{}, path: "path"},

			want: &model.EventIntent{Invalid: true},
		},
		{
			name: "error creating internal",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype", Enabled: true, Rules: map[string]*config.EventingTrigger{"rule": {Type: utils.EventFileDelete}}}},
			args: args{ctx: context.Background(), meta: map[string]interface{}{}, path: "path"},
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
			name: "file intent request handled",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype", Enabled: true, Rules: map[string]*config.EventingTrigger{"rule": {Type: utils.EventFileDelete}}}},
			args: args{ctx: context.Background(), meta: map[string]interface{}{}, path: "path"},
			crudMockArgs: []mockArgs{
				{
					method:         "InternalCreate",
					args:           []interface{}{mock.Anything, "dbtype", "abc", utils.TableEventingLogs, mock.Anything, false},
					paramsReturned: []interface{}{nil},
				},
			},
			want: &model.EventIntent{Docs: []*model.EventDocument{{Type: utils.EventFileDelete, RuleName: "rule", Timestamp: time.Now().Format(time.RFC3339Nano), Payload: `{"meta":{},"path":"path"}`, Status: "intent"}}},
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

			got, err := tt.m.DeleteFileIntentHook(context.Background(), tt.args.path, tt.args.meta)
			if (err != nil) != tt.wantErr {
				t.Errorf("Module.DeleteFileIntentHook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && tt.want != nil {
				invalid := got.Invalid

				if !reflect.DeepEqual(invalid, tt.want.Invalid) {
					t.Errorf("Module.DeleteFileIntentHook() = %v, want %v", invalid, tt.want.Invalid)
				}

				docs := got.Docs
				if docs != nil {
					if !reflect.DeepEqual(docs[0].Type, tt.want.Docs[0].Type) {
						t.Errorf("Module.DeleteFileIntentHook() = %v, want %v", docs[0].Type, tt.want.Docs[0].Type)
					}
					if !reflect.DeepEqual(docs[0].RuleName, tt.want.Docs[0].RuleName) {
						t.Errorf("Module.DeleteFileIntentHook() = %v, want %v", docs[0].RuleName, tt.want.Docs[0].RuleName)
					}
					t1, _ := time.Parse(time.RFC3339, docs[0].Timestamp)
					t1 = t1.Truncate(time.Second)
					t2, _ := time.Parse(time.RFC3339, tt.want.Docs[0].Timestamp)
					t2 = t2.Truncate(time.Second)
					if !t1.Equal(t2) {
						t.Errorf("Module.DeleteFileIntentHook() = %v, want %v", docs[0].Timestamp, tt.want.Docs[0].Timestamp)
					}
					if !reflect.DeepEqual(docs[0].Payload, tt.want.Docs[0].Payload) {
						t.Errorf("Module.DeleteFileIntentHook() = %v, want %v", docs[0].Payload, tt.want.Docs[0].Payload)
					}
					if !reflect.DeepEqual(docs[0].Status, tt.want.Docs[0].Status) {
						t.Errorf("Module.DeleteFileIntentHook() = %v, want %v", docs[0].Status, tt.want.Docs[0].Status)
					}
				}
			}

			mockSyncman.AssertExpectations(t)
			mockCrud.AssertExpectations(t)
		})
	}
}

func TestModule_HookStage(t *testing.T) {
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
			name: "event is staged",
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
					method:         "GetAssignedSpaceCloudID",
					args:           []interface{}{mock.Anything, "abc", 50},
					paramsReturned: []interface{}{"url", nil},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockCrud := mockCrudInterface{}
			mockAuth := mockAuthEventingInterface{}
			mockSyncman := mockSyncmanEventingInterface{}

			for _, m := range tt.crudMockArgs {
				mockCrud.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.syncMockArgs {
				mockSyncman.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.authMockArgs {
				mockAuth.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			tt.m.crud = &mockCrud
			tt.m.syncMan = &mockSyncman
			tt.m.auth = &mockAuth

			tt.m.HookStage(context.Background(), tt.args.intent, tt.args.err)

			mockCrud.AssertExpectations(t)
			mockSyncman.AssertExpectations(t)
			mockAuth.AssertExpectations(t)
		})
	}
}
