package routing

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

func TestRouting_DeleteProjectRoutes(t *testing.T) {
	type fields struct {
		routes config.Routes
	}
	type args struct {
		project string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   config.Routes
	}{
		// TODO: Add test cases.
		{
			name: "deleteProject",
			fields: fields{
				routes: config.Routes{
					&config.Route{
						ID:      "1234",
						Project: "test1",
						Source:  config.RouteSource{},
					},
					&config.Route{
						ID:      "12345",
						Project: "test1",
						Source:  config.RouteSource{},
					},
				},
			},
			args: args{
				project: "test1",
			},
			want: config.Routes{
				&config.Route{
					ID:      "12345",
					Project: "test1",
					Source:  config.RouteSource{},
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
				helpers.Logger.LogInfo(helpers.GetRequestID(nil), "Routing.deleteProjectRoutes()", nil)

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
		routes config.Routes
	}
	type args struct {
		project string
		routes  config.Routes
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   config.Routes
	}{
		{
			name: "set project rule",
			fields: fields{
				routes: config.Routes{},
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
			want: config.Routes{
				&config.Route{
					ID:      "12345",
					Project: "test1",
					Source:  config.RouteSource{},
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
				helpers.Logger.LogInfo(helpers.GetRequestID(nil), "Routing.addProjectRoutes()", nil)

				a, _ := json.MarshalIndent(tt.fields.routes, "", " ")
				log.Printf("got= %s", string(a))

				a, _ = json.MarshalIndent(tt.want, "", " ")
				log.Printf("want = %s", string(a))
			}
		})
	}
}
