package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/spaceuptech/space-cli/cmd/model"
	"github.com/spaceuptech/space-cli/cmd/utils"
)

func TestGetAuthProviders(t *testing.T) {
	srv := serverMock()
	defer srv.Close()
	utils.LoginStart("local-admin", "1YSU0YzJaWBu", srv.URL)

	type args struct {
		project     string
		commandName string
		params      map[string]string
	}
	tests := []struct {
		name    string
		args    args
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
			want:    []*model.SpecObject{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAuthProviders(tt.args.project, tt.args.commandName, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAuthProviders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAuthProviders() = %v-%v-%v-%v, want %v-%v-%v-%v", got, len(got), reflect.TypeOf(got), cap(got), tt.want, len(tt.want), reflect.TypeOf(tt.want), cap(tt.want))
				t.Errorf("%v", cmp.Equal(got, tt.want))
			}
		})
	}
}

func serverMock() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/v1/config/login", usersMock)
	handler.HandleFunc("/v1/config/projects/myproject/user-management/provider", users)

	srv := httptest.NewServer(handler)

	return srv
}

func usersMock(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("mock server responding"))
}

type provider struct {
	ID      string `json:"id" yaml:"id"`
	Enabled bool   `json:"enabled" yaml:"enabled"`
	Secret  string `json:"secret" yaml:"secret"`
}

type response struct {
	Error  string        `json:"error,omitempty"`
	Result []interface{} `json:"result,omitempty"`
}

func users(w http.ResponseWriter, r *http.Request) {
	p := provider{
		ID:      "local-admin",
		Enabled: true,
		Secret:  "hello",
	}
	providers := []interface{}{p}
	res := response{
		Result: providers,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_ = json.NewEncoder(w).Encode(res)
}
