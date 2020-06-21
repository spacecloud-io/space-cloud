package tmpl

import (
	"encoding/json"
	"reflect"
	"testing"
	"text/template"

	"github.com/spaceuptech/space-cloud/gateway/modules/auth"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func Test_goTemplate(t *testing.T) {
	type args struct {
		tmpl   string
		claims interface{}
		params interface{}
		format string
		token  string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "valid - empty params and claims",
			args: args{tmpl: `{"foo": "bar"}`, params: map[string]interface{}{}, format: "json"},
			want: map[string]interface{}{"foo": "bar"},
		},
		{
			name: "valid - empty params and claims",
			args: args{tmpl: `foo: bar`, params: map[string]interface{}{}, format: "yaml"},
			want: map[string]interface{}{"foo": "bar"},
		},
		{
			name: "valid - use params",
			args: args{
				tmpl:   `{"foo": "{{index . "args" "abc"}}"}`,
				params: map[string]interface{}{"abc": "bar"},
				format: "json",
			},
			want: map[string]interface{}{"foo": "bar"},
		},
		{
			name: "valid - use params with token",
			args: args{
				tmpl:   `{"foo": "{{index . "args" "abc"}}", "token": "{{index . "token"}}"}`,
				params: map[string]interface{}{"abc": "bar"},
				format: "json",
				token:  "jwt token",
			},
			want: map[string]interface{}{"foo": "bar", "token": "jwt token"},
		},
		{
			name: "valid - use params",
			args: args{
				tmpl:   `{"foo": "{{index . "args" "abc"}}"}`,
				params: map[string]interface{}{"abc": "bar"},
				format: "string",
			},
			want: `{"foo": "bar"}`,
		},
		{
			name: "valid - use params (nested objects)",
			args: args{
				tmpl:   `{"foo": "{{index . "args" "a" "b"}}"}`,
				params: map[string]interface{}{"a": map[string]interface{}{"b": "bar"}},
				format: "json",
			},
			want: map[string]interface{}{"foo": "bar"},
		},
		{
			name: "valid - use params (nested objects) without index",
			args: args{
				tmpl:   `{"foo": "{{.args.a.b}}"}`,
				params: map[string]interface{}{"a": map[string]interface{}{"b": "bar"}},
				format: "json",
			},
			want: map[string]interface{}{"foo": "bar"},
		},
		{
			name: "valid - use marshal function json",
			args: args{
				tmpl:   `{"foo": {{ marshalJSON (index . "args" "a")}}}`,
				params: map[string]interface{}{"a": map[string]interface{}{"b": "bar"}},
				format: "json",
			},
			want: map[string]interface{}{"foo": map[string]interface{}{"b": "bar"}},
		},
		{
			name: "valid - use params (nested objects and arrays)",
			args: args{
				tmpl:   `{"foo": "{{index . "args" "a" "b" 0}}"}`,
				params: map[string]interface{}{"a": map[string]interface{}{"b": []interface{}{"bar"}}},
				format: "json",
			},
			want: map[string]interface{}{"foo": "bar"},
		},
		{
			name: "valid - trying loops in yaml (noargs might need something this complex)",
			args: args{
				tmpl: `
{{ range $i, $value := index . "args" "array" }}
{{ index $value "p1" }}: {{ index $value "p2" }}
{{ end }}
`,
				params: map[string]interface{}{"array": []interface{}{
					map[string]interface{}{"p1": "foo1", "p2": "bar1"},
					map[string]interface{}{"p1": "foo2", "p2": "bar2"},
					map[string]interface{}{"p1": "foo3", "p2": "bar3"},
				}},
				format: "yaml",
			},
			want: map[string]interface{}{"foo1": "bar1", "foo2": "bar2", "foo3": "bar3"},
		},
		{
			name: "valid - real world use case (in yaml)",
			args: args{
				tmpl: `
query: "mutation { update_clusters(where: $where, set: $set) @db { status error }"
variables:
  where:
    owner_id: "{{ index . "auth" "id" }}"
    cluster_id: "{{ index . "args" "cluster" }}"
  set:
    session_id: ""
    cluster_key: "{{ generateId }}"
`,
				params: map[string]interface{}{"cluster": "cluster 1"},
				claims: auth.TokenClaims{"id": "1"},
				format: "yaml",
			},
			want: map[string]interface{}{
				"query": "mutation { update_clusters(where: $where, set: $set) @db { status error }",
				"variables": map[string]interface{}{
					"where": map[string]interface{}{"owner_id": "1", "cluster_id": "cluster 1"},
					"set":   map[string]interface{}{"session_id": "", "cluster_key": "generated id"},
				},
			},
		},
		{
			name: "valid - real world use case (in json)",
			args: args{
				tmpl: `
{
	"query": "mutation { update_clusters(where: $where, set: $set) @db { status error }",
	"variables": {
		"where": {
			"owner_id": "{{ index . "auth" "id" }}",
			"cluster_id": "{{ index . "args" "cluster" }}"
		},
		"set": {
			"session_id": "",
			"cluster_key": "{{ generateId }}"
		}
	}
}
`,
				params: map[string]interface{}{"cluster": "cluster 1"},
				claims: auth.TokenClaims{"id": "1"},
				format: "json",
			},
			want: map[string]interface{}{
				"query": "mutation { update_clusters(where: $where, set: $set) @db { status error }",
				"variables": map[string]interface{}{
					"where": map[string]interface{}{"owner_id": "1", "cluster_id": "cluster 1"},
					"set":   map[string]interface{}{"session_id": "", "cluster_key": "generated id"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := template.New(tt.name).Funcs(template.FuncMap{
				"hash":       utils.HashString,
				"generateId": func() string { return "generated id" },
				"add":        func(a, b int) int { return a + b },
				"marshalJSON": func(a interface{}) (string, error) {
					data, err := json.Marshal(a)
					return string(data), err
				},
			}).Parse(tt.args.tmpl)
			if err != nil {
				t.Errorf("GoTemplate() error = %v, could not pass template", err)
				return
			}

			got, err := GoTemplate("", "", tmpl, tt.args.format, tt.args.token, tt.args.claims, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GoTemplate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
