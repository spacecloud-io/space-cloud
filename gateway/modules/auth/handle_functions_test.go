package auth

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud"
)

func TestIsFuncCallAuthorised(t *testing.T) {
	var authMatchQuery = []struct {
		module                                      *Module
		testName, project, token, service, function string
		secretKeys                                  []*config.Secret
		params                                      interface{}
		result                                      map[string]interface{}
		IsErrExpected, CheckResult                  bool
	}{
		{
			testName: "Successful Test allow(Internal Services)", project: "project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
			module: &Module{fileRules: []*config.FileRule{&config.FileRule{
				Prefix: string(os.PathSeparator),
				Rule:   map[string]*config.Rule{"read": &config.Rule{Rule: "allow"}},
			},
			},
				funcRules: &config.ServicesModule{
					InternalServices: config.Services{
						"service": &config.Service{
							Endpoints: map[string]*config.Endpoint{
								"ep": {
									Rule: &config.Rule{Rule: "allow"},
								},
							},
						},
					},
				},
				project: "project"},
			service: "service", secretKeys: []*config.Secret{{IsPrimary: true, Secret: "mySecretkey"}},
			function:      "ep",
			IsErrExpected: false,
		},
		{
			testName: "Invalid Project Details(Services)", project: "project1", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
			module: &Module{fileRules: []*config.FileRule{&config.FileRule{
				Prefix: string(os.PathSeparator),
				Rule:   map[string]*config.Rule{"read": &config.Rule{Rule: "allow"}},
			},
			},
				funcRules: &config.ServicesModule{
					Services: config.Services{
						"service": &config.Service{
							Endpoints: map[string]*config.Endpoint{
								"ep": {
									Rule: &config.Rule{Rule: "allow"},
								},
							},
						},
					},
				},
				project: "project"},
			service: "service", secretKeys: []*config.Secret{{IsPrimary: true, Secret: "mySecretkey"}},
			function:      "ep",
			IsErrExpected: true,
		},
		{
			testName: "Test Case with rule deny", project: "project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
			module: &Module{fileRules: []*config.FileRule{&config.FileRule{
				Prefix: string(os.PathSeparator),
				Rule:   map[string]*config.Rule{"read": &config.Rule{Rule: "allow"}},
			},
			},
				funcRules: &config.ServicesModule{
					InternalServices: config.Services{
						"service": &config.Service{
							Endpoints: map[string]*config.Endpoint{
								"ep": {
									Rule: &config.Rule{Rule: "deny"},
								},
							},
						},
					},
				},
				project: "project"},
			service: "service", secretKeys: []*config.Secret{{IsPrimary: true, Secret: "mySecretkey"}},
			function:      "ep",
			IsErrExpected: true,
		},
		{
			testName: "Test Case-Successfully parse token", project: "project", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc",
			module: &Module{fileRules: []*config.FileRule{&config.FileRule{
				Prefix: string(os.PathSeparator),
				Rule:   map[string]*config.Rule{"read": &config.Rule{Rule: "match", Eval: "==", F1: 1, F2: 1, Type: "number"}},
			},
			}, funcRules: &config.ServicesModule{
				InternalServices: config.Services{
					"service": &config.Service{
						Endpoints: map[string]*config.Endpoint{
							"ep": {
								Rule: &config.Rule{Rule: "match", Eval: "==", F1: 1, F2: 1, Type: "number"},
							},
						},
					},
				},
			},
				project: "project"},
			service: "service", secretKeys: []*config.Secret{{IsPrimary: true, Secret: "mySecretkey"}},
			function:      "ep",
			IsErrExpected: false,
			CheckResult:   true,
			result:        map[string]interface{}{"token1": "token1value", "token2": "token2value"},
		},
	}
	authModule := Init("1", &crud.Module{})
	for _, test := range authMatchQuery {
		t.Run(test.testName, func(t *testing.T) {
			if er := authModule.SetConfig("project", test.secretKeys, "", config.Crud{}, &config.FileStore{}, test.module.funcRules, &config.Eventing{}); er != nil {
				t.Errorf("error setting config of auth module  - %s", er.Error())
			}
			auth, err := (authModule).IsFuncCallAuthorised(context.Background(), test.project, test.service, test.function, test.token, test.params)
			if (err != nil) != test.IsErrExpected {
				t.Error("Got Error-", err, "Want Error-", test.IsErrExpected)
			}
			// check result if TokenClaims is returned after parsing token and matching rule
			if test.CheckResult && !reflect.DeepEqual(test.result, auth) {
				t.Error("Got Result-", auth, "Wanted Result-", test.result)
			}
		})
	}
}
