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
		r       routeMapping
		args    args
		want    *config.Route
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "host match and url prefix match",
			r: routeMapping{
				"test": config.Routes{
					&config.Route{
						ID: "1234",
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
			},
			args: args{
				host: "spaceuptech.com",
				url:  "/abc/xyz",
			},
			want: &config.Route{
				ID: "1234",
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
			r: routeMapping{
				"test": config.Routes{
					&config.Route{
						ID: "1234",
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
						ID: "567",
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
			},
			args: args{
				host: "spaceuptech.com",
				url:  "/abc/xyz/123",
			},
			want: &config.Route{
				ID: "567",
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
			r: routeMapping{
				"test": config.Routes{
					&config.Route{
						ID: "1234",
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
			},
			args: args{
				host: "spaceuptech.com",
				url:  "/abc/xyz",
			},
			want: &config.Route{
				ID: "1234",
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
			r: routeMapping{
				"test": config.Routes{
					&config.Route{
						ID: "1234",
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
			r: routeMapping{
				"test": config.Routes{
					&config.Route{
						ID: "1234",
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
			r: routeMapping{
				"test": config.Routes{
					&config.Route{
						ID: "1234",
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
			},
			args: args{
				host: "",
				url:  "/abc",
			},
			want: &config.Route{
				ID: "1234",
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
			r: routeMapping{
				"test": config.Routes{
					&config.Route{
						ID: "1234",
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
			r: routeMapping{
				"test": config.Routes{
					&config.Route{
						ID: "1234",
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
			r: routeMapping{
				"test": config.Routes{
					&config.Route{
						ID: "1234",
						Source: config.RouteSource{
							Hosts:   []string{"spaceuptech.com"},
							Methods: []string{"*"},
							URL:     "/abc",
							Type:    config.RoutePrefix,
						},
					},
				},
			},
			args: args{
				host:   "spaceuptech.com",
				method: http.MethodPatch,
				url:    "/abc/xyz",
			},
			want: &config.Route{
				ID: "1234",
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
			r: routeMapping{
				"test": config.Routes{
					&config.Route{
						ID: "1234",
						Source: config.RouteSource{
							Hosts:   []string{"spaceuptech.com"},
							Methods: []string{http.MethodPatch, http.MethodGet},
							URL:     "/abc",
							Type:    config.RoutePrefix,
						},
					},
				},
			},
			args: args{
				host:   "spaceuptech.com",
				method: http.MethodPatch,
				url:    "/abc/xyz",
			},
			want: &config.Route{
				ID: "1234",
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.selectRoute(tt.args.host, tt.args.method, tt.args.url)
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
		r    routeMapping
		want routeMapping
		args args
	}{
		// TODO: Add test cases.
		{
			name: "delete project rule",
			r: routeMapping{
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
			want: routeMapping{
				"test2": config.Routes{
					&config.Route{
						ID:     "12345",
						Source: config.RouteSource{},
					},
				},
			},
			args: args{project: "test1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.r.deleteProjectRoutes(tt.args.project)

			if !reflect.DeepEqual(tt.want, tt.r) {
				t.Errorf("Routing.deleteProjectRoutes(): wanted - %v; got - %v", tt.want, tt.r)

				a, _ := json.MarshalIndent(tt.r, "", " ")
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
		r    routeMapping
		want routeMapping
		args args
	}{
		// TODO: Add test cases.
		{
			name: "add project rule",
			r: routeMapping{
				"test1": config.Routes{
					&config.Route{
						ID:     "12345",
						Source: config.RouteSource{},
					},
				},
			},
			args: args{
				project: "test2",
				routes: config.Routes{
					&config.Route{
						ID:     "1234",
						Source: config.RouteSource{URL: "/get"},
					},
					&config.Route{
						ID:     "1234",
						Source: config.RouteSource{URL: "/"},
					},
					&config.Route{
						ID:     "1234",
						Source: config.RouteSource{URL: "/add"},
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
				"test2": config.Routes{
					&config.Route{
						ID:     "1234",
						Source: config.RouteSource{URL: "/get"},
					},
					&config.Route{
						ID:     "1234",
						Source: config.RouteSource{URL: "/add"},
					},
					&config.Route{
						ID:     "1234",
						Source: config.RouteSource{URL: "/"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.addProjectRoutes(tt.args.project, tt.args.routes)

			if !reflect.DeepEqual(tt.r, tt.want) {
				t.Errorf("Routing.addProjectRoutes(): wanted - %v; got - %v", tt.want, tt.r)

				a, _ := json.MarshalIndent(tt.r, "", " ")
				log.Printf("got= %s", string(a))

				a, _ = json.MarshalIndent(tt.want, "", " ")
				log.Printf("want = %s", string(a))
			}
		})
	}
}
