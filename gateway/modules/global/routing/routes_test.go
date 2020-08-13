package routing

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

func Test_routeMapping_selectRoute(t *testing.T) {
	type args struct {
		host   string
		method string
		url    string
	}
	tests := []struct {
		name    string
		r       config.Routes
		args    args
		want    *config.Route
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "host match and url prefix match",
			r: config.Routes{
				&config.Route{
					ID:      "1234",
					Project: "test",
					Source: config.RouteSource{
						Hosts: []string{"spaceuptech.com"},
						URL:   "/abc",
						Type:  config.RoutePrefix,
					},
					Targets: []config.RouteTarget{{
						Host: "git.com",
						Port: 8080,
					}},
				},
			},
			args: args{
				host: "spaceuptech.com",
				url:  "/abc/xyz",
			},
			want: &config.Route{
				ID:      "1234",
				Project: "test",
				Source: config.RouteSource{
					Hosts: []string{"spaceuptech.com"},
					URL:   "/abc",
					Type:  config.RoutePrefix,
				},
				Targets: []config.RouteTarget{{
					Host: "git.com",
					Port: 8080,
				}},
			},
			wantErr: false,
		},
		{
			name: "does not confuses exact with prefix",
			r: config.Routes{
				&config.Route{
					ID:      "1234",
					Project: "test",
					Source: config.RouteSource{
						Hosts: []string{"spaceuptech.com"},
						URL:   "/abc/xyz",
						Type:  config.RouteExact,
					},
					Targets: []config.RouteTarget{{
						Host: "git.com",
						Port: 8080,
					}},
				},
				&config.Route{
					ID:      "567",
					Project: "test",
					Source: config.RouteSource{
						Hosts: []string{"spaceuptech.com"},
						URL:   "/abc",
						Type:  config.RoutePrefix,
					},
					Targets: []config.RouteTarget{{
						Host: "git.com",
						Port: 8080,
					}},
				},
			},
			args: args{
				host: "spaceuptech.com",
				url:  "/abc/xyz/123",
			},
			want: &config.Route{
				ID:      "567",
				Project: "test",
				Source: config.RouteSource{
					Hosts: []string{"spaceuptech.com"},
					URL:   "/abc",
					Type:  config.RoutePrefix,
				},
				Targets: []config.RouteTarget{{
					Host: "git.com",
					Port: 8080,
				}},
			},
			wantErr: false,
		},
		{
			name: "host match and url prefix match",
			r: config.Routes{
				&config.Route{
					ID:      "1234",
					Project: "test",
					Source: config.RouteSource{
						Hosts: []string{"spaceuptech.com"},
						URL:   "/abc/xyz",
						Type:  config.RouteExact,
					},
					Targets: []config.RouteTarget{{
						Host: "git.com",
						Port: 8080,
					}},
				},
			},
			args: args{
				host: "spaceuptech.com",
				url:  "/abc/xyz",
			},
			want: &config.Route{
				ID:      "1234",
				Project: "test",
				Source: config.RouteSource{
					Hosts: []string{"spaceuptech.com"},
					URL:   "/abc/xyz",
					Type:  config.RouteExact,
				},
				Targets: []config.RouteTarget{{
					Host: "git.com",
					Port: 8080,
				}},
			},
			wantErr: false,
		},

		{
			name: "host match and url match Type not provided",
			r: config.Routes{
				&config.Route{
					ID:      "1234",
					Project: "test",
					Source: config.RouteSource{
						Hosts: []string{"spaceuptech.com"},
						URL:   "/abc",
					},
					Targets: []config.RouteTarget{{
						Host: "git.com",
						Port: 8080,
					}},
				},
			},
			args: args{
				host: "spaceuptech.com",
				url:  "/abc/xyz",
			},
			want:    nil,
			wantErr: true,
		},

		{
			name: "host not match",
			r: config.Routes{
				&config.Route{
					ID:      "1234",
					Project: "test",
					Source: config.RouteSource{
						Hosts: []string{"spaceuptech.com"},
						URL:   "/abc",
						Type:  config.RoutePrefix,
					},
					Targets: []config.RouteTarget{{
						Host: "git.com",
						Port: 8080,
					}},
				},
			},
			args: args{
				host: "spaceupTech.com",
				url:  "/abc",
			},
			want:    nil,
			wantErr: true,
		},

		{
			name: "Route not present but host contains *",
			r: config.Routes{
				&config.Route{
					ID:      "1234",
					Project: "test",
					Source: config.RouteSource{
						Hosts: []string{"spaceuptech.com", "*"},
						URL:   "/abc",
						Type:  config.RoutePrefix,
					},
					Targets: []config.RouteTarget{{
						Host: "git.com",
						Port: 8080,
					}},
				},
			},
			args: args{
				host: "",
				url:  "/abc",
			},
			want: &config.Route{
				ID:      "1234",
				Project: "test",
				Source: config.RouteSource{
					Hosts: []string{"spaceuptech.com", "*"},
					URL:   "/abc",
					Type:  config.RoutePrefix,
				},
				Targets: []config.RouteTarget{{
					Host: "git.com",
					Port: 8080,
				}},
			},
			wantErr: false,
		},
		{
			name: "url not match",
			r: config.Routes{
				&config.Route{
					ID:      "1234",
					Project: "test",
					Source: config.RouteSource{
						Hosts: []string{"spaceuptech.com"},
						URL:   "/abc",
						Type:  config.RouteExact,
					},
					Targets: []config.RouteTarget{{
						Host: "git.com",
						Port: 8080,
					}},
				},
			},
			args: args{
				host: "spaceuptech.com",
				url:  "/abc/xyz",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "methods not match",
			r: config.Routes{
				&config.Route{
					ID:      "1234",
					Project: "test",
					Source: config.RouteSource{
						Hosts:   []string{"spaceuptech.com"},
						Methods: []string{http.MethodGet, http.MethodDelete},
						URL:     "/abc",
						Type:    config.RoutePrefix,
					},
					Targets: []config.RouteTarget{{
						Host: "git.com",
						Port: 8080,
					}},
				},
			},
			args: args{
				host:   "spaceuptech.com",
				method: http.MethodPatch,
				url:    "/abc/xyz",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "* method",
			r: config.Routes{
				&config.Route{
					ID:      "1234",
					Project: "test",
					Source: config.RouteSource{
						Hosts:   []string{"spaceuptech.com"},
						Methods: []string{"*"},
						URL:     "/abc",
						Type:    config.RoutePrefix,
					},
				},
			},
			args: args{
				host:   "spaceuptech.com",
				method: http.MethodPatch,
				url:    "/abc/xyz",
			},
			want: &config.Route{
				ID:      "1234",
				Project: "test",
				Source: config.RouteSource{
					Hosts:   []string{"spaceuptech.com"},
					Methods: []string{"*"},
					URL:     "/abc",
					Type:    config.RoutePrefix,
				},
			},
			wantErr: false,
		},
		{
			name: "valid method",
			r: config.Routes{
				&config.Route{
					ID:      "1234",
					Project: "test",
					Source: config.RouteSource{
						Hosts:   []string{"spaceuptech.com"},
						Methods: []string{http.MethodPatch, http.MethodGet},
						URL:     "/abc",
						Type:    config.RoutePrefix,
					},
				},
			},
			args: args{
				host:   "spaceuptech.com",
				method: http.MethodPatch,
				url:    "/abc/xyz",
			},
			want: &config.Route{
				ID:      "1234",
				Project: "test",
				Source: config.RouteSource{
					Hosts:   []string{"spaceuptech.com"},
					Methods: []string{http.MethodPatch, http.MethodGet},
					URL:     "/abc",
					Type:    config.RoutePrefix,
				},
			},
			wantErr: false,
		},
	}
	routeObj := New()
	for _, tt := range tests {
		routeObj.routes = tt.r
		t.Run(tt.name, func(t *testing.T) {
			got, err := routeObj.selectRoute(tt.args.host, tt.args.method, tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("routeMapping.selectRoute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("routeMapping.selectRoute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_routeMapping_deleteProjectRoutes(t *testing.T) {
	type args struct {
		project string
	}
	tests := []struct {
		name string
		r    config.Routes
		want config.Routes
		args args
	}{
		// TODO: Add test cases.
		{
			name: "delete project rule",
			r: config.Routes{
				&config.Route{
					ID:      "1234",
					Project: "test1",
					Source:  config.RouteSource{},
				},
				&config.Route{
					ID:      "12345",
					Project: "test2",
					Source:  config.RouteSource{},
				},
			},
			want: config.Routes{
				&config.Route{
					ID:      "12345",
					Project: "test2",
					Source:  config.RouteSource{},
				},
			},
			args: args{project: "test1"},
		},
	}
	routeObj := New()
	for _, tt := range tests {
		routeObj.routes = tt.r
		t.Run(tt.name, func(t *testing.T) {
			routeObj.deleteProjectRoutes(tt.args.project)

			if !reflect.DeepEqual(tt.want, routeObj.routes) {
				t.Errorf("Routing.deleteProjectRoutes(): wanted - %v; got - %v", tt.want, routeObj.routes)

				a, _ := json.MarshalIndent(routeObj.routes, "", " ")
				log.Printf("got= %s", string(a))

				a, _ = json.MarshalIndent(tt.want, "", " ")
				log.Printf("want = %s", string(a))
			}
		})
	}
}

func Test_routeMapping_addProjectRoutes(t *testing.T) {
	type args struct {
		project string
		routes  config.Routes
	}
	tests := []struct {
		name string
		r    config.Routes
		want config.Routes
		args args
	}{
		// TODO: Add test cases.
		{
			name: "add project rule",
			r: config.Routes{
				&config.Route{
					ID:      "12345",
					Project: "test1",
					Source:  config.RouteSource{},
				},
			},
			args: args{
				project: "test2",
				routes: config.Routes{
					&config.Route{
						ID:      "1234",
						Project: "test2",
						Source:  config.RouteSource{URL: "/get"},
					},
					&config.Route{
						ID:      "1234",
						Project: "test2",
						Source:  config.RouteSource{URL: "/"},
					},
					&config.Route{
						ID:      "1234",
						Project: "test2",
						Source:  config.RouteSource{URL: "/add"},
					},
				},
			},
			want: config.Routes{
				&config.Route{
					ID:      "1234",
					Project: "test2",
					Source:  config.RouteSource{URL: "/get"},
				},
				&config.Route{
					ID:      "1234",
					Project: "test2",
					Source:  config.RouteSource{URL: "/add"},
				},
				&config.Route{
					ID:      "1234",
					Project: "test2",
					Source:  config.RouteSource{URL: "/"},
				},
				&config.Route{
					ID:      "12345",
					Project: "test1",
					Source:  config.RouteSource{},
				},
			},
		},
	}
	routeObj := New()

	for _, tt := range tests {
		routeObj.routes = tt.r
		t.Run(tt.name, func(t *testing.T) {
			routeObj.addProjectRoutes(tt.args.project, tt.args.routes)

			if !reflect.DeepEqual(routeObj.routes, tt.want) {
				t.Errorf("Routing.addProjectRoutes(): wanted - %v; got - %v", tt.want, routeObj.routes)

				a, _ := json.MarshalIndent(routeObj.routes, "", " ")
				log.Printf("got= %s", string(a))

				a, _ = json.MarshalIndent(tt.want, "", " ")
				log.Printf("want = %s", string(a))
			}
		})
	}
}
