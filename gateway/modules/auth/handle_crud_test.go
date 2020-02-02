package auth

import (
	"context"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud"
	"github.com/spaceuptech/space-cloud/gateway/modules/schema"
)

func TestIsCreateOpAuthorised(t *testing.T) {
	var authMatchQuery = []struct {
		module                       *Module
		testName, dbType, col, token string
		project                      string
		IsErrExpected                bool
		rule                         *config.Crud
		value                        model.CreateRequest
		status                       int
	}{
		{
			testName: "Successful Test allow", dbType: "mongo", col: "tweet", project: "project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
			value: model.CreateRequest{
				Document:  map[string]interface{}{"exp": 12},
				Operation: "one",
			},
			IsErrExpected: false,
			status:        200,
		},
		{
			testName: "Unsuccessful Test-Unauthenticated Crud Request", dbType: "pongo", col: "tweet", project: "project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
			value: model.CreateRequest{
				Document:  map[string]interface{}{"exp": 12},
				Operation: "one",
			},
			IsErrExpected: true,
			status:        401,
		},
		{
			testName: "Unsuccessful Test-Incorrect Rule", dbType: "mongo", col: "tweet", project: "projec", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
			value: model.CreateRequest{
				Document:  []interface{}{map[string]interface{}{"exp": 12}},
				Operation: "one",
			},
			IsErrExpected: true,
			status:        403,
		},
		{
			testName: "Successful Test-Batch of Operations", dbType: "mongo", col: "tweet", project: "project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
			value: model.CreateRequest{
				Document:  []interface{}{map[string]interface{}{"exp": 12}},
				Operation: "all",
			},
			IsErrExpected: false,
			status:        200,
		},
	}
	project := "project"
	rule := config.Crud{"mongo": &config.CrudStub{Collections: map[string]*config.TableRule{"tweet": {Rules: map[string]*config.Rule{"create": {Rule: "allow", Eval: "Eval", Type: "Type", DB: "mongo", Col: "tweet", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}
	s := schema.Init(crud.Init(false), false)
	s.SetConfig(rule, project)
	auth := Init("1", &crud.Module{}, s, false)
	auth.SetConfig(project, "", rule, &config.FileStore{}, &config.ServicesModule{}, &config.Eventing{})
	for _, test := range authMatchQuery {
		t.Run(test.testName, func(t *testing.T) {
			r, err := (auth).IsCreateOpAuthorised(context.Background(), test.project, test.dbType, test.col, test.token, &test.value)
			if (err != nil) != test.IsErrExpected {
				t.Error("Received Error-", err, "Wanted Error-", test.IsErrExpected)
			}
			if !reflect.DeepEqual(r, test.status) {
				t.Error("Received Status Code-", r, "Expected Status-", test.status)
			}
		})
	}
}

func TestIsAggregateOpAuthorised(t *testing.T) {
	var authMatchQuery = []struct {
		module                       *Module
		testName, dbType, col, token string
		project                      string
		IsErrExpected                bool
		rule                         *config.Crud
		status                       int
		value                        model.AggregateRequest
	}{
		{
			testName: "Successful Test allow", dbType: "mongo", col: "tweet", project: "project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
			value: model.AggregateRequest{
				Pipeline:  "12",
				Operation: "one",
			},
			IsErrExpected: false,
			status:        200,
		},
		{
			testName: "Unsuccessful Test-Unauthenticated Crud Request", dbType: "pongo", col: "tweet", project: "project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
			value: model.AggregateRequest{
				Pipeline:  map[string]interface{}{"exp": 12},
				Operation: "one",
			},
			IsErrExpected: true,
			status:        401,
		},
		{
			testName: "Unsuccessful Test-Incorrect Rule", dbType: "mongo", col: "tweet", project: "projec", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
			value: model.AggregateRequest{
				Pipeline:  []interface{}{map[string]interface{}{"exp": 12}},
				Operation: "one",
			},
			IsErrExpected: true,
			status:        403,
		},
		{
			testName: "Successful Test-Batch of Operations", dbType: "mongo", col: "tweet", project: "project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
			value: model.AggregateRequest{
				Pipeline:  []interface{}{map[string]interface{}{"exp": 12}},
				Operation: "all",
			},
			IsErrExpected: false,
			status:        200,
		},
	}
	project := "project"
	rule := config.Crud{"mongo": &config.CrudStub{Collections: map[string]*config.TableRule{"tweet": {Rules: map[string]*config.Rule{"aggr": {Rule: "allow", Eval: "Eval", Type: "Type", DB: "mongo", Col: "tweet", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}
	s := schema.Init(crud.Init(false), false)
	s.SetConfig(rule, project)
	auth := Init("1", &crud.Module{}, s, false)
	auth.SetConfig(project, "", rule, &config.FileStore{}, &config.ServicesModule{}, &config.Eventing{})
	for _, test := range authMatchQuery {
		t.Run(test.testName, func(t *testing.T) {
			r, err := (auth).IsAggregateOpAuthorised(context.Background(), test.project, test.dbType, test.col, test.token, &test.value)
			if (err != nil) != test.IsErrExpected {
				t.Error("Success GoErr", err, "Want Error", test.IsErrExpected)
			}
			if !reflect.DeepEqual(r, test.status) {
				t.Error("Received Status Code-", r, "Expected Status-", test.status)
			}
		})
	}
}

func TestPostProcessMethod(t *testing.T) {
	var authMatchQuery = []struct {
		module        *Module
		testName      string
		postProcess   *PostProcess
		result        interface{}
		finalResult   interface{}
		IsErrExpected bool
	}{
		{
			testName: "remove from object", IsErrExpected: false,
			result:      map[string]interface{}{"age": 10},
			finalResult: map[string]interface{}{},
			postProcess: &PostProcess{[]PostProcessAction{PostProcessAction{Action: "remove", Field: "res.age"}}},
		}, {
			testName: "deep remove from object", IsErrExpected: false,
			result:      map[string]interface{}{"k1": map[string]interface{}{"k2": "val"}},
			finalResult: map[string]interface{}{"k1": map[string]interface{}{}},
			postProcess: &PostProcess{[]PostProcessAction{PostProcessAction{Action: "remove", Field: "res.k1.k2"}}},
		}, {
			testName: "deep remove from object 2", IsErrExpected: false,
			result:      map[string]interface{}{"k1": map[string]interface{}{"k12": "val"}, "k2": "v2"},
			finalResult: map[string]interface{}{"k2": "v2"},
			postProcess: &PostProcess{[]PostProcessAction{PostProcessAction{Action: "remove", Field: "res.k1"}}},
		}, {
			testName: "remove from array (single element)", IsErrExpected: false,
			result:      []interface{}{map[string]interface{}{"age": 10}},
			finalResult: []interface{}{map[string]interface{}{}},
			postProcess: &PostProcess{[]PostProcessAction{PostProcessAction{Action: "remove", Field: "res.age"}}},
		}, {
			testName: "remove from array (multiple elements)", IsErrExpected: false,
			result:      []interface{}{map[string]interface{}{"age": 10, "yo": "haha"}, map[string]interface{}{"age": 10}, map[string]interface{}{"yes": 11}},
			finalResult: []interface{}{map[string]interface{}{"yo": "haha"}, map[string]interface{}{}, map[string]interface{}{"yes": 11}},
			postProcess: &PostProcess{[]PostProcessAction{PostProcessAction{Action: "remove", Field: "res.age"}}},
		}, {
			testName: "Unsuccessful Test Case-remove", IsErrExpected: true,
			result:      map[string]interface{}{"key": "value"},
			finalResult: map[string]interface{}{"key": "value"},
			postProcess: &PostProcess{[]PostProcessAction{PostProcessAction{Action: "remove", Field: "response.age", Value: nil}}},
		}, {
			testName: "force into object", IsErrExpected: false,
			result:      map[string]interface{}{},
			finalResult: map[string]interface{}{"k1": "v1"},
			postProcess: &PostProcess{[]PostProcessAction{PostProcessAction{Action: "force", Field: "res.k1", Value: "v1"}}},
		}, {
			testName: "force into array (single)", IsErrExpected: false,
			result:      []interface{}{map[string]interface{}{}},
			finalResult: []interface{}{map[string]interface{}{"k1": "v1"}},
			postProcess: &PostProcess{[]PostProcessAction{PostProcessAction{Action: "force", Field: "res.k1", Value: "v1"}}},
		}, {
			testName: "force into array (multiple)", IsErrExpected: false,
			result:      []interface{}{map[string]interface{}{}, map[string]interface{}{"k2": "v2"}, map[string]interface{}{"k1": "v2"}},
			finalResult: []interface{}{map[string]interface{}{"k1": "v1"}, map[string]interface{}{"k2": "v2", "k1": "v1"}, map[string]interface{}{"k1": "v1"}},
			postProcess: &PostProcess{[]PostProcessAction{PostProcessAction{Action: "force", Field: "res.k1", Value: "v1"}}},
		}, {
			testName: "Unsuccessful Test Case-force", IsErrExpected: true,
			result:      map[string]interface{}{"res": map[string]interface{}{"age": 12}},
			finalResult: map[string]interface{}{"res": map[string]interface{}{"age": 12}},
			postProcess: &PostProcess{[]PostProcessAction{PostProcessAction{Action: "force", Field: "resp.age", Value: "1234"}}},
		}, {
			testName: "Unsuccessful Test Case-neither force nor remove", IsErrExpected: true,
			result:      map[string]interface{}{"res": map[string]interface{}{"age": 12}},
			finalResult: map[string]interface{}{"res": map[string]interface{}{"age": 12}},
			postProcess: &PostProcess{[]PostProcessAction{PostProcessAction{Action: "forced", Field: "res.age", Value: "1234"}}},
		},
		{testName: "Unsuccessful Test Case-invalid result", IsErrExpected: true,
			result:      1234,
			finalResult: 1234,
			postProcess: &PostProcess{[]PostProcessAction{PostProcessAction{Action: "forced", Field: "res.age", Value: "1234"}}},
		},
		{testName: "Unsuccessful Test Case-slice of interface as result", IsErrExpected: true,
			result:      []interface{}{1234, "suyash"},
			finalResult: []interface{}{1234, "suyash"},
			postProcess: &PostProcess{[]PostProcessAction{PostProcessAction{Action: "forced", Field: "res.age", Value: "1234"}}},
		},
	}
	project := "project"
	rule := config.Crud{"mongo": &config.CrudStub{Collections: map[string]*config.TableRule{"tweet": {Rules: map[string]*config.Rule{"aggr": {Rule: "allow", Eval: "Eval", Type: "Type", DB: "mongo", Col: "tweet", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}
	s := schema.Init(crud.Init(false), false)
	s.SetConfig(rule, project)
	auth := Init("1", &crud.Module{}, s, false)
	auth.SetConfig(project, "", rule, &config.FileStore{}, &config.ServicesModule{}, &config.Eventing{})
	for _, test := range authMatchQuery {
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

func TestIsReadOpAuthorised(t *testing.T) {
	var authMatchQuery = []struct {
		module                       *Module
		testName, dbType, col, token string
		project                      string
		IsErrExpected                bool
		rule                         *config.Crud
		value                        model.ReadRequest
		status                       int
	}{
		{
			testName: "Successful Test allow", dbType: "mongo", col: "tweet", project: "project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
			value: model.ReadRequest{
				Find:      map[string]interface{}{"exp": 12},
				Operation: "one",
			},
			IsErrExpected: false,
			status:        200,
		},
		{
			testName: "Unsuccessful Test-Unauthenticated Crud Request", dbType: "pongo", col: "tweet", project: "project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
			value: model.ReadRequest{
				Find:      map[string]interface{}{"exp": 12},
				Operation: "one",
			},
			IsErrExpected: true,
			status:        401,
		},
		{
			testName: "Unsuccessful Test-Unauthorized Crud Request", dbType: "mongo", col: "tweet", project: "projec", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
			value: model.ReadRequest{
				Find:      map[string]interface{}{"exp": 12},
				Operation: "one",
			},
			IsErrExpected: true,
			status:        403,
		},
	}
	project := "project"
	rule := config.Crud{"mongo": &config.CrudStub{Collections: map[string]*config.TableRule{"tweet": {Rules: map[string]*config.Rule{"read": {Rule: "allow", Eval: "Eval", Type: "Type", DB: "mongo", Col: "tweet", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}
	s := schema.Init(crud.Init(false), false)
	s.SetConfig(rule, project)
	auth := Init("1", &crud.Module{}, s, false)
	auth.SetConfig(project, "", rule, &config.FileStore{}, &config.ServicesModule{}, &config.Eventing{})
	for _, test := range authMatchQuery {
		t.Run(test.testName, func(t *testing.T) {
			_, r, err := (auth).IsReadOpAuthorised(context.Background(), test.project, test.dbType, test.col, test.token, &test.value)
			if (err != nil) != test.IsErrExpected {
				t.Error("Success GoErr", err, "Want Error", test.IsErrExpected)
			}
			if !reflect.DeepEqual(r, test.status) {
				t.Error("Received Status Code-", r, "Expected Status-", test.status)
			}
		})
	}
}

func TestIsDeleteOpAuthorised(t *testing.T) {
	var authMatchQuery = []struct {
		module                       *Module
		testName, dbType, col, token string
		project                      string
		IsErrExpected                bool
		rule                         *config.Crud
		value                        model.DeleteRequest
		status                       int
	}{
		{
			testName: "Successful Test allow", dbType: "mongo", col: "tweet", project: "project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
			value: model.DeleteRequest{
				Find:      map[string]interface{}{"exp": 12},
				Operation: "one",
			},
			IsErrExpected: false,
			status:        200,
		},
		{
			testName: "Unsuccessful Test-Unauthenticated Crud Request", dbType: "pongo", col: "tweet", project: "project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
			value: model.DeleteRequest{
				Find:      map[string]interface{}{"exp": 12},
				Operation: "one",
			},
			IsErrExpected: true,
			status:        401,
		},
		{
			testName: "Unsuccessful Test-Unauthorized Crud Request", dbType: "mongo", col: "tweet", project: "projec", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
			value: model.DeleteRequest{
				Find:      map[string]interface{}{"exp": 12},
				Operation: "one",
			},
			IsErrExpected: true,
			status:        403,
		},
	}
	project := "project"
	rule := config.Crud{"mongo": &config.CrudStub{Collections: map[string]*config.TableRule{"tweet": {Rules: map[string]*config.Rule{"delete": {Rule: "allow", Eval: "Eval", Type: "Type", DB: "mongo", Col: "tweet", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}
	s := schema.Init(crud.Init(false), false)
	s.SetConfig(rule, project)
	auth := Init("1", &crud.Module{}, s, false)
	auth.SetConfig(project, "", rule, &config.FileStore{}, &config.ServicesModule{}, &config.Eventing{})

	for _, test := range authMatchQuery {
		t.Run(test.testName, func(t *testing.T) {
			r, err := (auth).IsDeleteOpAuthorised(context.Background(), test.project, test.dbType, test.col, test.token, &test.value)
			if (err != nil) != test.IsErrExpected {
				t.Error("Success GoErr", err, "Want Error", test.IsErrExpected)
			}
			if !reflect.DeepEqual(r, test.status) {
				t.Error("Received Status Code-", r, "Expected Status-", test.status)
			}
		})
	}
}

func TestIsUpdateOpAuthorised(t *testing.T) {
	var authMatchQuery = []struct {
		module                       *Module
		testName, dbType, col, token string
		project                      string
		IsErrExpected                bool
		rule                         *config.Crud
		value                        model.UpdateRequest
		status                       int
	}{
		{
			testName: "Successful Test allow", dbType: "mongo", col: "tweet", project: "project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
			value: model.UpdateRequest{
				Update:    map[string]interface{}{"exp": 12},
				Operation: "one",
			},
			IsErrExpected: false,
			status:        200,
		},
		{
			testName: "Unsuccessful Test-Unauthenticated Crud Request", dbType: "pongo", col: "tweet", project: "project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
			value: model.UpdateRequest{
				Update:    map[string]interface{}{"exp": 12},
				Operation: "one",
			},
			IsErrExpected: true,
			status:        401,
		},
		{
			testName: "Unsuccessful Test-Unauthorized Crud Request", dbType: "mongo", col: "tweet", project: "projec", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
			value: model.UpdateRequest{
				Update:    map[string]interface{}{"exp": 12},
				Operation: "one",
			},
			IsErrExpected: true,
			status:        403,
		},
	}
	project := "project"
	rule := config.Crud{"mongo": &config.CrudStub{Collections: map[string]*config.TableRule{"tweet": {Rules: map[string]*config.Rule{"update": {Rule: "allow", Eval: "Eval", Type: "Type", DB: "mongo", Col: "tweet", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}
	s := schema.Init(crud.Init(false), false)
	s.SetConfig(rule, project)
	auth := Init("1", &crud.Module{}, s, false)
	auth.SetConfig(project, "", rule, &config.FileStore{}, &config.ServicesModule{}, &config.Eventing{})
	for _, test := range authMatchQuery {
		t.Run(test.testName, func(t *testing.T) {
			r, err := (auth).IsUpdateOpAuthorised(context.Background(), test.project, test.dbType, test.col, test.token, &test.value)
			if (err != nil) != test.IsErrExpected {
				t.Error("Success GoErr", err, "Want Error", test.IsErrExpected)
			}
			if !reflect.DeepEqual(r, test.status) {
				t.Error("Received Status Code-", r, "Expected Status-", test.status)
			}
		})
	}
}
