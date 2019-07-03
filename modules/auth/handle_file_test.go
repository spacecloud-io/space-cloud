package auth

import (
	"os"
	"testing"

	"github.com/spaceuptech/space-cloud/config"
)

func TestGetFileRule(t *testing.T) {

	var ps = string(os.PathSeparator)

	fileRule := &config.FileRule{
		Prefix: ps,
		Rule:   map[string]*config.Rule{"rule": &config.Rule{Rule: "allow"}},
	}

	fileRule1 := &config.FileRule{
		Prefix: ps + "folder",
		Rule:   map[string]*config.Rule{"rule": &config.Rule{Rule: "allow"}},
	}

	var mod = []struct {
		module   *Module
		testName string
		path     string
	}{
		//Successful Tests
		{testName: "Success", path: ps, module: &Module{fileRules: []*config.FileRule{fileRule, fileRule, fileRule}}},
		{testName: "Success", path: ps + "folder", module: &Module{fileRules: []*config.FileRule{fileRule, fileRule, fileRule}}},
		{testName: "Success", path: ps + "folder", module: &Module{fileRules: []*config.FileRule{fileRule1, fileRule1, fileRule1}}},
		{testName: "Success", path: ps + "folder" + ps + "file", module: &Module{fileRules: []*config.FileRule{fileRule1, fileRule1, fileRule1}}},

		//Error Tests
		{testName: "Fail", path: ps + "NewFolder" + ps + "file", module: &Module{fileRules: []*config.FileRule{fileRule1, fileRule1, fileRule1}}},
		{testName: "Fail", path: ps + ".." + ps + "folder" + ps + "file", module: &Module{fileRules: []*config.FileRule{fileRule, fileRule, fileRule}}},

	}

	for _, test := range mod {
		t.Run(test.testName, func(t *testing.T) {

			data, rules, err1 := (test.module).getFileRule(test.path)
			if test.testName == "Success" {
				if err1 != nil {
					t.Error(data, rules, err1)
				}
			} else {
				if err1 == nil {
					t.Error(data, rules, err1)
				}
			}
		})
	}
}
