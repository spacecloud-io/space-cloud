package routing

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"text/template"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

func TestRouting_modifyResponse(t *testing.T) {
	type args struct {
		res           string
		headers       http.Header
		resTmpl       string
		globalHeaders config.Headers
		routeHeaders  config.Headers
		auth          interface{}
	}
	tests := []struct {
		name        string
		args        args
		want        string
		wantHeaders http.Header
		wantErr     bool
	}{
		{
			name: "no modification in response even when content type header is set",
			args: args{
				res:     `{"abc": "xyz"}`,
				headers: map[string][]string{"Content-Type": {"application/json"}},
			},
			want:        `{"abc": "xyz"}`,
			wantHeaders: map[string][]string{"Content-Type": {"application/json"}},
		},
		{
			name: "set headers and no body - only global",
			args: args{
				res:           `{"abc": "xyz"}`,
				headers:       map[string][]string{},
				globalHeaders: config.Headers{{Key: "abc", Value: "xyz"}},
			},
			want:        `{"abc": "xyz"}`,
			wantHeaders: map[string][]string{"Abc": {"xyz"}},
		},
		{
			name: "set headers and no body - only route",
			args: args{
				res:          `{"abc": "xyz"}`,
				headers:      map[string][]string{},
				routeHeaders: config.Headers{{Key: "abc", Value: "xyz"}},
			},
			want:        `{"abc": "xyz"}`,
			wantHeaders: map[string][]string{"Abc": {"xyz"}},
		},
		{
			name: "set headers and no body - both global & route",
			args: args{
				res:           `{"abc": "xyz"}`,
				headers:       map[string][]string{},
				globalHeaders: config.Headers{{Key: "xyz", Value: "abc"}},
				routeHeaders:  config.Headers{{Key: "abc", Value: "xyz"}},
			},
			want:        `{"abc": "xyz"}`,
			wantHeaders: map[string][]string{"Abc": {"xyz"}, "Xyz": {"abc"}},
		},
		{
			name: "mutate body",
			args: args{
				res:     `{"abc":"xyz"}`,
				headers: map[string][]string{"Content-Type": {"application/json"}},
				resTmpl: `{"res": {{marshalJSON .args}}}`,
			},
			want:        `{"res":{"abc":"xyz"}}`,
			wantHeaders: map[string][]string{"Content-Type": {"application/json"}, "Content-Length": {"21"}},
		},
		{
			name: "mutate body - invalid res payload",
			args: args{
				res:     `{"abc":"xyz"}as`,
				headers: map[string][]string{"Content-Type": {"application/json"}},
				resTmpl: `{"res": {{marshalJSON .args}}}`,
			},
			want:        `{"abc":"xyz"}as`,
			wantHeaders: map[string][]string{"Content-Type": {"application/json"}},
		},
		{
			name: "mutate body - no json header",
			args: args{
				res:     `{"abc": "xyz"}`,
				headers: map[string][]string{},
				resTmpl: `{"res": {{marshalJSON .args}}}`,
			},
			want:        `{"abc": "xyz"}`,
			wantHeaders: map[string][]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make an instance of the routing module
			r := &Routing{goTemplates: make(map[string]*template.Template), globalConfig: &config.GlobalRoutesConfig{ResponseHeaders: tt.args.globalHeaders}}
			if tt.args.resTmpl != "" {
				err := r.createGoTemplate("response", "p", "id", tt.args.resTmpl)
				if err != nil {
					t.Error("Unable to parse template", err)
					return
				}
			}
			// Make an instance of the route
			modify := struct {
				Tmpl            config.TemplatingEngine `json:"template,omitempty" yaml:"template,omitempty" mapstructure:"template"`
				ReqTmpl         string                  `json:"requestTemplate" yaml:"requestTemplate" mapstructure:"requestTemplate"`
				ResTmpl         string                  `json:"responseTemplate" yaml:"responseTemplate" mapstructure:"responseTemplate"`
				OpFormat        string                  `json:"outputFormat,omitempty" yaml:"outputFormat,omitempty" mapstructure:"outputFormat"`
				RequestHeaders  config.Headers          `json:"headers" yaml:"headers" mapstructure:"headers"`
				ResponseHeaders config.Headers          `json:"resHeaders" yaml:"resHeaders" mapstructure:"resHeaders"`
			}{Tmpl: config.TemplatingEngineGo, ResTmpl: tt.args.resTmpl, OpFormat: "json", ResponseHeaders: tt.args.routeHeaders}
			route := &config.Route{ID: "id", Project: "p", Modify: modify}

			// Make an instance of the response object
			res := &http.Response{Body: ioutil.NopCloser(bytes.NewBuffer([]byte(tt.args.res))), Header: tt.args.headers}

			if err := r.modifyResponse(context.Background(), res, route, "", tt.args.auth); (err != nil) != tt.wantErr {
				t.Errorf("modifyResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check if response body is correct
			data, _ := ioutil.ReadAll(res.Body)
			if got := string(data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("modifyResponse()[body] = %v, want %v", got, tt.want)
				return
			}

			// Check if response headers is correct
			if !reflect.DeepEqual(res.Header, tt.wantHeaders) {
				t.Errorf("modifyResponse()[header] = %v, want %v", res.Header, tt.wantHeaders)
				return
			}
		})
	}
}

func Test_makeQueryArguments(t *testing.T) {
	type args struct {
		url     string
		headers map[string][]string
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "valid test case",
			args: args{url: "/abc/xyz?foo=bar&bool=true", headers: map[string][]string{"a1": {"b1"}, "a2": {"b2"}}},
			want: map[string]interface{}{
				"path":      "/abc/xyz",
				"pathArray": []interface{}{"abc", "xyz"},
				"params":    map[string]interface{}{"foo": "bar", "bool": "true"},
				"headers":   map[string]interface{}{"A1": "b1", "A2": "b2"},
			},
		},
		{
			name: "url param without value",
			args: args{url: "/abc/xyz?foo=bar&bool", headers: map[string][]string{"a1": {"b1"}, "a2": {"b2"}}},
			want: map[string]interface{}{
				"path":      "/abc/xyz",
				"pathArray": []interface{}{"abc", "xyz"},
				"params":    map[string]interface{}{"foo": "bar", "bool": ""},
				"headers":   map[string]interface{}{"A1": "b1", "A2": "b2"},
			},
		},
		{
			name: "multiple headers only the first one becomes available",
			args: args{url: "/abc/xyz?foo=bar&bool=true", headers: map[string][]string{"a1": {"b11", "b12", "b13"}, "a2": {"b2"}}},
			want: map[string]interface{}{
				"path":      "/abc/xyz",
				"pathArray": []interface{}{"abc", "xyz"},
				"params":    map[string]interface{}{"foo": "bar", "bool": "true"},
				"headers":   map[string]interface{}{"A1": "b11", "A2": "b2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.Parse(tt.args.url)
			if err != nil {
				t.Errorf("makeQueryArguments() - error while passing url - %v", err)
				return
			}
			header := http.Header{}
			for k, array := range tt.args.headers {
				for _, v := range array {
					header.Add(k, v)
				}
			}
			if got := makeQueryArguments(&http.Request{Header: header, URL: u}); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("makeQueryArguments() = %v, want %v", got, tt.want)
			}
		})
	}
}
