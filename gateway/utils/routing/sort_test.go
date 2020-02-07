package routing

import (
	"encoding/json"
	"github.com/spaceuptech/space-cloud/gateway/config"
	"log"
	"reflect"
	"testing"
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
						Id: "1",
						Source: config.RouteSource{
							URL: "/",
						},
					},
					{
						Id: "2",
						Source: config.RouteSource{
							URL: "/api",
						},
					},
					{
						Id: "3",
						Source: config.RouteSource{
							URL: "/api/good/create/yes",
						},
					},
					{
						Id: "4",
						Source: config.RouteSource{
							URL: "/api/good",
						},
					},
					{
						Id: "5",
						Source: config.RouteSource{
							URL: "/api/abc",
						},
					},
					{
						Id: "6",
						Source: config.RouteSource{
							URL: "/api/abc/create/yes",
						},
					},
					{
						Id: "7",
						Source: config.RouteSource{
							URL: "/api/abc/yes",
						},
					},
				},
			},
			want: args{
				rules: []*config.Route{
					{
						Id: "6",
						Source: config.RouteSource{
							URL: "/api/abc/create/yes",
						},
					},
					{
						Id: "3",
						Source: config.RouteSource{
							URL: "/api/good/create/yes",
						},
					},
					{
						Id: "7",
						Source: config.RouteSource{
							URL: "/api/abc/yes",
						},
					},
					{
						Id: "5",
						Source: config.RouteSource{
							URL: "/api/abc",
						},
					},
					{
						Id: "4",
						Source: config.RouteSource{
							URL: "/api/good",
						},
					},
					{
						Id: "2",
						Source: config.RouteSource{
							URL: "/api",
						},
					},
					{
						Id: "1",
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
