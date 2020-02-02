package auth

import (
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud"
	"github.com/spaceuptech/space-cloud/gateway/modules/schema"
	"golang.org/x/net/context"
)

func TestMatch_Rule(t *testing.T) {
	var testCases = []struct {
		name          string
		IsErrExpected bool
		project       string
		rule          *config.Rule
		args          map[string]interface{}
		auth          map[string]interface{}
	}{
		{
			name: "internal sc user", IsErrExpected: false, project: "default",
			rule: &config.Rule{Rule: "deny", Eval: "!=", Type: "string", F1: "interfaceString1", F2: "interfaceString2", DB: "mongo", Col: "default", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}},
			args: map[string]interface{}{"string1": "interface1", "string2": "interface2"},
			auth: map[string]interface{}{"id": "internal-sc-user", "roll": "1234"},
		}, {
			name: "allow", IsErrExpected: false, project: "default",
			rule: &config.Rule{Rule: "allow", Eval: "!=", Type: "string", F1: "interfaceString1", F2: "interfaceString2", DB: "mongo", Col: "default", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}},
			args: map[string]interface{}{"string1": "interface1", "string2": "interface2"},
			auth: map[string]interface{}{"id": "internal-sc-user", "roll": "1234"},
		}, {
			name: "deny rule", IsErrExpected: true, project: "default",
			rule: &config.Rule{Rule: "deny", Eval: "!=", Type: "string", F1: "interfaceString1", F2: "interfaceString2", DB: "mongo", Col: "default", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}},
			args: map[string]interface{}{"string1": "interface1", "string2": "interface2"},
			auth: map[string]interface{}{"id": "internal-sc", "roll": "1234"},
		},
		{name: "Invalid matchand rule", IsErrExpected: true, project: "default",
			rule: &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "deny", Eval: "!="},
				{Rule: "allow", Eval: "!="}}},
			args: map[string]interface{}{"string1": "interface1", "string2": "interface2"},
			auth: map[string]interface{}{"id": "internal-sc", "roll": "1234"},
		},
		{name: "Valid matchand rule", IsErrExpected: false, project: "default",
			rule: &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "allow", Eval: "!="},
				{Rule: "allow", Eval: "!="}}},
			args: map[string]interface{}{"string1": "interface1", "string2": "interface2"},
			auth: map[string]interface{}{"id": "internal-sc", "roll": "1234"},
		},
		{name: "match rule-string", IsErrExpected: false, project: "default",
			rule: &config.Rule{Rule: "match", Eval: "!=", Type: "string", F1: "interfaceString1", F2: "interfaceString2", DB: "mongo", Col: "default", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}},
			args: map[string]interface{}{"string1": "interface1", "string2": "interface2"},
			auth: map[string]interface{}{"id": "internal-sc", "roll": "1234"},
		},
		{name: "match rule-invalid type", IsErrExpected: true, project: "default",
			rule: &config.Rule{Rule: "match", Eval: "!=", Type: "integer", F1: "interfaceString1", F2: "interfaceString2", DB: "mongo", Col: "default", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}},
			args: map[string]interface{}{"string1": "interface1", "string2": "interface2"},
			auth: map[string]interface{}{"id": "internal-sc", "roll": "1234"},
		},
		{name: "match rule-number", IsErrExpected: false, project: "default",
			rule: &config.Rule{Rule: "match", Eval: "!=", Type: "number", F1: 0, F2: "args.string1", DB: "mongo", Col: "default"},
			args: map[string]interface{}{"args": map[string]interface{}{"string1": 1}},
			auth: map[string]interface{}{"id": "internal-sc", "roll": "1234"},
		},
		{name: "match rule-number(indirectly passing arguments)", IsErrExpected: false, project: "default",
			rule: &config.Rule{Rule: "match", Eval: "==", Type: "number", F1: "args.string1", F2: "args.string1", DB: "mongo", Col: "default"},
			args: map[string]interface{}{"args": map[string]interface{}{"string1": 1}},
			auth: map[string]interface{}{"id": "internal-sc", "roll": "1234"},
		},
		{name: "match rule-invalid boolean fields", IsErrExpected: true, project: "default",
			rule: &config.Rule{Rule: "match", Eval: "!=", Type: "bool", F1: "interfaceString1", F2: "interfaceString2", DB: "mongo", Col: "default", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}},
			args: map[string]interface{}{"args": map[string]interface{}{"string1": "interface1"}},
			auth: map[string]interface{}{"id": "internal-sc", "roll": "1234"},
		},
		{name: "match rule-valid boolean fields", IsErrExpected: false, project: "default",
			rule: &config.Rule{Rule: "match", Eval: "!=", Type: "bool", F1: true, F2: "args.k1"},
			args: map[string]interface{}{"args": map[string]interface{}{"k1": false}},
			auth: map[string]interface{}{"id": "internal-sc", "roll": "1234"},
		},
		{name: "Invalid match or rule", IsErrExpected: true, project: "default",
			rule: &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "deny"}, {Rule: "deny"}}},
			args: map[string]interface{}{"string1": "interface1", "string2": "interface2"},
			auth: map[string]interface{}{"id": "internal-sc", "roll": "1234"},
		},
		{name: "match or rule", IsErrExpected: false, project: "default",
			rule: &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "allow"}, {Rule: "deny"}}},
			args: map[string]interface{}{"string1": "interface1", "string2": "interface2"},
			auth: map[string]interface{}{"id": "internal-sc", "roll": "1234"},
		},
		{name: "webhook rule", IsErrExpected: false, project: "default",
			rule: &config.Rule{Rule: "webhook", Eval: "!=", Type: "string", F1: "interfaceString1", F2: "interfaceString2", DB: "mongo", Col: "default"},
			args: map[string]interface{}{"args": map[string]interface{}{"token": "interface1"}},
			auth: map[string]interface{}{"id": "internal-sc", "roll": "1234"},
		},
		{name: "query rule", IsErrExpected: true, project: "default",
			rule: &config.Rule{Rule: "query", Type: "string", DB: "mongo", Col: "default", Find: map[string]interface{}{"age": 12}},
			args: map[string]interface{}{"args": map[string]interface{}{"age": 12}},
			auth: map[string]interface{}{"id": "internal-sc", "roll": "1234"},
		},
		{name: "match-force rule", IsErrExpected: false, project: "default",
			rule: &config.Rule{Rule: "force", Eval: "!=", Type: "string", Field: "args.token", DB: "mongo", Col: "default"},
			args: map[string]interface{}{"args": map[string]interface{}{"token": "interface1"}},
			auth: map[string]interface{}{"id": "internal-sc", "roll": "1234"},
		},
		{name: "match-remove rule", IsErrExpected: false, project: "default",
			rule: &config.Rule{Rule: "remove", Eval: "!=", Type: "string", F1: "interfaceString1", F2: "interfaceString2", DB: "mongo", Col: "default"},
			args: map[string]interface{}{"args": map[string]interface{}{"token": "interface1"}},
			auth: map[string]interface{}{"id": "internal-sc", "roll": "1234"},
		},
	}
	auth := Init("1", &crud.Module{}, &schema.Schema{}, false)
	rule := config.Crud{"mongo": &config.CrudStub{Collections: map[string]*config.TableRule{"default": {Rules: map[string]*config.Rule{"update": {Rule: "query", Eval: "Eval", Type: "Type", DB: "mongo", Col: "default"}}}}}}
	auth.makeHttpRequest = func(ctx context.Context, method, url, token, scToken string, params, vPtr interface{}) error {
		return nil
	}
	auth.SetConfig("default", "", rule, &config.FileStore{}, &config.ServicesModule{}, &config.Eventing{})
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			_, err := (auth).matchRule(context.Background(), test.project, test.rule, test.args, test.auth)
			if (err != nil) != test.IsErrExpected {
				t.Error("| Got This ", err, "| Wanted Error |", test.IsErrExpected)
			}
		})
	}
}
func TestMatchForce_Rule(t *testing.T) {
	var testCases = []struct {
		name          string
		IsErrExpected bool
		IsSkipable    bool
		result        *PostProcess
		rule          *config.Rule
		args          map[string]interface{}
	}{
		{name: "res directly passing value", IsErrExpected: false, IsSkipable: false,
			result: &PostProcess{[]PostProcessAction{PostProcessAction{Action: "force", Field: "res.age", Value: "1234"}}},
			rule:   &config.Rule{Rule: "force", Value: "1234", Field: "res.age"},
			args:   map[string]interface{}{"string1": "interface1", "string2": "interface2"},
		},
		{name: "Scope not present for given variable", IsErrExpected: true, IsSkipable: true,
			rule: &config.Rule{Rule: "force", Value: "1234", Field: "args.age"},
			args: map[string]interface{}{"string": "interface1", "string2": "interface2"},
		},
		{name: "res indirectly passing value", IsErrExpected: false, IsSkipable: false,
			result: &PostProcess{[]PostProcessAction{PostProcessAction{Action: "force", Field: "res.age", Value: "1234"}}},
			rule:   &config.Rule{Rule: "force", Value: "args.string2", Field: "res.age"},
			args:   map[string]interface{}{"args": map[string]interface{}{"string1": "interface1", "string2": "1234"}},
		},
		{name: "Incorrect Rule Field Test Case", IsErrExpected: true, IsSkipable: true,
			rule: &config.Rule{Rule: "force", Value: "args.string2", Field: "arg.string1"},
			args: map[string]interface{}{"args": map[string]interface{}{"string1": "interface1", "string2": "interface2"}},
		},
		{name: "Valid args", IsErrExpected: false, IsSkipable: false, result: &PostProcess{},
			rule: &config.Rule{Rule: "force", Value: "1234", Field: "args.string1"},
			args: map[string]interface{}{"args": map[string]interface{}{"string1": "interface1", "string2": "interface2"}},
		},
	}
	auth := Init("1", &crud.Module{}, &schema.Schema{}, false)
	auth.makeHttpRequest = func(ctx context.Context, method, url, token, scToken string, params, vPtr interface{}) error {
		return nil
	}
	auth.SetConfig("default", "", config.Crud{}, &config.FileStore{}, &config.ServicesModule{}, &config.Eventing{})
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			r, err := matchForce(test.rule, test.args)
			if (err != nil) != test.IsErrExpected {
				t.Error("| Got This ", err, "| Wanted Error |", test.IsErrExpected)
			}
			//check return value if post process is appended
			if !test.IsSkipable {
				if !reflect.DeepEqual(r, test.result) {
					t.Error("| Got This ", r, "| Wanted Result |", test.result)
				}
			}
		})
	}
}

func TestMatchRemove_Rule(t *testing.T) {
	var testCases = []struct {
		name          string
		IsErrExpected bool
		rule          *config.Rule
		result        *PostProcess
		IsSkipable    bool
		args          map[string]interface{}
	}{
		{name: "res", IsErrExpected: false,
			IsSkipable: false,
			rule:       &config.Rule{Rule: "remove", Value: "12", Fields: []string{"res.age"}},
			args:       map[string]interface{}{"res": map[string]interface{}{"age": "12"}},
			result:     &PostProcess{[]PostProcessAction{PostProcessAction{Action: "remove", Field: "res.age", Value: nil}}},
		},
		{name: "invalid field provided", IsErrExpected: true, IsSkipable: true,
			rule: &config.Rule{Rule: "remove", Type: "string", Fields: []string{"args:age"}},
			args: map[string]interface{}{"string": "interface1", "string2": "interface2"},
		},
		{name: "scope not present", IsErrExpected: true, IsSkipable: true,
			rule: &config.Rule{Rule: "remove", Type: "string", Fields: []string{"args.age"}},
			args: map[string]interface{}{"string": "interface1", "string2": "interface2"},
		},
		{name: "remove multiple args", IsErrExpected: false, IsSkipable: false,
			result: &PostProcess{},
			rule:   &config.Rule{Rule: "remove", Type: "number", Fields: []string{"args.age", "args.exp"}},
			args:   map[string]interface{}{"args": map[string]interface{}{"age": 10, "exp": 10}},
		},
		{name: "invalid map value to another map", IsErrExpected: true, IsSkipable: true,
			rule: &config.Rule{Rule: "remove", Type: "number", Fields: []string{"args.age.exp"}},
			args: map[string]interface{}{"args": map[string]interface{}{"age": 10, "exp": 10}},
		},
		{name: "cannot find property of map", IsErrExpected: true, IsSkipable: true,
			rule: &config.Rule{Rule: "remove", Type: "number", Fields: []string{"args.aged.exp"}},
			args: map[string]interface{}{"args": map[string]interface{}{"age": 10, "exp": 10}},
		},
		{name: "invalid prefix", IsErrExpected: true, IsSkipable: true,
			rule: &config.Rule{Rule: "remove", Type: "number", Fields: []string{"arg.age.exp"}},
			args: map[string]interface{}{"args": map[string]interface{}{"age": 10, "exp": 10}},
		},
	}
	auth := Init("1", &crud.Module{}, &schema.Schema{}, false)
	auth.makeHttpRequest = func(ctx context.Context, method, url, token, scToken string, params, vPtr interface{}) error {
		return nil
	}
	auth.SetConfig("default", "", config.Crud{}, &config.FileStore{}, &config.ServicesModule{}, &config.Eventing{})
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			r, err := matchRemove(test.rule, test.args)
			if (err != nil) != test.IsErrExpected {
				t.Error("| Got This ", err, "| Wanted Error |", test.IsErrExpected)
			}
			//check return value if post process is appended
			if !test.IsSkipable {
				if !reflect.DeepEqual(r, test.result) {
					t.Error("| Got This ", r, "| Wanted Result |", test.result)
				}
			}
		})
	}
}
