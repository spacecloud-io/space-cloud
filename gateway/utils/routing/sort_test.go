package routing

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

func Test_sortFileRule(t *testing.T) {
	type args struct {
		rules []*config.Route
	}
	tests := []struct {
		name string
		args args
		want args
	}{
		{
			name: "unsorted routing rules",
			args: args{
				rules: []*config.Route{
					{
						ID: "1",
						Source: config.RouteSource{
							URL: "/",
						},
					},
					{
						ID: "2",
						Source: config.RouteSource{
							URL: "/api",
						},
					},
					{
						ID: "3",
						Source: config.RouteSource{
							URL: "/api/good/create/yes",
						},
					},
					{
						ID: "4",
						Source: config.RouteSource{
							URL: "/api/good",
						},
					},
					{
						ID: "5",
						Source: config.RouteSource{
							URL: "/api/abc",
						},
					},
					{
						ID: "6",
						Source: config.RouteSource{
							URL: "/api/abc/create/yes",
						},
					},
					{
						ID: "7",
						Source: config.RouteSource{
							URL: "/api/abc/yes",
						},
					},
				},
			},
			want: args{
				rules: []*config.Route{
					{
						ID: "6",
						Source: config.RouteSource{
							URL: "/api/abc/create/yes",
						},
					},
					{
						ID: "3",
						Source: config.RouteSource{
							URL: "/api/good/create/yes",
						},
					},
					{
						ID: "7",
						Source: config.RouteSource{
							URL: "/api/abc/yes",
						},
					},
					{
						ID: "5",
						Source: config.RouteSource{
							URL: "/api/abc",
						},
					},
					{
						ID: "4",
						Source: config.RouteSource{
							URL: "/api/good",
						},
					},
					{
						ID: "2",
						Source: config.RouteSource{
							URL: "/api",
						},
					},
					{
						ID: "1",
						Source: config.RouteSource{
							URL: "/",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortRoutes(tt.args.rules)
			if !reflect.DeepEqual(tt.args.rules, tt.want.rules) {
				t.Errorf("error sorting routes unablt to sort routes properly")
			}
			val, _ := json.MarshalIndent(tt.args.rules, "", " ")
			log.Println(string(val))
		})
	}
}

// /api/abc/qwer
// /api/abc/*/yus
// /api/abc
// /api/xyz
// /api/*/qwer45t/yus
// /api
// /
