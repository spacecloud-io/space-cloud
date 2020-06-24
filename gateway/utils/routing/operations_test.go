package routing

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

func TestRouting_DeleteProjectRoutes(t *testing.T) {
	type fields struct {
		routes routeMapping
	}
	type args struct {
		project string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   routeMapping
	}{
		// TODO: Add test cases.
		{
			name: "deleteProject",
			fields: fields{
				routes: routeMapping{
					"test1": config.Routes{
						&config.Route{
							ID:     "1234",
							Source: config.RouteSource{},
						},
					},
					"test2": config.Routes{
						&config.Route{
							ID:     "12345",
							Source: config.RouteSource{},
						},
					},
				},
			},
			args: args{
				project: "test1",
			},
			want: routeMapping{
				"test2": config.Routes{
					&config.Route{
						ID:     "12345",
						Source: config.RouteSource{},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Routing{
				routes: tt.fields.routes,
			}
			r.DeleteProjectRoutes(tt.args.project)
			if !reflect.DeepEqual(tt.want, tt.fields.routes) {
				log.Println("Routing.deleteProjectRoutes()")

				a, _ := json.MarshalIndent(tt.fields.routes, "", " ")
				log.Printf("got= %s", string(a))

				a, _ = json.MarshalIndent(tt.want, "", " ")
				log.Printf("want = %s", string(a))
			}
		})
	}
}

func TestRouting_SetProjectRoutes(t *testing.T) {
	type fields struct {
		routes routeMapping
	}
	type args struct {
		project string
		routes  config.Routes
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   routeMapping
	}{
		{
			name: "set project rule",
			fields: fields{
				routes: routeMapping{},
			},
			args: args{
				project: "test1",
				routes: config.Routes{
					&config.Route{
						ID:     "12345",
						Source: config.RouteSource{},
					},
				},
			},
			want: routeMapping{
				"test1": config.Routes{
					&config.Route{
						ID:     "12345",
						Source: config.RouteSource{},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Routing{
				routes: tt.fields.routes,
			}
			_ = r.SetProjectRoutes(tt.args.project, tt.args.routes)
			if !reflect.DeepEqual(tt.fields.routes, tt.want) {
				log.Println("Routing.addProjectRoutes()")

				a, _ := json.MarshalIndent(tt.fields.routes, "", " ")
				log.Printf("got= %s", string(a))

				a, _ = json.MarshalIndent(tt.want, "", " ")
				log.Printf("want = %s", string(a))
			}
		})
	}
}
