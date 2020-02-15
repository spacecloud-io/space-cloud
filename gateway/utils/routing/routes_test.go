package routing

import (
	"encoding/json"
	"log"
	"reflect"
	"sync"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

func Test_routeMapping_selectRoute(t *testing.T) {
	type args struct {
		host string
		url  string
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
						Id: "1234",
						Source: config.RouteSource{
							Hosts: []string{"spaceuptech.com"},
							URL:   "/abc",
							Type:  config.RoutePrefix,
						},
						Destination: config.RouteDestination{
							Host: "git.com",
							Port: "8080",
						},
					},
				},
			},
			args: args{
				host: "spaceuptech.com",
				url:  "/abc/xyz",
			},
			want: &config.Route{
				Id: "1234",
				Source: config.RouteSource{
					Hosts: []string{"spaceuptech.com"},
					URL:   "/abc",
					Type:  config.RoutePrefix,
				},
				Destination: config.RouteDestination{
					Host: "git.com",
					Port: "8080",
				},
			},
			wantErr: false,
		},
		{
			name: "2 routes present",
			r: routeMapping{
				"test": config.Routes{
					&config.Route{
						Id: "1234",
						Source: config.RouteSource{
							Hosts: []string{"spaceuptech.com"},
							URL:   "/abc/xyz",
							Type:  config.RouteExact,
						},
						Destination: config.RouteDestination{
							Host: "git.com",
							Port: "8080",
						},
					},
					&config.Route{
						Id: "567",
						Source: config.RouteSource{
							Hosts: []string{"spaceuptech.com"},
							URL:   "/abc",
							Type:  config.RoutePrefix,
						},
						Destination: config.RouteDestination{
							Host: "git.com",
							Port: "8080",
						},
					},
				},
			},
			args: args{
				host: "spaceuptech.com",
				url:  "/abc/xyz/123",
			},
			want: &config.Route{
				Id: "567",
				Source: config.RouteSource{
					Hosts: []string{"spaceuptech.com"},
					URL:   "/abc",
					Type:  config.RoutePrefix,
				},
				Destination: config.RouteDestination{
					Host: "git.com",
					Port: "8080",
				},
			},
			wantErr: false,
		},
		{
			name: "host match and url prefix match",
			r: routeMapping{
				"test": config.Routes{
					&config.Route{
						Id: "1234",
						Source: config.RouteSource{
							Hosts: []string{"spaceuptech.com"},
							URL:   "/abc/xyz",
							Type:  config.RouteExact,
						},
						Destination: config.RouteDestination{
							Host: "git.com",
							Port: "8080",
						},
					},
				},
			},
			args: args{
				host: "spaceuptech.com",
				url:  "/abc/xyz",
			},
			want: &config.Route{
				Id: "1234",
				Source: config.RouteSource{
					Hosts: []string{"spaceuptech.com"},
					URL:   "/abc/xyz",
					Type:  config.RouteExact,
				},
				Destination: config.RouteDestination{
					Host: "git.com",
					Port: "8080",
				},
			},
			wantErr: false,
		},

		{
			name: "host match and url match Type not provided",
			r: routeMapping{
				"test": config.Routes{
					&config.Route{
						Id: "1234",
						Source: config.RouteSource{
							Hosts: []string{"spaceuptech.com"},
							URL:   "/abc",
						},
						Destination: config.RouteDestination{
							Host: "git.com",
							Port: "8080",
						},
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
						Id: "1234",
						Source: config.RouteSource{
							Hosts: []string{"spaceuptech.com"},
							URL:   "/abc",
							Type:  config.RoutePrefix,
						},
						Destination: config.RouteDestination{
							Host: "git.com",
							Port: "8080",
						},
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
						Id: "1234",
						Source: config.RouteSource{
							Hosts: []string{"spaceuptech.com", "*"},
							URL:   "/abc",
							Type:  config.RoutePrefix,
						},
						Destination: config.RouteDestination{
							Host: "git.com",
							Port: "8080",
						},
					},
				},
			},
			args: args{
				host: "",
				url:  "/abc",
			},
			want: &config.Route{
				Id: "1234",
				Source: config.RouteSource{
					Hosts: []string{"spaceuptech.com", "*"},
					URL:   "/abc",
					Type:  config.RoutePrefix,
				},
				Destination: config.RouteDestination{
					Host: "git.com",
					Port: "8080",
				},
			},
			wantErr: false,
		},

		{
			name: "url not match",
			r: routeMapping{
				"test": config.Routes{
					&config.Route{
						Id: "1234",
						Source: config.RouteSource{
							Hosts: []string{"spaceuptech.com"},
							URL:   "/abc",
							Type:  config.RouteExact,
						},
						Destination: config.RouteDestination{
							Host: "git.com",
							Port: "8080",
						},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.selectRoute(tt.args.host, tt.args.url)
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

func TestRouting_selectRoute(t *testing.T) {
	type fields struct {
		lock   sync.RWMutex
		routes routeMapping
	}
	type args struct {
		host string
		url  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *config.Route
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "host match and url prefix match",
			fields: fields{
				routes: routeMapping{
					"test": config.Routes{
						&config.Route{
							Id: "1234",
							Source: config.RouteSource{
								Hosts: []string{"spaceuptech.com"},
								URL:   "/abc",
								Type:  config.RoutePrefix,
							},
							Destination: config.RouteDestination{
								Host: "git.com",
								Port: "8080",
							},
						},
					},
				},
			},
			args: args{
				host: "spaceuptech.com",
				url:  "/abc/xyz",
			},
			want: &config.Route{
				Id: "1234",
				Source: config.RouteSource{
					Hosts: []string{"spaceuptech.com"},
					URL:   "/abc",
					Type:  config.RoutePrefix,
				},
				Destination: config.RouteDestination{
					Host: "git.com",
					Port: "8080",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Routing{
				lock:   tt.fields.lock,
				routes: tt.fields.routes,
			}
			got, err := r.selectRoute(tt.args.host, tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Routing.selectRoute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Routing.selectRoute() = %v, want %v", got, tt.want)
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
			want: routeMapping{
				"test2": config.Routes{
					&config.Route{
						Id:          "12345",
						Source:      config.RouteSource{},
						Destination: config.RouteDestination{},
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
				log.Println("Routing.deleteProjectRoutes()")

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
						Id:          "12345",
						Source:      config.RouteSource{},
						Destination: config.RouteDestination{},
					},
				},
			},
			args: args{
				project: "test2",
				routes: config.Routes{
					&config.Route{
						Id:          "1234",
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
				"test2": config.Routes{
					&config.Route{
						Id:          "1234",
						Source:      config.RouteSource{},
						Destination: config.RouteDestination{},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.addProjectRoutes(tt.args.project, tt.args.routes)

			if !reflect.DeepEqual(tt.r, tt.want) {
				log.Println("Routing.addProjectRoutes()")

				a, _ := json.MarshalIndent(tt.r, "", " ")
				log.Printf("got= %s", string(a))

				a, _ = json.MarshalIndent(tt.want, "", " ")
				log.Printf("want = %s", string(a))
			}
		})
	}
}
