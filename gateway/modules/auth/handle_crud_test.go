package auth

import (
	"context"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud"
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
			testName: "Unsuccessful Test-Unauthenticated crud Request", dbType: "pongo", col: "tweet", project: "project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
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
	dbRules := config.DatabaseRules{
		"chicago--project--db-rule--mongo-tweet-rule": &config.DatabaseRule{Rules: map[string]*config.Rule{"create": {Rule: "allow", Eval: "Eval", Type: "Type", DB: "mongo", Col: "tweet", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}},
	}
	auth := Init("chicago", "1", &crud.Module{}, nil)
	if er := auth.SetConfig(context.TODO(), "local", &config.ProjectConfig{ID: project, Secrets: []*config.Secret{}}, dbRules, config.DatabasePreparedQueries{}, config.FileStoreRules{}, config.Services{}, config.EventingRules{}, config.SecurityFunctions{}); er != nil {
		t.Errorf("error setting config of auth module  - %s", er.Error())
	}
	for _, test := range authMatchQuery {
		t.Run(test.testName, func(t *testing.T) {
			_, err := (auth).IsCreateOpAuthorised(context.Background(), test.project, test.dbType, test.col, test.token, &test.value)
			if (err != nil) != test.IsErrExpected {
				t.Error("Received Error-", err, "Wanted Error-", test.IsErrExpected)
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
			testName: "Unsuccessful Test-Unauthenticated crud Request", dbType: "pongo", col: "tweet", project: "project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
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
	dbRules := config.DatabaseRules{
		"chicago--project--db-rule--mongo-tweet-rule":   &config.DatabaseRule{Rules: map[string]*config.Rule{"aggr": {Rule: "allow", Eval: "Eval", Type: "Type", DB: "mongo", Col: "tweet", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}},
		"chicago--project--db-rule--mongo-default-rule": &config.DatabaseRule{Rules: map[string]*config.Rule{"aggr": {Rule: "allow"}, "read": {Rule: "allow"}, "update": {Rule: "allow"}, "delete": {Rule: "allow"}}},
	}
	auth := Init("chicago", "1", &crud.Module{}, nil)
	if er := auth.SetConfig(context.TODO(), "local", &config.ProjectConfig{ID: project, Secrets: []*config.Secret{}}, dbRules, config.DatabasePreparedQueries{}, config.FileStoreRules{}, config.Services{}, config.EventingRules{}, config.SecurityFunctions{}); er != nil {
		t.Errorf("error setting config of auth module  - %s", er.Error())
	}
	for _, test := range authMatchQuery {
		t.Run(test.testName, func(t *testing.T) {
			_, err := (auth).IsAggregateOpAuthorised(context.Background(), test.project, test.dbType, test.col, test.token, &test.value)
			if (err != nil) != test.IsErrExpected {
				t.Error("Success GoErr", err, "Want Error", test.IsErrExpected)
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
			testName: "Unsuccessful Test-Unauthenticated crud Request", dbType: "pongo", col: "tweet", project: "project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
			value: model.ReadRequest{
				Find:      map[string]interface{}{"exp": 12},
				Operation: "one",
			},
			IsErrExpected: true,
			status:        401,
		},
		{
			testName: "Unsuccessful Test-Unauthorized crud Request", dbType: "mongo", col: "tweet", project: "projec", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
			value: model.ReadRequest{
				Find:      map[string]interface{}{"exp": 12},
				Operation: "one",
			},
			IsErrExpected: true,
			status:        403,
		},
	}
	project := "project"
	dbRules := config.DatabaseRules{
		"chicago--project--db-rule--mongo-tweet-rule": &config.DatabaseRule{Rules: map[string]*config.Rule{"read": {Rule: "allow", Eval: "Eval", Type: "Type", DB: "mongo", Col: "tweet", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}},
	}
	auth := Init("chicago", "1", &crud.Module{}, nil)
	if er := auth.SetConfig(context.TODO(), "local", &config.ProjectConfig{ID: project, Secrets: []*config.Secret{}}, dbRules, config.DatabasePreparedQueries{}, config.FileStoreRules{}, config.Services{}, config.EventingRules{}, config.SecurityFunctions{}); er != nil {
		t.Errorf("error setting config of auth module  - %s", er.Error())
	}
	for _, test := range authMatchQuery {
		t.Run(test.testName, func(t *testing.T) {
			_, _, err := (auth).IsReadOpAuthorised(context.Background(), test.project, test.dbType, test.col, test.token, &test.value, model.ReturnWhereStub{})
			if (err != nil) != test.IsErrExpected {
				t.Error("Success GoErr", err, "Want Error", test.IsErrExpected)
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
			testName: "Unsuccessful Test-Unauthenticated crud Request", dbType: "pongo", col: "tweet", project: "project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
			value: model.DeleteRequest{
				Find:      map[string]interface{}{"exp": 12},
				Operation: "one",
			},
			IsErrExpected: true,
			status:        401,
		},
		{
			testName: "Unsuccessful Test-Unauthorized crud Request", dbType: "mongo", col: "tweet", project: "projec", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
			value: model.DeleteRequest{
				Find:      map[string]interface{}{"exp": 12},
				Operation: "one",
			},
			IsErrExpected: true,
			status:        403,
		},
	}
	project := "project"
	dbRules := config.DatabaseRules{
		"chicago--project--db-rule--mongo-tweet-rule": &config.DatabaseRule{Rules: map[string]*config.Rule{"delete": {Rule: "allow", Eval: "Eval", Type: "Type", DB: "mongo", Col: "tweet", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}},
	}
	auth := Init("chicago", "1", &crud.Module{}, nil)
	if er := auth.SetConfig(context.TODO(), "local", &config.ProjectConfig{ID: project, Secrets: []*config.Secret{}}, dbRules, config.DatabasePreparedQueries{}, config.FileStoreRules{}, config.Services{}, config.EventingRules{}, config.SecurityFunctions{}); er != nil {
		t.Errorf("error setting config of auth module  - %s", er.Error())
	}
	for _, test := range authMatchQuery {
		t.Run(test.testName, func(t *testing.T) {
			_, err := (auth).IsDeleteOpAuthorised(context.Background(), test.project, test.dbType, test.col, test.token, &test.value)
			if (err != nil) != test.IsErrExpected {
				t.Error("Success GoErr", err, "Want Error", test.IsErrExpected)
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
			testName: "Unsuccessful Test-Unauthenticated crud Request", dbType: "pongo", col: "tweet", project: "project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
			value: model.UpdateRequest{
				Update:    map[string]interface{}{"exp": 12},
				Operation: "one",
			},
			IsErrExpected: true,
			status:        401,
		},
		{
			testName: "Unsuccessful Test-Unauthorized crud Request", dbType: "mongo", col: "tweet", project: "projec", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
			value: model.UpdateRequest{
				Update:    map[string]interface{}{"exp": 12},
				Operation: "one",
			},
			IsErrExpected: true,
			status:        403,
		},
	}
	project := "project"
	dbRules := config.DatabaseRules{
		"chicago--project--db-rule--mongo-tweet-rule": &config.DatabaseRule{Rules: map[string]*config.Rule{"update": {Rule: "allow", Eval: "Eval", Type: "Type", DB: "mongo", Col: "tweet", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}},
	}
	auth := Init("chicago", "1", &crud.Module{}, nil)
	if er := auth.SetConfig(context.TODO(), "local", &config.ProjectConfig{ID: project, Secrets: []*config.Secret{}}, dbRules, config.DatabasePreparedQueries{}, config.FileStoreRules{}, config.Services{}, config.EventingRules{}, config.SecurityFunctions{}); er != nil {
		t.Errorf("error setting config of auth module  - %s", er.Error())
	}
	for _, test := range authMatchQuery {
		t.Run(test.testName, func(t *testing.T) {
			_, err := (auth).IsUpdateOpAuthorised(context.Background(), test.project, test.dbType, test.col, test.token, &test.value)
			if (err != nil) != test.IsErrExpected {
				t.Error("Success GoErr", err, "Want Error", test.IsErrExpected)
			}
		})
	}
}

func TestIsPreparedQueryAuthorised(t *testing.T) {
	var authMatchQuery = []struct {
		module                      *Module
		testName, dbType, id, token string
		project                     string
		IsErrExpected               bool
		rule                        *config.Crud
		value                       model.PreparedQueryRequest
		status                      int
	}{
		{
			testName: "Successful Test allow", dbType: "mongo", id: "tweet", project: "project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
			value: model.PreparedQueryRequest{
				Params: map[string]interface{}{"exp": 12},
			},
			IsErrExpected: false,
			status:        200,
		},
		{
			testName: "Unsuccessful Test-Unauthenticated crud Request", dbType: "pongo", id: "tweet", project: "project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
			value: model.PreparedQueryRequest{
				Params: map[string]interface{}{"exp": 12},
			},
			IsErrExpected: true,
			status:        401,
		},
		{
			testName: "Unsuccessful Test-Unauthorized crud Request", dbType: "mongo", id: "tweet", project: "projec", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
			value: model.PreparedQueryRequest{
				Params: map[string]interface{}{"exp": 12},
			},
			IsErrExpected: true,
			status:        403,
		},
	}
	project := "project"
	dbPrepRules := config.DatabasePreparedQueries{
		config.GenerateResourceID("chicago", "project", config.ResourceDatabasePreparedQuery, "mongo", "tweet"): &config.DatbasePreparedQuery{Rule: &config.Rule{Rule: "allow", Eval: "Eval", Type: "Type", DB: "mongo", Col: "tweet", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
	}
	auth := Init("chicago", "1", &crud.Module{}, nil)
	if er := auth.SetConfig(context.TODO(), "local", &config.ProjectConfig{ID: project, Secrets: []*config.Secret{}}, config.DatabaseRules{}, dbPrepRules, config.FileStoreRules{}, config.Services{}, config.EventingRules{}, config.SecurityFunctions{}); er != nil {
		t.Errorf("error setting config of auth module  - %s", er.Error())
	}
	for _, test := range authMatchQuery {
		t.Run(test.testName, func(t *testing.T) {
			_, _, err := (auth).IsPreparedQueryAuthorised(context.Background(), test.project, test.dbType, test.id, test.token, &test.value)
			if (err != nil) != test.IsErrExpected {
				t.Error("Success GoErr", err, "Want Error", test.IsErrExpected)
			}
		})
	}
}

func Test_authenticatePreparedQueryRequest(t *testing.T) {
	tests := []struct {
		name                        string
		module                      *Module
		dbAlias, id, token, project string
		wantRule                    *config.Rule
		wantAuth                    map[string]interface{}
		wantErr                     bool
	}{
		// TODO: Add test cases.
		{
			name: "Successful Test for authenticate Prepared Query Request", dbAlias: "mongo", id: "tweet", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
			module:   &Module{clusterID: "chicago", dbPrepQueryRules: config.DatabasePreparedQueries{config.GenerateResourceID("chicago", "project", config.ResourceDatabasePreparedQuery, "mongo", "tweet"): &config.DatbasePreparedQuery{ID: "tweet", Rule: &config.Rule{Rule: "allow", Eval: "Eval", Type: "Type", DB: "mongo", Col: "tweet", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
			project:  "project",
			wantRule: &config.Rule{Rule: "allow", Eval: "Eval", Type: "Type", DB: "mongo", Col: "tweet", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}},
			wantAuth: nil,
			wantErr:  false,
		},
		{
			name: "Unsuccessful Test-authenticate Prepared Query Request", dbAlias: "pongo", id: "tweet", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
			module:   &Module{clusterID: "chicago", dbPrepQueryRules: config.DatabasePreparedQueries{config.GenerateResourceID("chicago", "project", config.ResourceDatabasePreparedQuery, "mongo", "tweet"): &config.DatbasePreparedQuery{ID: "tweet", Rule: &config.Rule{Rule: "allow", Eval: "Eval", Type: "Type", DB: "mongo", Col: "tweet", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
			project:  "project",
			wantRule: nil,
			wantAuth: nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRule, gotAuth, err := (tt.module).authenticatePreparedQueryRequest(context.Background(), tt.project, tt.dbAlias, tt.id, tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("Module.authenticatePreparedQueryRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRule, tt.wantRule) {
				t.Errorf("Module.authenticatePreparedQueryRequest() gotRule = %v, want %v", gotRule, tt.wantRule)
			}
			if !reflect.DeepEqual(gotAuth, tt.wantAuth) {
				t.Errorf("Module.authenticatePreparedQueryRequest() gotAuth = %v, want %v", gotAuth, tt.wantAuth)
			}
		})
	}
}

func Test_getPrepareQueryRule(t *testing.T) {
	tests := []struct {
		name        string
		module      *Module
		dbAlias, id string
		project     string
		want        *config.Rule
		wantErr     bool
	}{
		{
			name: "Successful Test to get Prepare Query Rule", dbAlias: "mongo", id: "tweet",
			module:  &Module{clusterID: "chicago", dbPrepQueryRules: config.DatabasePreparedQueries{config.GenerateResourceID("chicago", "project", config.ResourceDatabasePreparedQuery, "mongo", "tweet"): &config.DatbasePreparedQuery{ID: "tweet", Rule: &config.Rule{Rule: "allow", Eval: "Eval", Type: "Type", DB: "mongo", Col: "tweet", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
			project: "project",
			want:    &config.Rule{Rule: "allow", Eval: "Eval", Type: "Type", DB: "mongo", Col: "tweet", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}},
			wantErr: false,
		},
		{
			name: "Unsuccessful Test- Prepared Query Rule Request", dbAlias: "pongo", id: "tweet",
			module:  &Module{clusterID: "chicago", dbPrepQueryRules: config.DatabasePreparedQueries{config.GenerateResourceID("chicago", "project", config.ResourceDatabasePreparedQuery, "mongo", "tweet"): &config.DatbasePreparedQuery{ID: "tweet", Rule: &config.Rule{Rule: "allow", Eval: "Eval", Type: "Type", DB: "mongo", Col: "tweet", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
			project: "project",
			want:    nil,
			wantErr: true,
		},
		{
			name: "Successful Test to get default Prepare Query Rule", dbAlias: "mongo", id: "tweet",
			module:  &Module{clusterID: "chicago", dbPrepQueryRules: config.DatabasePreparedQueries{config.GenerateResourceID("chicago", "project", config.ResourceDatabasePreparedQuery, "mongo", "tweet"): &config.DatbasePreparedQuery{ID: "tweet", Rule: &config.Rule{Rule: "allow", Eval: "Eval", Type: "Type", DB: "mongo", Col: "default", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
			project: "project",
			want:    &config.Rule{Rule: "allow", Eval: "Eval", Type: "Type", DB: "mongo", Col: "default", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := (tt.module).getPrepareQueryRule(context.Background(), tt.project, tt.dbAlias, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Module.getPrepareQueryRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Module.getPrepareQueryRule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModule_getCrudRule(t *testing.T) {
	type fields struct {
		rules map[string]*config.DatabaseRule
	}
	type args struct {
		dbAlias string
		col     string
		project string
		query   model.OperationType
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *config.Rule
		wantErr bool
	}{
		{
			name:   "valid col",
			fields: fields{rules: map[string]*config.DatabaseRule{config.GenerateResourceID("chicago", "project", config.ResourceDatabaseRule, "db", "col", "rule"): {Rules: map[string]*config.Rule{"op": {Type: "allow"}}}}},
			args:   args{project: "project", query: "op", dbAlias: "db", col: "col"},
			want:   &config.Rule{Type: "allow"},
		},
		{
			name:   "valid default - wrong op",
			fields: fields{rules: map[string]*config.DatabaseRule{config.GenerateResourceID("chicago", "project", config.ResourceDatabaseRule, "db", "default", "rule"): {Rules: map[string]*config.Rule{"op": {Type: "default"}}}}},
			args:   args{project: "project", query: "op", dbAlias: "db", col: "col"},
			want:   &config.Rule{Type: "default"},
		},
		{
			name:    "wrong db",
			fields:  fields{rules: map[string]*config.DatabaseRule{config.GenerateResourceID("chicago", "project", config.ResourceDatabaseRule, "db", "col", "rule"): {Rules: map[string]*config.Rule{"op": {Type: "allow"}}}}},
			args:    args{project: "project", query: "op", dbAlias: "db-bad", col: "col"},
			wantErr: true,
		},
		{
			name:    "wrong col",
			fields:  fields{rules: map[string]*config.DatabaseRule{config.GenerateResourceID("chicago", "project", config.ResourceDatabaseRule, "db", "col", "rule"): {Rules: map[string]*config.Rule{"op": {Type: "allow"}}}}},
			args:    args{project: "project", query: "op", dbAlias: "db", col: "col-bad"},
			wantErr: true,
		},
		{
			name:    "invalid default - wrong op",
			fields:  fields{rules: map[string]*config.DatabaseRule{config.GenerateResourceID("chicago", "project", config.ResourceDatabaseRule, "db", "col", "rule"): {Rules: map[string]*config.Rule{}}, "default": {Rules: map[string]*config.Rule{"op": {Type: "default"}}}}},
			args:    args{project: "project", query: "op-bad", dbAlias: "db", col: "col"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Module{
				clusterID: "chicago",
				dbRules:   tt.fields.rules,
			}
			got, err := m.getCrudRule(context.Background(), tt.args.project, tt.args.dbAlias, tt.args.col, tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCrudRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCrudRule() got = %v, want %v", got, tt.want)
			}
		})
	}
}
