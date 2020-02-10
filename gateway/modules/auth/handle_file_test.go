package auth

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func TestGetFileRule(t *testing.T) {

	var ps = "/"

	fileRule := &config.FileRule{
		Prefix: ps,
		Rule:   map[string]*config.Rule{"rule": &config.Rule{Rule: "allow"}},
	}
	fileRule1 := &config.FileRule{
		Prefix: ps + "folder",
		Rule:   map[string]*config.Rule{"rule": &config.Rule{Rule: "allow"}},
	}
	fileRule2 := &config.FileRule{
		Prefix: ps + "folder/:suyash",
		Rule:   map[string]*config.Rule{"rule": &config.Rule{Rule: "allow"}},
	}
	fileRule3 := &config.FileRule{
		Prefix: ps + "folder/suyash",
		Rule:   map[string]*config.Rule{"rule": &config.Rule{Rule: "deny"}},
	}

	var mod = []struct {
		module        *Module
		IsErrExpected bool
		testName      string
		path          string
		pathParams    map[string]interface{}
		result        *config.FileRule
	}{
		{
			testName: "Valid Test Case-Basic Path", IsErrExpected: false,
			result:     &config.FileRule{Name: "", Prefix: "/", Rule: map[string]*config.Rule{"rule": {Rule: "allow"}}},
			pathParams: map[string]interface{}{}, path: ps,
			module: &Module{fileRules: []*config.FileRule{fileRule, fileRule, fileRule}},
		},
		{
			testName: "Test Case-local file store type", IsErrExpected: false, path: ps,
			module:     &Module{fileRules: []*config.FileRule{fileRule, fileRule, fileRule}, fileStoreType: "local"},
			pathParams: map[string]interface{}{},
			result:     &config.FileRule{Name: "", Prefix: "/", Rule: map[string]*config.Rule{"rule": {Rule: "allow"}}},
		},
		{
			testName: "Valid Test Case-File Rule with folder specified", IsErrExpected: false, path: ps + "folder",
			module:     &Module{fileRules: []*config.FileRule{fileRule1, fileRule1, fileRule1}},
			pathParams: map[string]interface{}{},
			result:     &config.FileRule{Name: "", Prefix: "/folder", Rule: map[string]*config.Rule{"rule": {Rule: "allow"}}},
		},
		{
			testName: "Valid Test Case-Folder with variable mentioned", IsErrExpected: false, path: ps + "folder/:suyash",
			module:     &Module{fileRules: []*config.FileRule{fileRule2, fileRule2, fileRule2}},
			result:     &config.FileRule{Name: "", Prefix: "/folder/:suyash", Rule: map[string]*config.Rule{"rule": {Rule: "allow"}}},
			pathParams: map[string]interface{}{"suyash": ":suyash"},
		},
		{
			testName: "Test case-Rule and Actual Path do not match", IsErrExpected: true, path: ps + "folder" + ps + "file",
			module: &Module{fileRules: []*config.FileRule{fileRule3, fileRule3, fileRule3}},
		},
		{
			testName: "Invalid Path Test Case", IsErrExpected: true, path: ps + "NewFolder" + ps + "file",
			module: &Module{fileRules: []*config.FileRule{fileRule1, fileRule1, fileRule1}},
		},
		{
			testName: "Invalid Case-Provided path should be absolute", IsErrExpected: true, path: ps + ".." + ps + "folder" + ps + "file",
			module: &Module{fileRules: []*config.FileRule{fileRule, fileRule, fileRule}},
		},
	}

	for _, test := range mod {
		t.Run(test.testName, func(t *testing.T) {

			data, rules, err1 := (test.module).getFileRule(test.path)
			if (err1 != nil) != test.IsErrExpected {
				t.Error(data, rules, err1)
			}
			if !test.IsErrExpected {
				if !reflect.DeepEqual(rules, test.result) {
					t.Errorf("getFileRule():Wanted Rule%v,Got Rule%v", test.result, rules)
				}
				//check if valid path paramters are returned
				if !reflect.DeepEqual(data, test.pathParams) {
					t.Errorf("getFileRule():Wanted Path Parameters%v,Got Path Parameters%v", test.pathParams, data)
				}
			}
		})
	}
}

func TestIsFileOpAuthorised(t *testing.T) {
	var authMatchQuery = []struct {
		module                         *Module
		testName, project, token, path string
		op                             utils.FileOpType
		args                           map[string]interface{}
		IsErrExpected                  bool
		result                         *model.PostProcess
	}{
		{
			testName: "Successful Test allow", project: "project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
			module: &Module{fileRules: []*config.FileRule{&config.FileRule{
				Prefix: string(os.PathSeparator),
				Rule:   map[string]*config.Rule{"read": &config.Rule{Rule: "allow"}},
			}},
				project: "project"}, result: &model.PostProcess{},
			IsErrExpected: false, op: "read", args: map[string]interface{}{"age": 12}, path: string(os.PathSeparator),
		},
		{
			testName: "Test Case-Invalid Project Details", project: "projec", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
			module: &Module{fileRules: []*config.FileRule{&config.FileRule{
				Prefix: string(os.PathSeparator),
				Rule:   map[string]*config.Rule{"read": &config.Rule{Rule: "allow"}},
			}},
				project: "project"},
			IsErrExpected: true, op: "read", args: map[string]interface{}{"age": 12}, path: string(os.PathSeparator),
		},
		{
			testName: "Test Case-Not able to parse token", project: "project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
			module: &Module{fileRules: []*config.FileRule{&config.FileRule{
				Prefix: string(os.PathSeparator),
				Rule:   map[string]*config.Rule{"read": &config.Rule{Rule: "allowed"}},
			}},
				project: "project"},
			IsErrExpected: true, op: "read", args: map[string]interface{}{"age": 12}, path: string(os.PathSeparator),
		},
		{
			testName: "Test Case-invalid file rule", project: "project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
			module: &Module{fileRules: []*config.FileRule{&config.FileRule{
				Prefix: string(os.PathSeparator) + "folder",
				Rule:   map[string]*config.Rule{"read": &config.Rule{Rule: "allowed"}},
			}},
				project: "project"},
			IsErrExpected: true, op: "read", args: map[string]interface{}{"age": 12}, path: string(os.PathSeparator),
		},
		{
			testName: "Invalid Test Case-Fields do not match", project: "project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
			module: &Module{fileRules: []*config.FileRule{&config.FileRule{
				Prefix: string(os.PathSeparator),
				Rule:   map[string]*config.Rule{"read": &config.Rule{Rule: "deny"}},
			}},
				project: "project",
				secret:  "mySecretkey",
			},
			IsErrExpected: true, op: "read", args: map[string]interface{}{"params": "age"}, path: string(os.PathSeparator),
		},
	}
	for _, test := range authMatchQuery {
		t.Run(test.testName, func(t *testing.T) {
			result, err := (test.module).IsFileOpAuthorised(context.Background(), test.project, test.token, test.path, test.op, test.args)
			if (err != nil) != test.IsErrExpected {
				t.Error("Got Error-", err, "Want Error-", test.IsErrExpected)
			}
			//check Post Process Result if match rule function is called
			if !test.IsErrExpected && !reflect.DeepEqual(result, test.result) {
				t.Error("Got Result-", result, "Wanted Result-", test.result)
			}
		})
	}
}
