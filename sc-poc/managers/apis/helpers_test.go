package apis

import (
	"reflect"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func Test_replacePathParam(t *testing.T) {
	type args struct {
		path    string
		param   string
		indexes map[string]string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 map[string]string
	}{
		{
			name: "No path param",
			args: args{
				path:    "/foo/bar",
				param:   "project",
				indexes: map[string]string{},
			},
			want:  "/foo/bar",
			want1: map[string]string{},
		},
		{
			name: "single path param",
			args: args{
				path:    "/foo/{bar}",
				param:   "bar",
				indexes: map[string]string{},
			},
			want:  "/foo/*",
			want1: map[string]string{"0": "bar"},
		},
		{
			name: "2 path params",
			args: args{
				path:    "/foo/{bar}/abc/{bar}/xyz",
				param:   "bar",
				indexes: map[string]string{},
			},
			want:  "/foo/*/abc/*/xyz",
			want1: map[string]string{"0": "bar", "1": "bar"},
		},
		{
			name: "Multiple path params",
			args: args{
				path:    "/foo/{bar}/abc/{bar}/xyz/{bar}/{bar}",
				param:   "bar",
				indexes: map[string]string{},
			},
			want:  "/foo/*/abc/*/xyz/*/*",
			want1: map[string]string{"0": "bar", "1": "bar", "2": "bar", "3": "bar"},
		},
		{
			name: "Mixed path params",
			args: args{
				path:    "/foo/{bar}/abc/{proj}/xyz/{proj}/{bar}",
				param:   "bar",
				indexes: map[string]string{},
			},
			want:  "/foo/*/abc/{proj}/xyz/{proj}/*",
			want1: map[string]string{"0": "bar", "3": "bar"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := replacePathParam(tt.args.path, tt.args.param, tt.args.indexes)
			if got != tt.want {
				t.Errorf("replacePathParam() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("replacePathParam() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_sanitizeUrl(t *testing.T) {
	type args struct {
		api *API
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 []string
	}{
		{
			name: "No path params",
			args: args{
				api: &API{
					Path: "/a/b/c",
					OpenAPI: &OpenAPI{
						PathDef: &openapi3.PathItem{Parameters: openapi3.NewParameters()},
					},
				},
			},
			want:  "/a/b/c",
			want1: []string{},
		},
		{
			name: "Mixed path params",
			args: args{
				api: &API{
					Path: "/a/{b}/{c}/d/{b}/{e}",
					OpenAPI: &OpenAPI{
						PathDef: &openapi3.PathItem{Parameters: openapi3.Parameters{{
							Value: &openapi3.Parameter{In: "path", Name: "b"},
						}, {
							Value: &openapi3.Parameter{In: "path", Name: "c"},
						}, {
							Value: &openapi3.Parameter{In: "path", Name: "e"},
						}}}},
				},
			},
			want:  "/a/*/*/d/*/*",
			want1: []string{"b", "c", "b", "e"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := sanitizeURL(tt.args.api)
			if got != tt.want {
				t.Errorf("sanitizeUrl() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("sanitizeUrl() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_getPathParams(t *testing.T) {
	type args struct {
		ogURL       string
		receivedURL string
		indexes     []string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "no path params",
			args: args{
				ogURL:       "/a/b/c",
				receivedURL: "/a/b/c",
				indexes:     []string{},
			},
			want: map[string]string{},
		},
		{
			name: "single path params",
			args: args{
				ogURL:       "/a/*/c",
				receivedURL: "/a/b/c",
				indexes:     []string{"project"},
			},
			want: map[string]string{"project": "b"},
		},
		{
			name: "multiple path params",
			args: args{
				ogURL:       "/a/*/c/*",
				receivedURL: "/a/b/c/d",
				indexes:     []string{"project", "env"},
			},
			want: map[string]string{"project": "b", "env": "d"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPathParams(tt.args.ogURL, tt.args.receivedURL, tt.args.indexes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPathParams() = %v, want %v", got, tt.want)
			}
		})
	}
}
