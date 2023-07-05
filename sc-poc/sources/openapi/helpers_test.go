package openapi

import (
	"reflect"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

// TestCreateRequestParams tests the createRequestParams function with different inputs and outputs
func TestCreateRequestParams(t *testing.T) {
	type args struct {
		urlPath       string
		reqBodySchema *openapi3.Schema
		params        openapi3.Parameters
		vars          map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   map[string]interface{}
		want2   map[string]string
		err     error
		wantErr bool
	}{
		{
			name: "simple query parameter",
			args: args{
				urlPath:       "/users",
				reqBodySchema: &openapi3.Schema{},
				params: openapi3.Parameters{
					{Value: &openapi3.Parameter{Name: "name", In: "query", Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}}}},
				},
				vars: map[string]interface{}{
					"name": "Alice",
				},
			},
			want:  "/users?name=Alice",
			want1: map[string]interface{}{},
			want2: map[string]string{},
		},
		{
			name: "complex path parameter",
			args: args{
				urlPath:       "/books/{id}",
				reqBodySchema: &openapi3.Schema{},
				params: openapi3.Parameters{
					{Value: &openapi3.Parameter{Name: "id", In: "path", Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "integer"}}}},
				},
				vars: map[string]interface{}{
					"id": 42,
				},
			},
			want:  "/books/42",
			want1: map[string]interface{}{},
			want2: map[string]string{},
		},
		{
			name: "mixed path and query parameters",
			args: args{
				urlPath:       "/users/{id}/posts",
				reqBodySchema: &openapi3.Schema{},
				params: openapi3.Parameters{
					{Value: &openapi3.Parameter{Name: "id", In: "path", Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "integer"}}}},
					{Value: &openapi3.Parameter{Name: "title", In: "query", Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}}}},
					{Value: &openapi3.Parameter{Name: "tags", In: "query", Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "array", Items: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}}}}}},
				},
				vars: map[string]interface{}{
					"id":    42,
					"title": "Hello World",
					"tags":  []string{"golang", "testing"},
				},
			},
			want:  "/users/42/posts?tags=%5B%22golang%22%2C%22testing%22%5D&title=Hello+World",
			want1: map[string]interface{}{},
			want2: map[string]string{},
		},
		{
			name: "request body parameter",
			args: args{
				urlPath: "/users",
				reqBodySchema: &openapi3.Schema{
					Type: "object",
					Properties: map[string]*openapi3.SchemaRef{
						"name": {Value: &openapi3.Schema{Type: "string"}},
						"age":  {Value: &openapi3.Schema{Type: "integer"}},
					},
				},
				params: openapi3.Parameters{},
				vars: map[string]interface{}{
					"name": "Charlie",
					"age":  30,
				},
			},
			want: "/users",
			want1: map[string]interface{}{
				"name": "Charlie",
				"age":  30,
			},
			want2: map[string]string{},
		},
		{
			name: "query parameter in request body",
			args: args{
				urlPath: "/search",
				reqBodySchema: &openapi3.Schema{
					Type: "object",
					Properties: map[string]*openapi3.SchemaRef{
						"query": {Value: &openapi3.Schema{Type: "string"}},
					},
				},
				params: openapi3.Parameters{
					{Value: &openapi3.Parameter{Name: "query", In: "query", Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}}}},
				},
				vars: map[string]interface{}{
					"query": "golang",
				},
			},
			want: "/search?query=golang",
			want1: map[string]interface{}{
				"query": "golang",
			},
			want2: map[string]string{},
		},
		{
			name: "header parameter",
			args: args{
				urlPath:       "/users",
				reqBodySchema: &openapi3.Schema{},
				params: openapi3.Parameters{
					{Value: &openapi3.Parameter{Name: "Authorization", In: "header", Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}}}},
				},
				vars: map[string]interface{}{
					"Authorization": "Bearer token",
				},
			},
			want:  "/users",
			want1: map[string]interface{}{},
			want2: map[string]string{
				"Authorization": "Bearer token",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, err := createRequestParams(tt.args.urlPath, tt.args.reqBodySchema, tt.args.params, tt.args.vars)
			if err != nil && !tt.wantErr {
				t.Errorf("createRequestParams() err = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("createRequestParams() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("createRequestParams() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("createRequestParams() got2 = %v, want %v", got2, tt.want2)
			}

		})
	}
}
