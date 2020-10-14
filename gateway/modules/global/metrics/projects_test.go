package metrics

import (
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/managers/admin"
	"github.com/spaceuptech/space-cloud/gateway/managers/syncman"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func TestModule_generateMetricsRequest(t *testing.T) {
	type args struct {
		project *config.Project
		ssl     *config.SSL
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want2 map[string]interface{}
		want3 map[string]interface{}
	}{
		{
			name: "valid config",
			args: args{
				project: &config.Project{
					ProjectConfig: &config.ProjectConfig{ID: "project"},
					DatabaseConfigs: config.DatabaseConfigs{
						"": &config.DatabaseConfig{DbAlias: "postgres", Type: "postgres", Enabled: true},
					},
					DatabaseSchemas: config.DatabaseSchemas{
						"resourceID1": &config.DatabaseSchema{
							DbAlias: "postgres",
							Table:   "table1",
						},
						"resourceID2": &config.DatabaseSchema{
							DbAlias: "postgres",
							Table:   "event_logs",
						},
						"resourceID3": &config.DatabaseSchema{
							DbAlias: "postgres",
							Table:   "invocation_logs",
						},
						"resourceID4": &config.DatabaseSchema{
							DbAlias: "postgres",
							Table:   "default",
						},
					},
					Auths: config.Auth{
						"": &config.AuthStub{ID: "email", Enabled: true},
					},
					RemoteService: config.Services{
						"": &config.Service{},
					},
					FileStoreConfig: &config.FileStoreConfig{
						Enabled:   true,
						StoreType: "local",
					},
					FileStoreRules: config.FileStoreRules{
						"": &config.FileRule{ID: "file"},
					},
					EventingConfig: &config.EventingConfig{Enabled: true, DBAlias: "postgres"},
					EventingTriggers: config.EventingTriggers{
						"": &config.EventingTrigger{Type: "type"},
					},
					LetsEncrypt: &config.LetsEncrypt{WhitelistedDomains: []string{"www.google.com"}},
					IngressRoutes: config.IngressRoutes{
						"": &config.Route{ID: "route"},
					},
				},
				ssl: &config.SSL{
					Enabled: true,
				},
			},
			want: "clusterID--project",
			want2: map[string]interface{}{
				"nodes":        1,
				"os":           runtime.GOOS,
				"is_prod":      false,
				"version":      utils.BuildVersion,
				"distribution": "ce",
				"last_updated": time.Now().UnixNano() / int64(time.Millisecond),
				"ssl_enabled":  true,
				"project":      "project",
				"crud": map[string]interface{}{
					"postgres": map[string]interface{}{
						"tables": 1,
					},
				},
				"databases": map[string][]string{
					"databases": {"postgres"},
				},
				"file_store_store_type": "local",
				"file_store_rules":      1,
				"auth": map[string]interface{}{
					"providers": 1,
				},
				"services":     1,
				"lets_encrypt": 1,
				"routes":       1,
				"total_events": 1,
			},
			want3: map[string]interface{}{"start_time": ""},
		},
	}
	m, _ := New("", "", false, admin.New("", "clusterID", false, &config.AdminUser{}), &syncman.Manager{}, false)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, got2, got3 := m.generateMetricsRequest(tt.args.project, tt.args.ssl)
			if got != tt.want {
				t.Errorf("generateMetricsRequest() got = %v, want %v", got, tt.want)
			}
			for key, value := range tt.want2 {
				if key == "last_updated" {
					continue
				}
				gotValue, ok := got2[key]
				if !ok {
					t.Errorf("createCrudDocuments() key = %s doesn't exist in result", key)
					continue
				}
				if !reflect.DeepEqual(gotValue, value) {
					t.Errorf("createCrudDocuments() got value = %v %T want = %v %T for key (%s)", gotValue, gotValue, value, value, key)
				}
			}
			if _, ok := got3["start_time"]; !ok {
				t.Errorf("generateMetricsRequest() got3 = %v, want %v", got3, tt.want3)
			}

		})
	}
}
