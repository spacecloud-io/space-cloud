package syncman

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

func TestManager_SetEventingRule(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx      context.Context
		project  string
		ruleName string
		value    *config.EventingTrigger
	}
	tests := []struct {
		name            string
		s               *Manager
		args            args
		modulesMockArgs []mockArgs
		storeMockArgs   []mockArgs
		wantErr         bool
	}{
		{
			name:    "unable to get project config",
			s:       &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}, EventingTriggers: config.EventingTriggers{config.GenerateResourceID("chicago", "1", config.ResourceEventingTrigger, "rule"): &config.EventingTrigger{ID: "rule"}}}}}},
			args:    args{ctx: context.Background(), project: "2", ruleName: "rule", value: &config.EventingTrigger{ID: "rule"}},
			wantErr: true,
		},
		{
			name: "unable to set eventing config",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}, EventingTriggers: config.EventingTriggers{config.GenerateResourceID("chicago", "1", config.ResourceEventingTrigger, "rule"): &config.EventingTrigger{ID: "rule"}}}}}},
			args: args{ctx: context.Background(), project: "1", ruleName: "rule", value: &config.EventingTrigger{ID: "rule"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingTriggerConfig",
					args:           []interface{}{mock.Anything, config.EventingTriggers{config.GenerateResourceID("chicago", "1", config.ResourceEventingTrigger, "rule"): &config.EventingTrigger{ID: "rule"}}},
					paramsReturned: []interface{}{errors.New("error setting eventing module config")},
				},
			},
			wantErr: true,
		},
		{
			name: "unable to set project",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}, EventingTriggers: config.EventingTriggers{config.GenerateResourceID("chicago", "1", config.ResourceEventingTrigger, "rule"): &config.EventingTrigger{ID: "rule"}}}}}},
			args: args{ctx: context.Background(), project: "1", ruleName: "rule", value: &config.EventingTrigger{ID: "rule"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingTriggerConfig",
					args:           []interface{}{mock.Anything, config.EventingTriggers{config.GenerateResourceID("chicago", "1", config.ResourceEventingTrigger, "rule"): &config.EventingTrigger{ID: "rule"}}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{mock.Anything, config.GenerateResourceID("chicago", "1", config.ResourceEventingTrigger, "rule"), &config.EventingTrigger{ID: "rule"}},
					paramsReturned: []interface{}{errors.New("Invalid config file type")},
				},
			},
			wantErr: true,
		},
		{
			name: "eventing rules are set",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}, EventingTriggers: config.EventingTriggers{config.GenerateResourceID("chicago", "1", config.ResourceEventingTrigger, "rule"): &config.EventingTrigger{ID: "rule"}}}}}},
			args: args{ctx: context.Background(), project: "1", ruleName: "rule", value: &config.EventingTrigger{ID: "rule"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingTriggerConfig",
					args:           []interface{}{mock.Anything, config.EventingTriggers{config.GenerateResourceID("chicago", "1", config.ResourceEventingTrigger, "rule"): &config.EventingTrigger{ID: "rule"}}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{mock.Anything, config.GenerateResourceID("chicago", "1", config.ResourceEventingTrigger, "rule"), &config.EventingTrigger{ID: "rule"}},
					paramsReturned: []interface{}{nil},
				},
			},
		},
		{
			name: "eventing rules are set when config rules are nil",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}}}}},
			args: args{ctx: context.Background(), project: "1", ruleName: "rule", value: &config.EventingTrigger{ID: "rule"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingTriggerConfig",
					args:           []interface{}{mock.Anything, config.EventingTriggers{config.GenerateResourceID("chicago", "1", config.ResourceEventingTrigger, "rule"): &config.EventingTrigger{ID: "rule"}}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{mock.Anything, config.GenerateResourceID("chicago", "1", config.ResourceEventingTrigger, "rule"), &config.EventingTrigger{ID: "rule"}},
					paramsReturned: []interface{}{nil},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockModules := mockModulesInterface{}
			mockStore := mockStoreInterface{}

			for _, m := range tt.modulesMockArgs {
				mockModules.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.storeMockArgs {
				mockStore.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			tt.s.modules = &mockModules
			tt.s.store = &mockStore

			if _, err := tt.s.SetEventingRule(tt.args.ctx, tt.args.project, tt.args.ruleName, tt.args.value, model.RequestParams{}); (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetEventingRule() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockModules.AssertExpectations(t)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestManager_SetDeleteEventingRule(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx      context.Context
		project  string
		ruleName string
	}
	tests := []struct {
		name            string
		s               *Manager
		args            args
		modulesMockArgs []mockArgs
		storeMockArgs   []mockArgs
		wantErr         bool
	}{
		{
			name:    "unable to get project config",
			s:       &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}, EventingTriggers: config.EventingTriggers{config.GenerateResourceID("chicago", "1", config.ResourceEventingTrigger, "rule"): &config.EventingTrigger{ID: "rule"}}}}}},
			args:    args{ctx: context.Background(), project: "2", ruleName: "rule"},
			wantErr: true,
		},
		{
			name: "unable to set eventing config",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}, EventingTriggers: config.EventingTriggers{config.GenerateResourceID("chicago", "1", config.ResourceEventingTrigger, "rule"): &config.EventingTrigger{ID: "rule"}}}}}},
			args: args{ctx: context.Background(), project: "1", ruleName: "rule"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingTriggerConfig",
					args:           []interface{}{mock.Anything, config.EventingTriggers{}},
					paramsReturned: []interface{}{errors.New("error setting eventing module config")},
				},
			},
			wantErr: true,
		},
		{
			name: "unable to delete project",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}, EventingTriggers: config.EventingTriggers{config.GenerateResourceID("chicago", "1", config.ResourceEventingTrigger, "rule"): &config.EventingTrigger{ID: "rule"}}}}}},
			args: args{ctx: context.Background(), project: "1", ruleName: "rule"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingTriggerConfig",
					args:           []interface{}{mock.Anything, config.EventingTriggers{}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "DeleteResource",
					args:           []interface{}{mock.Anything, config.GenerateResourceID("chicago", "1", config.ResourceEventingTrigger, "rule")},
					paramsReturned: []interface{}{errors.New("Invalid config file type")},
				},
			},
			wantErr: true,
		},
		{
			name: "eventing rule deleted successfully",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}, EventingTriggers: config.EventingTriggers{config.GenerateResourceID("chicago", "1", config.ResourceEventingTrigger, "rule"): &config.EventingTrigger{ID: "rule"}}}}}},
			args: args{ctx: context.Background(), project: "1", ruleName: "rule"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingTriggerConfig",
					args:           []interface{}{mock.Anything, config.EventingTriggers{}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "DeleteResource",
					args:           []interface{}{mock.Anything, config.GenerateResourceID("chicago", "1", config.ResourceEventingTrigger, "rule")},
					paramsReturned: []interface{}{nil},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockModules := mockModulesInterface{}
			mockStore := mockStoreInterface{}

			for _, m := range tt.modulesMockArgs {
				mockModules.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.storeMockArgs {
				mockStore.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			tt.s.modules = &mockModules
			tt.s.store = &mockStore

			if _, err := tt.s.SetDeleteEventingRule(tt.args.ctx, tt.args.project, tt.args.ruleName, model.RequestParams{}); (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetDeleteEventingRule() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockModules.AssertExpectations(t)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestManager_SetEventingSchema(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx     context.Context
		project string
		evType  string
		schema  string
	}
	tests := []struct {
		name            string
		s               *Manager
		args            args
		modulesMockArgs []mockArgs
		storeMockArgs   []mockArgs
		wantErr         bool
	}{
		{
			name:    "unable to get project config",
			s:       &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}, EventingSchemas: config.EventingSchemas{config.GenerateResourceID("chicago", "1", config.ResourceEventingSchema, "evType"): &config.EventingSchema{ID: "evType", Schema: "type evType {id: String!}"}}}}}},
			args:    args{ctx: context.Background(), evType: "evType", project: "2", schema: "type evType {id: String!}"},
			wantErr: true,
		},
		{
			name: "unable to set eventing config",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}, EventingSchemas: config.EventingSchemas{}}}}},
			args: args{ctx: context.Background(), evType: "evType", project: "1", schema: "type evType {id: String!}"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingSchemaConfig",
					args:           []interface{}{mock.Anything, config.EventingSchemas{config.GenerateResourceID("chicago", "1", config.ResourceEventingSchema, "evType"): &config.EventingSchema{ID: "evType", Schema: "type evType {id: String!}"}}},
					paramsReturned: []interface{}{errors.New("error setting eventing module config")},
				},
			},
			wantErr: true,
		},
		{
			name: "unable to set project",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}, EventingSchemas: config.EventingSchemas{}}}}},
			args: args{ctx: context.Background(), evType: "evType", project: "1", schema: "type evType {id: String!}"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingSchemaConfig",
					args:           []interface{}{mock.Anything, config.EventingSchemas{config.GenerateResourceID("chicago", "1", config.ResourceEventingSchema, "evType"): &config.EventingSchema{ID: "evType", Schema: "type evType {id: String!}"}}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{mock.Anything, config.GenerateResourceID("chicago", "1", config.ResourceEventingSchema, "evType"), &config.EventingSchema{ID: "evType", Schema: "type evType {id: String!}"}},
					paramsReturned: []interface{}{errors.New("Invalid config file type")},
				},
			},
			wantErr: true,
		},
		{
			name: "eventing schema is set",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}, EventingSchemas: config.EventingSchemas{}}}}},
			args: args{ctx: context.Background(), evType: "evType", project: "1", schema: "type evType {id: String!}"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingSchemaConfig",
					args:           []interface{}{mock.Anything, config.EventingSchemas{config.GenerateResourceID("chicago", "1", config.ResourceEventingSchema, "evType"): &config.EventingSchema{ID: "evType", Schema: "type evType {id: String!}"}}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{mock.Anything, config.GenerateResourceID("chicago", "1", config.ResourceEventingSchema, "evType"), &config.EventingSchema{ID: "evType", Schema: "type evType {id: String!}"}},
					paramsReturned: []interface{}{nil},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockModules := mockModulesInterface{}
			mockStore := mockStoreInterface{}

			for _, m := range tt.modulesMockArgs {
				mockModules.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.storeMockArgs {
				mockStore.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			tt.s.modules = &mockModules
			tt.s.store = &mockStore

			if _, err := tt.s.SetEventingSchema(tt.args.ctx, tt.args.project, tt.args.evType, tt.args.schema, model.RequestParams{}); (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetEventingSchema() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockModules.AssertExpectations(t)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestManager_SetDeleteEventingSchema(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx     context.Context
		project string
		evType  string
	}
	tests := []struct {
		name            string
		s               *Manager
		args            args
		modulesMockArgs []mockArgs
		storeMockArgs   []mockArgs
		wantErr         bool
	}{
		{
			name:    "unable to get project config",
			s:       &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}, EventingSchemas: config.EventingSchemas{config.GenerateResourceID("chicago", "1", config.ResourceEventingSchema, "evType"): &config.EventingSchema{ID: "evType", Schema: "type evType {id: String!}"}}}}}},
			args:    args{ctx: context.Background(), evType: "evType", project: "2"},
			wantErr: true,
		},
		{
			name: "unable to set eventing config",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}, EventingSchemas: config.EventingSchemas{config.GenerateResourceID("chicago", "1", config.ResourceEventingSchema, "evType"): &config.EventingSchema{ID: "evType", Schema: "type evType {id: String!}"}}}}}},
			args: args{ctx: context.Background(), evType: "evType", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingSchemaConfig",
					args:           []interface{}{mock.Anything, config.EventingSchemas{}},
					paramsReturned: []interface{}{errors.New("error setting eventing module config")},
				},
			},
			wantErr: true,
		},
		{
			name: "unable to set project",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}, EventingSchemas: config.EventingSchemas{config.GenerateResourceID("chicago", "1", config.ResourceEventingSchema, "evType"): &config.EventingSchema{ID: "evType", Schema: "type evType {id: String!}"}}}}}},
			args: args{ctx: context.Background(), evType: "evType", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingSchemaConfig",
					args:           []interface{}{mock.Anything, config.EventingSchemas{}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "DeleteResource",
					args:           []interface{}{mock.Anything, config.GenerateResourceID("chicago", "1", config.ResourceEventingSchema, "evType")},
					paramsReturned: []interface{}{errors.New("Invalid config file type")},
				},
			},
			wantErr: true,
		},
		{
			name: "eventing schema is deleted",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}, EventingSchemas: config.EventingSchemas{config.GenerateResourceID("chicago", "1", config.ResourceEventingSchema, "evType"): &config.EventingSchema{ID: "evType", Schema: "type evType {id: String!}"}}}}}},
			args: args{ctx: context.Background(), evType: "evType", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingSchemaConfig",
					args:           []interface{}{mock.Anything, config.EventingSchemas{}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "DeleteResource",
					args:           []interface{}{mock.Anything, config.GenerateResourceID("chicago", "1", config.ResourceEventingSchema, "evType")},
					paramsReturned: []interface{}{nil},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockModules := mockModulesInterface{}
			mockStore := mockStoreInterface{}

			for _, m := range tt.modulesMockArgs {
				mockModules.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.storeMockArgs {
				mockStore.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			tt.s.modules = &mockModules
			tt.s.store = &mockStore

			if _, err := tt.s.SetDeleteEventingSchema(tt.args.ctx, tt.args.project, tt.args.evType, model.RequestParams{}); (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetDeleteEventingSchema() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockModules.AssertExpectations(t)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestManager_SetEventingSecurityRules(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx     context.Context
		project string
		evType  string
		rule    *config.Rule
	}
	tests := []struct {
		name            string
		s               *Manager
		args            args
		modulesMockArgs []mockArgs
		storeMockArgs   []mockArgs
		wantErr         bool
	}{
		{
			name:    "unable to get project config",
			s:       &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}, EventingRules: config.EventingRules{config.GenerateResourceID("chicago", "1", config.ResourceEventingRule, "evType"): &config.Rule{Type: "evType", ID: "evType"}}}}}},
			args:    args{ctx: context.Background(), evType: "evType", project: "2", rule: &config.Rule{}},
			wantErr: true,
		},
		{
			name: "unable to set eventing config",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}, EventingRules: config.EventingRules{config.GenerateResourceID("chicago", "1", config.ResourceEventingRule, "evType"): &config.Rule{Type: "evType", ID: "evType"}}}}}},
			args: args{ctx: context.Background(), evType: "evType", project: "1", rule: &config.Rule{}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingRuleConfig",
					args:           []interface{}{mock.Anything, config.EventingRules{config.GenerateResourceID("chicago", "1", config.ResourceEventingRule, "evType"): &config.Rule{ID: "evType"}}},
					paramsReturned: []interface{}{errors.New("error setting eventing module config")},
				},
			},
			wantErr: true,
		},
		{
			name: "unable to delete resource",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}, EventingRules: config.EventingRules{config.GenerateResourceID("chicago", "1", config.ResourceEventingRule, "evType"): &config.Rule{Type: "evType", ID: "evType"}}}}}},
			args: args{ctx: context.Background(), evType: "evType", project: "1", rule: &config.Rule{ID: "evType"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingRuleConfig",
					args:           []interface{}{mock.Anything, config.EventingRules{config.GenerateResourceID("chicago", "1", config.ResourceEventingRule, "evType"): &config.Rule{ID: "evType"}}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{mock.Anything, config.GenerateResourceID("chicago", "1", config.ResourceEventingRule, "evType"), &config.Rule{ID: "evType"}},
					paramsReturned: []interface{}{errors.New("Invalid config file type")},
				},
			},
			wantErr: true,
		},
		{
			name: "security rules empty in config and they are set",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}}}}},
			args: args{ctx: context.Background(), evType: "evType", project: "1", rule: &config.Rule{ID: "evType"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingRuleConfig",
					args:           []interface{}{mock.Anything, config.EventingRules{config.GenerateResourceID("chicago", "1", config.ResourceEventingRule, "evType"): &config.Rule{ID: "evType"}}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{mock.Anything, config.GenerateResourceID("chicago", "1", config.ResourceEventingRule, "evType"), &config.Rule{ID: "evType"}},
					paramsReturned: []interface{}{nil},
				},
			},
		},
		{
			name: "security rules are set",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}, EventingRules: config.EventingRules{config.GenerateResourceID("chicago", "1", config.ResourceEventingRule, "evType"): &config.Rule{Type: "evType", ID: "evType"}}}}}},
			args: args{ctx: context.Background(), evType: "evType", project: "1", rule: &config.Rule{ID: "evType"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingRuleConfig",
					args:           []interface{}{mock.Anything, config.EventingRules{config.GenerateResourceID("chicago", "1", config.ResourceEventingRule, "evType"): &config.Rule{ID: "evType"}}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetResource",
					args:           []interface{}{mock.Anything, config.GenerateResourceID("chicago", "1", config.ResourceEventingRule, "evType"), &config.Rule{ID: "evType"}},
					paramsReturned: []interface{}{nil},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockModules := mockModulesInterface{}
			mockStore := mockStoreInterface{}

			for _, m := range tt.modulesMockArgs {
				mockModules.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.storeMockArgs {
				mockStore.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			tt.s.modules = &mockModules
			tt.s.store = &mockStore

			if _, err := tt.s.SetEventingSecurityRules(tt.args.ctx, tt.args.project, tt.args.evType, tt.args.rule, model.RequestParams{}); (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetEventingSecurityRules() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockModules.AssertExpectations(t)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestManager_SetDeleteEventingSecurityRules(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx     context.Context
		project string
		evType  string
	}
	tests := []struct {
		name            string
		s               *Manager
		args            args
		modulesMockArgs []mockArgs
		storeMockArgs   []mockArgs
		wantErr         bool
	}{
		{
			name:    "unable to get project config",
			s:       &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}, EventingRules: config.EventingRules{config.GenerateResourceID("chicago", "1", config.ResourceEventingRule, "evType"): &config.Rule{Type: "evType", ID: "evType"}}}}}},
			args:    args{ctx: context.Background(), evType: "evType", project: "2"},
			wantErr: true,
		},
		{
			name: "unable to set eventing config",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}, EventingRules: config.EventingRules{config.GenerateResourceID("chicago", "1", config.ResourceEventingRule, "evType"): &config.Rule{Type: "evType", ID: "evType"}}}}}},
			args: args{ctx: context.Background(), evType: "evType", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingRuleConfig",
					args:           []interface{}{mock.Anything, config.EventingRules{}},
					paramsReturned: []interface{}{errors.New("error setting eventing module config")},
				},
			},
			wantErr: true,
		},
		{
			name: "unable to set project",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}, EventingRules: config.EventingRules{config.GenerateResourceID("chicago", "1", config.ResourceEventingRule, "evType"): &config.Rule{Type: "evType", ID: "evType"}}}}}},
			args: args{ctx: context.Background(), evType: "evType", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingRuleConfig",
					args:           []interface{}{mock.Anything, config.EventingRules{}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "DeleteResource",
					args:           []interface{}{mock.Anything, config.GenerateResourceID("chicago", "1", config.ResourceEventingRule, "evType")},
					paramsReturned: []interface{}{errors.New("Invalid config file type")},
				},
			},
			wantErr: true,
		},
		{
			name: "eventing security rules deleted",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}, EventingRules: config.EventingRules{config.GenerateResourceID("chicago", "1", config.ResourceEventingRule, "evType"): &config.Rule{Type: "evType", ID: "evType"}}}}}},
			args: args{ctx: context.Background(), evType: "evType", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingRuleConfig",
					args:           []interface{}{mock.Anything, config.EventingRules{}},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "DeleteResource",
					args:           []interface{}{mock.Anything, config.GenerateResourceID("chicago", "1", config.ResourceEventingRule, "evType")},
					paramsReturned: []interface{}{nil},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockModules := mockModulesInterface{}
			mockStore := mockStoreInterface{}

			for _, m := range tt.modulesMockArgs {
				mockModules.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.storeMockArgs {
				mockStore.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			tt.s.modules = &mockModules
			tt.s.store = &mockStore

			if _, err := tt.s.SetDeleteEventingSecurityRules(tt.args.ctx, tt.args.project, tt.args.evType, model.RequestParams{}); (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetDeleteEventingSecurityRules() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockModules.AssertExpectations(t)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestManager_GetEventingTriggerRules(t *testing.T) {
	type args struct {
		ctx     context.Context
		project string
		id      string
	}
	tests := []struct {
		name    string
		s       *Manager
		args    args
		want    []interface{}
		wantErr bool
	}{
		{
			name:    "unable to get project config",
			s:       &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}, EventingTriggers: config.EventingTriggers{config.GenerateResourceID("chicago", "1", config.ResourceEventingTrigger, "rule"): &config.EventingTrigger{ID: "rule"}}}}}},
			args:    args{ctx: context.Background(), id: "rule", project: "2"},
			wantErr: true,
		},
		{
			name:    "id not present in config",
			s:       &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}, EventingTriggers: config.EventingTriggers{config.GenerateResourceID("chicago", "1", config.ResourceEventingTrigger, "rule"): &config.EventingTrigger{ID: "rule"}}}}}},
			args:    args{ctx: context.Background(), id: "notRule", project: "1"},
			wantErr: true,
		},
		{
			name: "got trigger rule",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}, EventingTriggers: config.EventingTriggers{config.GenerateResourceID("chicago", "1", config.ResourceEventingTrigger, "rule"): &config.EventingTrigger{ID: "rule"}}}}}},
			args: args{ctx: context.Background(), id: "rule", project: "1"},
			want: []interface{}{&config.EventingTrigger{ID: "rule"}},
		},
		{
			name: "id is empty and got all trigger rules",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}, EventingTriggers: config.EventingTriggers{config.GenerateResourceID("chicago", "1", config.ResourceEventingTrigger, "rule"): &config.EventingTrigger{ID: "rule"}}}}}},
			args: args{ctx: context.Background(), id: "*", project: "1"},
			want: []interface{}{&config.EventingTrigger{ID: "rule"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got, err := tt.s.GetEventingTriggerRules(tt.args.ctx, tt.args.project, tt.args.id, model.RequestParams{})
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.GetEventingTriggerRules() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.GetEventingTriggerRules() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_GetEventingSchema(t *testing.T) {
	type args struct {
		ctx     context.Context
		project string
		id      string
	}
	tests := []struct {
		name    string
		s       *Manager
		args    args
		want    []interface{}
		wantErr bool
	}{
		{
			name:    "unable to get project config",
			s:       &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}, EventingSchemas: config.EventingSchemas{"": &config.EventingSchema{ID: "id"}}}}}},
			args:    args{ctx: context.Background(), id: "id", project: "2"},
			wantErr: true,
		},
		{
			name:    "id not present in config",
			s:       &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}, EventingSchemas: config.EventingSchemas{"": &config.EventingSchema{ID: "id"}}}}}},
			args:    args{ctx: context.Background(), id: "notId", project: "1"},
			wantErr: true,
		},
		{
			name: "got schema",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}, EventingSchemas: config.EventingSchemas{config.GenerateResourceID("chicago", "1", config.ResourceEventingSchema, "id"): &config.EventingSchema{ID: "id"}}}}}},
			args: args{ctx: context.Background(), id: "id", project: "1"},
			want: []interface{}{&config.EventingSchema{ID: "id"}},
		},
		{
			name: "id empty and got schemas",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}, EventingSchemas: config.EventingSchemas{config.GenerateResourceID("chicago", "1", config.ResourceEventingSchema, "id"): &config.EventingSchema{ID: "id"}}}}}},
			args: args{ctx: context.Background(), id: "*", project: "1"},
			want: []interface{}{&config.EventingSchema{ID: "id"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got, err := tt.s.GetEventingSchema(tt.args.ctx, tt.args.project, tt.args.id, model.RequestParams{})
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.GetEventingSchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.GetEventingSchema() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_GetEventingSecurityRules(t *testing.T) {
	type args struct {
		ctx     context.Context
		project string
		id      string
	}
	tests := []struct {
		name    string
		s       *Manager
		args    args
		want    []interface{}
		wantErr bool
	}{
		{
			name:    "unable to get project config",
			s:       &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}, EventingSchemas: config.EventingSchemas{"": &config.EventingSchema{ID: "id"}}}}}},
			args:    args{ctx: context.Background(), id: "id", project: "2"},
			wantErr: true,
		},
		{
			name:    "id not present in config",
			s:       &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}, EventingSchemas: config.EventingSchemas{"": &config.EventingSchema{ID: "id"}}}}}},
			args:    args{ctx: context.Background(), id: "notId", project: "1"},
			wantErr: true,
		},
		{
			name: "got security rule",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}, EventingRules: config.EventingRules{config.GenerateResourceID("chicago", "1", config.ResourceEventingRule, "id"): &config.Rule{ID: "id"}}}}}},
			args: args{ctx: context.Background(), id: "id", project: "1"},
			want: []interface{}{&config.Rule{ID: "id"}},
		},
		{
			name: "id empty and got security rules",
			s:    &Manager{clusterID: "chicago", projectConfig: &config.Config{Projects: config.Projects{"1": &config.Project{ProjectConfig: &config.ProjectConfig{ID: "1"}, EventingConfig: &config.EventingConfig{}, EventingRules: config.EventingRules{config.GenerateResourceID("chicago", "1", config.ResourceEventingRule, "id"): &config.Rule{ID: "id"}}}}}},
			args: args{ctx: context.Background(), id: "*", project: "1"},
			want: []interface{}{&config.Rule{ID: "id"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got, err := tt.s.GetEventingSecurityRules(tt.args.ctx, tt.args.project, tt.args.id, model.RequestParams{})
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.GetEventingSecurityRules() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.GetEventingSecurityRules() = %v, want %v", got, tt.want)
			}
		})
	}
}
