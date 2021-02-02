package auth

import (
	"context"
	"encoding/base64"
	"errors"
	"reflect"
	"testing"

	"github.com/getlantern/deepcopy"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	jwtUtils "github.com/spaceuptech/space-cloud/gateway/utils/jwt"
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
			rule: &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "allow"},
				{Rule: "allow"}}},
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
			rule: &config.Rule{Rule: "remove", Fields: []interface{}{"args.token"}},
			args: map[string]interface{}{"args": map[string]interface{}{"token": "interface1"}},
			auth: map[string]interface{}{"id": "internal-sc", "roll": "1234"},
		},
		{name: "match-encrypt rule", IsErrExpected: false, project: "default",
			rule: &config.Rule{Rule: "encrypt", Fields: []interface{}{"args.username"}},
			args: map[string]interface{}{"args": map[string]interface{}{"username": "username1"}},
			auth: map[string]interface{}{"id": "internal-sc", "roll": "1234"},
		},
		{name: "match-decrypt rule", IsErrExpected: false, project: "default",
			rule: &config.Rule{Rule: "decrypt", Fields: []interface{}{"args.username"}},
			args: map[string]interface{}{"args": map[string]interface{}{"username": base64.StdEncoding.EncodeToString([]byte("username1"))}},
			auth: map[string]interface{}{"id": "internal-sc", "roll": "1234"},
		},
		{
			name: "match-hash rule", IsErrExpected: false, project: "default",
			rule: &config.Rule{Rule: "hash", Fields: []interface{}{"args.password"}},
			args: map[string]interface{}{"args": map[string]interface{}{"password": "password"}},
			auth: map[string]interface{}{"id": "internal-sc", "roll": "1234"},
		},
		{
			name:          "Match global security function rule",
			IsErrExpected: false,
			project:       "default",
			rule: &config.Rule{
				Rule:                 "function",
				SecurityFunctionName: "force-and-remove",
				FnBlockVariables: map[string]string{
					"role":        "args.auth.role",
					"id":          "args.find.id",
					"custom-role": "super-user",
					"operation":   "args.op",
				},
			},
			args: map[string]interface{}{
				"args": map[string]interface{}{
					"op": "one",
					"auth": map[string]interface{}{
						"role":    "admin",
						"user-id": "UL6VUwwGEFTgxzoZPy9g",
					},
					"find": map[string]interface{}{
						"id": 1500,
					},
				},
			},
			auth: map[string]interface{}{
				"role":    "admin",
				"user-id": "UL6VUwwGEFTgxzoZPy9g",
			},
		},
	}
	auth := Init("chicago", "1", &crud.Module{}, nil)
	dbRules := config.DatabaseRules{config.GenerateResourceID("chicago", "project", config.ResourceDatabaseRule, ""): &config.DatabaseRule{Rules: map[string]*config.Rule{"update": {Rule: "query", Eval: "Eval", Type: "Type", DB: "mongo", Col: "default"}}}}
	securityFunctions := config.SecurityFunctions{
		config.GenerateResourceID("chicago", "default", config.ResourceSecurityFunction, "force-and-remove"): &config.SecurityFunction{
			ID: "force-and-remove",
			Rule: &config.Rule{
				Name: "Main and rule",
				Rule: "and",
				Clauses: []*config.Rule{
					{
						Name: "match role rule",
						Rule: "match",
						Type: "string",
						Eval: "==",
						F1:   "args.role",
						F2:   "admin",
					},
					{
						Name: "match id rule",
						Rule: "match",
						Type: "number",
						Eval: "!=",
						F1:   "args.id",
						F2:   100,
					},
					{
						Name: "match custom-role rule",
						Rule: "match",
						Type: "string",
						Eval: "==",
						F1:   "args.custom-role",
						F2:   "super-user",
					},
				},
			},
			Variables: []string{"role", "id", "operation", "custom-role"},
		},
	}
	auth.makeHTTPRequest = func(ctx context.Context, method, url, token, scToken string, params, vPtr interface{}) error {
		return nil
	}
	err := auth.SetConfig(context.TODO(),
		"local",
		&config.ProjectConfig{ID: "default", AESKey: "Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=", Secrets: []*config.Secret{{IsPrimary: true, Secret: "mySecretKey"}}},
		dbRules,
		config.DatabasePreparedQueries{},
		config.FileStoreRules{},
		config.Services{},
		config.EventingRules{},
		securityFunctions)
	if err != nil {
		t.Errorf("Unable to set auth config %s", err.Error())
		return
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			_, err := (auth).matchRule(context.Background(), test.project, test.rule, test.args, test.auth, model.ReturnWhereStub{})
			if (err != nil) != test.IsErrExpected {
				t.Error("| Got This ", err, "| Wanted Error |", test.IsErrExpected)
			}
		})
	}
}
func TestMatchForce_Rule(t *testing.T) {
	m := Module{}
	emptyAuth := make(map[string]interface{})
	var testCases = []struct {
		name             string
		isErrExpected    bool
		checkPostProcess bool
		checkArgs        bool
		result           *model.PostProcess
		rule             *config.Rule
		args             map[string]interface{}
		wantedargs       map[string]interface{}
	}{
		{name: "res directly passing value", isErrExpected: false, checkPostProcess: true, checkArgs: false,
			result: &model.PostProcess{PostProcessAction: []model.PostProcessAction{{Action: "force", Field: "res.age", Value: "1234"}}},
			rule:   &config.Rule{Rule: "force", Value: "1234", Field: "res.age"},
			args:   map[string]interface{}{"string1": "interface1", "string2": "interface2"},
		},
		{name: "Scope not present for given variable", isErrExpected: true, checkPostProcess: false, checkArgs: false,
			rule: &config.Rule{Rule: "force", Value: "1234", Field: "args.age"},
			args: map[string]interface{}{"string": "interface1", "string2": "interface2"},
		},
		{name: "res indirectly passing value", isErrExpected: false, checkPostProcess: true, checkArgs: false,
			result: &model.PostProcess{PostProcessAction: []model.PostProcessAction{{Action: "force", Field: "res.age", Value: "1234"}}},
			rule:   &config.Rule{Rule: "force", Value: "args.string2", Field: "res.age"},
			args:   map[string]interface{}{"args": map[string]interface{}{"string1": "interface1", "string2": "1234"}},
		},
		{name: "Incorrect Rule Field Test Case", isErrExpected: true, checkPostProcess: false, checkArgs: false,
			rule: &config.Rule{Rule: "force", Value: "args.string2", Field: "arg.string1"},
			args: map[string]interface{}{"args": map[string]interface{}{"string1": "interface1", "string2": "interface2"}},
		},
		{name: "Valid args", isErrExpected: false, checkPostProcess: false, checkArgs: true,
			rule:       &config.Rule{Rule: "force", Value: "1234", Field: "args.string1"},
			args:       map[string]interface{}{"args": map[string]interface{}{"string1": "interface1", "string2": "interface2"}},
			wantedargs: map[string]interface{}{"args": map[string]interface{}{"string1": "1234", "string2": "interface2"}},
		},
		{
			name: "rule clause - allow",
			rule: &config.Rule{Rule: "force", Clause: &config.Rule{Rule: "allow"}},
		},
		{
			name: "rule clause - deny",
			rule: &config.Rule{Rule: "force", Clause: &config.Rule{Rule: "deny"}},
		},
	}
	auth := Init("chicago", "1", &crud.Module{}, nil)
	auth.makeHTTPRequest = func(ctx context.Context, method, url, token, scToken string, params, vPtr interface{}) error {
		return nil
	}
	err := auth.SetConfig(context.TODO(), "local", &config.ProjectConfig{ID: "project", Secrets: []*config.Secret{{IsPrimary: true, Secret: "mySecretKey"}}}, config.DatabaseRules{}, config.DatabasePreparedQueries{}, config.FileStoreRules{}, config.Services{}, config.EventingRules{}, config.SecurityFunctions{})
	if err != nil {
		t.Errorf("Unable to set auth config %s", err.Error())
		return
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			r, err := m.matchForce(context.Background(), "testID", test.rule, test.args, emptyAuth)
			if (err != nil) != test.isErrExpected {
				t.Error("| Got This ", err, "| Wanted Error |", test.isErrExpected)
			}
			// check return value if post process is appended
			if test.checkPostProcess {
				if !reflect.DeepEqual(r, test.result) {
					t.Error("| Got This ", r, "| Wanted Result |", test.result)
				}
			}
			if test.checkArgs {
				if !reflect.DeepEqual(test.args, test.wantedargs) {
					t.Error("| Got This ", test.args, "| Wanted Result |", test.wantedargs)
				}
			}
		})
	}
}

func TestMatchRemove_Rule(t *testing.T) {
	m := Module{}
	emptyAuth := make(map[string]interface{})
	var testCases = []struct {
		name             string
		isErrExpected    bool
		checkArgs        bool
		rule             *config.Rule
		result           *model.PostProcess
		checkPostProcess bool
		args             map[string]interface{}
		wantedargs       map[string]interface{}
	}{
		{name: "res", isErrExpected: false,
			checkPostProcess: true, checkArgs: false,
			rule:   &config.Rule{Rule: "remove", Fields: []interface{}{"res.age"}},
			args:   map[string]interface{}{"res": map[string]interface{}{"age": "12"}},
			result: &model.PostProcess{PostProcessAction: []model.PostProcessAction{{Action: "remove", Field: "res.age", Value: nil}}},
		},
		{
			name:             "Provide values to remove fields from args object",
			isErrExpected:    false,
			checkPostProcess: true,
			checkArgs:        false,
			rule:             &config.Rule{Rule: "remove", Fields: "args.auth.obj"},
			args:             map[string]interface{}{"res": map[string]interface{}{"age": "12"}, "args": map[string]interface{}{"auth": map[string]interface{}{"obj": []interface{}{"res.age"}}}},
			result:           &model.PostProcess{PostProcessAction: []model.PostProcessAction{{Action: "remove", Field: "res.age", Value: nil}}},
		},
		{name: "invalid field provided", isErrExpected: true, checkPostProcess: false, checkArgs: false,
			rule: &config.Rule{Rule: "remove", Fields: []string{"args:age"}},
			args: map[string]interface{}{"string": "interface1", "string2": "interface2"},
		},
		{name: "scope not present", isErrExpected: true, checkPostProcess: false, checkArgs: false,
			rule: &config.Rule{Rule: "remove", Fields: []string{"args.age"}},
			args: map[string]interface{}{"string": "interface1", "string2": "interface2"},
		},
		{name: "remove multiple args", isErrExpected: false, checkPostProcess: false, checkArgs: true,
			rule:       &config.Rule{Rule: "remove", Fields: []interface{}{"args.age", "args.exp"}},
			args:       map[string]interface{}{"args": map[string]interface{}{"age": 10, "exp": 10}},
			wantedargs: map[string]interface{}{"args": map[string]interface{}{}},
		},
		{name: "invalid map value to another map", isErrExpected: true, checkPostProcess: false, checkArgs: false,
			rule: &config.Rule{Rule: "remove", Fields: []string{"args.age.exp"}},
			args: map[string]interface{}{"args": map[string]interface{}{"age": 10, "exp": 10}},
		},
		{name: "cannot find property of map", isErrExpected: true, checkPostProcess: false, checkArgs: false,
			rule: &config.Rule{Rule: "remove", Fields: []string{"args.aged.exp"}},
			args: map[string]interface{}{"args": map[string]interface{}{"age": 10, "exp": 10}},
		},
		{name: "invalid prefix", isErrExpected: true, checkPostProcess: false, checkArgs: false,
			rule: &config.Rule{Rule: "remove", Fields: []string{"arg.age.exp"}},
			args: map[string]interface{}{"args": map[string]interface{}{"age": 10, "exp": 10}},
		},
		{
			name:             "Invalid value provide for get fields",
			isErrExpected:    true,
			checkPostProcess: false,
			checkArgs:        false,
			rule:             &config.Rule{Rule: "remove", Fields: 1},
			args:             map[string]interface{}{"args": map[string]interface{}{"age": 10, "exp": 10}},
		},
		{
			name:             "Throw error if fields field contains a value which is not string",
			isErrExpected:    true,
			checkPostProcess: false,
			checkArgs:        false,
			rule:             &config.Rule{Rule: "remove", Fields: []interface{}{1}},
			args:             map[string]interface{}{"args": map[string]interface{}{"age": 10, "exp": 10}},
		},
		{
			name: "rule clause - allow",
			rule: &config.Rule{Rule: "force", Clause: &config.Rule{Rule: "allow"}},
		},
		{
			name: "rule clause - deny",
			rule: &config.Rule{Rule: "force", Clause: &config.Rule{Rule: "deny"}},
		},
	}
	auth := Init("chicago", "1", &crud.Module{}, nil)
	auth.makeHTTPRequest = func(ctx context.Context, method, url, token, scToken string, params, vPtr interface{}) error {
		return nil
	}
	err := auth.SetConfig(context.TODO(), "local", &config.ProjectConfig{ID: "project", Secrets: []*config.Secret{{IsPrimary: true, Secret: "mySecretKey"}}}, config.DatabaseRules{}, config.DatabasePreparedQueries{}, config.FileStoreRules{}, config.Services{}, config.EventingRules{}, config.SecurityFunctions{})
	if err != nil {
		t.Errorf("Unable to set auth config %s", err.Error())
		return
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			r, err := m.matchRemove(context.Background(), "testID", test.rule, test.args, emptyAuth)
			if (err != nil) != test.isErrExpected {
				t.Error("| Got This ", err, "| Wanted Error |", test.isErrExpected)
			}
			// check return value if post process is appended
			if test.checkPostProcess {
				if !reflect.DeepEqual(r, test.result) {
					t.Error("| Got This ", r, "| Wanted Result |", test.result)
				}
			}
			if test.checkArgs {
				if !reflect.DeepEqual(test.args, test.wantedargs) {
					t.Error("| Got This ", test.args, "| Wanted Result |", test.wantedargs)
				}
			}
		})
	}
}

func base64DecodeString(key string) []byte {
	decodedKey, _ := base64.StdEncoding.DecodeString(key)
	return decodedKey
}

func TestModule_matchEncrypt(t *testing.T) {
	type args struct {
		rule *config.Rule
		args map[string]interface{}
	}
	tests := []struct {
		name         string
		m            *Module
		args         args
		want         *model.PostProcess
		shouldChange bool
		wantErr      bool
	}{
		{
			name:    "invalid field",
			m:       &Module{aesKey: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=")},
			args:    args{rule: &config.Rule{Rule: "encrypt", Fields: []string{"args.abc"}}, args: map[string]interface{}{"args": map[string]interface{}{"username": "username1"}}},
			wantErr: true,
		},
		{
			name:    "invalid value type",
			m:       &Module{aesKey: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g")},
			args:    args{rule: &config.Rule{Rule: "encrypt", Fields: []string{"args.username"}}, args: map[string]interface{}{"args": map[string]interface{}{"username": 10}}},
			wantErr: true,
		},
		{
			name:    "invalid key",
			m:       &Module{aesKey: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g")},
			args:    args{rule: &config.Rule{Rule: "encrypt", Fields: []string{"args.username"}}, args: map[string]interface{}{"args": map[string]interface{}{"username": "username1"}}},
			wantErr: true,
		},
		{
			name:    "invalid field prefix",
			m:       &Module{aesKey: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=")},
			args:    args{rule: &config.Rule{Rule: "encrypt", Fields: []string{"abc.username"}}, args: map[string]interface{}{"abc": map[string]interface{}{"username": "username1"}}},
			wantErr: true,
		},
		{
			name:         "valid args",
			m:            &Module{aesKey: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=")},
			args:         args{rule: &config.Rule{Rule: "encrypt", Fields: []interface{}{"args.username"}}, args: map[string]interface{}{"args": map[string]interface{}{"username": "username1"}}},
			want:         &model.PostProcess{},
			shouldChange: true,
		},
		{
			name:         "valid res",
			m:            &Module{aesKey: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=")},
			args:         args{rule: &config.Rule{Rule: "encrypt", Fields: []interface{}{"res.username"}}, args: map[string]interface{}{"res": map[string]interface{}{"username": "username1"}}},
			want:         &model.PostProcess{PostProcessAction: []model.PostProcessAction{{Action: "encrypt", Field: "res.username"}}},
			shouldChange: false,
		},
		{
			name: "Provide values to encrypt fields from args object",
			m:    &Module{aesKey: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=")},
			args: args{rule: &config.Rule{Rule: "encrypt", Fields: "args.auth.obj"}, args: map[string]interface{}{"args": map[string]interface{}{"auth": map[string]interface{}{"obj": []interface{}{"res.username"}}}, "res": map[string]interface{}{"username": "username1"}}},
			want: &model.PostProcess{PostProcessAction: []model.PostProcessAction{{Action: "encrypt", Field: "res.username"}}},
		},
		{
			name:         "valid args with rule allow",
			m:            &Module{aesKey: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=")},
			args:         args{rule: &config.Rule{Rule: "encrypt", Fields: []interface{}{"args.username"}, Clause: &config.Rule{Rule: "allow"}}, args: map[string]interface{}{"args": map[string]interface{}{"username": "username1"}}},
			want:         &model.PostProcess{},
			shouldChange: true,
		},
		{
			name:         "valid args with rule deny",
			m:            &Module{aesKey: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=")},
			args:         args{rule: &config.Rule{Rule: "encrypt", Fields: []string{"args.username"}, Clause: &config.Rule{Rule: "deny"}}, args: map[string]interface{}{"args": map[string]interface{}{"username": "username1"}}},
			want:         &model.PostProcess{},
			shouldChange: false,
		},
		{
			name:    "Invalid value provide for get fields",
			m:       &Module{aesKey: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=")},
			args:    args{rule: &config.Rule{Rule: "encrypt", Fields: 1}, args: map[string]interface{}{"abc": map[string]interface{}{"username": "username1"}}},
			wantErr: true,
		},
		{
			name:    "Throw error if fields field contains a value which is not string",
			m:       &Module{aesKey: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=")},
			args:    args{rule: &config.Rule{Rule: "encrypt", Fields: []interface{}{"res.username", 1}}, args: map[string]interface{}{"res": map[string]interface{}{"username": "username1"}}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := map[string]interface{}{}
			_ = deepcopy.Copy(&args, tt.args.args)

			got, err := tt.m.matchEncrypt(context.Background(), "", tt.args.rule, args, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("Module.matchEncrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.shouldChange == reflect.DeepEqual(args, tt.args.args) {
				t.Errorf("Module.matchEncrypt() args = %v, ogArgs = %v, shouldChange = %v", args, tt.args.args, tt.shouldChange)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Module.matchEncrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModule_matchDecrypt(t *testing.T) {
	type args struct {
		rule *config.Rule
		args map[string]interface{}
	}
	tests := []struct {
		name         string
		m            *Module
		args         args
		want         *model.PostProcess
		shouldChange bool
		wantErr      bool
	}{
		{
			name:    "invalid field",
			m:       &Module{aesKey: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=")},
			args:    args{rule: &config.Rule{Rule: "decrypt", Fields: []string{"args.abc"}}, args: map[string]interface{}{"args": map[string]interface{}{"username": "username1"}}},
			wantErr: true,
		},
		{
			name:    "invalid value type",
			m:       &Module{aesKey: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g")},
			args:    args{rule: &config.Rule{Rule: "decrypt", Fields: []string{"args.username"}}, args: map[string]interface{}{"args": map[string]interface{}{"username": 10}}},
			wantErr: true,
		},
		{
			name:    "invalid key",
			m:       &Module{aesKey: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g")},
			args:    args{rule: &config.Rule{Rule: "decrypt", Fields: []string{"args.username"}}, args: map[string]interface{}{"args": map[string]interface{}{"username": "username1"}}},
			wantErr: true,
		},
		{
			name:         "valid args",
			m:            &Module{aesKey: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=")},
			args:         args{rule: &config.Rule{Rule: "decrypt", Fields: []interface{}{"args.username"}}, args: map[string]interface{}{"args": map[string]interface{}{"username": "BXioRN4GyvZs"}}},
			want:         &model.PostProcess{},
			shouldChange: true,
		},
		{
			name: "valid res",
			m:    &Module{aesKey: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=")},
			args: args{rule: &config.Rule{Rule: "decrypt", Fields: []interface{}{"res.username"}}, args: map[string]interface{}{"res": map[string]interface{}{"username": "username1"}}},
			want: &model.PostProcess{PostProcessAction: []model.PostProcessAction{{Action: "decrypt", Field: "res.username"}}},
		},
		{
			name: "Provide values to decrypt fields from args object",
			m:    &Module{aesKey: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=")},
			args: args{rule: &config.Rule{Rule: "encrypt", Fields: "args.auth.obj"}, args: map[string]interface{}{"args": map[string]interface{}{"auth": map[string]interface{}{"obj": []interface{}{"res.username"}}}, "res": map[string]interface{}{"username": "username1"}}},
			want: &model.PostProcess{PostProcessAction: []model.PostProcessAction{{Action: "decrypt", Field: "res.username"}}},
		},
		{
			name:    "invalid field prefix",
			m:       &Module{aesKey: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=")},
			args:    args{rule: &config.Rule{Rule: "decrypt", Fields: []string{"abc.username"}}, args: map[string]interface{}{"abc": map[string]interface{}{"username": "username1"}}},
			wantErr: true,
		},
		{
			name:         "valid args with rule allow",
			m:            &Module{aesKey: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=")},
			args:         args{rule: &config.Rule{Rule: "decrypt", Fields: []interface{}{"args.username"}, Clause: &config.Rule{Rule: "allow"}}, args: map[string]interface{}{"args": map[string]interface{}{"username": "BXioRN4GyvZs"}}},
			want:         &model.PostProcess{},
			shouldChange: true,
		},
		{
			name:         "valid args with rule deny",
			m:            &Module{aesKey: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=")},
			args:         args{rule: &config.Rule{Rule: "decrypt", Fields: []string{"args.username"}, Clause: &config.Rule{Rule: "deny"}}, args: map[string]interface{}{"args": map[string]interface{}{"username": "BXioRN4GyvZs"}}},
			want:         &model.PostProcess{},
			shouldChange: false,
		},
		{
			name:    "Invalid value provide for get fields",
			m:       &Module{aesKey: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=")},
			args:    args{rule: &config.Rule{Rule: "encrypt", Fields: 1}, args: map[string]interface{}{"abc": map[string]interface{}{"username": "username1"}}},
			wantErr: true,
		},
		{
			name:    "Throw error if fields field contains a value which is not string",
			m:       &Module{aesKey: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=")},
			args:    args{rule: &config.Rule{Rule: "encrypt", Fields: []interface{}{"res.username", 1}}, args: map[string]interface{}{"res": map[string]interface{}{"username": "username1"}}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := map[string]interface{}{}
			_ = deepcopy.Copy(&args, tt.args.args)

			got, err := tt.m.matchDecrypt(context.Background(), "", tt.args.rule, tt.args.args, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("Module.matchDecrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.shouldChange == reflect.DeepEqual(args, tt.args.args) {
				t.Errorf("Module.matchDecrypt() args = %v, ogArgs = %v, shouldChange = %v", args, tt.args.args, tt.shouldChange)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Module.matchDecrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_matchHash(t *testing.T) {
	type args struct {
		rule *config.Rule
		args map[string]interface{}
	}
	tests := []struct {
		name         string
		m            *Module
		args         args
		want         *model.PostProcess
		shouldChange bool
		wantErr      bool
	}{
		{
			name:    "invalid field",
			m:       &Module{aesKey: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=")},
			args:    args{rule: &config.Rule{Rule: "hash", Fields: []string{"args.abc"}}, args: map[string]interface{}{"args": map[string]interface{}{"password": "password"}}},
			wantErr: true,
		},
		{
			name:         "valid args",
			args:         args{rule: &config.Rule{Rule: "hash", Fields: []interface{}{"args.password"}}, args: map[string]interface{}{"args": map[string]interface{}{"password": "password"}}},
			want:         &model.PostProcess{},
			shouldChange: true,
		},
		{
			name: "valid res",
			m:    &Module{aesKey: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=")},
			args: args{rule: &config.Rule{Rule: "hash", Fields: []interface{}{"res.password"}}, args: map[string]interface{}{"res": map[string]interface{}{"password": "password"}}},
			want: &model.PostProcess{PostProcessAction: []model.PostProcessAction{{Action: "hash", Field: "res.password"}}},
		},
		{
			name: "Provide values to hash fields from args object",
			m:    &Module{aesKey: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=")},
			args: args{rule: &config.Rule{Rule: "encrypt", Fields: "args.auth.obj"}, args: map[string]interface{}{"args": map[string]interface{}{"auth": map[string]interface{}{"obj": []interface{}{"res.username"}}}, "res": map[string]interface{}{"username": "username1"}}},
			want: &model.PostProcess{PostProcessAction: []model.PostProcessAction{{Action: "hash", Field: "res.username"}}},
		},
		{
			name:    "invalid value type",
			m:       &Module{aesKey: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=")},
			args:    args{rule: &config.Rule{Rule: "hash", Fields: []string{"args.password"}}, args: map[string]interface{}{"args": map[string]interface{}{"password": 123456}}},
			wantErr: true,
		},
		{
			name:    "invalid field prefix",
			m:       &Module{aesKey: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=")},
			args:    args{rule: &config.Rule{Rule: "hash", Fields: []string{"abc.password"}}, args: map[string]interface{}{"abc": map[string]interface{}{"password": "password"}}},
			wantErr: true,
		},
		{
			name:         "valid args with rule allow",
			args:         args{rule: &config.Rule{Rule: "hash", Fields: []interface{}{"args.password"}, Clause: &config.Rule{Rule: "allow"}}, args: map[string]interface{}{"args": map[string]interface{}{"password": "password"}}},
			want:         &model.PostProcess{},
			shouldChange: true,
		},
		{
			name:         "valid args with rule deny",
			args:         args{rule: &config.Rule{Rule: "hash", Fields: []string{"args.password"}, Clause: &config.Rule{Rule: "deny"}}, args: map[string]interface{}{"args": map[string]interface{}{"password": "password"}}},
			want:         &model.PostProcess{},
			shouldChange: false,
		},
		{
			name:    "Invalid value provide for get fields",
			m:       &Module{aesKey: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=")},
			args:    args{rule: &config.Rule{Rule: "encrypt", Fields: 1}, args: map[string]interface{}{"abc": map[string]interface{}{"username": "username1"}}},
			wantErr: true,
		},
		{
			name:    "Throw error if fields field contains a value which is not string",
			m:       &Module{aesKey: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=")},
			args:    args{rule: &config.Rule{Rule: "encrypt", Fields: []interface{}{"res.username", 1}}, args: map[string]interface{}{"res": map[string]interface{}{"username": "username1"}}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := map[string]interface{}{}
			_ = deepcopy.Copy(&args, tt.args.args)

			m := &Module{}
			got, err := m.matchHash(context.Background(), "", tt.args.rule, tt.args.args, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("matchHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.shouldChange == reflect.DeepEqual(args, tt.args.args) {
				t.Errorf("Module.matchDecrypt() args = %v, ogArgs = %v, shouldChange = %v", args, tt.args.args, tt.shouldChange)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("matchHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModule_matchFunc(t *testing.T) {
	type params struct {
		url               string
		claims            map[string]interface{}
		params            interface{}
		shouldRequestFail bool
	}
	type args struct {
		rule       *config.Rule
		httpParams params
		args       map[string]interface{}
	}
	tests := []struct {
		name    string
		m       *Module
		args    args
		want    *model.PostProcess
		wantErr bool
	}{
		{
			name: "Normal webhook call validation should pass",
			m:    &Module{jwt: jwtUtils.New(), aesKey: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=")},
			args: args{
				httpParams: params{
					url:    "http://localhost/validate",
					claims: map[string]interface{}{},
					params: map[string]interface{}{"auth": map[string]interface{}{"role": "admin"}, "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.mcLuntEPgBDN1U_ywGLpC7L0--iD7OwX6eqjEWUo4oo"},
				},
				rule: &config.Rule{Rule: "webhook", URL: "http://localhost/validate"},
				args: map[string]interface{}{"args": map[string]interface{}{"auth": map[string]interface{}{"role": "admin"}, "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.mcLuntEPgBDN1U_ywGLpC7L0--iD7OwX6eqjEWUo4oo"}}},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Normal webhook call validation should fail",
			m:    &Module{jwt: jwtUtils.New(), aesKey: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=")},
			args: args{
				httpParams: params{
					shouldRequestFail: true,
				},
				rule: &config.Rule{Rule: "webhook", URL: "http://localhost/validate"},
				args: map[string]interface{}{"args": map[string]interface{}{"auth": map[string]interface{}{}, "token": "loremparis"}},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Normal webhook call with custom claims",
			m:    &Module{jwt: jwtUtils.New(), aesKey: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=")},
			args: args{
				httpParams: params{
					url:    "http://localhost/validate",
					claims: map[string]interface{}{"id": "4f42ofgrlg34o", "name": "tony"},
					params: map[string]interface{}{"auth": map[string]interface{}{"id": "4f42ofgrlg34o"}, "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRmNDJvZmdybGczNG8iLCJuYW1lIjoidG9ueSJ9.8ckk11R62P6KXly3d5Fx2I2NAaSukqTA0Zx-FreebpE"},
				},
				rule: &config.Rule{Rule: "webhook", URL: "http://localhost/validate", Claims: "{\"id\": \"{{ .auth.id }}\", \"name\": \"tony\"}"},
				args: map[string]interface{}{"args": map[string]interface{}{"auth": map[string]interface{}{"id": "4f42ofgrlg34o"}, "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRmNDJvZmdybGczNG8iLCJuYW1lIjoidG9ueSJ9.8ckk11R62P6KXly3d5Fx2I2NAaSukqTA0Zx-FreebpE"}}},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Normal webhook call with custom claims and template",
			m:    &Module{jwt: jwtUtils.New(), aesKey: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=")},
			args: args{
				httpParams: params{
					url:    "http://localhost/validate",
					claims: map[string]interface{}{"id": "4f42ofgrlg34o", "name": "tony"},
					params: map[string]interface{}{"service": "http://localhost:9000/"},
				},
				rule: &config.Rule{Rule: "webhook", URL: "http://localhost/validate", Claims: "{\"id\": \"{{ .args.auth.id }}\", \"name\": \"tony\"}", ReqTmpl: `{"service":"{{.args.auth.serviceAddr}}"}`, Template: config.TemplatingEngineGo},
				args: map[string]interface{}{"args": map[string]interface{}{"auth": map[string]interface{}{"id": "4f42ofgrlg34o", "serviceAddr": "http://localhost:9000/"}, "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRmNDJvZmdybGczNG8ifQ.kUx6dq6qHDYX2HyPthgthi3Rch5UPUrqm5Io1cbKVr0"}}},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var HTTPCall utils.TypeMakeHTTPRequest = func(ctx context.Context, method, url, token, scToken string, params, vPtr interface{}) error {
				if tt.args.httpParams.shouldRequestFail {
					return errors.New("Cannot make http request to web hook url")
				}
				if tt.args.httpParams.url != url {
					t.Errorf("matchFunc() Url mis match in makeHTTPRequest wanted (%s) got (%s)", tt.args.httpParams.url, url)
					return nil
				}
				claims, err := tt.m.jwt.ParseToken(ctx, token)
				if err != nil {
					t.Errorf("matchFunc() cannot parse token (%s)", token)
					return nil
				}
				delete(claims, "exp")
				if !reflect.DeepEqual(tt.args.httpParams.claims, claims) {
					t.Errorf("matchFunc() token claims mis match in makeHTTPRequest wanted (%s) got (%s)", tt.args.httpParams.claims, claims)
					return nil
				}

				if !reflect.DeepEqual(tt.args.httpParams.params, params) {
					t.Errorf("matchFunc() Url body params mis match in makeHTTPRequest wanted (%s) got (%s)", tt.args.httpParams.params, params)
					return nil
				}
				return nil
			}
			_ = tt.m.SetConfig(context.TODO(), "local", &config.ProjectConfig{ID: "project", AESKey: string(tt.m.aesKey), Secrets: []*config.Secret{{IsPrimary: true, Alg: config.HS256, Secret: "some-secret"}}}, config.DatabaseRules{}, config.DatabasePreparedQueries{}, config.FileStoreRules{}, config.Services{}, config.EventingRules{}, config.SecurityFunctions{})
			if err := tt.m.matchFunc(context.Background(), tt.args.rule, HTTPCall, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("matchFunc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
