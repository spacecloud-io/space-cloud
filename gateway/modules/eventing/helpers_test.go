package eventing

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/auth"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func TestModule_selectRule(t *testing.T) {
	type args struct {
		name   string
		evType string
	}
	tests := []struct {
		name    string
		m       *Module
		args    args
		want    *config.EventingTrigger
		wantErr bool
	}{
		{
			name: "event type is an internal type",
			m:    &Module{config: &config.Eventing{InternalRules: map[string]*config.EventingTrigger{"some-rule": {Type: "DB_INSERT", URL: "abc"}}}},
			args: args{name: "some-rule", evType: "DB_INSERT"},
			want: &config.EventingTrigger{Type: "DB_INSERT", URL: "abc"},
		},
		{
			name: "event type is found in rules",
			m:    &Module{config: &config.Eventing{Rules: map[string]*config.EventingTrigger{"some-rule": {Type: "event"}}}},
			args: args{name: "some-rule", evType: "event"},
			want: &config.EventingTrigger{Type: "event"},
		},
		{
			name: "event type is found in internal rules",
			m:    &Module{config: &config.Eventing{InternalRules: map[string]*config.EventingTrigger{"some-rule": {Type: "event"}}}},
			args: args{name: "some-rule", evType: "event"},
			want: &config.EventingTrigger{Type: "event"},
		},
		{
			name:    "event type is not found",
			m:       &Module{config: &config.Eventing{}},
			args:    args{name: "some-rule", evType: "event"},
			want:    &config.EventingTrigger{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.selectRule(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Module.selectRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Module.selectRule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModule_validate(t *testing.T) {
	authModule := auth.Init("chicago", "1", &crud.Module{}, nil)
	err := authModule.SetConfig(context.TODO(), "local", &config.ProjectConfig{ID: "project", Secrets: []*config.Secret{{IsPrimary: true, Secret: "mySecretkey"}}}, config.DatabaseRules{}, config.DatabasePreparedQueries{}, config.FileStoreRules{}, config.Services{}, config.EventingRules{config.GenerateResourceID("chicago", "project", config.ResourceEventingRule, "event"): &config.Rule{ID: "event", Rule: "authenticated"}}, config.SecurityFunctions{})
	if err != nil {
		t.Fatalf("error setting config (%s)", err.Error())
	}
	type args struct {
		ctx     context.Context
		project string
		token   string
		event   *model.QueueEventRequest
	}
	tests := []struct {
		name    string
		m       *Module
		args    args
		wantErr bool
	}{
		{
			name:    "event type is an internal type",
			m:       &Module{config: &config.Eventing{Rules: map[string]*config.EventingTrigger{"some-rule": {Type: "DB_INSERT"}}}},
			args:    args{event: &model.QueueEventRequest{Type: "DB_INSERT", Delay: 0, Payload: "something", Options: make(map[string]string)}},
			wantErr: true,
		},
		{
			name:    "invalid project details",
			m:       &Module{auth: &auth.Module{}},
			args:    args{ctx: context.Background(), project: "some-project", event: &model.QueueEventRequest{Type: "event", Delay: 0, Payload: "something", Options: make(map[string]string)}},
			wantErr: true,
		},
		{
			name:    "invalid token",
			m:       &Module{auth: &auth.Module{}},
			args:    args{ctx: context.Background(), token: "token", event: &model.QueueEventRequest{Type: "event", Delay: 0, Payload: "something", Options: make(map[string]string)}},
			wantErr: true,
		},
		{
			name: "event type not in schemas",
			m: &Module{
				auth: authModule,
				config: &config.Eventing{
					SecurityRules: map[string]*config.Rule{"event": {Rule: "authenticated"}},
					Schemas:       map[string]config.SchemaObject{"event": {Schema: "some-schema"}}}},
			args: args{ctx: context.Background(), project: "project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc", event: &model.QueueEventRequest{Type: "event", Delay: 0, Payload: "some-schema", Options: make(map[string]string)}},
		},
		{
			name: "no schema given",
			m: &Module{
				schemas: map[string]model.Fields{"event": {}},
				auth:    authModule,
				config: &config.Eventing{
					SecurityRules: map[string]*config.Rule{
						"event": {
							Rule: "authenticated",
						}},
					Schemas: map[string]config.SchemaObject{"event": {Schema: "type event {id: ID! title: String}"}}}},
			args: args{ctx: context.Background(), project: "project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc", event: &model.QueueEventRequest{Type: "event", Delay: 0, Payload: make(map[string]interface{}), Options: make(map[string]string)}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.auth = authModule

			if err := tt.m.validate(context.Background(), tt.args.project, tt.args.token, tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("Module.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_isOptionsValid(t *testing.T) {
	type args struct {
		ruleOptions     map[string]string
		providedOptions map[string]string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "key present in providedOptions",
			args: args{ruleOptions: map[string]string{"key1": "value1", "key2": "value2"}, providedOptions: map[string]string{"key1": "value1", "key2": "value2"}},
			want: true,
		},
		{
			name: "key not present in providedOptions",
			args: args{ruleOptions: map[string]string{"key1": "value1", "key2": "value2"}, providedOptions: map[string]string{"key": "value1", "key2": "value2"}},
			want: false,
		},
		{
			name: "value not present in providedOptions",
			args: args{ruleOptions: map[string]string{"key1": "value1", "key2": "value2"}, providedOptions: map[string]string{"key1": "value", "key2": "value2"}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isOptionsValid(tt.args.ruleOptions, tt.args.providedOptions); got != tt.want {
				t.Errorf("isOptionsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertToArray(t *testing.T) {
	type args struct {
		eventDocs []*model.EventDocument
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{
			name: "conversion takes place",
			args: args{eventDocs: []*model.EventDocument{{ID: "ID", BatchID: "BatchID", Type: "DB_INSERT", RuleName: "encrypt", Token: 123, Payload: "payload", Status: "ok", Remark: "Remark", Timestamp: "", EventTimestamp: ""}}},
			want: []interface{}{0: map[string]interface{}{"_id": "ID", "batchid": "BatchID", "payload": "payload", "remark": "Remark", "rule_name": "encrypt", "status": "ok", "token": 123, "type": "DB_INSERT", "event_ts": "", "ts": ""}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertToArray(tt.args.eventDocs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertToArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModule_getMatchingRules(t *testing.T) {
	type args struct {
		req *model.QueueEventRequest
	}
	tests := []struct {
		name string
		m    *Module
		args args
		want []*config.EventingTrigger
	}{
		{
			name: "rule type is not equal to name",
			m: &Module{config: &config.Eventing{
				Rules:         map[string]*config.EventingTrigger{"some-rule": {Type: "rule", Options: map[string]string{"option": "value"}}},
				InternalRules: map[string]*config.EventingTrigger{"some-internal-rule": {Type: "internalRule", Options: map[string]string{"option": "value"}}}}},
			args: args{req: &model.QueueEventRequest{Type: "name", Options: map[string]string{"option": "value"}}},
			want: []*config.EventingTrigger{},
		},
		{
			name: "rule options are not valid",
			m: &Module{config: &config.Eventing{
				Rules:         map[string]*config.EventingTrigger{"some-rule": {Type: "rule", Options: map[string]string{"option": "value"}}},
				InternalRules: map[string]*config.EventingTrigger{"some-internal-rule": {Type: "internalRule", Options: map[string]string{"option": "value"}}}}},
			args: args{req: &model.QueueEventRequest{Type: "rule", Options: map[string]string{"wrong-option": "value"}}},
			want: []*config.EventingTrigger{},
		},
		{
			name: "rule matching in Rules",
			m: &Module{config: &config.Eventing{
				Rules:         map[string]*config.EventingTrigger{"some-rule": {Type: "rule", Options: map[string]string{"option": "value"}}},
				InternalRules: map[string]*config.EventingTrigger{"some-internal-rule": {Type: "internalRule", Options: map[string]string{"option": "value"}}}}},
			args: args{req: &model.QueueEventRequest{Type: "rule", Options: map[string]string{"option": "value"}}},
			want: []*config.EventingTrigger{{Type: "rule", TriggerType: "external", Retries: 0, Timeout: 0, Options: map[string]string{"option": "value"}}},
		},
		{
			name: "rule matching in InternalRules",
			m: &Module{config: &config.Eventing{
				Rules:         map[string]*config.EventingTrigger{"some-rule": {Type: "rule", Options: map[string]string{"option": "value"}}},
				InternalRules: map[string]*config.EventingTrigger{"some-internal-rule": {Type: "internalRule", Options: map[string]string{"option": "value"}}}}},
			args: args{req: &model.QueueEventRequest{Type: "internalRule", Options: map[string]string{"option": "value"}}},
			want: []*config.EventingTrigger{{Type: "internalRule", TriggerType: "internal", Retries: 0, Timeout: 0, Options: map[string]string{"option": "value"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.getMatchingRules(context.Background(), tt.args.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Module.getMatchingRules() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getCreateRows(t *testing.T) {
	type args struct {
		doc interface{}
		op  string
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{
			name: "op is utils one",
			args: args{doc: []interface{}{}, op: utils.One},
			want: []interface{}{[]interface{}{}},
		},
		{
			name: "op is utils all",
			args: args{doc: []interface{}{}, op: utils.All},
			want: []interface{}{},
		},
		{
			name: "default case",
			args: args{doc: []interface{}{}, op: "notOneorAll"},
			want: []interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCreateRows(tt.args.doc, tt.args.op); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCreateRows() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModule_generateProcessedEventRequest(t *testing.T) {
	type args struct {
		eventID string
	}
	tests := []struct {
		name string
		m    *Module
		args args
		want *model.UpdateRequest
	}{
		{
			name: "update request returned",
			m:    &Module{},
			args: args{eventID: "eventID"},
			want: &model.UpdateRequest{
				Find:      map[string]interface{}{"_id": "eventID"},
				Operation: utils.All,
				Update: map[string]interface{}{
					"$set": map[string]interface{}{"status": utils.EventStatusProcessed},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.generateProcessedEventRequest(tt.args.eventID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Module.generateProcessedEventRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModule_generateFailedEventRequest(t *testing.T) {
	type args struct {
		eventID string
		remark  string
	}
	tests := []struct {
		name string
		m    *Module
		args args
		want *model.UpdateRequest
	}{
		{
			name: "update request returned",
			m:    &Module{},
			args: args{eventID: "eventID", remark: "remark"},
			want: &model.UpdateRequest{
				Find:      map[string]interface{}{"_id": "eventID"},
				Operation: utils.All,
				Update: map[string]interface{}{
					"$set": map[string]interface{}{"status": utils.EventStatusFailed, "remark": "remark"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.generateFailedEventRequest(tt.args.eventID, tt.args.remark); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Module.generateFailedEventRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModule_generateStageEventRequest(t *testing.T) {
	type args struct {
		eventID string
	}
	tests := []struct {
		name string
		m    *Module
		args args
		want *model.UpdateRequest
	}{
		{
			name: "update request returned",
			m:    &Module{},
			args: args{eventID: "eventID"},
			want: &model.UpdateRequest{
				Find:      map[string]interface{}{"_id": "eventID"},
				Operation: utils.All,
				Update: map[string]interface{}{
					"$set": map[string]interface{}{"status": utils.EventStatusStaged},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.generateStageEventRequest(tt.args.eventID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Module.generateStageEventRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModule_generateCancelEventRequest(t *testing.T) {
	type args struct {
		eventID string
	}
	tests := []struct {
		name string
		m    *Module
		args args
		want *model.UpdateRequest
	}{
		{
			name: "update request returned",
			m:    &Module{},
			args: args{eventID: "eventID"},
			want: &model.UpdateRequest{
				Find:      map[string]interface{}{"_id": "eventID"},
				Operation: utils.All,
				Update: map[string]interface{}{
					"$set": map[string]interface{}{"status": utils.EventStatusCancelled},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.generateCancelEventRequest(tt.args.eventID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Module.generateCancelEventRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModule_getSpaceCloudIDFromBatchID(t *testing.T) {
	type args struct {
		batchID string
	}
	tests := []struct {
		name string
		m    *Module
		args args
		want string
	}{
		{
			name: "got spaceCloud id",
			m:    &Module{},
			args: args{batchID: "some--id"},
			want: "id",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.getSpaceCloudIDFromBatchID(tt.args.batchID); got != tt.want {
				t.Errorf("Module.getSpaceCloudIDFromBatchID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModule_generateBatchID(t *testing.T) {
	tests := []struct {
		name string
		m    *Module
		want string
	}{
		{
			name: "ID is generated",
			m:    &Module{nodeID: "abc"},
			want: "abc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.generateBatchID()
			got = strings.Split(got, "--")[1]
			if got != tt.want {
				t.Errorf("Module.generateBatchID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModule_transmitEvents(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		eventToken int
		eventDocs  []*model.EventDocument
	}
	tests := []struct {
		name          string
		m             *Module
		args          args
		syncMockArgs  []mockArgs
		adminMockArgs []mockArgs
		authMockArgs  []mockArgs
	}{
		{
			name: "error getting url",
			m:    &Module{project: "abc"},
			args: args{eventToken: 50, eventDocs: []*model.EventDocument{{ID: "ID", BatchID: "BatchID", Type: "DB_INSERT", RuleName: "encrypt", Token: 50, Payload: "payload", Status: "ok", Remark: "Remark"}}},
			syncMockArgs: []mockArgs{
				{
					method:         "GetAssignedSpaceCloudID",
					args:           []interface{}{mock.Anything, "abc", 50},
					paramsReturned: []interface{}{"", errors.New("some error")},
				},
			},
		},
		{
			name: "error getting token",
			m:    &Module{project: "abc"},
			args: args{eventToken: 50, eventDocs: []*model.EventDocument{{ID: "ID", BatchID: "BatchID", Type: "DB_INSERT", RuleName: "encrypt", Token: 50, Payload: "payload", Status: "ok", Remark: "Remark"}}},
			syncMockArgs: []mockArgs{
				{
					method:         "GetAssignedSpaceCloudID",
					args:           []interface{}{mock.Anything, "abc", 50},
					paramsReturned: []interface{}{"url", nil},
				},
			},
		},
		{
			name: "error getting scToken",
			m:    &Module{project: "abc"},
			args: args{eventToken: 50, eventDocs: []*model.EventDocument{{ID: "ID", BatchID: "BatchID", Type: "DB_INSERT", RuleName: "encrypt", Token: 50, Payload: "payload", Status: "ok", Remark: "Remark"}}},
			syncMockArgs: []mockArgs{
				{
					method:         "GetAssignedSpaceCloudID",
					args:           []interface{}{mock.Anything, "abc", 50},
					paramsReturned: []interface{}{"url", nil},
				},
			},
		},
		{
			name: "error making http request",
			m:    &Module{project: "abc"},
			args: args{eventToken: 50, eventDocs: []*model.EventDocument{{ID: "ID", BatchID: "BatchID", Type: "DB_INSERT", RuleName: "encrypt", Token: 50, Payload: "payload", Status: "ok", Remark: "Remark"}}},
			syncMockArgs: []mockArgs{
				{
					method:         "GetAssignedSpaceCloudID",
					args:           []interface{}{mock.Anything, "abc", 50},
					paramsReturned: []interface{}{"url", nil},
				},
			},
		},
		{
			name: "event is transmitted",
			m:    &Module{project: "abc"},
			args: args{eventToken: 50, eventDocs: []*model.EventDocument{{ID: "ID", BatchID: "BatchID", Type: "DB_INSERT", RuleName: "encrypt", Token: 50, Payload: "payload", Status: "ok", Remark: "Remark"}}},
			syncMockArgs: []mockArgs{
				{
					method:         "GetAssignedSpaceCloudID",
					args:           []interface{}{mock.Anything, "abc", 50},
					paramsReturned: []interface{}{"url", nil},
				},
			},
			authMockArgs: []mockArgs{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockAuth := mockAuthEventingInterface{}
			mockSyncman := mockSyncmanEventingInterface{}

			for _, m := range tt.syncMockArgs {
				mockSyncman.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.authMockArgs {
				mockAuth.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			tt.m.syncMan = &mockSyncman
			tt.m.auth = &mockAuth

			tt.m.transmitEvents(tt.args.eventToken, tt.args.eventDocs)

			mockSyncman.AssertExpectations(t)
			mockAuth.AssertExpectations(t)
		})
	}
}

func TestModule_generateQueueEventRequestRaw(t *testing.T) {
	type args struct {
		token      int
		name       string
		eventDocID string
		batchID    string
		status     string
		event      *model.QueueEventRequest
	}
	tests := []struct {
		name string
		m    *Module
		args args
		want *model.EventDocument
	}{
		{
			name: "error parsing timestamp since it was not provided",
			m:    &Module{},
			args: args{token: 50, name: "rule", eventDocID: "eventDocID", batchID: "batchID", status: "ok", event: &model.QueueEventRequest{Delay: int64(0), IsSynchronous: false, Options: make(map[string]string), Payload: "payload", Type: "DB_INSERT"}},
			want: &model.EventDocument{ID: "eventDocID", BatchID: "batchID", Type: "DB_INSERT", RuleName: "rule", Token: 50, Timestamp: time.Now().Format(time.RFC3339Nano), Payload: "\"payload\"", Status: "ok"},
		},
		{
			name: "error parsing timestamp since it was provided",
			m:    &Module{},
			args: args{token: 50, name: "rule", eventDocID: "eventDocID", batchID: "batchID", status: "ok", event: &model.QueueEventRequest{Timestamp: "incorrectTimestamp", Delay: int64(0), IsSynchronous: false, Options: make(map[string]string), Payload: "payload", Type: "DB_INSERT"}},
			want: &model.EventDocument{ID: "eventDocID", BatchID: "batchID", Type: "DB_INSERT", RuleName: "rule", Token: 50, Timestamp: time.Now().Format(time.RFC3339Nano), Payload: "\"payload\"", Status: "ok"},
		},
		{
			name: "event timestamp > timestamp",
			m:    &Module{},
			args: args{token: 50, name: "rule", eventDocID: "eventDocID", batchID: "batchID", status: "ok", event: &model.QueueEventRequest{Timestamp: "5020-03-31T16:16:26+05:30", Delay: int64(0), IsSynchronous: false, Options: make(map[string]string), Payload: "payload", Type: "DB_INSERT"}},
			want: &model.EventDocument{ID: "eventDocID", BatchID: "batchID", Type: "DB_INSERT", RuleName: "rule", Token: 50, Timestamp: "5020-03-31T16:16:26+05:30", Payload: "\"payload\"", Status: "ok"},
		},
		{
			name: "event delay > 0",
			m:    &Module{},
			args: args{token: 50, name: "rule", eventDocID: "eventDocID", batchID: "batchID", status: "ok", event: &model.QueueEventRequest{Timestamp: time.Now().Format(time.RFC3339Nano), Delay: int64(10), IsSynchronous: false, Options: make(map[string]string), Payload: "payload", Type: "DB_INSERT"}},
			want: &model.EventDocument{ID: "eventDocID", BatchID: "batchID", Type: "DB_INSERT", RuleName: "rule", Token: 50, Timestamp: time.Now().Format(time.RFC3339Nano), Payload: "\"payload\"", Status: "ok"},
		},
		{
			name: "event request generated",
			m:    &Module{},
			args: args{token: 50, name: "rule", eventDocID: "eventDocID", batchID: "batchID", status: "ok", event: &model.QueueEventRequest{Timestamp: time.Now().Format(time.RFC3339Nano), Delay: int64(0), IsSynchronous: false, Options: make(map[string]string), Payload: "payload", Type: "DB_INSERT"}},
			want: &model.EventDocument{ID: "eventDocID", BatchID: "batchID", Type: "DB_INSERT", RuleName: "rule", Token: 50, Timestamp: time.Now().Format(time.RFC3339Nano), Payload: "\"payload\"", Status: "ok"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.generateQueueEventRequestRaw(context.Background(), tt.args.token, &config.EventingTrigger{ID: tt.args.name}, tt.args.eventDocID, tt.args.batchID, tt.args.status, tt.args.event)
			if got != nil {
				if !reflect.DeepEqual(got.BatchID, tt.want.BatchID) {
					t.Errorf("Module.generateQueueEventRequest() = %v, want %v", got, tt.want)
				}
				if !reflect.DeepEqual(got.ID, tt.want.ID) {
					t.Errorf("Module.generateQueueEventRequest() = %v, want %v", got, tt.want)
				}
				if !reflect.DeepEqual(got.Type, tt.want.Type) {
					t.Errorf("Module.generateQueueEventRequest() = %v, want %v", got, tt.want)
				}
				if !reflect.DeepEqual(got.RuleName, tt.want.RuleName) {
					t.Errorf("Module.generateQueueEventRequest() = %v, want %v", got, tt.want)
				}
				if !reflect.DeepEqual(got.Token, tt.want.Token) {
					t.Errorf("Module.generateQueueEventRequest() = %v, want %v", got, tt.want)
				}
				if !reflect.DeepEqual(got.Payload, tt.want.Payload) {
					t.Errorf("Module.generateQueueEventRequest() = %v, want %v", got, tt.want)
				}
				if !reflect.DeepEqual(got.Status, tt.want.Status) {
					t.Errorf("Module.generateQueueEventRequest() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestModule_batchRequestsRaw(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx        context.Context
		eventDocID string
		token      int
		requests   []*model.QueueEventRequest
		batchID    string
	}
	tests := []struct {
		name            string
		m               *Module
		args            args
		crudMockArgs    []mockArgs
		syncmanMockArgs []mockArgs
		adminMockArgs   []mockArgs
		authMockArgs    []mockArgs
		wantErr         bool
	}{
		{
			name: "internalCreate error",
			m:    &Module{config: &config.Eventing{DBAlias: mock.Anything, Rules: map[string]*config.EventingTrigger{"some-rule": {Type: utils.EventDBCreate, ID: mock.Anything, Options: map[string]string{}, URL: mock.Anything, Retries: 3}}}, project: mock.Anything},
			args: args{ctx: context.Background(), eventDocID: mock.Anything, token: 50, batchID: mock.Anything, requests: []*model.QueueEventRequest{{Type: utils.EventDBCreate, Delay: 0, IsSynchronous: false, Options: map[string]string{}, Payload: "payload", Timestamp: time.Now().Format(time.RFC3339Nano)}}},
			crudMockArgs: []mockArgs{
				{
					method:         "InternalCreate",
					args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, utils.TableEventingLogs, mock.Anything, false},
					paramsReturned: []interface{}{errors.New("some error")},
				},
			},
			wantErr: true,
		},
		{
			name: "requests are batched",
			m:    &Module{config: &config.Eventing{DBAlias: mock.Anything, Rules: map[string]*config.EventingTrigger{"some-rule": {Type: utils.EventDBCreate, ID: mock.Anything, Options: map[string]string{}, URL: mock.Anything, Retries: 3}}}, project: mock.Anything},
			args: args{ctx: context.Background(), eventDocID: mock.Anything, token: 50, batchID: mock.Anything, requests: []*model.QueueEventRequest{{Type: utils.EventDBCreate, Delay: 0, IsSynchronous: false, Options: map[string]string{}, Payload: "payload", Timestamp: time.Now().Format(time.RFC3339Nano)}}},
			crudMockArgs: []mockArgs{
				{
					method:         "InternalCreate",
					args:           []interface{}{mock.Anything, mock.Anything, mock.Anything, utils.TableEventingLogs, mock.Anything, false},
					paramsReturned: []interface{}{nil},
				},
			},
			syncmanMockArgs: []mockArgs{
				{
					method:         "GetAssignedSpaceCloudID",
					args:           []interface{}{mock.Anything, mock.Anything, 50},
					paramsReturned: []interface{}{mock.Anything, nil},
				},
			},
			adminMockArgs: []mockArgs{},
			authMockArgs:  []mockArgs{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockCrud := mockCrudInterface{}
			mockSyncman := mockSyncmanEventingInterface{}
			mockAuth := mockAuthEventingInterface{}

			for _, m := range tt.crudMockArgs {
				mockCrud.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.syncmanMockArgs {
				mockSyncman.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.authMockArgs {
				mockAuth.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			tt.m.crud = &mockCrud
			tt.m.syncMan = &mockSyncman
			tt.m.auth = &mockAuth

			if err := tt.m.batchRequestsRaw(context.Background(), tt.args.token, tt.args.requests, tt.args.batchID); (err != nil) != tt.wantErr {
				t.Errorf("Module.batchRequests() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockCrud.AssertExpectations(t)
			mockSyncman.AssertExpectations(t)
			mockAuth.AssertExpectations(t)
		})
	}
}
