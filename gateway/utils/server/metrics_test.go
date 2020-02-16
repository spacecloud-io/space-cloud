package server

import (
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

func Test_getProjectInfo(t *testing.T) {
	type args struct {
		config *config.Modules
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "successful test",
			args: args{
				config: &config.Modules{
					Crud: config.Crud{
						"mongo": &config.CrudStub{
							Type:    "mongo",
							Enabled: true,
							Collections: map[string]*config.TableRule{
								"collection1": &config.TableRule{},
								"collection2": &config.TableRule{},
							},
						},
					},
				},
			},
			want: map[string]interface{}{
				"auth":      []string{},
				"crud":      map[string]interface{}{"dbs": []string{"mongo"}, "collections": 2},
				"functions": map[string]interface{}{"enabled": false, "services": 0, "functions": 0},
				"fileStore": map[string]interface{}{"enabled": false, "storeTypes": []string{}, "rules": 0},
				"static":    map[string]interface{}{"routes": 0, "internalRoutes": 0},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getProjectInfo(tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getProjectInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
