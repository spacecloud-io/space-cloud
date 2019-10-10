package server

import (
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/config"
)

type TestProjectInfoStub struct {
	modules *config.Modules
	static  *config.Static
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
				Services: config.Services{
					"service1": &config.Service{
						Functions: map[string]config.Function{
							"func1": config.Function{
								Rule: &config.Rule{},
							},
						},
					},
					"service2": &config.Service{
						Functions: map[string]config.Function{
							"func1": config.Function{Rule: &config.Rule{}},
							"func2": config.Function{Rule: &config.Rule{}},
							"func3": config.Function{Rule: &config.Rule{}},
						},
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
			Auth: map[string]*config.AuthStub{
				"email": &config.AuthStub{
					Enabled: true,
				},
			},
		},
		static: &config.Static{
			Routes: []*config.StaticRoute{&config.StaticRoute{}, &config.StaticRoute{}, &config.StaticRoute{}},
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
				"enabled":    true,
				"storeTypes": []string{"amazon-s3"},
				"rules":      2,
			},
			"static": map[string]interface{}{
				"routes":         3,
				"internalRoutes": 0,
			},
			"realtime": map[string]interface{}{
				"enabled": true,
			},
			"auth": []string{"email"},
		},
	})

	for i, testCase := range test {
		res := getProjectInfo(testCase.modules, testCase.static)
		eq := reflect.DeepEqual(testCase.result, res)
		if i < trueCases && !eq {
			t.Error(i+1, ":", "Incorrect Match - Actual: ", res, " Expected: ", testCase.result)
		} else if i > trueCases && eq {
			t.Error(i+1, ":", "Incorrect Match")
		}
	}
}
