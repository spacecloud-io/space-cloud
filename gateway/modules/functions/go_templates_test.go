package functions

import (
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
			name: "valid - use params",
			args: args{
				tmpl:   `{"foo": "{{index . "abc"}}"}`,
				params: map[string]interface{}{"abc": "bar"},
				format: "json",
			},
			want: map[string]interface{}{"foo": "bar"},
		},
		{
			name: "valid - use params (nested objects)",
			args: args{
				tmpl:   `{"foo": "{{index . "a" "b"}}"}`,
				params: map[string]interface{}{"a": map[string]interface{}{"b": "bar"}},
				format: "json",
			},
			want: map[string]interface{}{"foo": "bar"},
		},
		{
			name: "valid - use params (nested objects and arrays)",
			args: args{
				tmpl:   `{"foo": "{{index . "a" "b" 0}}"}`,
				params: map[string]interface{}{"a": map[string]interface{}{"b": []interface{}{"bar"}}},
				format: "json",
			},
			want: map[string]interface{}{"foo": "bar"},
		},
		{
			name: "valid - trying loops in yaml (nobody might need something this complex)",
			args: args{
				tmpl: `
{{ range $i, $value := index . "array" }}
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
    cluster_id: "{{ index . "cluster" }}"
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
			"cluster_id": "{{ index . "cluster" }}"
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
			}).Parse(tt.args.tmpl)
			if err != nil {
				t.Errorf("goTemplate() error = %v, could not pass template", err)
				return
			}

			got, err := goTemplate(tmpl, tt.args.format, tt.args.claims, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("goTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("goTemplate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
