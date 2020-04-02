package eventing

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

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
			m:    &Module{config: &config.Eventing{Rules: map[string]config.EventingRule{"rule": config.EventingRule{Type: "someType", Options: map[string]string{"col": "col", "db": "db"}}}}},
			args: args{token: 50, eventDocID: "eventid", batchID: "batchid", dbAlias: "db", col: "col", rows: nil},
			want: nil,
		},
		{
			name: "rows are nil",
			m:    &Module{config: &config.Eventing{Rules: map[string]config.EventingRule{"rule": config.EventingRule{Type: utils.EventDBCreate, Options: map[string]string{"col": "col", "db": "db"}}}}},
			args: args{token: 50, eventDocID: "eventid", batchID: "batchid", dbAlias: "db", col: "col", rows: nil},
			want: []*model.EventDocument{},
		},
		{
			name: "eventing is not possible",
			m:    &Module{config: &config.Eventing{Rules: map[string]config.EventingRule{"rule": config.EventingRule{Type: utils.EventDBCreate, Options: map[string]string{"col": "col", "db": "db"}}}}},
			args: args{token: 50, eventDocID: "eventid", batchID: "batchid", dbAlias: "db", col: "col", rows: []interface{}{map[string]interface{}{"key": "value"}}},
			schemaMockArgs: []mockArgs{
				mockArgs{
					method:         "CheckIfEventingIsPossible",
					args:           []interface{}{"db", "col", map[string]interface{}{"key": "value"}, false},
					paramsReturned: []interface{}{nil, false},
				},
			},
			want: nil,
		},
		{
			name: "doc are created",
			m:    &Module{config: &config.Eventing{Rules: map[string]config.EventingRule{"rule": config.EventingRule{Type: utils.EventDBCreate, Options: map[string]string{"col": "col", "db": "db"}}}}},
			args: args{token: 50, batchID: "batchid", dbAlias: "db", col: "col", rows: []interface{}{map[string]interface{}{"key": "value"}}},
			schemaMockArgs: []mockArgs{
				mockArgs{
					method:         "CheckIfEventingIsPossible",
					args:           []interface{}{"db", "col", map[string]interface{}{"key": "value"}, false},
					paramsReturned: []interface{}{map[string]interface{}{}, true},
				},
			},
			want: []*model.EventDocument{&model.EventDocument{BatchID: "batchid", Type: utils.EventDBCreate, RuleName: "rule", Token: 50, Timestamp: time.Now().Format(time.RFC3339), Payload: string(data), Status: utils.EventStatusIntent}},
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
			m:     &Module{config: &config.Eventing{Rules: map[string]config.EventingRule{"rule": config.EventingRule{Type: "nottype", Options: map[string]string{"col": "col", "db": "db"}}}}},
			args:  args{token: 50, eventType: "type", batchID: "batchid", dbAlias: "db", col: "col", find: map[string]interface{}{"key": "value"}},
			want:  nil,
			want1: false,
		},
		{
			name: "eventing is not possible",
			m:    &Module{config: &config.Eventing{Rules: map[string]config.EventingRule{"rule": config.EventingRule{Type: "type", Options: map[string]string{"col": "col", "db": "db"}}}}},
			args: args{token: 50, eventType: "type", batchID: "batchid", dbAlias: "db", col: "col", find: map[string]interface{}{"key": "value"}},
			schemaMockArgs: []mockArgs{
				mockArgs{
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
			m:    &Module{config: &config.Eventing{Rules: map[string]config.EventingRule{"rule": config.EventingRule{Type: "type", Options: map[string]string{"col": "col", "db": "db"}}}}},
			args: args{token: 50, eventType: "type", batchID: "batchid", dbAlias: "db", col: "col", find: map[string]interface{}{"key": "value"}},
			schemaMockArgs: []mockArgs{
				mockArgs{
					method:         "CheckIfEventingIsPossible",
					args:           []interface{}{"db", "col", map[string]interface{}{"key": "value"}, true},
					paramsReturned: []interface{}{map[string]interface{}{}, true},
				},
			},
			want:  []*model.EventDocument{&model.EventDocument{BatchID: "batchid", Type: "type", RuleName: "rule", Token: 50, Timestamp: time.Now().Format(time.RFC3339), Payload: string(data), Status: utils.EventStatusIntent}},
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
