package routing

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

func Test_getHostAndURL(t *testing.T) {
	type args struct {
		request *http.Request
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		// TODO: Add test cases.
		{
			name: "timepass",
			args: args{
				request: &http.Request{
					Host: "www.golang.com:abc",
					URL: &url.URL{
						Path: "/abc",
					},
				},
			},
			want:  "www.golang.com",
			want1: "/abc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getHostAndURL(tt.args.request)
			if got != tt.want {
				t.Errorf("getHostAndURL() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getHostAndURL() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_rewriteURL(t *testing.T) {
	type args struct {
		url   string
		route *config.Route
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "rewrite url",
			args: args{
				url: "/abc/xyz",
				route: &config.Route{
					ID: "1234",
					Source: config.RouteSource{
						URL:        "/abc",
						RewriteURL: "/v1/abc",
					},
				},
			},
			want: "/v1/abc/xyz",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rewriteURL(tt.args.url, tt.args.route); got != tt.want {
				t.Errorf("rewriteURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setRequest(t *testing.T) {
	type args struct {
		request *http.Request
		route   *config.Route
		url     string
	}
	tests := []struct {
		name string
		args args
		want *http.Request
	}{
		{
			name: "set request",
			args: args{
				url: "/abc",
				route: &config.Route{
					Targets: []config.RouteTarget{{
						Host:   "spacecloud.com",
						Port:   8080,
						Weight: 100,
					}},
				},
				request: &http.Request{
					URL: &url.URL{},
				},
			},
			want: &http.Request{
				Host: "spacecloud.com",
				URL: &url.URL{
					Host:   "spacecloud.com:8080",
					Path:   "/abc",
					Scheme: "http",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = setRequest(tt.args.request, tt.args.route, tt.args.url)
			if !reflect.DeepEqual(tt.args.request, tt.want) {
				t.Errorf("Routing.addProjectRoutes(): wanted - %v; got - %v", tt.want, tt.args.request)

				a, _ := json.MarshalIndent(tt.args.request, "", " ")
				log.Printf("got= %s", string(a))

				a, _ = json.MarshalIndent(tt.want, "", " ")
				log.Printf("want = %s", string(a))
			}
		})
	}
}
