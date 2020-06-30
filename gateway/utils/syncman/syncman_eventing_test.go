package syncman

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/stretchr/testify/mock"
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
		value    config.EventingRule
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
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{Rules: map[string]config.EventingRule{"rule": {}}}}}}}},
			args:    args{ctx: context.Background(), project: "2", ruleName: "rule", value: config.EventingRule{}},
			wantErr: true,
		},
		{
			name: "unable to set eventing config",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{Rules: map[string]config.EventingRule{"rule": {}}}}}}}},
			args: args{ctx: context.Background(), project: "1", ruleName: "rule", value: config.EventingRule{}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{errors.New("error setting eventing module config")},
				},
			},
			wantErr: true,
		},
		{
			name: "unable to set project",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{Rules: map[string]config.EventingRule{"rule": {}}}}}}}},
			args: args{ctx: context.Background(), project: "1", ruleName: "rule", value: config.EventingRule{}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{mock.Anything, mock.Anything},
					paramsReturned: []interface{}{errors.New("Invalid config file type")},
				},
			},
			wantErr: true,
		},
		{
			name: "eventing rules are set",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{Rules: map[string]config.EventingRule{"rule": {}}}}}}}},
			args: args{ctx: context.Background(), project: "1", ruleName: "rule", value: config.EventingRule{}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
		},
		{
			name: "eventing config is not set when config rules are nil",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{}}}}}},
			args: args{ctx: context.Background(), project: "1", ruleName: "rule", value: config.EventingRule{}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{errors.New("error setting eventing module config")},
				},
			},
			wantErr: true,
		},
		{
			name: "eventing rules are not set when config rules are nil",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{}}}}}},
			args: args{ctx: context.Background(), project: "1", ruleName: "rule", value: config.EventingRule{}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{mock.Anything, mock.Anything},
					paramsReturned: []interface{}{errors.New("Invalid config file type")},
				},
			},
			wantErr: true,
		},
		{
			name: "eventing rules are set when config rules are nil",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{}}}}}},
			args: args{ctx: context.Background(), project: "1", ruleName: "rule", value: config.EventingRule{}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{mock.Anything, mock.Anything},
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

			if err := tt.s.SetEventingRule(tt.args.ctx, tt.args.project, tt.args.ruleName, tt.args.value); (err != nil) != tt.wantErr {
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
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{Rules: map[string]config.EventingRule{"rule": {}}}}}}}},
			args:    args{ctx: context.Background(), project: "2", ruleName: "rule"},
			wantErr: true,
		},
		{
			name: "unable to set eventing config",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{Rules: map[string]config.EventingRule{"rule": {}}}}}}}},
			args: args{ctx: context.Background(), project: "1", ruleName: "rule"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{errors.New("error setting eventing module config")},
				},
			},
			wantErr: true,
		},
		{
			name: "unable to set project",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{Rules: map[string]config.EventingRule{"rule": {}}}}}}}},
			args: args{ctx: context.Background(), project: "1", ruleName: "rule"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{mock.Anything, mock.Anything},
					paramsReturned: []interface{}{errors.New("Invalid config file type")},
				},
			},
			wantErr: true,
		},
		{
			name: "eventing rule deleted succesfully",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{Rules: map[string]config.EventingRule{"rule": {}}}}}}}},
			args: args{ctx: context.Background(), project: "1", ruleName: "rule"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{mock.Anything, mock.Anything},
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

			if err := tt.s.SetDeleteEventingRule(tt.args.ctx, tt.args.project, tt.args.ruleName); (err != nil) != tt.wantErr {
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
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{Schemas: map[string]config.SchemaObject{"evType": {ID: "evType", Schema: "type evType {id: String!}"}}}}}}}},
			args:    args{ctx: context.Background(), evType: "evType", project: "2", schema: "type evType {id: String!}"},
			wantErr: true,
		},
		{
			name: "schemas empty in config and unable to set eventing config",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{}}}}}},
			args: args{ctx: context.Background(), evType: "evType", project: "1", schema: "type evType {id: String!}"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{errors.New("error setting eventing module config")},
				},
			},
			wantErr: true,
		},
		{
			name: "unable to set eventing config",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{Schemas: map[string]config.SchemaObject{"evType": {ID: "evType", Schema: "type evType {id: String!}"}}}}}}}},
			args: args{ctx: context.Background(), evType: "evType", project: "1", schema: "type evType {id: String!}"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{errors.New("error setting eventing module config")},
				},
			},
			wantErr: true,
		},
		{
			name: "unable to set project",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{Schemas: map[string]config.SchemaObject{"evType": {ID: "evType", Schema: "type evType {id: String!}"}}}}}}}},
			args: args{ctx: context.Background(), evType: "evType", project: "1", schema: "type evType {id: String!}"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{mock.Anything, mock.Anything},
					paramsReturned: []interface{}{errors.New("Invalid config file type")},
				},
			},
			wantErr: true,
		},
		{
			name: "eventing schema is set",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{Schemas: map[string]config.SchemaObject{"evType": {ID: "evType", Schema: "type evType {id: String!}"}}}}}}}},
			args: args{ctx: context.Background(), evType: "evType", project: "1", schema: "type evType {id: String!}"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{mock.Anything, mock.Anything},
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

			if err := tt.s.SetEventingSchema(tt.args.ctx, tt.args.project, tt.args.evType, tt.args.schema); (err != nil) != tt.wantErr {
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
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{Schemas: map[string]config.SchemaObject{"evType": {ID: "evType", Schema: "type evType {id: String!}"}}}}}}}},
			args:    args{ctx: context.Background(), evType: "evType", project: "2"},
			wantErr: true,
		},
		{
			name: "unable to set eventing config",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{Schemas: map[string]config.SchemaObject{"evType": {ID: "evType", Schema: "type evType {id: String!}"}}}}}}}},
			args: args{ctx: context.Background(), evType: "evType", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{errors.New("error setting eventing module config")},
				},
			},
			wantErr: true,
		},
		{
			name: "unable to set project",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{Schemas: map[string]config.SchemaObject{"evType": {ID: "evType", Schema: "type evType {id: String!}"}}}}}}}},
			args: args{ctx: context.Background(), evType: "evType", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{mock.Anything, mock.Anything},
					paramsReturned: []interface{}{errors.New("Invalid config file type")},
				},
			},
			wantErr: true,
		},
		{
			name: "eventing schema is deleted",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{Schemas: map[string]config.SchemaObject{"evType": {ID: "evType", Schema: "type evType {id: String!}"}}}}}}}},
			args: args{ctx: context.Background(), evType: "evType", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{mock.Anything, mock.Anything},
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

			if err := tt.s.SetDeleteEventingSchema(tt.args.ctx, tt.args.project, tt.args.evType); (err != nil) != tt.wantErr {
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
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{SecurityRules: map[string]*config.Rule{"evType": {}}}}}}}},
			args:    args{ctx: context.Background(), evType: "evType", project: "2", rule: &config.Rule{}},
			wantErr: true,
		},
		{
			name: "unable to set eventing config",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{SecurityRules: map[string]*config.Rule{"evType": {}}}}}}}},
			args: args{ctx: context.Background(), evType: "evType", project: "1", rule: &config.Rule{}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{errors.New("error setting eventing module config")},
				},
			},
			wantErr: true,
		},
		{
			name: "unable to set project",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{SecurityRules: map[string]*config.Rule{"evType": {}}}}}}}},
			args: args{ctx: context.Background(), evType: "evType", project: "1", rule: &config.Rule{}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{mock.Anything, mock.Anything},
					paramsReturned: []interface{}{errors.New("Invalid config file type")},
				},
			},
			wantErr: true,
		},
		{
			name: "security rules empty in config and they are set",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{}}}}}},
			args: args{ctx: context.Background(), evType: "evType", project: "1", rule: &config.Rule{}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
		},
		{
			name: "security rules are set",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{SecurityRules: map[string]*config.Rule{"evType": {}}}}}}}},
			args: args{ctx: context.Background(), evType: "evType", project: "1", rule: &config.Rule{}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{mock.Anything, mock.Anything},
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

			if err := tt.s.SetEventingSecurityRules(tt.args.ctx, tt.args.project, tt.args.evType, tt.args.rule); (err != nil) != tt.wantErr {
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
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{SecurityRules: map[string]*config.Rule{"evType": {}}}}}}}},
			args:    args{ctx: context.Background(), evType: "evType", project: "2"},
			wantErr: true,
		},
		{
			name: "unable to set eventing config",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{SecurityRules: map[string]*config.Rule{"evType": {}}}}}}}},
			args: args{ctx: context.Background(), evType: "evType", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{errors.New("error setting eventing module config")},
				},
			},
			wantErr: true,
		},
		{
			name: "unable to set project",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{SecurityRules: map[string]*config.Rule{"evType": {}}}}}}}},
			args: args{ctx: context.Background(), evType: "evType", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{mock.Anything, mock.Anything},
					paramsReturned: []interface{}{errors.New("Invalid config file type")},
				},
			},
			wantErr: true,
		},
		{
			name: "eventing security rules deleted",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{SecurityRules: map[string]*config.Rule{"evType": {}}}}}}}},
			args: args{ctx: context.Background(), evType: "evType", project: "1"},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetEventingConfig",
					args:           []interface{}{"1", mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{mock.Anything, mock.Anything},
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

			if err := tt.s.SetDeleteEventingSecurityRules(tt.args.ctx, tt.args.project, tt.args.evType); (err != nil) != tt.wantErr {
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
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{Rules: map[string]config.EventingRule{"rule": {}}}}}}}},
			args:    args{ctx: context.Background(), id: "rule", project: "2"},
			wantErr: true,
		},
		{
			name:    "id not present in config",
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{Rules: map[string]config.EventingRule{"rule": {}}}}}}}},
			args:    args{ctx: context.Background(), id: "notRule", project: "1"},
			wantErr: true,
		},
		{
			name: "got trigger rule",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{Rules: map[string]config.EventingRule{"rule": {}}}}}}}},
			args: args{ctx: context.Background(), id: "rule", project: "1"},
			want: []interface{}{config.EventingRule{}},
		},
		{
			name: "id is empty and got all trigger rules",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{Rules: map[string]config.EventingRule{"rule": {}}}}}}}},
			args: args{ctx: context.Background(), id: "*", project: "1"},
			want: []interface{}{config.EventingRule{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.GetEventingTriggerRules(tt.args.ctx, tt.args.project, tt.args.id)
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
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{Schemas: map[string]config.SchemaObject{"id": {}}}}}}}},
			args:    args{ctx: context.Background(), id: "id", project: "2"},
			wantErr: true,
		},
		{
			name:    "id not present in config",
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{Schemas: map[string]config.SchemaObject{"id": {}}}}}}}},
			args:    args{ctx: context.Background(), id: "notId", project: "1"},
			wantErr: true,
		},
		{
			name: "got schema",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{Schemas: map[string]config.SchemaObject{"id": {}}}}}}}},
			args: args{ctx: context.Background(), id: "id", project: "1"},
			want: []interface{}{config.SchemaObject{}},
		},
		{
			name: "id empty and got schemas",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{Schemas: map[string]config.SchemaObject{"id": {}}}}}}}},
			args: args{ctx: context.Background(), id: "*", project: "1"},
			want: []interface{}{config.SchemaObject{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.GetEventingSchema(tt.args.ctx, tt.args.project, tt.args.id)
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
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{SecurityRules: map[string]*config.Rule{"id": {}}}}}}}},
			args:    args{ctx: context.Background(), id: "id", project: "2"},
			wantErr: true,
		},
		{
			name:    "id not present in config",
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{SecurityRules: map[string]*config.Rule{"id": {}}}}}}}},
			args:    args{ctx: context.Background(), id: "notId", project: "1"},
			wantErr: true,
		},
		{
			name: "got security rule",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{SecurityRules: map[string]*config.Rule{"id": {}}}}}}}},
			args: args{ctx: context.Background(), id: "id", project: "1"},
			want: []interface{}{&config.Rule{}},
		},
		{
			name: "id empty and got security rules",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Eventing: config.Eventing{SecurityRules: map[string]*config.Rule{"id": {}}}}}}}},
			args: args{ctx: context.Background(), id: "*", project: "1"},
			want: []interface{}{&config.Rule{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.GetEventingSecurityRules(tt.args.ctx, tt.args.project, tt.args.id)
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
