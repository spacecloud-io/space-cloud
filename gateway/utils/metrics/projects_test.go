package metrics

import (
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
	"github.com/spaceuptech/space-cloud/gateway/utils/syncman"
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
					ID: "project",
					Modules: &config.Modules{
						Crud: map[string]*config.CrudStub{"db": {
							Type: "postgres",
							Collections: map[string]*config.TableRule{
								"table1": {},
							},
							Enabled: true,
						}},
						Auth: config.Auth{"auth": {
							ID:      "email",
							Enabled: true,
						}},
						Services: &config.ServicesModule{Services: map[string]*config.Service{"service": {}}},
						FileStore: &config.FileStore{
							Enabled:   true,
							StoreType: "local",
							Rules:     []*config.FileRule{{ID: "file"}},
						},
						Eventing: config.Eventing{
							Enabled: true,
							DBAlias: "db",
							Rules: map[string]config.EventingRule{
								"type": {
									Type: "type",
								},
							},
						},
						LetsEncrypt: config.LetsEncrypt{
							ID:                 "letsEncrypt",
							WhitelistedDomains: []string{"1"},
						},
						Routes: config.Routes{{
							ID: "route",
						}},
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
					"db": map[string]interface{}{
						"tables": -2,
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
					t.Errorf("createCrudDocuments() got value = %v %T want = %v %T", gotValue, gotValue, value, value)
				}
			}
			if _, ok := got3["start_time"]; !ok {
				t.Errorf("generateMetricsRequest() got3 = %v, want %v", got3, tt.want3)
			}

		})
	}
}
