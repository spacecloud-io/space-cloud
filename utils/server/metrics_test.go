package server

import (
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/config"
)

type TestProjectInfoStub struct {
	modules *config.Modules
	result  map[string]interface{}
}

func Test(t *testing.T) {
	trueCases := 1
	test := []*TestProjectInfoStub{}
	test = append(test, &TestProjectInfoStub{
		modules: &config.Modules{
			Crud: map[string]*config.CrudStub{
				"mongo": &config.CrudStub{
					Enabled: true,
					Collections: map[string]*config.TableRule{
						"collection1": &config.TableRule{},
						"collection2": &config.TableRule{},
					},
				},
				"sql-mysql": &config.CrudStub{
					Enabled: true,
					Collections: map[string]*config.TableRule{
						"collection1": &config.TableRule{},
					},
				},
			},
			Functions: &config.Functions{
				Enabled: true,
				Rules: map[string]map[string]*config.Rule{
					"service 1": map[string]*config.Rule{
						"func1": &config.Rule{},
					},
					"service 2": map[string]*config.Rule{
						"func1": &config.Rule{},
						"func2": &config.Rule{},
						"func3": &config.Rule{},
					},
				},
			},
			FileStore: &config.FileStore{
				Enabled:   true,
				StoreType: "amazon-s3",
				Rules: []*config.FileRule{
					&config.FileRule{},
					&config.FileRule{},
				},
			},
			Static: &config.Static{
				Enabled: true,
				Routes: []*config.StaticRoute{
					&config.StaticRoute{},
					&config.StaticRoute{},
					&config.StaticRoute{},
				},
			},
			Realtime: &config.Realtime{
				Enabled: true,
			},
		},
		result: map[string]interface{}{
			"crud": map[string]interface{}{
				"dbs":         []string{"mongo", "sql-mysql"},
				"collections": 3,
			},
			"functions": map[string]interface{}{
				"enabled":   true,
				"services":  2,
				"functions": 4,
			},
			"fileStore": map[string]interface{}{
				"enabled":   true,
				"storeType": "amazon-s3",
				"rules":     2,
			},
			"static": map[string]interface{}{
				"enabled": true,
				"routes":  3,
			},
			"realtime": map[string]interface{}{
				"enabled": true,
			},
			"auth": []string{},
		},
	})

	for i, testCase := range test {
		res := getProjectInfo(testCase.modules)
		eq := reflect.DeepEqual(testCase.result, res)
		if i < trueCases && !eq {
			t.Error(i+1, ":", "Incorrect Match - Actual: ", res, " Expected: ", testCase.result)
		} else if i > trueCases && eq {
			t.Error(i+1, ":", "Incorrect Match")
		}
	}
}
