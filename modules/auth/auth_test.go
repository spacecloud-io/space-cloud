package auth

import (
	"errors"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/modules/crud"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/utils"
)

func TestGetRule(t *testing.T) {
	var authGetRule = []struct {
		dbType, col, testName string
		query                 utils.OperationType
		wantThis              *config.Rule
		authModuleRules       config.Crud
	}{
		// success condition
		{testName: "Successful Test", dbType: "my-sql", col: "collectionName", query: "rule1", wantThis: &config.Rule{Rule: "Rule", Eval: "Eval", Type: "Type", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}, authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "Rule", Eval: "Eval", Type: "Type", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}},
		// error conditiion
		{testName: "Error : Nothing is Provided"},
	}
	successTestCases := 0
	authModule := Init(&crud.Module{})
	for i, test := range authGetRule {
		t.Run(test.testName, func(t *testing.T) {
			(*authModule).rules = test.authModuleRules
			gotThisRule, err := authModule.getRule(test.dbType, test.col, test.query)
			if i <= successTestCases {
				if !reflect.DeepEqual(gotThisRule, test.wantThis) || err != nil {
					t.Error("Success Test ", "Got This ", gotThisRule, "Wanted This ", test.wantThis, "Error ", err)
				}
			} else {

				if (gotThisRule != nil && reflect.DeepEqual(gotThisRule, test.wantThis)) || err == nil {
					t.Error("Error Test", "Got This ", gotThisRule, "Wanted This ", test.wantThis, "Error ", err)
				}
			}

		})
	}
}

// todo : this test generates a tokenstring even if object is empty is this the behaviour we want
func TestCreateToken(t *testing.T) {
	var authCreateToken = []struct {
		testName, wantThis, secretKey string
		object                        map[string]interface{}
	}{
		// success test
		{testName: "Successful Test", secretKey: "mySecretkey", wantThis: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc", object: map[string]interface{}{"token1": "token1value", "token2": "token2value"}},
		// error test
		// {testName: "Error Test : nothing is provided "},
	}
	successTestCases := 0
	authModule := Init(&crud.Module{})
	for i, test := range authCreateToken {
		t.Run(test.testName, func(t *testing.T) {
			authModule.SetSecret(test.secretKey)
			tokenString, err := authModule.CreateToken(test.object)
			if i <= successTestCases {
				if (test.wantThis != tokenString) || err != nil {
					t.Error("Success Test ", "Got This ", tokenString, "Wanted This ", test.wantThis, "Error ", err)
				}
			}

		})
	}
}

func TestGetAuthObj(t *testing.T) {
	var authGetAuthObj = []struct {
		testName, secretKey, token string
		wantThis                   map[string]interface{}
	}{
		{testName: "Successful Test", secretKey: "mySecretkey", wantThis: map[string]interface{}{"token1": "token1value", "token2": "token2value"}, token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc"},
		// error test
		{testName: "Error Test secret key is not provide", wantThis: map[string]interface{}{"token1": "token1value", "token2": "token2value"}, token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc"},
		{testName: "Successful Test", secretKey: "mySecretkey", wantThis: map[string]interface{}{"token1": "token1value", "token2": "token2value"}, token: "abcdefghijkleyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc"},
	}
	successTestCases := 0
	authModule := Init(&crud.Module{})
	for i, test := range authGetAuthObj {
		t.Run(test.testName, func(t *testing.T) {
			authModule.SetSecret(test.secretKey)
			getAuthObect, err := authModule.GetAuthObj(test.token)
			if i <= successTestCases {
				if (!reflect.DeepEqual(getAuthObect, test.wantThis)) || err != nil {
					t.Error("Success Test ", "Got This ", getAuthObect, "Wanted This ", test.wantThis, "Error ", err)
				}
			} else {

				if (reflect.DeepEqual(getAuthObect, test.wantThis)) || err == nil {
					t.Error("Error Test", "Got This ", getAuthObect, "Wanted This ", test.wantThis, "Error ", err)
				}
			}

		})
	}
}

func TestIsAuthorized(t *testing.T) {
	var authIsAuthorized = []struct {
		testName, dbType, col string
		project               string
		query                 utils.OperationType
		err                   error
		args                  map[string]interface{}
		authModuleRules       config.Crud
	}{
		{testName: "Successful Test allow", dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "allow", Eval: "Eval", Type: "Type", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}},
		{testName: "Successful Test authenticated", dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "authenticated", Eval: "Eval", Type: "Type", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}},
		{testName: "Successful Test and ==", dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "and", Eval: "==", Type: "string", F1: "interfaceString1", F2: "interfaceString1", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}}}},
		{testName: "Successful Test and !=", dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "and", Eval: "!=", Type: "string", F1: "interfaceString1", F2: "interfaceString2", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}}}},
		{testName: "Successful Test and ==", dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "and", Eval: "==", Type: "number", F1: 1.0, F2: 1.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}}}},
		{testName: "Successful Test and <=", dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "and", Eval: "<=", Type: "number", F1: 2.0, F2: 3.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}}}},
		{testName: "Successful Test and >=", dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "Rule", Eval: ">=", Type: "number", F1: 3.0, F2: 2.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}}}},
		{testName: "Successful Test and <", dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "Rule", Eval: "<", Type: "number", F1: 1.0, F2: 2.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}}}},
		{testName: "Successful Test and >", dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "Rule", Eval: ">", Type: "number", F1: 3.0, F2: 2.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}}}},
		{testName: "Successful Test and !=", dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "Rule", Eval: "!=", Type: "number", F1: 3.0, F2: 2.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}}}},
		{testName: "Successful Test and ==", dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "Rule", Eval: "==", Type: "bool", F1: true, F2: true, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}}}},
		{testName: "Successful Test and !=", dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "Rule", Eval: "!=", Type: "bool", F1: true, F2: false, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}}}},

		{testName: "Successful Test or ==", dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "and", Eval: "==", Type: "string", F1: "interfaceString1", F2: "interfaceString1", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}}}},
		{testName: "Successful Test or !=", dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "and", Eval: "!=", Type: "string", F1: "interfaceString1", F2: "interfaceString2", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}}}},
		{testName: "Successful Test or ==", dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "and", Eval: "==", Type: "number", F1: 1.0, F2: 1.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}}}},
		{testName: "Successful Test or <=", dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "and", Eval: "<=", Type: "number", F1: 2.0, F2: 3.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}}}},
		{testName: "Successful Test or >=", dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "Rule", Eval: ">=", Type: "number", F1: 3.0, F2: 2.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}}}},
		{testName: "Successful Test or <", dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "Rule", Eval: "<", Type: "number", F1: 1.0, F2: 2.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}}}},
		{testName: "Successful Test or >", dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "Rule", Eval: ">", Type: "number", F1: 3.0, F2: 2.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}}}},
		{testName: "Successful Test or !=", dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "Rule", Eval: "!=", Type: "number", F1: 3.0, F2: 2.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}}}},
		{testName: "Successful Test or ==", dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "Rule", Eval: "==", Type: "bool", F1: true, F2: true, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}}}},
		{testName: "Successful Test or !=", dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "Rule", Eval: "!=", Type: "bool", F1: true, F2: false, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}}}},

		{testName: "Successful Test match ==", dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "match", Eval: "==", Type: "string", F1: "interfaceString1", F2: "interfaceString1", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}},
		{testName: "Successful Test match !=", dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "match", Eval: "!=", Type: "string", F1: "interfaceString1", F2: "interfaceString2", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}},
		{testName: "Successful Test match ==", dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "match", Eval: "==", Type: "number", F1: 1.0, F2: 1.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}},
		{testName: "Successful Test match <=", dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "match", Eval: "<=", Type: "number", F1: 1.0, F2: 1.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}},
		{testName: "Successful Test match >=", dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "match", Eval: ">=", Type: "number", F1: 1.0, F2: 1.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}},
		{testName: "Successful Test match <", dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "match", Eval: "<", Type: "number", F1: 1.0, F2: 2.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}},
		{testName: "Successful Test match >", dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "match", Eval: ">", Type: "number", F1: 2.0, F2: 1.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}},
		{testName: "Successful Test match !=", dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "match", Eval: "!=", Type: "number", F1: 2.0, F2: 1.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}},
		{testName: "Successful Test match ==", dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "match", Eval: "==", Type: "bool", F1: true, F2: true, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}},
		{testName: "Successful Test match !=", dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "match", Eval: "!=", Type: "bool", F1: false, F2: true, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}},

		// {testName: "Successful Test", dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "and", Clauses: []*config.Rule{{}}}}}}}}},

		// {testName: "Successful Test", dbType: "my-sql", col: "collectionName", project: "project", query: "rule1"},
		// error test
		{testName: "Error Test deny", err: ErrIncorrectMatch, dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "deny", Eval: "!=", Type: "bool", F1: false, F2: true, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}},
		{testName: "Error Test default", err: ErrIncorrectMatch, dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "", Eval: "!=", Type: "bool", F1: false, F2: true, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}},
		{testName: "Error Test or ==", err: ErrIncorrectMatch, dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "and", Eval: "==", Type: "", F1: "interfaceString1", F2: "interfaceString1", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}}}},
		{testName: "Error Test and ==", err: ErrIncorrectMatch, dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "and", Eval: "==", Type: "", F1: "interfaceString1", F2: "interfaceString1", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}}}},
		{testName: "Error : Nothing is Provided for get rule func", project: "project", err: ErrRuleNotFound},
		{testName: "Error Test and == f1 is int instead of string", err: ErrIncorrectRuleFieldType, dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "and", Eval: "==", Type: "string", F1: 1, F2: "interfaceString1", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}}}},
		{testName: "Error Test and == f2 is int instead of string", err: ErrIncorrectRuleFieldType, dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "and", Eval: "==", Type: "string", F1: "interfaceString1", F2: 1, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}}}},
		{testName: "Error Test and eval is not provide string", err: ErrIncorrectMatch, dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "and", Eval: "", Type: "string", F1: "interfaceString1", F2: "interfaceString2", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}}}},
		{testName: "Error Test and == f1 number is not provided", err: errors.New("Store: Cloud not load value"), dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "and", Eval: "==", Type: "number", F2: 1.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}}}},
		{testName: "Error Test and == f2 number is not provided", err: errors.New("Store: Cloud not load value"), dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "and", Eval: "==", Type: "number", F1: 1.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}}}},
		{testName: "Error Test and eval is not provided number", err: ErrIncorrectRuleFieldType, dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "and", Eval: "", Type: "number", F1: 1.0, F2: 1.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}}}},
		{testName: "Error Test and == bool f1 bool is not provided", err: errors.New("Store: Cloud not load value"), dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "Rule", Eval: "==", Type: "bool", F2: true, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}}}},
		{testName: "Error Test and == bool f2 bool is not provided", err: errors.New("Store: Cloud not load value"), dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "Rule", Eval: "==", Type: "bool", F1: true, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}}}},
		{testName: "Error Test and == bool eval is not provided", err: ErrIncorrectRuleFieldType, dbType: "my-sql", col: "collectionName", project: "project", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "Rule", Eval: "", Type: "bool", F1: true, F2: true, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}}}},
	}
	successTestCases := 31
	authModule := Init(&crud.Module{})
	for i, test := range authIsAuthorized {
		t.Run(test.testName, func(t *testing.T) {
			authModule.rules = test.authModuleRules
			authModule.project = "project"
			err := authModule.IsAuthorized(test.project, test.dbType, test.col, test.query, test.args)
			if i <= successTestCases {
				if !reflect.DeepEqual(err, test.err) {
					t.Error("Success Test ", "| Got This | ", err, "| Wanted This |", test.err)
				}
			} else {
				if !reflect.DeepEqual(err, test.err) {
					t.Error("Error Test", "| Got This |", err, "| Wanted This |", test.err)
				}
			}

		})
	}
}

func TestIsAuthenticated(t *testing.T) {
	var authIsAuthenticated = []struct {
		testName, dbType, col, token, secretKey string
		query                                   utils.OperationType
		wantThis                                map[string]interface{}
		authModuleRules                         config.Crud
	}{
		{testName: "Successful Test", secretKey: "mySecretkey", wantThis: map[string]interface{}{"token1": "token1value", "token2": "token2value"}, token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc", dbType: "my-sql", col: "collectionName", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "Rule", Eval: "Eval", Type: "Type", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}},
		{testName: "Successful Test rule is allow", secretKey: "mySecretkey", wantThis: map[string]interface{}{}, token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc", dbType: "my-sql", col: "collectionName", query: "rule1", authModuleRules: config.Crud{"my-sql": &config.CrudStub{Collections: map[string]*config.TableRule{"collectionName": &config.TableRule{Rules: map[string]*config.Rule{"rule1": &config.Rule{Rule: "allow", Eval: "Eval", Type: "Type", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}}}}},

		// error test
	}
	successTestCases := 1
	authModule := Init(&crud.Module{})
	for i, test := range authIsAuthenticated {
		t.Run(test.testName, func(t *testing.T) {
			authModule.SetConfig("project", test.secretKey, test.authModuleRules, nil)
			getAuthObect, err := authModule.IsAuthenticated(test.token, test.dbType, test.col, test.query)
			if i <= successTestCases {
				if (!reflect.DeepEqual(getAuthObect, test.wantThis)) || err != nil {
					t.Error("Success Test ", "Got This ", getAuthObect, "Wanted This ", test.wantThis, "Error ", err)
				}
			} else {

				if (reflect.DeepEqual(getAuthObect, test.wantThis)) || err == nil {
					t.Error("Error Test", "Got This ", getAuthObect, "Wanted This ", test.wantThis, "Error ", err)
				}
			}

		})
	}
}
