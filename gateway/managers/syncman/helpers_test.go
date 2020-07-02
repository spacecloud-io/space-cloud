package syncman

import (
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

func Test_calcTokens(t *testing.T) {
	type args struct {
		n      int
		tokens int
		i      int
	}
	tests := []struct {
		name      string
		args      args
		wantStart int
		wantEnd   int
	}{
		{name: "test1", args: args{n: 7, tokens: 100, i: 0}, wantStart: 0, wantEnd: 14},
		{name: "test2", args: args{n: 7, tokens: 100, i: 4}, wantStart: 60, wantEnd: 74},
		{name: "test3", args: args{n: 7, tokens: 100, i: 5}, wantStart: 75, wantEnd: 89},
		{name: "test4", args: args{n: 7, tokens: 100, i: 6}, wantStart: 90, wantEnd: 99},
		{name: "test5", args: args{n: 1, tokens: 100, i: 0}, wantStart: 0, wantEnd: 99},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStart, gotEnd := calcTokens(tt.args.n, tt.args.tokens, tt.args.i)
			if gotStart != tt.wantStart {
				t.Errorf("calcTokens() gotStart = %v, want %v", gotStart, tt.wantStart)
			}
			if gotEnd != tt.wantEnd {
				t.Errorf("calcTokens() gotEnd = %v, want %v", gotEnd, tt.wantEnd)
			}
		})
	}
}

func TestManager_setProjectConfig(t *testing.T) {
	type args struct {
		conf *config.Project
	}
	tests := []struct {
		name string
		s    *Manager
		args args
		want []*config.Project
	}{
		{
			name: "empty config passed as parameter",
			s:    &Manager{projectConfig: config.GenerateEmptyConfig()},
			args: args{conf: &config.Project{}},
			want: []*config.Project{{}},
		},
		{
			name: "config id doesn't match project id",
			s:    &Manager{projectConfig: config.GenerateEmptyConfig()},
			args: args{conf: &config.Project{ID: "someID"}},
			want: []*config.Project{{ID: "someID"}},
		},
		{
			name: "config id matches a project id",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "notSomeID"}, {ID: "someID"}}}},
			args: args{conf: &config.Project{ID: "someID"}},
			want: []*config.Project{{ID: "notSomeID"}, {ID: "someID"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.setProjectConfig(tt.args.conf)

			if !reflect.DeepEqual(tt.s.projectConfig.Projects, tt.want) {
				t.Errorf("Got: %v, Want: %v", tt.s.projectConfig.Projects, tt.want)
			}
		})
	}
}

func Test_remove(t *testing.T) {
	type args struct {
		s []*config.Project
		i int
	}
	tests := []struct {
		name string
		args args
		want []*config.Project
	}{
		{
			name: "project is removed",
			args: args{s: []*config.Project{{ID: "id1"}, {ID: "id2"}, {ID: "id3"}}, i: 1},
			want: []*config.Project{{ID: "id1"}, {ID: "id3"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := remove(tt.args.s, tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("remove() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_delete(t *testing.T) {
	type args struct {
		projectID string
	}
	tests := []struct {
		name string
		s    *Manager
		args args
		want []*config.Project
	}{
		{
			name: "project ID does not match",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "id1"}, {ID: "id2"}}}},
			args: args{projectID: "notMatching"},
			want: []*config.Project{{ID: "id1"}, {ID: "id2"}},
		},
		{
			name: "project ID does matches and project config is deleted",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "id1"}, {ID: "id2"}}}},
			args: args{projectID: "id1"},
			want: []*config.Project{{ID: "id2"}, {ID: "id2"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.delete(tt.args.projectID)

			if !reflect.DeepEqual(tt.s.projectConfig.Projects, tt.want) {
				t.Errorf("Got: %v, Want: %v", tt.s.projectConfig.Projects, tt.want)
			}
		})
	}
}

func Test_scServices_Len(t *testing.T) {
	tests := []struct {
		name string
		a    scServices
		want int
	}{
		{
			name: "length is returned",
			a:    scServices{{id: "id1"}, {id: "id2"}},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Len(); got != tt.want {
				t.Errorf("scServices.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_scServices_Swap(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		a    scServices
		args args
		want scServices
	}{
		{
			name: "swap successful",
			a:    scServices{{id: "id1"}, {id: "id2"}},
			args: args{0, 1},
			want: scServices{{id: "id2"}, {id: "id1"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.a.Swap(tt.args.i, tt.args.j)
			if !reflect.DeepEqual(tt.a, tt.want) {
				t.Errorf("Got: %v, Want: %v", tt.a, tt.want)
			}
		})
	}
}

func Test_scServices_Less(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		a    scServices
		args args
		want bool
	}{
		{
			name: "true condition",
			a:    scServices{{id: "1"}, {id: "2"}},
			args: args{0, 1},
			want: true,
		},
		{
			name: "false condition",
			a:    scServices{{id: "2"}, {id: "1"}},
			args: args{0, 1},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Less(tt.args.i, tt.args.j); got != tt.want {
				t.Errorf("scServices.Less() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calcIndex(t *testing.T) {
	type args struct {
		token       int
		totalTokens int
		n           int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "index calculated",
			args: args{token: 100, totalTokens: 1000, n: 100},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calcIndex(tt.args.token, tt.args.totalTokens, tt.args.n); got != tt.want {
				t.Errorf("calcIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_GetGatewayIndex(t *testing.T) {
	tests := []struct {
		name string
		s    *Manager
		want int
	}{
		{
			name: "services is empty",
			s:    &Manager{services: []*service{}},
			want: 0,
		},
		{
			name: "service id does not match node id",
			s:    &Manager{services: []*service{{id: "1"}}, nodeID: "not1"},
			want: 0,
		},
		{
			name: "got gateway index",
			s:    &Manager{services: []*service{{id: "0"}, {id: "1"}}, nodeID: "1"},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.GetGatewayIndex(); got != tt.want {
				t.Errorf("Manager.GetGatewayIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_GetNodeID(t *testing.T) {
	tests := []struct {
		name string
		s    *Manager
		want string
	}{
		{
			name: "nodeID returned",
			s:    &Manager{nodeID: "nodeID"},
			want: "nodeID",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.GetNodeID(); got != tt.want {
				t.Errorf("Manager.GetNodeID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_getConfigWithoutLock(t *testing.T) {
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
			name:    "project not present in the state",
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "id1"}, {ID: "id2"}}}},
			args:    args{projectID: "someID"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "project id matches",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "id1"}, {ID: "id2"}}}},
			args: args{projectID: "id1"},
			want: &config.Project{ID: "id1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.getConfigWithoutLock(tt.args.projectID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.getConfigWithoutLock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.getConfigWithoutLock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_GetSpaceCloudURLFromID(t *testing.T) {
	type args struct {
		nodeID string
	}
	tests := []struct {
		name    string
		s       *Manager
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "service does not exist with given node id",
			s:       &Manager{services: []*service{{id: "someID", addr: "someAddr"}}},
			args:    args{nodeID: "notSomeID"},
			want:    "",
			wantErr: true,
		},
		{
			name: "got space cloud url from id",
			s:    &Manager{services: []*service{{id: "someID", addr: "spacecloud.com"}}},
			args: args{nodeID: "someID"},
			want: "spacecloud.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.GetSpaceCloudURLFromID(tt.args.nodeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.GetSpaceCloudURLFromID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Manager.GetSpaceCloudURLFromID() = %v, want %v", got, tt.want)
			}
		})
	}
}
