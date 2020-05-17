package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cli/cmd/model"
	"github.com/spaceuptech/space-cli/cmd/utils"
)

func TestGetAuthProviders(t *testing.T) {

	type args struct {
		project     string
		commandName string
		params      map[string]string
	}
	type provider struct {
		ID      string `json:"id" yaml:"id"`
		Enabled bool   `json:"enabled" yaml:"enabled"`
		Secret  string `json:"secret" yaml:"secret"`
	}
	tests := []struct {
		name    string
		args    args
		url     map[string]interface{}
		want    []*model.SpecObject
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test-1",
			args: args{
				project:     "myproject",
				commandName: "auth-providers",
				params:      map[string]string{},
			},
			url: map[string]interface{}{
				"/v1/config/login": "mock server responding",
				"/v1/config/projects/myproject/user-management/provider": model.Response{
					Result: []interface{}{provider{
						ID:      "local-admin",
						Enabled: true,
						Secret:  "hello",
					},
					},
				},
			},
			want:    []*model.SpecObject{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := serverMock(tt.url)
			defer srv.Close()
			utils.LoginStart("local-admin", "1YSU0YzJaWBu", srv.URL)
			got, err := GetAuthProviders(tt.args.project, tt.args.commandName, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAuthProviders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAuthProviders() = %v-%v-%v-%v, want %v-%v-%v-%v", got, len(got), reflect.TypeOf(got), cap(got), tt.want, len(tt.want), reflect.TypeOf(tt.want), cap(tt.want))
			}
		})
	}
}

func serverMock(url map[string]interface{}) *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/", usersMock(url))

	srv := httptest.NewServer(handler)

	return srv
}

func usersMock(url map[string]interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := url[r.URL.Path]
		logrus.Errorf("res=%v", res)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(res)
	}
}
