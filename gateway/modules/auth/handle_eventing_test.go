package auth

import (
	"context"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

func TestModule_getEventingRule(t *testing.T) {
	type args struct {
		eventType string
	}
	tests := []struct {
		name    string
		m       *Module
		args    args
		want    *config.Rule
		wantErr bool
	}{
		{
			name: "rule found",
			m:    &Module{eventingRules: map[string]*config.Rule{"some-type": &config.Rule{}}},
			args: args{eventType: "some-type"},
			want: &config.Rule{},
		},
		{
			name:    "rule not found",
			m:       &Module{eventingRules: map[string]*config.Rule{"some-type2": &config.Rule{}}},
			args:    args{eventType: "some-type1"},
			wantErr: true,
		},
		{
			name: "default rule found",
			m:    &Module{eventingRules: map[string]*config.Rule{"default": &config.Rule{Rule: "allow"}}},
			args: args{eventType: "some-type"},
			want: &config.Rule{Rule: "allow"},
		},
		{
			name:    "empty rules",
			m:       &Module{},
			args:    args{eventType: "some-type"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.getEventingRule(tt.args.eventType)
			if (err != nil) != tt.wantErr {
				t.Errorf("Module.getEventingRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Module.getEventingRule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModule_IsEventingOpAuthorised(t *testing.T) {
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
			name: "got rule",
			m:    &Module{project: "some-project", eventingRules: map[string]*config.Rule{"some-type": &config.Rule{Rule: "authenticated"}}, secrets: []*config.Secret{{IsPrimary: true, Secret: "mySecretkey"}}},
			args: args{ctx: context.Background(), project: "some-project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc", event: &model.QueueEventRequest{Type: "some-type", Delay: 0, Payload: "something", Options: make(map[string]string)}},
		},
		{
			name:    "did not get rule",
			m:       &Module{eventingRules: map[string]*config.Rule{"some-type": &config.Rule{}}, secrets: []*config.Secret{{IsPrimary: true, Secret: "mySecretkey"}}},
			args:    args{ctx: context.Background(), project: "some-project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc", event: &model.QueueEventRequest{Type: "some-type1", Delay: 0, Payload: "something", Options: make(map[string]string)}},
			wantErr: true,
		},
		{
			name: "valid project details",
			m:    &Module{project: "some-project", eventingRules: map[string]*config.Rule{"some-type": &config.Rule{Rule: "allow"}}, secrets: []*config.Secret{{IsPrimary: true, Secret: "mySecretkey"}}},
			args: args{ctx: context.Background(), project: "some-project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc", event: &model.QueueEventRequest{Type: "some-type", Delay: 0, Payload: "something", Options: make(map[string]string)}},
		},
		{
			name:    "invalid project details",
			m:       &Module{project: "some-project", eventingRules: map[string]*config.Rule{"some-type": &config.Rule{Rule: "allow"}}, secrets: []*config.Secret{{IsPrimary: true, Secret: "mySecretkey"}}},
			args:    args{ctx: context.Background(), project: "some-project1", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc", event: &model.QueueEventRequest{Type: "some-type", Delay: 0, Payload: "something", Options: make(map[string]string)}},
			wantErr: true,
		},
		{
			name:    "did not get auth",
			m:       &Module{project: "some-project", eventingRules: map[string]*config.Rule{"some-type": &config.Rule{Rule: "authenticated"}}, secrets: []*config.Secret{{IsPrimary: true, Secret: "mySecretkey"}}},
			args:    args{ctx: context.Background(), project: "some-project", token: "token", event: &model.QueueEventRequest{Type: "some-type", Delay: 0, Payload: "something", Options: make(map[string]string)}},
			wantErr: true,
		},
		{
			name:    "rules did not match",
			m:       &Module{project: "some-project", eventingRules: map[string]*config.Rule{"some-type": &config.Rule{Rule: "deny"}}, secrets: []*config.Secret{{IsPrimary: true, Secret: "mySecretkey"}}},
			args:    args{ctx: context.Background(), project: "some-project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc", event: &model.QueueEventRequest{Type: "some-type", Delay: 0, Payload: "something", Options: make(map[string]string)}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.IsEventingOpAuthorised(tt.args.ctx, tt.args.project, tt.args.token, tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("Module.IsEventingOpAuthorised() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
