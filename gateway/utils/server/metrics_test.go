package server

import (
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

type TestProjectInfoStub struct {
	projects []*config.Project
	static   *config.Static
	result   map[string]interface{}
}

func Test(t *testing.T) {
	trueCases := 1
	test := []*TestProjectInfoStub{}
	test = append(test, &TestProjectInfoStub{
		projects: []*config.Project{
			&config.Project{
				Modules: &config.Modules{
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
				},
			},
		},
		result: map[string]interface{}{
			"crud": map[string]interface{}{
				"dbs":         []string{"mongo", "sql-mysql"},
				"collections": 3,
			},
			"functions": map[string]interface{}{
				"enabled":   false,
				"services":  0,
				"functions": 0,
			},
			"fileStore": map[string]interface{}{
				"enabled":    false,
				"storeTypes": []string{},
				"rules":      0,
			},
			"auth": []string{},
		},
	})

	for i, testCase := range test {
		res := getProjectInfo(testCase.projects)
		eq := reflect.DeepEqual(testCase.result, res)
		if i < trueCases && !eq {
			t.Error(i+1, ":", "Incorrect Match - Actual: ", res, " Expected: ", testCase.result)
		} else if i > trueCases && eq {
			t.Error(i+1, ":", "Incorrect Match")
		}
	}
}
