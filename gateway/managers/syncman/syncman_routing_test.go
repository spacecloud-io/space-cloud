package syncman

import (
	"context"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

func TestManager_GetProjectRoutes(t *testing.T) {
	type args struct {
		ctx     context.Context
		project string
	}
	tests := []struct {
		name    string
		s       *Manager
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Project config not found",
			s: &Manager{
				clusterID: "chicago",
				projectConfig: &config.Config{
					Projects: config.Projects{
						"myproject": &config.Project{
							ProjectConfig: &config.ProjectConfig{ID: "myproject"},
						},
					},
				},
			},
			args:    args{ctx: context.Background(), project: "test"},
			wantErr: true,
		},
		{
			name: "Get all project routes",
			s: &Manager{
				clusterID: "chicago",
				projectConfig: &config.Config{
					Projects: config.Projects{
						"myproject": &config.Project{
							ProjectConfig: &config.ProjectConfig{ID: "myproject"},
							IngressRoutes: config.IngressRoutes{
								"resourceID1": &config.Route{ID: "route1"},
							},
						},
					},
				},
			},
			args: args{ctx: context.Background(), project: "myproject"},
			want: config.Routes{{ID: "route1"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got, err := tt.s.GetProjectRoutes(tt.args.ctx, tt.args.project)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.GetProjectRoutes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.GetProjectRoutes() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
			name: "Project config not found",
			s: &Manager{
				clusterID: "chicago",
				projectConfig: &config.Config{
					Projects: config.Projects{
						"myproject": &config.Project{
							ProjectConfig: &config.ProjectConfig{ID: "myproject"},
						},
					},
				},
			},
			args:    args{ctx: context.Background(), project: "test", routeID: "route1"},
			wantErr: true,
		},
		{
			name: "Get all project routes",
			s: &Manager{
				clusterID: "chicago",
				projectConfig: &config.Config{
					Projects: config.Projects{
						"myproject": &config.Project{
							ProjectConfig: &config.ProjectConfig{ID: "myproject"},
							IngressRoutes: config.IngressRoutes{
								"resourceID1": &config.Route{ID: "route1"},
							},
						},
					},
				},
			},
			args: args{ctx: context.Background(), project: "myproject", routeID: "*"},
			want: []interface{}{&config.Route{ID: "route1"}},
		},
		{
			name: "Get specific project route",
			s: &Manager{
				clusterID: "chicago",
				projectConfig: &config.Config{
					Projects: config.Projects{
						"myproject": &config.Project{
							ProjectConfig: &config.ProjectConfig{ID: "myproject"},
							IngressRoutes: config.IngressRoutes{
								config.GenerateResourceID("chicago", "myproject", config.ResourceIngressRoute, "route1"): &config.Route{ID: "route1"},
								config.GenerateResourceID("chicago", "myproject", config.ResourceIngressRoute, "route2"): &config.Route{ID: "route2"},
							},
						},
					},
				},
			},
			args: args{ctx: context.Background(), project: "myproject", routeID: "route1"},
			want: []interface{}{&config.Route{ID: "route1"}},
		},
		{
			name: "Throw error when specific route id is asked but it doesn't exists in config",
			s: &Manager{
				clusterID: "chicago",
				projectConfig: &config.Config{
					Projects: config.Projects{
						"myproject": &config.Project{
							ProjectConfig: &config.ProjectConfig{ID: "myproject"},
							IngressRoutes: config.IngressRoutes{
								"resourceID1": &config.Route{ID: "route1"},
								"resourceID2": &config.Route{ID: "route2"},
							},
						},
					},
				},
			},
			args:    args{ctx: context.Background(), project: "myproject", routeID: "route3"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.integrationMan = &mockIntegrationManager{skip: true}
			_, got, err := tt.s.GetIngressRouting(tt.args.ctx, tt.args.project, tt.args.routeID, model.RequestParams{})
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.GetIngressRouting() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.GetIngressRouting() = %v type %v, want %v type %v", got, reflect.TypeOf(got), tt.want, reflect.TypeOf(tt.want))
			}
		})
	}
}
