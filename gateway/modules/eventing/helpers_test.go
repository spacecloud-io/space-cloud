package eventing

import (
	"context"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/auth"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud"
	"github.com/spaceuptech/space-cloud/gateway/modules/schema"
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
		want    config.EventingRule
		wantErr bool
	}{
		{
			name: "event type is an internal type",
			m:    &Module{config: &config.Eventing{Rules: map[string]config.EventingRule{"some-rule": config.EventingRule{Type: "DB_INSERT"}}}},
			args: args{name: "some-rule", evType: "DB_INSERT"},
			want: config.EventingRule{Type: "DB_INSERT", Retries: 3, Timeout: 5000},
		},
		{
			name: "event type is found in rules",
			m:    &Module{config: &config.Eventing{Rules: map[string]config.EventingRule{"some-rule": config.EventingRule{Type: "event"}}}},
			args: args{name: "some-rule", evType: "event"},
			want: config.EventingRule{Type: "event"},
		},
		{
			name: "event type is found in internal rules",
			m:    &Module{config: &config.Eventing{InternalRules: map[string]config.EventingRule{"some-rule": config.EventingRule{Type: "event"}}}},
			args: args{name: "some-rule", evType: "event"},
			want: config.EventingRule{Type: "event"},
		},
		{
			name:    "event type is not found",
			m:       &Module{config: &config.Eventing{}},
			args:    args{name: "some-rule", evType: "event"},
			want:    config.EventingRule{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.selectRule(tt.args.name, tt.args.evType)
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
	authModule := auth.Init("1", &crud.Module{}, &schema.Schema{}, false)
	err := authModule.SetConfig("project", "mySecretkey", config.Crud{}, &config.FileStore{}, &config.ServicesModule{}, &config.Eventing{SecurityRules: map[string]*config.Rule{"event": &config.Rule{Rule: "authenticated"}}})
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
			name: "event type is an internal type",
			m:    &Module{config: &config.Eventing{Rules: map[string]config.EventingRule{"some-rule": config.EventingRule{Type: "DB_INSERT"}}}},
			args: args{event: &model.QueueEventRequest{Type: "DB_INSERT", Delay: 0, Timestamp: 0, Payload: "something", Options: make(map[string]string)}},
		},
		{
			name:    "invalid project details",
			m:       &Module{auth: &auth.Module{}},
			args:    args{ctx: context.Background(), project: "some-project", event: &model.QueueEventRequest{Type: "event", Delay: 0, Timestamp: 0, Payload: "something", Options: make(map[string]string)}},
			wantErr: true,
		},
		{
			name:    "invalid token",
			m:       &Module{auth: &auth.Module{}},
			args:    args{ctx: context.Background(), token: "token", event: &model.QueueEventRequest{Type: "event", Delay: 0, Timestamp: 0, Payload: "something", Options: make(map[string]string)}},
			wantErr: true,
		},
		{
			name: "event type not in schemas",
			m:    &Module{auth: authModule, config: &config.Eventing{SecurityRules: map[string]*config.Rule{"event": {Rule: "authenticated"}}, Schemas: map[string]config.SchemaObject{"event": config.SchemaObject{Schema: "some-schema"}}}},
			args: args{ctx: context.Background(), project: "project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc", event: &model.QueueEventRequest{Type: "event", Delay: 0, Timestamp: 0, Payload: "some-schema", Options: make(map[string]string)}},
		},
		{
			name: "no schema given",
			m:    &Module{schemas: map[string]schema.Fields{"event": {}}, auth: authModule, config: &config.Eventing{SecurityRules: map[string]*config.Rule{"event": &config.Rule{Rule: "authenticated"}}, Schemas: map[string]config.SchemaObject{"event": config.SchemaObject{Schema: "type event {id: ID! title: String}"}}}},
			args: args{ctx: context.Background(), project: "project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc", event: &model.QueueEventRequest{Type: "event", Delay: 0, Timestamp: 0, Payload: make(map[string]interface{}), Options: make(map[string]string)}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.auth = authModule

			if err := tt.m.validate(tt.args.ctx, tt.args.project, tt.args.token, tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("Module.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
