package routing

import (
	"encoding/json"
	"log"
	"reflect"
	"sync"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

func TestRouting_DeleteProjectRoutes(t *testing.T) {
	type fields struct {
		lock   sync.RWMutex
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
				lock: sync.RWMutex{},
				routes: routeMapping{
					"test1": config.Routes{
						&config.Route{
							Id:          "1234",
							Source:      config.RouteSource{},
							Destination: config.RouteDestination{},
						},
					},
					"test2": config.Routes{
						&config.Route{
							Id:          "12345",
							Source:      config.RouteSource{},
							Destination: config.RouteDestination{},
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
						Id:          "12345",
						Source:      config.RouteSource{},
						Destination: config.RouteDestination{},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Routing{
				lock:   tt.fields.lock,
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
		lock   sync.RWMutex
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
		// TODO: Add test cases.
		{
			name: "set project rule",
			fields: fields{
				lock:   sync.RWMutex{},
				routes: routeMapping{},
			},
			args: args{
				project: "test1",
				routes: config.Routes{
					&config.Route{
						Id:          "12345",
						Source:      config.RouteSource{},
						Destination: config.RouteDestination{},
					},
				},
			},
			want: routeMapping{
				"test1": config.Routes{
					&config.Route{
						Id:          "12345",
						Source:      config.RouteSource{},
						Destination: config.RouteDestination{},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Routing{
				lock:   tt.fields.lock,
				routes: tt.fields.routes,
			}
			r.SetProjectRoutes(tt.args.project, tt.args.routes)
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
