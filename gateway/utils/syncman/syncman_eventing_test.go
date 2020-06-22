package syncman

import (
	"context"
	"errors"
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
