package syncman

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

func TestManager_GetEventSource(t *testing.T) {
	tests := []struct {
		name string
		s    *Manager
		want string
	}{
		{
			name: "got event source",
			s:    &Manager{nodeID: "nodeID"},
			want: "sc-nodeID",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.GetEventSource(); got != tt.want {
				t.Errorf("Manager.GetEventSource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_GetClusterID(t *testing.T) {
	tests := []struct {
		name string
		s    *Manager
		want string
	}{
		{
			name: "got cluster id",
			s:    &Manager{clusterID: "clusterID"},
			want: "clusterID",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.GetClusterID(); got != tt.want {
				t.Errorf("Manager.GetClusterID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_GetNodesInCluster(t *testing.T) {
	tests := []struct {
		name string
		s    *Manager
		want int
	}{
		{
			name: "length of services is 0",
			s:    &Manager{services: []*service{}},
			want: 1,
		},
		{
			name: "length of services is returned",
			s:    &Manager{services: []*service{{id: "1"}, {id: "2"}}},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.GetNodesInCluster(); got != tt.want {
				t.Errorf("Manager.GetNodesInCluster() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_GetAssignedSpaceCloudURL(t *testing.T) {
	type args struct {
		ctx     context.Context
		project string
		token   int
	}
	tests := []struct {
		name    string
		s       *Manager
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "got assigned space cloud url",
			s:    &Manager{storeType: "kube", services: []*service{{id: "1", addr: "some.com"}}},
			args: args{ctx: context.Background(), project: "project", token: 0},
			want: "http://some.com/v1/api/project/eventing/process",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.GetAssignedSpaceCloudURL(tt.args.ctx, tt.args.project, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.GetAssignedSpaceCloudURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Manager.GetAssignedSpaceCloudURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_GetSpaceCloudNodeURLs(t *testing.T) {
	type args struct {
		project string
	}
	tests := []struct {
		name string
		s    *Manager
		args args
		want []string
	}{
		{
			name: "got space cloud urls",
			s:    &Manager{services: []*service{{id: "1", addr: "some.com"}}},
			args: args{project: "project"},
			want: []string{"http://some.com/v1/api/project/realtime/process"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.GetSpaceCloudNodeURLs(tt.args.project); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.GetSpaceCloudNodeURLs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_GetRealtimeURL(t *testing.T) {
	type args struct {
		project string
	}
	tests := []struct {
		name string
		s    *Manager
		args args
		want string
	}{
		{
			name: "got realtime url",
			s:    &Manager{port: 4122},
			args: args{project: "project"},
			want: "http://localhost:4122/v1/api/project/realtime/handle",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.GetRealtimeURL(tt.args.project); got != tt.want {
				t.Errorf("Manager.GetRealtimeURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_GetAssignedTokens(t *testing.T) {
	tests := []struct {
		name      string
		s         *Manager
		wantStart int
		wantEnd   int
	}{
		{
			name:      "got assigned tokens",
			s:         &Manager{storeType: "kube", services: []*service{{id: "1"}, {id: "2"}}},
			wantStart: 0,
			wantEnd:   49,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStart, gotEnd := tt.s.GetAssignedTokens()
			if gotStart != tt.wantStart {
				t.Errorf("Manager.GetAssignedTokens() gotStart = %v, want %v", gotStart, tt.wantStart)
			}
			if gotEnd != tt.wantEnd {
				t.Errorf("Manager.GetAssignedTokens() gotEnd = %v, want %v", gotEnd, tt.wantEnd)
			}
		})
	}
}

func TestManager_ApplyProjectConfig(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx     context.Context
		project *config.Project
	}
	tests := []struct {
		name            string
		s               *Manager
		args            args
		adminMockArgs   []mockArgs
		modulesMockArgs []mockArgs
		storeMockArgs   []mockArgs
		want            int
		wantErr         bool
	}{
		{
			name: "sync operation is not valid",
			s:    &Manager{projectConfig: config.GenerateEmptyConfig()},
			args: args{ctx: context.Background(), project: &config.Project{}},
			adminMockArgs: []mockArgs{
				{
					method:         "ValidateSyncOperation",
					args:           []interface{}{config.GenerateEmptyConfig(), &config.Project{}},
					paramsReturned: []interface{}{false},
				},
			},
			want:    http.StatusInternalServerError,
			wantErr: true,
		},
		{
			name: "could not get internal access token",
			s:    &Manager{projectConfig: config.GenerateEmptyConfig()},
			args: args{ctx: context.Background(), project: &config.Project{}},
			adminMockArgs: []mockArgs{
				{
					method:         "ValidateSyncOperation",
					args:           []interface{}{config.GenerateEmptyConfig(), &config.Project{}},
					paramsReturned: []interface{}{true},
				},
				{
					method:         "GetInternalAccessToken",
					paramsReturned: []interface{}{"", errors.New("could not generate signed string for token")},
				},
			},
			want:    http.StatusInternalServerError,
			wantErr: true,
		},
		{
			name: "project exists already and store type kube and can not set project",
			s:    &Manager{storeType: "kube", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1"}}}},
			args: args{ctx: context.Background(), project: &config.Project{ID: "1"}},
			adminMockArgs: []mockArgs{
				{
					method:         "ValidateSyncOperation",
					args:           []interface{}{mock.Anything, mock.Anything},
					paramsReturned: []interface{}{true},
				},
				{
					method:         "GetInternalAccessToken",
					paramsReturned: []interface{}{"token", nil},
				},
			},
			modulesMockArgs: []mockArgs{
				{
					method: "SetProjectConfig",
					args:   []interface{}{mock.Anything, mock.Anything, mock.Anything},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{context.Background(), mock.Anything},
					paramsReturned: []interface{}{errors.New("error marshalling project config")},
				},
			},
			want:    http.StatusInternalServerError,
			wantErr: true,
		},
		{
			name: "project doesn't exist and store type kube and can not set project",
			s:    &Manager{storeType: "kube", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1"}}}},
			args: args{ctx: context.Background(), project: &config.Project{ID: "2"}},
			adminMockArgs: []mockArgs{
				{
					method:         "ValidateSyncOperation",
					args:           []interface{}{mock.Anything, mock.Anything},
					paramsReturned: []interface{}{true},
				},
				{
					method:         "GetInternalAccessToken",
					paramsReturned: []interface{}{"token", nil},
				},
			},
			modulesMockArgs: []mockArgs{
				{
					method: "SetProjectConfig",
					args:   []interface{}{mock.Anything, mock.Anything, mock.Anything},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{context.Background(), mock.Anything},
					paramsReturned: []interface{}{errors.New("error marshalling project config")},
				},
			},
			want:    http.StatusInternalServerError,
			wantErr: true,
		},
		{
			name: "project exists already and store type kube and project is set",
			s:    &Manager{storeType: "kube", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1"}}}},
			args: args{ctx: context.Background(), project: &config.Project{ID: "1"}},
			adminMockArgs: []mockArgs{
				{
					method:         "ValidateSyncOperation",
					args:           []interface{}{mock.Anything, mock.Anything},
					paramsReturned: []interface{}{true},
				},
				{
					method:         "GetInternalAccessToken",
					paramsReturned: []interface{}{"token", nil},
				},
			},
			modulesMockArgs: []mockArgs{
				{
					method: "SetProjectConfig",
					args:   []interface{}{mock.Anything, mock.Anything, mock.Anything},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{context.Background(), mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "project doesn't exist and store type kube and project is set",
			s:    &Manager{storeType: "kube", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1"}}}},
			args: args{ctx: context.Background(), project: &config.Project{ID: "2"}},
			adminMockArgs: []mockArgs{
				{
					method:         "ValidateSyncOperation",
					args:           []interface{}{mock.Anything, mock.Anything},
					paramsReturned: []interface{}{true},
				},
				{
					method:         "GetInternalAccessToken",
					paramsReturned: []interface{}{"token", nil},
				},
			},
			modulesMockArgs: []mockArgs{
				{
					method: "SetProjectConfig",
					args:   []interface{}{mock.Anything, mock.Anything, mock.Anything},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{context.Background(), mock.Anything},
					paramsReturned: []interface{}{nil},
				},
			},
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockAdmin := mockAdminSyncmanInterface{}
			mockModules := mockModulesInterface{}
			mockStore := mockStoreInterface{}

			for _, m := range tt.adminMockArgs {
				mockAdmin.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.modulesMockArgs {
				mockModules.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.storeMockArgs {
				mockStore.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			tt.s.adminMan = &mockAdmin
			tt.s.modules = &mockModules
			tt.s.store = &mockStore

			got, err := tt.s.ApplyProjectConfig(tt.args.ctx, tt.args.project)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.ApplyProjectConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Manager.ApplyProjectConfig() = %v, want %v", got, tt.want)
			}

			mockAdmin.AssertExpectations(t)
			mockModules.AssertExpectations(t)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestManager_setProject(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx     context.Context
		project *config.Project
	}
	tests := []struct {
		name          string
		s             *Manager
		args          args
		storeMockArgs []mockArgs
		wantErr       bool
	}{
		{
			name: "store type kube and couldn't set project",
			s:    &Manager{storeType: "kube", projectConfig: &config.Config{Projects: []*config.Project{{ID: "2"}}}},
			args: args{ctx: context.Background(), project: &config.Project{ID: "1"}},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{mock.Anything, mock.Anything},
					paramsReturned: []interface{}{errors.New("could not marshall project config")},
				},
			},
			wantErr: true,
		},
		{
			name: "store type kube and project is set",
			s:    &Manager{storeType: "kube", projectConfig: &config.Config{Projects: []*config.Project{{ID: "2"}}}},
			args: args{ctx: context.Background(), project: &config.Project{ID: "1"}},
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

			mockStore := mockStoreInterface{}

			for _, m := range tt.storeMockArgs {
				mockStore.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			tt.s.store = &mockStore

			if err := tt.s.setProject(tt.args.ctx, tt.args.project); (err != nil) != tt.wantErr {
				t.Errorf("Manager.setProject() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestManager_SetProjectGlobalConfig(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx     context.Context
		project *config.Project
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
			name: "couldn't set global config",
			s:    &Manager{storeType: "kube"},
			args: args{ctx: context.Background(), project: &config.Project{Secrets: []*config.Secret{{Secret: "secret"}}, AESKey: "aeskey", ContextTimeGraphQL: 10, DockerRegistry: "registry", ID: "1", Name: "name"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetGlobalConfig",
					args:           []interface{}{"name", []*config.Secret{{Secret: "secret"}}, "aeskey"},
					paramsReturned: []interface{}{errors.New("couldn't decode aeskey")},
				},
			},
			wantErr: true,
		},
		{
			name: "could not get config",
			s:    &Manager{storeType: "kube", projectConfig: &config.Config{Projects: []*config.Project{{ID: "2"}}}},
			args: args{ctx: context.Background(), project: &config.Project{Secrets: []*config.Secret{{Secret: "secret"}}, AESKey: "aeskey", ContextTimeGraphQL: 10, DockerRegistry: "registry", ID: "1", Name: "name"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetGlobalConfig",
					args:           []interface{}{"name", []*config.Secret{{Secret: "secret"}}, "aeskey"},
					paramsReturned: []interface{}{nil},
				},
			},
			wantErr: true,
		},
		{
			name: "global project config is set",
			s:    &Manager{storeType: "kube", projectConfig: &config.Config{Projects: []*config.Project{{ID: "1"}}}},
			args: args{ctx: context.Background(), project: &config.Project{Secrets: []*config.Secret{{Secret: "secret"}}, AESKey: "aeskey", ContextTimeGraphQL: 10, DockerRegistry: "registry", ID: "1", Name: "name"}},
			modulesMockArgs: []mockArgs{
				{
					method:         "SetGlobalConfig",
					args:           []interface{}{"name", []*config.Secret{{Secret: "secret"}}, "aeskey"},
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

			if err := tt.s.SetProjectGlobalConfig(tt.args.ctx, tt.args.project); (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetProjectGlobalConfig() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockModules.AssertExpectations(t)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestManager_SetProjectConfig(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx     context.Context
		project *config.Project
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
			name: "project config is not set",
			s:    &Manager{storeType: "kube", projectConfig: config.GenerateEmptyConfig()},
			args: args{ctx: context.Background(), project: &config.Project{ID: "1"}},
			modulesMockArgs: []mockArgs{
				{
					method: "SetProjectConfig",
					args:   []interface{}{mock.Anything, mock.Anything, mock.Anything},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "SetProject",
					args:           []interface{}{mock.Anything, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to get db config")},
				},
			},
			wantErr: true,
		},
		{
			name: "project config is set",
			s:    &Manager{storeType: "kube", projectConfig: config.GenerateEmptyConfig()},
			args: args{ctx: context.Background(), project: &config.Project{ID: "1"}},
			modulesMockArgs: []mockArgs{
				{
					method: "SetProjectConfig",
					args:   []interface{}{mock.Anything, mock.Anything, mock.Anything},
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

			if err := tt.s.SetProjectConfig(tt.args.ctx, tt.args.project); (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetProjectConfig() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockModules.AssertExpectations(t)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestManager_DeleteProjectConfig(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		ctx       context.Context
		projectID string
	}
	tests := []struct {
		name            string
		s               *Manager
		args            args
		adminMockArgs   []mockArgs
		modulesMockArgs []mockArgs
		storeMockArgs   []mockArgs
		wantErr         bool
	}{
		{
			name: "could not get internal access token",
			s:    &Manager{storeType: "local"},
			args: args{ctx: context.Background(), projectID: "project"},
			adminMockArgs: []mockArgs{
				{
					method:         "GetInternalAccessToken",
					paramsReturned: []interface{}{"", errors.New("could not generate signed string for token")},
				},
			},
			wantErr: true,
		},
		{
			name: "store type kube and couldn't delete project",
			s:    &Manager{storeType: "kube", projectConfig: &config.Config{Projects: []*config.Project{{ID: "notProject"}}}},
			args: args{ctx: context.Background(), projectID: "project"},
			adminMockArgs: []mockArgs{
				{
					method:         "GetInternalAccessToken",
					paramsReturned: []interface{}{"token", nil},
				},
			},
			modulesMockArgs: []mockArgs{
				{
					method: "SetProjectConfig",
					args:   []interface{}{mock.Anything, mock.Anything, mock.Anything},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "DeleteProject",
					args:           []interface{}{mock.Anything, "project"},
					paramsReturned: []interface{}{errors.New("unable to get config map")},
				},
			},
			wantErr: true,
		},
		{
			name: "store type kube and project is deleted",
			s:    &Manager{storeType: "kube", projectConfig: &config.Config{Projects: []*config.Project{{ID: "notProject"}}}},
			args: args{ctx: context.Background(), projectID: "project"},
			adminMockArgs: []mockArgs{
				{
					method:         "GetInternalAccessToken",
					paramsReturned: []interface{}{"token", nil},
				},
			},
			modulesMockArgs: []mockArgs{
				{
					method: "SetProjectConfig",
					args:   []interface{}{mock.Anything, mock.Anything, mock.Anything},
				},
			},
			storeMockArgs: []mockArgs{
				{
					method:         "DeleteProject",
					args:           []interface{}{mock.Anything, "project"},
					paramsReturned: []interface{}{nil},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockAdmin := mockAdminSyncmanInterface{}
			mockModules := mockModulesInterface{}
			mockStore := mockStoreInterface{}

			for _, m := range tt.adminMockArgs {
				mockAdmin.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.modulesMockArgs {
				mockModules.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.storeMockArgs {
				mockStore.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			tt.s.adminMan = &mockAdmin
			tt.s.modules = &mockModules
			tt.s.store = &mockStore

			if err := tt.s.DeleteProjectConfig(tt.args.ctx, tt.args.projectID); (err != nil) != tt.wantErr {
				t.Errorf("Manager.DeleteProjectConfig() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockAdmin.AssertExpectations(t)
			mockModules.AssertExpectations(t)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestManager_GetProjectConfig(t *testing.T) {
	type args struct {
		projectID string
	}
	tests := []struct {
		name    string
		s       *Manager
		args    args
		want    []interface{}
		wantErr bool
	}{
		{
			name:    "project not present in state",
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1"}, {ID: "2"}}}},
			args:    args{projectID: "3"},
			want:    []interface{}{},
			wantErr: true,
		},
		{
			name: "projectID is *",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1"}, {ID: "2"}}}},
			args: args{projectID: "*"},
			want: []interface{}{config.Project{ID: "1"}, config.Project{ID: "2"}},
		},
		{
			name: "projectID matches an existing project's id",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1"}, {ID: "2"}}}},
			args: args{projectID: "1"},
			want: []interface{}{config.Project{ID: "1"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.GetProjectConfig(tt.args.projectID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.GetProjectConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.GetProjectConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_GetConfig(t *testing.T) {
	type args struct {
		projectID string
	}
	tests := []struct {
		name    string
		s       *Manager
		args    args
		want    *config.Project
		wantErr bool
	}{
		{
			name:    "project not present in state",
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1"}}}},
			args:    args{projectID: "2"},
			wantErr: true,
		},
		{
			name: "projectID matches an existing project's ID",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1"}}}},
			args: args{projectID: "1"},
			want: &config.Project{ID: "1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.GetConfig(tt.args.projectID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.GetConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.GetConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
