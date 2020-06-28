package auth

import (
	"context"
	"crypto/aes"
	"encoding/base64"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud"
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
			rule: &config.Rule{Rule: "remove", Eval: "!=", Type: "string", F1: "interfaceString1", F2: "interfaceString2", DB: "mongo", Col: "default"},
			args: map[string]interface{}{"args": map[string]interface{}{"token": "interface1"}},
			auth: map[string]interface{}{"id": "internal-sc", "roll": "1234"},
		},
		{name: "match-encrypt rule", IsErrExpected: false, project: "default",
			rule: &config.Rule{Rule: "encrypt", Fields: []string{"args.username"}},
			args: map[string]interface{}{"args": map[string]interface{}{"username": "username1"}},
			auth: map[string]interface{}{"id": "internal-sc", "roll": "1234"},
		},
		{name: "match-decrypt rule", IsErrExpected: false, project: "default",
			rule: &config.Rule{Rule: "decrypt", Fields: []string{"args.username"}},
			args: map[string]interface{}{"args": map[string]interface{}{"username": base64.StdEncoding.EncodeToString([]byte("username1"))}},
			auth: map[string]interface{}{"id": "internal-sc", "roll": "1234"},
		},
		{
			name: "match-hash rule", IsErrExpected: false, project: "default",
			rule: &config.Rule{Rule: "hash", Fields: []string{"args.password"}},
			args: map[string]interface{}{"args": map[string]interface{}{"password": "password"}},
			auth: map[string]interface{}{"id": "internal-sc", "roll": "1234"},
		},
	}
	auth := Init("1", &crud.Module{})
	rule := config.Crud{"mongo": &config.CrudStub{Collections: map[string]*config.TableRule{"default": {Rules: map[string]*config.Rule{"update": {Rule: "query", Eval: "Eval", Type: "Type", DB: "mongo", Col: "default"}}}}}}
	auth.makeHTTPRequest = func(ctx context.Context, method, url, token, scToken string, params, vPtr interface{}) error {
		return nil
	}
	err := auth.SetConfig("default", []*config.Secret{{IsPrimary: true, Secret: "mySecretKey"}}, "Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=", rule, &config.FileStore{}, &config.ServicesModule{}, &config.Eventing{})
	if err != nil {
		t.Errorf("Unable to set auth config %s", err.Error())
		return
	}
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
			result: &model.PostProcess{PostProcessAction: []model.PostProcessAction{model.PostProcessAction{Action: "force", Field: "res.age", Value: "1234"}}},
			rule:   &config.Rule{Rule: "force", Value: "1234", Field: "res.age"},
			args:   map[string]interface{}{"string1": "interface1", "string2": "interface2"},
		},
		{name: "Scope not present for given variable", isErrExpected: true, checkPostProcess: false, checkArgs: false,
			rule: &config.Rule{Rule: "force", Value: "1234", Field: "args.age"},
			args: map[string]interface{}{"string": "interface1", "string2": "interface2"},
		},
		{name: "res indirectly passing value", isErrExpected: false, checkPostProcess: true, checkArgs: false,
			result: &model.PostProcess{PostProcessAction: []model.PostProcessAction{model.PostProcessAction{Action: "force", Field: "res.age", Value: "1234"}}},
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
	auth := Init("1", &crud.Module{})
	auth.makeHTTPRequest = func(ctx context.Context, method, url, token, scToken string, params, vPtr interface{}) error {
		return nil
	}
	_ = auth.SetConfig("default", []*config.Secret{}, "Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=", config.Crud{}, &config.FileStore{}, &config.ServicesModule{}, &config.Eventing{})
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
			rule:   &config.Rule{Rule: "remove", Fields: []string{"res.age"}},
			args:   map[string]interface{}{"res": map[string]interface{}{"age": "12"}},
			result: &model.PostProcess{PostProcessAction: []model.PostProcessAction{model.PostProcessAction{Action: "remove", Field: "res.age", Value: nil}}},
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
			rule:       &config.Rule{Rule: "remove", Fields: []string{"args.age", "args.exp"}},
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
			name: "rule clause - allow",
			rule: &config.Rule{Rule: "force", Clause: &config.Rule{Rule: "allow"}},
		},
		{
			name: "rule clause - deny",
			rule: &config.Rule{Rule: "force", Clause: &config.Rule{Rule: "deny"}},
		},
	}
	auth := Init("1", &crud.Module{})
	auth.makeHTTPRequest = func(ctx context.Context, method, url, token, scToken string, params, vPtr interface{}) error {
		return nil
	}
	_ = auth.SetConfig("default", []*config.Secret{}, "Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=", config.Crud{}, &config.FileStore{}, &config.ServicesModule{}, &config.Eventing{})
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
		name    string
		m       *Module
		args    args
		want    *model.PostProcess
		wantErr bool
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
			name: "valid res",
			m:    &Module{aesKey: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=")},
			args: args{rule: &config.Rule{Rule: "encrypt", Fields: []string{"res.username"}}, args: map[string]interface{}{"res": map[string]interface{}{"username": "username1"}}},
			want: &model.PostProcess{PostProcessAction: []model.PostProcessAction{model.PostProcessAction{Action: "encrypt", Field: "res.username"}}},
		},
		{
			name:    "invalid field prefix",
			m:       &Module{aesKey: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=")},
			args:    args{rule: &config.Rule{Rule: "encrypt", Fields: []string{"abc.username"}}, args: map[string]interface{}{"abc": map[string]interface{}{"username": "username1"}}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.matchEncrypt(tt.args.rule, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Module.matchEncrypt() error = %v, wantErr %v", err, tt.wantErr)
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
		name    string
		m       *Module
		args    args
		want    *model.PostProcess
		wantErr bool
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
			name: "valid res",
			m:    &Module{aesKey: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=")},
			args: args{rule: &config.Rule{Rule: "decrypt", Fields: []string{"res.username"}}, args: map[string]interface{}{"res": map[string]interface{}{"username": "username1"}}},
			want: &model.PostProcess{PostProcessAction: []model.PostProcessAction{model.PostProcessAction{Action: "decrypt", Field: "res.username"}}},
		},
		{
			name:    "invalid field prefix",
			m:       &Module{aesKey: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g=")},
			args:    args{rule: &config.Rule{Rule: "decrypt", Fields: []string{"abc.username"}}, args: map[string]interface{}{"abc": map[string]interface{}{"username": "username1"}}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.matchDecrypt(tt.args.rule, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Module.matchDecrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Module.matchDecrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decryptAESCFB(t *testing.T) {
	type args struct {
		dst []byte
		src []byte
		key []byte
		iv  []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "invalid key",
			args:    args{dst: make([]byte, len("username1")), src: []byte("username1"), key: []byte("invalidKey"), iv: []byte("invalidKey123456")[:aes.BlockSize]},
			wantErr: true,
		},
		{
			name: "decryption takes place",
			args: args{dst: make([]byte, len("username1")), src: []byte{5, 120, 168, 68, 222, 6, 202, 246, 108}, key: base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g="), iv: []byte(base64DecodeString("Olw6AhA/GzSxfhwKLxO7JJsUL6VUwwGEFTgxzoZPy9g="))[:aes.BlockSize]},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := decryptAESCFB(tt.args.dst, tt.args.src, tt.args.key, tt.args.iv); (err != nil) != tt.wantErr {
				t.Errorf("decryptAESCFB() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && reflect.DeepEqual(tt.args.dst, tt.args.src) {
				t.Errorf("decryptAESCFB() decryption did not take place")
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
		name    string
		args    args
		want    *model.PostProcess
		wantErr bool
	}{
		{
			name:    "invalid field",
			args:    args{rule: &config.Rule{Rule: "hash", Fields: []string{"args.abc"}}, args: map[string]interface{}{"args": map[string]interface{}{"password": "password"}}},
			wantErr: true,
		},
		{
			name: "valid res",
			args: args{rule: &config.Rule{Rule: "hash", Fields: []string{"res.password"}}, args: map[string]interface{}{"res": map[string]interface{}{"password": "password"}}},
			want: &model.PostProcess{PostProcessAction: []model.PostProcessAction{model.PostProcessAction{Action: "hash", Field: "res.password"}}},
		},
		{
			name:    "invalid value type",
			args:    args{rule: &config.Rule{Rule: "hash", Fields: []string{"args.password"}}, args: map[string]interface{}{"args": map[string]interface{}{"password": 123456}}},
			wantErr: true,
		},
		{
			name:    "invalid field prefix",
			args:    args{rule: &config.Rule{Rule: "hash", Fields: []string{"abc.password"}}, args: map[string]interface{}{"abc": map[string]interface{}{"password": "password"}}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := matchHash(tt.args.rule, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("matchHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("matchHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
