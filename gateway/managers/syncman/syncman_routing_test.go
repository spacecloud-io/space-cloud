package syncman

import (
	"context"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

func TestManager_GetIngressRouting(t *testing.T) {
	type args struct {
		ctx     context.Context
		project string
		routeID string
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
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Routes: config.Routes{{ID: "1"}, {ID: "2"}}}}}}},
			args:    args{ctx: context.Background(), project: "2", routeID: "1"},
			wantErr: true,
		},
		{
			name: "routeID empty and got all routes",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Routes: config.Routes{{ID: "1"}}}}}}},
			args: args{ctx: context.Background(), project: "1", routeID: "*"},
			want: []interface{}{&config.Route{ID: "1"}},
		},
		{
			name:    "route id not present in config",
			s:       &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Routes: config.Routes{{ID: "1"}}}}}}},
			args:    args{ctx: context.Background(), project: "1", routeID: "2"},
			wantErr: true,
		},
		{
			name: "got ingress route",
			s:    &Manager{projectConfig: &config.Config{Projects: []*config.Project{{ID: "1", Modules: &config.Modules{Routes: config.Routes{{ID: "1"}}}}}}},
			args: args{ctx: context.Background(), project: "1", routeID: "1"},
			want: []interface{}{&config.Route{ID: "1"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.integrationMan = &mockIntegrationManager{skip: true}
			_, got, err := tt.s.GetIngressRouting(tt.args.ctx, tt.args.project, tt.args.routeID, model.RequestParams{})
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.GetIngressRouting() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.GetIngressRouting() = %v, want %v", got, tt.want)
			}
		})
	}
}
