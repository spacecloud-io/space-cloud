package syncman

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"

	"github.com/spaceuptech/space-cloud/gateway/config"
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
