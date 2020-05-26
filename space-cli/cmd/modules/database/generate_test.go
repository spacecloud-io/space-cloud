package database

import (
	"errors"
	"reflect"
	"testing"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spaceuptech/space-cli/cmd/model"
	"github.com/spaceuptech/space-cli/cmd/utils"
	"github.com/spaceuptech/space-cli/cmd/utils/input"
	"github.com/stretchr/testify/mock"
)

func Test_generateDBRule(t *testing.T) {
	// surveyReturnValue stores the values returned from the survey
	surveyReturnValue := ""
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	tests := []struct {
		name           string
		surveyMockArgs []mockArgs
		want           *model.SpecObject
		wantErr        bool
	}{
		{
			name: "error while surveying project ID",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{mock.Anything, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "error while surveying ID",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{mock.Anything, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "error while surveying dbAlias",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Collection Name"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{mock.Anything, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "error while surveying realtime enabled bool",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Collection Name"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter DB Alias"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{mock.Anything, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "db rule spec object generated",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "project"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Collection Name"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "col"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter DB Alias"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "db"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{mock.Anything, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "y"},
				},
			},
			want: &model.SpecObject{
				API:  "/v1/config/projects/{project}/database/{dbAlias}/collections/{col}/rules",
				Type: "db-rules",
				Meta: map[string]string{
					"dbAlias": "db",
					"col":     "col",
					"project": "project",
				},
				Spec: map[string]interface{}{
					"isRealtimeEnabled": true,
					"rules": map[string]interface{}{
						"create": map[string]interface{}{
							"rule": "allow",
						},
						"delete": map[string]interface{}{
							"rule": "allow",
						},
						"read": map[string]interface{}{
							"rule": "allow",
						},
						"update": map[string]interface{}{
							"rule": "allow",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockSurvey := utils.MockInputInterface{}

			for _, m := range tt.surveyMockArgs {
				mockSurvey.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			input.Survey = &mockSurvey

			got, err := generateDBRule()
			if (err != nil) != tt.wantErr {
				t.Errorf("generateDBRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateDBRule() = %v, want %v", got, tt.want)
			}

			mockSurvey.AssertExpectations(t)
		})
	}
}

func Test_generateDBConfig(t *testing.T) {
	// surveyReturnValue stores the values returned from the survey
	surveyReturnValue := ""
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	tests := []struct {
		name           string
		surveyMockArgs []mockArgs
		want           *model.SpecObject
		wantErr        bool
	}{
		{
			name: "error surveying project id",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{mock.Anything, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "error surveying dbtype",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{mock.Anything, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "error surveying connections string",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select database choice ", Options: []string{"mongo", "mysql", "postgres", "sqlserver", "embedded"}}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "postgres"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{mock.Anything, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "error surveying db alias",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select database choice ", Options: []string{"mongo", "mysql", "postgres", "sqlserver", "embedded"}}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "postgres"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Database Connection String ", Default: "postgres://postgres:mysecretpassword@localhost:5432/postgres?sslmode=disable"}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{mock.Anything, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "db config spec object created",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "project"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select database choice ", Options: []string{"mongo", "mysql", "postgres", "sqlserver", "embedded"}}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "postgres"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Database Connection String ", Default: "postgres://postgres:mysecretpassword@localhost:5432/postgres?sslmode=disable"}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil, "postgres://postgres:password@localhost:5432/postgres?sslmode=disable"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter DB Alias", Default: "postgres"}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil, "dbAlias"},
				},
			},
			want: &model.SpecObject{
				API:  "/v1/config/projects/{project}/database/{dbAlias}/config/{id}",
				Type: "db-config",
				Meta: map[string]string{
					"dbAlias": "dbAlias",
					"project": "project",
					"id":      "dbAlias" + "-config",
				},
				Spec: map[string]interface{}{
					"conn":      "postgres://postgres:password@localhost:5432/postgres?sslmode=disable",
					"enabled":   true,
					"isPrimary": false,
					"type":      "postgres",
				},
			},
		},
		{
			name: "dbtype default case",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{mock.Anything, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "default"},
				},
			},
			wantErr: true,
		},
		{
			name: "dbtype sqlserver case",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select database choice ", Options: []string{"mongo", "mysql", "postgres", "sqlserver", "embedded"}}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "sqlserver"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{mock.Anything, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "dbtype embedded case",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select database choice ", Options: []string{"mongo", "mysql", "postgres", "sqlserver", "embedded"}}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "embedded"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{mock.Anything, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "dbtype mongo case",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select database choice ", Options: []string{"mongo", "mysql", "postgres", "sqlserver", "embedded"}}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "mongo"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{mock.Anything, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "dbtype mysql case",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select database choice ", Options: []string{"mongo", "mysql", "postgres", "sqlserver", "embedded"}}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "mysql"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{mock.Anything, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), ""},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockSurvey := utils.MockInputInterface{}

			for _, m := range tt.surveyMockArgs {
				mockSurvey.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			input.Survey = &mockSurvey

			got, err := generateDBConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("generateDBConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateDBConfig() = %v, want %v", got, tt.want)
			}

			mockSurvey.AssertExpectations(t)
		})
	}
}

func Test_generateDBSchema(t *testing.T) {
	// surveyReturnValue stores the values returned from the survey
	surveyReturnValue := ""
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	tests := []struct {
		name           string
		surveyMockArgs []mockArgs
		want           *model.SpecObject
		wantErr        bool
	}{
		{
			name: "error surveying project",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{mock.Anything, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "error surveying collection",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project "}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{mock.Anything, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "error surveying DB Alias",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project "}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Collection "}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{mock.Anything, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "error surveying Schema",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project "}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Collection "}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter DB Alias"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{mock.Anything, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "db schema spec object created",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project "}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "project"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Collection "}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "col"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter DB Alias"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "db"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Schema"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "type col {id : ID @primary}"},
				},
			},
			want: &model.SpecObject{
				API:  "/v1/config/projects/{project}/database/{dbAlias}/collections/{col}/schema/mutate",
				Type: "db-schema",
				Meta: map[string]string{
					"dbAlias": "db",
					"project": "project",
					"col":     "col",
				},
				Spec: map[string]interface{}{
					"schema": "type col {id : ID @primary}",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockSurvey := utils.MockInputInterface{}

			for _, m := range tt.surveyMockArgs {
				mockSurvey.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			input.Survey = &mockSurvey

			got, err := generateDBSchema()
			if (err != nil) != tt.wantErr {
				t.Errorf("generateDBSchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateDBSchema() = %v, want %v", got, tt.want)
			}

			mockSurvey.AssertExpectations(t)
		})
	}
}
