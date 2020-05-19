package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud"
	"github.com/spaceuptech/space-cloud/gateway/modules/schema"
)

func hash(s string) string {
	h := sha256.New()
	_, _ = h.Write([]byte(s))
	hashed := hex.EncodeToString(h.Sum(nil))
	return hashed
}

func TestPostProcessMethod(t *testing.T) {
	var authMatchQuery = []struct {
		testName      string
		postProcess   *model.PostProcess
		result        interface{}
		finalResult   interface{}
		IsErrExpected bool
	}{
		{
			testName: "remove from object", IsErrExpected: false,
			result:      map[string]interface{}{"age": 10},
			finalResult: map[string]interface{}{},
			postProcess: &model.PostProcess{PostProcessAction: []model.PostProcessAction{model.PostProcessAction{Action: "remove", Field: "res.age"}}},
		}, {
			testName: "deep remove from object", IsErrExpected: false,
			result:      map[string]interface{}{"k1": map[string]interface{}{"k2": "val"}},
			finalResult: map[string]interface{}{"k1": map[string]interface{}{}},
			postProcess: &model.PostProcess{PostProcessAction: []model.PostProcessAction{model.PostProcessAction{Action: "remove", Field: "res.k1.k2"}}},
		}, {
			testName: "deep remove from object 2", IsErrExpected: false,
			result:      map[string]interface{}{"k1": map[string]interface{}{"k12": "val"}, "k2": "v2"},
			finalResult: map[string]interface{}{"k2": "v2"},
			postProcess: &model.PostProcess{PostProcessAction: []model.PostProcessAction{model.PostProcessAction{Action: "remove", Field: "res.k1"}}},
		}, {
			testName: "remove from array (single element)", IsErrExpected: false,
			result:      []interface{}{map[string]interface{}{"age": 10}},
			finalResult: []interface{}{map[string]interface{}{}},
			postProcess: &model.PostProcess{PostProcessAction: []model.PostProcessAction{model.PostProcessAction{Action: "remove", Field: "res.age"}}},
		}, {
			testName: "remove from array (multiple elements)", IsErrExpected: false,
			result:      []interface{}{map[string]interface{}{"age": 10, "yo": "haha"}, map[string]interface{}{"age": 10}, map[string]interface{}{"yes": 11}},
			finalResult: []interface{}{map[string]interface{}{"yo": "haha"}, map[string]interface{}{}, map[string]interface{}{"yes": 11}},
			postProcess: &model.PostProcess{PostProcessAction: []model.PostProcessAction{model.PostProcessAction{Action: "remove", Field: "res.age"}}},
		}, {
			testName: "Unsuccessful Test Case-remove", IsErrExpected: true,
			result:      map[string]interface{}{"key": "value"},
			finalResult: map[string]interface{}{"key": "value"},
			postProcess: &model.PostProcess{PostProcessAction: []model.PostProcessAction{model.PostProcessAction{Action: "remove", Field: "response.age", Value: nil}}},
		}, {
			testName: "force into object", IsErrExpected: false,
			result:      map[string]interface{}{},
			finalResult: map[string]interface{}{"k1": "v1"},
			postProcess: &model.PostProcess{PostProcessAction: []model.PostProcessAction{model.PostProcessAction{Action: "force", Field: "res.k1", Value: "v1"}}},
		}, {
			testName: "force into array (single)", IsErrExpected: false,
			result:      []interface{}{map[string]interface{}{}},
			finalResult: []interface{}{map[string]interface{}{"k1": "v1"}},
			postProcess: &model.PostProcess{PostProcessAction: []model.PostProcessAction{model.PostProcessAction{Action: "force", Field: "res.k1", Value: "v1"}}},
		}, {
			testName: "force into array (multiple)", IsErrExpected: false,
			result:      []interface{}{map[string]interface{}{}, map[string]interface{}{"k2": "v2"}, map[string]interface{}{"k1": "v2"}},
			finalResult: []interface{}{map[string]interface{}{"k1": "v1"}, map[string]interface{}{"k2": "v2", "k1": "v1"}, map[string]interface{}{"k1": "v1"}},
			postProcess: &model.PostProcess{PostProcessAction: []model.PostProcessAction{model.PostProcessAction{Action: "force", Field: "res.k1", Value: "v1"}}},
		}, {
			testName: "Unsuccessful Test Case-force", IsErrExpected: true,
			result:      map[string]interface{}{"res": map[string]interface{}{"age": 12}},
			finalResult: map[string]interface{}{"res": map[string]interface{}{"age": 12}},
			postProcess: &model.PostProcess{PostProcessAction: []model.PostProcessAction{model.PostProcessAction{Action: "force", Field: "resp.age", Value: "1234"}}},
		}, {
			testName: "Unsuccessful Test Case-neither force nor remove", IsErrExpected: true,
			result:      map[string]interface{}{"res": map[string]interface{}{"age": 12}},
			finalResult: map[string]interface{}{"res": map[string]interface{}{"age": 12}},
			postProcess: &model.PostProcess{PostProcessAction: []model.PostProcessAction{model.PostProcessAction{Action: "forced", Field: "res.age", Value: "1234"}}},
		},
		{testName: "Unsuccessful Test Case-invalid result", IsErrExpected: true,
			result:      1234,
			finalResult: 1234,
			postProcess: &model.PostProcess{PostProcessAction: []model.PostProcessAction{model.PostProcessAction{Action: "forced", Field: "res.age", Value: "1234"}}},
		},
		{testName: "Unsuccessful Test Case-slice of interface as result", IsErrExpected: true,
			result:      []interface{}{1234, "suyash"},
			finalResult: []interface{}{1234, "suyash"},
			postProcess: &model.PostProcess{PostProcessAction: []model.PostProcessAction{model.PostProcessAction{Action: "forced", Field: "res.age", Value: "1234"}}},
		},
		{
			testName:      "invalid field provided for encryption",
			postProcess:   &model.PostProcess{PostProcessAction: []model.PostProcessAction{model.PostProcessAction{Action: "encrypt", Field: "res.age"}}},
			result:        map[string]interface{}{"username": "username1"},
			finalResult:   map[string]interface{}{"username": "username1"},
			IsErrExpected: true,
		},
		{
			testName:      "invalid type of loaded value for encryption",
			postProcess:   &model.PostProcess{PostProcessAction: []model.PostProcessAction{model.PostProcessAction{Action: "encrypt", Field: "res.username"}}},
			result:        map[string]interface{}{"username": 10},
			finalResult:   map[string]interface{}{"username": 10},
			IsErrExpected: true,
		},
		{
			testName:    "valid key in encryption",
			postProcess: &model.PostProcess{PostProcessAction: []model.PostProcessAction{model.PostProcessAction{Action: "encrypt", Field: "res.username"}}},
			result:      map[string]interface{}{"username": "username1"},
			finalResult: map[string]interface{}{"username": base64.StdEncoding.EncodeToString([]byte{5, 120, 168, 68, 222, 6, 202, 246, 108})},
		},
		{
			testName:      "invalid key in encryption",
			postProcess:   &model.PostProcess{PostProcessAction: []model.PostProcessAction{model.PostProcessAction{Action: "encrypt", Field: "res.username"}}},
			result:        map[string]interface{}{"username": "username1"},
			finalResult:   map[string]interface{}{"username": "username1"},
			IsErrExpected: true,
		},
		{
			testName:      "invalid field provided for decryption",
			postProcess:   &model.PostProcess{PostProcessAction: []model.PostProcessAction{model.PostProcessAction{Action: "decrypt", Field: "res.age"}}},
			result:        map[string]interface{}{"username": "username1"},
			finalResult:   map[string]interface{}{"username": "username1"},
			IsErrExpected: true,
		},
		{
			testName:      "invalid type of loaded value for decryption",
			postProcess:   &model.PostProcess{PostProcessAction: []model.PostProcessAction{model.PostProcessAction{Action: "decrypt", Field: "res.username"}}},
			result:        map[string]interface{}{"username": 10},
			finalResult:   map[string]interface{}{"username": 10},
			IsErrExpected: true,
		},
		{
			testName:    "valid key in decryption",
			postProcess: &model.PostProcess{PostProcessAction: []model.PostProcessAction{model.PostProcessAction{Action: "decrypt", Field: "res.username"}}},
			result:      map[string]interface{}{"username": base64.StdEncoding.EncodeToString([]byte{5, 120, 168, 68, 222, 6, 202, 246, 108})},
			finalResult: map[string]interface{}{"username": "username1"},
		},
		{
			testName:      "invalid key in decryption",
			postProcess:   &model.PostProcess{PostProcessAction: []model.PostProcessAction{model.PostProcessAction{Action: "decrypt", Field: "res.username"}}},
			result:        map[string]interface{}{"username": string([]byte{5, 120, 168, 68, 222, 6, 202, 246, 108})},
			finalResult:   map[string]interface{}{"username": string([]byte{5, 120, 168, 68, 222, 6, 202, 246, 108})},
			IsErrExpected: true,
		},
		{
			testName:      "invalid field provided for hash",
			postProcess:   &model.PostProcess{PostProcessAction: []model.PostProcessAction{model.PostProcessAction{Action: "hash", Field: "res.age"}}},
			result:        map[string]interface{}{"password": "password"},
			finalResult:   map[string]interface{}{"password": "password"},
			IsErrExpected: true,
		},
		{
			testName:      "invalid type of loaded value for hash",
			postProcess:   &model.PostProcess{PostProcessAction: []model.PostProcessAction{model.PostProcessAction{Action: "hash", Field: "res.password"}}},
			result:        map[string]interface{}{"password": 10},
			finalResult:   map[string]interface{}{"password": 10},
			IsErrExpected: true,
		},
		{
			testName:    "valid hash",
			postProcess: &model.PostProcess{PostProcessAction: []model.PostProcessAction{model.PostProcessAction{Action: "hash", Field: "res.password"}}},
			result:      map[string]interface{}{"password": "password"},
			finalResult: map[string]interface{}{"password": hash("password")},
		},
	}
	project := "project"
	rule := config.Crud{"mongo": &config.CrudStub{Collections: map[string]*config.TableRule{"tweet": {Rules: map[string]*config.Rule{"aggr": {Rule: "allow", Eval: "Eval", Type: "Type", DB: "mongo", Col: "tweet", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}
	s := schema.Init(crud.Init())
	_ = s.SetConfig(rule, project)
	auth := Init("1", &crud.Module{})
	_ = auth.SetConfig(project, []*config.Secret{}, "Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=", rule, &config.FileStore{}, &config.ServicesModule{}, &config.Eventing{})
	for _, test := range authMatchQuery {
		if test.testName == "invalid key in encryption" || test.testName == "invalid key in decryption" {
			auth.aesKey = base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g")
		} else {
			auth.aesKey = base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=")
		}
		t.Run(test.testName, func(t *testing.T) {
			err := (auth).PostProcessMethod(test.postProcess, test.result)
			if (err != nil) != test.IsErrExpected {
				t.Error("Success GoErr", err, "Want Error", test.IsErrExpected)
				return
			}

			if !reflect.DeepEqual(test.result, test.finalResult) {
				t.Errorf("Error: got %v; wanted %v", test.result, test.finalResult)
				return
			}
		})
	}
}
