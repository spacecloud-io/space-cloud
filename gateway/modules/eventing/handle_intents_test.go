package eventing

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func TestModule_processIntents(t *testing.T) {
	timeValue, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z07:00")
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
			name: "eventing is not enable",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype", Enabled: false}},
			args: args{t: &timeValue},
		},
		{
			name: "error reading",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype", Enabled: true}},
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
					args:           []interface{}{mock.Anything, "dbtype", utils.TableEventingLogs, &model.ReadRequest{Operation: utils.All, Find: map[string]interface{}{"status": utils.EventStatusIntent, "token": map[string]interface{}{"$gte": 1, "$lte": 100}}}},
					paramsReturned: []interface{}{[]interface{}{}, errors.New("some error")},
				},
			},
		},
		{
			name: "mapstructure decode error",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype", Enabled: true}},
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
					args:           []interface{}{mock.Anything, "dbtype", utils.TableEventingLogs, &model.ReadRequest{Operation: utils.All, Find: map[string]interface{}{"status": utils.EventStatusIntent, "token": map[string]interface{}{"$gte": 1, "$lte": 100}}}},
					paramsReturned: []interface{}{[]interface{}{"key"}, nil},
				},
			},
		},
		{
			name: "time parsing error",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype", Enabled: true}},
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
					args:           []interface{}{mock.Anything, "dbtype", utils.TableEventingLogs, &model.ReadRequest{Operation: utils.All, Find: map[string]interface{}{"status": utils.EventStatusIntent, "token": map[string]interface{}{"$gte": 1, "$lte": 100}}}},
					paramsReturned: []interface{}{[]interface{}{&model.EventDocument{EventTimestamp: time.Now().Format(time.RFC1123), ID: "id"}}, nil},
				},
			},
		},
		{
			name: "no error",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype", Enabled: true}},
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
					args:           []interface{}{mock.Anything, "dbtype", utils.TableEventingLogs, &model.ReadRequest{Operation: utils.All, Find: map[string]interface{}{"status": utils.EventStatusIntent, "token": map[string]interface{}{"$gte": 1, "$lte": 100}}}},
					paramsReturned: []interface{}{[]interface{}{&model.EventDocument{EventTimestamp: time.Now().Format(time.RFC3339), ID: "id"}}, nil},
				},
			},
		},
	}
	for _, tt := range tests {

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

		t.Run(tt.name, func(t *testing.T) {
			tt.m.processIntents(tt.args.t)
		})

		mockSyncman.AssertExpectations(t)
		mockCrud.AssertExpectations(t)
	}
}

func TestModule_processIntent(t *testing.T) {
	var res interface{}
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		eventDoc *model.EventDocument
	}
	tests := []struct {
		name              string
		m                 *Module
		args              args
		crudMockArgs      []mockArgs
		syncMockArgs      []mockArgs
		adminMockArgs     []mockArgs
		authMockArgs      []mockArgs
		filestoreMockArgs []mockArgs
	}{
		{
			name: "default case",
			m:    &Module{},
			args: args{&model.EventDocument{ID: "id", Type: "default"}},
		},
		{
			name: "db create case with error in read but no error in internal update",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype"}},
			args: args{&model.EventDocument{ID: "id", Type: utils.EventDBCreate, Token: 50, Payload: `{"db": "db", "col": "col", "doc": {}, "find": {}}`}},
			crudMockArgs: []mockArgs{
				{
					method:         "Read",
					args:           []interface{}{mock.Anything, "db", "col", &model.ReadRequest{Operation: utils.One, Find: map[string]interface{}{}}},
					paramsReturned: []interface{}{map[string]interface{}{}, errors.New("some error")},
				},
				{
					method:         "InternalUpdate",
					args:           []interface{}{mock.Anything, "dbtype", "abc", utils.TableEventingLogs, &model.UpdateRequest{Find: map[string]interface{}{"_id": "id"}, Operation: utils.All, Update: map[string]interface{}{"$set": map[string]interface{}{"status": utils.EventStatusCancelled}}}},
					paramsReturned: []interface{}{nil},
				},
			},
		},
		{
			name: "db create case with error in read and error in internal update",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype"}},
			args: args{&model.EventDocument{ID: "id", Type: utils.EventDBCreate, Token: 50, Payload: `{"db": "db", "col": "col", "doc": {}, "find": {}}`}},
			crudMockArgs: []mockArgs{
				{
					method:         "Read",
					args:           []interface{}{mock.Anything, "db", "col", &model.ReadRequest{Operation: utils.One, Find: map[string]interface{}{}}},
					paramsReturned: []interface{}{map[string]interface{}{}, errors.New("some error")},
				},
				{
					method:         "InternalUpdate",
					args:           []interface{}{mock.Anything, "dbtype", "abc", utils.TableEventingLogs, &model.UpdateRequest{Find: map[string]interface{}{"_id": "id"}, Operation: utils.All, Update: map[string]interface{}{"$set": map[string]interface{}{"status": utils.EventStatusCancelled}}}},
					paramsReturned: []interface{}{errors.New("some error")},
				},
			},
		},
		{
			name: "db create case with no error in read but error in internal update",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype"}},
			args: args{&model.EventDocument{ID: "id", Type: utils.EventDBCreate, Token: 50, Payload: `{"db": "db", "col": "col", "doc": {}, "find": {}}`}},
			crudMockArgs: []mockArgs{
				{
					method:         "Read",
					args:           []interface{}{mock.Anything, "db", "col", &model.ReadRequest{Operation: utils.One, Find: map[string]interface{}{}}},
					paramsReturned: []interface{}{map[string]interface{}{}, nil},
				},
				{
					method:         "InternalUpdate",
					args:           []interface{}{mock.Anything, "dbtype", "abc", utils.TableEventingLogs, &model.UpdateRequest{Find: map[string]interface{}{"_id": "id"}, Operation: utils.All, Update: map[string]interface{}{"$set": map[string]interface{}{"status": utils.EventStatusStaged}}}},
					paramsReturned: []interface{}{errors.New("some error")},
				},
			},
		},
		{
			name: "db create case with no errors",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype"}},
			args: args{&model.EventDocument{ID: "id", Type: utils.EventDBCreate, Token: 50, Payload: `{"db": "db", "col": "col", "doc": {}, "find": {}}`}},
			crudMockArgs: []mockArgs{
				{
					method:         "Read",
					args:           []interface{}{mock.Anything, "db", "col", &model.ReadRequest{Operation: utils.One, Find: map[string]interface{}{}}},
					paramsReturned: []interface{}{map[string]interface{}{}, nil},
				},
				{
					method:         "InternalUpdate",
					args:           []interface{}{mock.Anything, "dbtype", "abc", utils.TableEventingLogs, &model.UpdateRequest{Find: map[string]interface{}{"_id": "id"}, Operation: utils.All, Update: map[string]interface{}{"$set": map[string]interface{}{"status": utils.EventStatusStaged}}}},
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
					args:           []interface{}{mock.Anything, "POST", "url", mock.Anything, mock.Anything, []*model.EventDocument{{ID: "id", Type: "DB_INSERT", Token: 50, Payload: `{"db": "db", "col": "col", "doc": {}, "find": {}}`, Status: "staged"}}, &res},
					paramsReturned: []interface{}{nil},
				},
			},
			adminMockArgs: []mockArgs{
				{
					method:         "GetInternalAccessToken",
					paramsReturned: []interface{}{mock.Anything, nil},
				},
			},
			authMockArgs: []mockArgs{
				{
					method:         "GetSCAccessToken",
					paramsReturned: []interface{}{mock.Anything, nil},
				},
			},
		},
		{
			name: "db update case with error while reading",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype"}},
			args: args{&model.EventDocument{ID: "id", Type: utils.EventDBUpdate, Token: 50, Payload: `{"db": "db", "col": "col", "doc": {}, "find": {}}`}},
			crudMockArgs: []mockArgs{
				{
					method:         "Read",
					args:           []interface{}{mock.Anything, "db", "col", &model.ReadRequest{Operation: utils.One, Find: map[string]interface{}{}}},
					paramsReturned: []interface{}{map[string]interface{}{}, errors.New("some error")},
				},
			},
		},
		{
			name: "db update case with error while updating internal",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype"}},
			args: args{&model.EventDocument{ID: "id", Type: utils.EventDBUpdate, Token: 50, Payload: `{"db": "db", "col": "col", "doc": {}, "find": {}}`}},
			crudMockArgs: []mockArgs{
				{
					method:         "Read",
					args:           []interface{}{mock.Anything, "db", "col", &model.ReadRequest{Operation: utils.One, Find: map[string]interface{}{}}},
					paramsReturned: []interface{}{map[string]interface{}{}, nil},
				},
				{
					method:         "InternalUpdate",
					args:           []interface{}{mock.Anything, "dbtype", "abc", utils.TableEventingLogs, mock.Anything},
					paramsReturned: []interface{}{errors.New("some error")},
				},
			},
		},
		{
			name: "db update case with no errors",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype"}},
			args: args{&model.EventDocument{ID: "id", Type: utils.EventDBUpdate, Token: 50, Payload: `{"db": "db", "col": "col", "doc": {}, "find": {}}`}},
			crudMockArgs: []mockArgs{
				{
					method:         "Read",
					args:           []interface{}{mock.Anything, "db", "col", &model.ReadRequest{Operation: utils.One, Find: map[string]interface{}{}}},
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
					args:           []interface{}{mock.Anything, "POST", "url", mock.Anything, mock.Anything, []*model.EventDocument{{ID: "id", Type: "DB_UPDATE", Token: 50, Timestamp: time.Now().Format(time.RFC3339), Payload: "{\"db\":\"db\",\"col\":\"col\",\"doc\":{},\"find\":{}}", Status: "staged"}}, &res},
					paramsReturned: []interface{}{nil},
				},
			},
			adminMockArgs: []mockArgs{
				{
					method:         "GetInternalAccessToken",
					paramsReturned: []interface{}{mock.Anything, nil},
				},
			},
			authMockArgs: []mockArgs{
				{
					method:         "GetSCAccessToken",
					paramsReturned: []interface{}{mock.Anything, nil},
				},
			},
		},
		{
			name: "db delete case with no error while reading",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype"}},
			args: args{&model.EventDocument{ID: "id", Type: utils.EventDBDelete, Token: 50, Payload: `{"db": "db", "col": "col", "doc": {}, "find": {}}`}},
			crudMockArgs: []mockArgs{
				{
					method:         "Read",
					args:           []interface{}{mock.Anything, "db", "col", &model.ReadRequest{Operation: utils.One, Find: map[string]interface{}{}}},
					paramsReturned: []interface{}{map[string]interface{}{}, nil},
				},
				{
					method:         "InternalUpdate",
					args:           []interface{}{mock.Anything, "dbtype", "abc", utils.TableEventingLogs, mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
		},
		{
			name: "db delete case with error while reading and updating internal",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype"}},
			args: args{&model.EventDocument{ID: "id", Type: utils.EventDBDelete, Token: 50, Payload: `{"db": "db", "col": "col", "doc": {}, "find": {}}`}},
			crudMockArgs: []mockArgs{
				{
					method:         "Read",
					args:           []interface{}{mock.Anything, "db", "col", &model.ReadRequest{Operation: utils.One, Find: map[string]interface{}{}}},
					paramsReturned: []interface{}{map[string]interface{}{}, errors.New("some error")},
				},
				{
					method:         "InternalUpdate",
					args:           []interface{}{mock.Anything, "dbtype", "abc", utils.TableEventingLogs, mock.Anything},
					paramsReturned: []interface{}{errors.New("some error")},
				},
			},
		},
		{
			name: "db delete case with error while reading",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype"}},
			args: args{&model.EventDocument{ID: "id", Type: utils.EventDBDelete, Token: 50, Payload: `{"db": "db", "col": "col", "doc": {}, "find": {}}`}},
			crudMockArgs: []mockArgs{
				{
					method:         "Read",
					args:           []interface{}{mock.Anything, "db", "col", &model.ReadRequest{Operation: utils.One, Find: map[string]interface{}{}}},
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
					args:           []interface{}{mock.Anything, "POST", "url", mock.Anything, mock.Anything, []*model.EventDocument{{ID: "id", Type: "DB_DELETE", Token: 50, Payload: "{\"db\": \"db\", \"col\": \"col\", \"doc\": {}, \"find\": {}}", Status: "staged"}}, &res},
					paramsReturned: []interface{}{nil},
				},
			},
			adminMockArgs: []mockArgs{
				{
					method:         "GetInternalAccessToken",
					paramsReturned: []interface{}{mock.Anything, nil},
				},
			},
			authMockArgs: []mockArgs{
				{
					method:         "GetSCAccessToken",
					paramsReturned: []interface{}{mock.Anything, nil},
				},
			},
		},
		{
			name: "file create case with error while getting internal access token",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype"}},
			args: args{&model.EventDocument{ID: "id", Type: utils.EventFileCreate, Token: 50, Payload: `{"meta": {"key": "value"}, "path": "path"}`}},
			authMockArgs: []mockArgs{
				{
					method:         "GetInternalAccessToken",
					paramsReturned: []interface{}{"", errors.New("some error")},
				},
			},
		},
		{
			name: "file create case with error while DoesExists and not while updating internal",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype"}},
			args: args{&model.EventDocument{ID: "id", Type: utils.EventFileCreate, Token: 50, Payload: `{"meta": {"key": "value"}, "path": "path"}`}},
			authMockArgs: []mockArgs{
				{
					method:         "GetInternalAccessToken",
					paramsReturned: []interface{}{"token", nil},
				},
			},
			filestoreMockArgs: []mockArgs{
				{
					method:         "DoesExists",
					args:           []interface{}{mock.Anything, "abc", "token", "path"},
					paramsReturned: []interface{}{errors.New("some error")},
				},
			},
			crudMockArgs: []mockArgs{
				{
					method:         "InternalUpdate",
					args:           []interface{}{mock.Anything, "dbtype", "abc", utils.TableEventingLogs, mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
		},
		{
			name: "file create case with error while DoesExists and while updating internal",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype"}},
			args: args{&model.EventDocument{ID: "id", Type: utils.EventFileCreate, Token: 50, Payload: `{"meta": {"key": "value"}, "path": "path"}`}},
			authMockArgs: []mockArgs{
				{
					method:         "GetInternalAccessToken",
					paramsReturned: []interface{}{"token", nil},
				},
			},
			filestoreMockArgs: []mockArgs{
				{
					method:         "DoesExists",
					args:           []interface{}{mock.Anything, "abc", "token", "path"},
					paramsReturned: []interface{}{errors.New("some error")},
				},
			},
			crudMockArgs: []mockArgs{
				{
					method:         "InternalUpdate",
					args:           []interface{}{mock.Anything, "dbtype", "abc", utils.TableEventingLogs, mock.Anything},
					paramsReturned: []interface{}{errors.New("some error")},
				},
			},
		},
		{
			name: "file create case with error while updating internal",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype"}},
			args: args{&model.EventDocument{ID: "id", Type: utils.EventFileCreate, Token: 50, Payload: `{"meta": {"key": "value"}, "path": "path"}`}},
			authMockArgs: []mockArgs{
				{
					method:         "GetInternalAccessToken",
					paramsReturned: []interface{}{"token", nil},
				},
			},
			filestoreMockArgs: []mockArgs{
				{
					method:         "DoesExists",
					args:           []interface{}{mock.Anything, "abc", "token", "path"},
					paramsReturned: []interface{}{nil},
				},
			},
			crudMockArgs: []mockArgs{
				{
					method:         "InternalUpdate",
					args:           []interface{}{mock.Anything, "dbtype", "abc", utils.TableEventingLogs, mock.Anything},
					paramsReturned: []interface{}{errors.New("some error")},
				},
			},
		},
		{
			name: "file create case with no errors",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype"}},
			args: args{&model.EventDocument{ID: "id", Type: utils.EventFileCreate, Token: 50, Payload: `{"meta": {"key": "value"}, "path": "path"}`}},
			authMockArgs: []mockArgs{
				{
					method:         "GetInternalAccessToken",
					paramsReturned: []interface{}{"token", nil},
				},
				{
					method:         "GetSCAccessToken",
					paramsReturned: []interface{}{mock.Anything, nil},
				},
			},
			filestoreMockArgs: []mockArgs{
				{
					method:         "DoesExists",
					args:           []interface{}{mock.Anything, "abc", "token", "path"},
					paramsReturned: []interface{}{nil},
				},
			},
			crudMockArgs: []mockArgs{
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
					args:           []interface{}{mock.Anything, "POST", "url", mock.Anything, mock.Anything, []*model.EventDocument{{ID: "id", Type: "FILE_CREATE", Token: 50, Payload: "{\"meta\": {\"key\": \"value\"}, \"path\": \"path\"}", Status: "staged"}}, &res},
					paramsReturned: []interface{}{nil},
				},
			},
			adminMockArgs: []mockArgs{
				{
					method:         "GetInternalAccessToken",
					paramsReturned: []interface{}{mock.Anything, nil},
				},
			},
		},
		{
			name: "file delete case with error while getting internal access token",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype"}},
			args: args{&model.EventDocument{ID: "id", Type: utils.EventFileDelete, Token: 50, Payload: `{"meta": {"key": "value"}, "path": "path"}`}},
			authMockArgs: []mockArgs{
				{
					method:         "GetInternalAccessToken",
					paramsReturned: []interface{}{"", errors.New("some error")},
				},
			},
		},
		{
			name: "file delete case with no error while DoesExists but while updating internal",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype"}},
			args: args{&model.EventDocument{ID: "id", Type: utils.EventFileDelete, Token: 50, Payload: `{"meta": {"key": "value"}, "path": "path"}`}},
			authMockArgs: []mockArgs{
				{
					method:         "GetInternalAccessToken",
					paramsReturned: []interface{}{"token", nil},
				},
			},
			filestoreMockArgs: []mockArgs{
				{
					method:         "DoesExists",
					args:           []interface{}{mock.Anything, "abc", "token", "path"},
					paramsReturned: []interface{}{nil},
				},
			},
			crudMockArgs: []mockArgs{
				{
					method:         "InternalUpdate",
					args:           []interface{}{mock.Anything, "dbtype", "abc", utils.TableEventingLogs, mock.Anything},
					paramsReturned: []interface{}{errors.New("some error")},
				},
			},
		},
		{
			name: "file delete case with error while DoesExists and also while updating internal",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype"}},
			args: args{&model.EventDocument{ID: "id", Type: utils.EventFileDelete, Token: 50, Payload: `{"meta": {"key": "value"}, "path": "path"}`}},
			authMockArgs: []mockArgs{
				{
					method:         "GetInternalAccessToken",
					paramsReturned: []interface{}{"token", nil},
				},
			},
			filestoreMockArgs: []mockArgs{
				{
					method:         "DoesExists",
					args:           []interface{}{mock.Anything, "abc", "token", "path"},
					paramsReturned: []interface{}{errors.New("some error")},
				},
			},
			crudMockArgs: []mockArgs{
				{
					method:         "InternalUpdate",
					args:           []interface{}{mock.Anything, "dbtype", "abc", utils.TableEventingLogs, mock.Anything},
					paramsReturned: []interface{}{errors.New("some error")},
				},
			},
		},
		{
			name: "file delete case with error while DoesExists but not while updating internal",
			m:    &Module{project: "abc", config: &config.Eventing{DBAlias: "dbtype"}},
			args: args{&model.EventDocument{ID: "id", Type: utils.EventFileDelete, Token: 50, Payload: `{"meta": {"key": "value"}, "path": "path"}`}},
			authMockArgs: []mockArgs{
				{
					method:         "GetInternalAccessToken",
					paramsReturned: []interface{}{"token", nil},
				},
				{
					method:         "GetSCAccessToken",
					paramsReturned: []interface{}{mock.Anything, nil},
				},
			},
			filestoreMockArgs: []mockArgs{
				{
					method:         "DoesExists",
					args:           []interface{}{mock.Anything, "abc", "token", "path"},
					paramsReturned: []interface{}{errors.New("some error")},
				},
			},
			crudMockArgs: []mockArgs{
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
					args:           []interface{}{mock.Anything, "POST", "url", mock.Anything, mock.Anything, []*model.EventDocument{{ID: "id", Type: "FILE_DELETE", Token: 50, Payload: "{\"meta\": {\"key\": \"value\"}, \"path\": \"path\"}", Status: "staged"}}, &res},
					paramsReturned: []interface{}{nil},
				},
			},
			adminMockArgs: []mockArgs{
				{
					method:         "GetInternalAccessToken",
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
			mockFileStore := mockFileStoreEventingInterface{}

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
			for _, m := range tt.filestoreMockArgs {
				mockFileStore.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			tt.m.crud = &mockCrud
			tt.m.syncMan = &mockSyncman
			tt.m.adminMan = &mockAdmin
			tt.m.auth = &mockAuth
			tt.m.fileStore = &mockFileStore

			tt.m.processIntent(tt.args.eventDoc)

			mockCrud.AssertExpectations(t)
			mockSyncman.AssertExpectations(t)
			mockAdmin.AssertExpectations(t)
			mockAuth.AssertExpectations(t)
			mockFileStore.AssertExpectations(t)
		})
	}
}
