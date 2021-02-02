package auth

import (
	"context"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	jwtUtils "github.com/spaceuptech/space-cloud/gateway/utils/jwt"
)

func TestModule_getEventingRule(t *testing.T) {
	type args struct {
		eventType string
		project   string
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
			m:    &Module{clusterID: "chicago", eventingRules: map[string]*config.Rule{config.GenerateResourceID("chicago", "project", config.ResourceEventingRule, "some-type"): &config.Rule{}}},
			args: args{eventType: "some-type", project: "project"},
			want: &config.Rule{},
		},
		{
			name:    "rule not found",
			m:       &Module{eventingRules: map[string]*config.Rule{"some-type2": &config.Rule{}}},
			args:    args{eventType: "some-type1", project: "project"},
			wantErr: true,
		},
		{
			name: "default rule found",
			m:    &Module{clusterID: "chicago", eventingRules: map[string]*config.Rule{config.GenerateResourceID("chicago", "project", config.ResourceEventingRule, "default"): &config.Rule{Rule: "allow"}}},
			args: args{eventType: "some-type", project: "project"},
			want: &config.Rule{Rule: "allow"},
		},
		{
			name:    "empty rules",
			m:       &Module{},
			args:    args{eventType: "some-type", project: "project"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.getEventingRule(context.Background(), tt.args.project, tt.args.eventType)
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
			m:    &Module{clusterID: "chicago", jwt: jwtUtils.New(), project: "some-project", eventingRules: map[string]*config.Rule{config.GenerateResourceID("chicago", "some-project", config.ResourceEventingRule, "some-type"): &config.Rule{Rule: "authenticated"}}},
			args: args{ctx: context.Background(), project: "some-project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc", event: &model.QueueEventRequest{Type: "some-type", Delay: 0, Payload: "something", Options: make(map[string]string)}},
		},
		{
			name:    "did not get rule",
			m:       &Module{jwt: jwtUtils.New(), eventingRules: map[string]*config.Rule{"some-type": &config.Rule{}}},
			args:    args{ctx: context.Background(), project: "some-project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc", event: &model.QueueEventRequest{Type: "some-type1", Delay: 0, Payload: "something", Options: make(map[string]string)}},
			wantErr: true,
		},
		{
			name: "valid project details",
			m:    &Module{clusterID: "chicago", jwt: jwtUtils.New(), project: "some-project", eventingRules: map[string]*config.Rule{config.GenerateResourceID("chicago", "some-project", config.ResourceEventingRule, "some-type"): &config.Rule{Rule: "allow"}}},
			args: args{ctx: context.Background(), project: "some-project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc", event: &model.QueueEventRequest{Type: "some-type", Delay: 0, Payload: "something", Options: make(map[string]string)}},
		},
		{
			name:    "invalid project details",
			m:       &Module{jwt: jwtUtils.New(), project: "some-project", eventingRules: map[string]*config.Rule{"some-type": &config.Rule{Rule: "allow"}}},
			args:    args{ctx: context.Background(), project: "some-project1", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc", event: &model.QueueEventRequest{Type: "some-type", Delay: 0, Payload: "something", Options: make(map[string]string)}},
			wantErr: true,
		},
		{
			name:    "did not get auth",
			m:       &Module{jwt: jwtUtils.New(), project: "some-project", eventingRules: map[string]*config.Rule{"some-type": &config.Rule{Rule: "authenticated"}}},
			args:    args{ctx: context.Background(), project: "some-project", token: "token", event: &model.QueueEventRequest{Type: "some-type", Delay: 0, Payload: "something", Options: make(map[string]string)}},
			wantErr: true,
		},
		{
			name:    "rules did not match",
			m:       &Module{jwt: jwtUtils.New(), project: "some-project", eventingRules: map[string]*config.Rule{"some-type": &config.Rule{Rule: "deny"}}},
			args:    args{ctx: context.Background(), project: "some-project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc", event: &model.QueueEventRequest{Type: "some-type", Delay: 0, Payload: "something", Options: make(map[string]string)}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = tt.m.SetConfig(context.TODO(), "local", &config.ProjectConfig{ID: tt.m.project, Secrets: []*config.Secret{{IsPrimary: true, Secret: "mySecretkey"}}}, config.DatabaseRules{}, config.DatabasePreparedQueries{}, config.FileStoreRules{}, config.Services{}, tt.m.eventingRules, config.SecurityFunctions{})
			if _, err := tt.m.IsEventingOpAuthorised(context.Background(), tt.args.project, tt.args.token, tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("Module.IsEventingOpAuthorised() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
