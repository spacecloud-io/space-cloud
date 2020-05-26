package eventing

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

func Test_generateEventingRule(t *testing.T) {
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
			name: "error surveying ruleType",
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
			name: "eventing rule spec object created",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "project"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{mock.Anything, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "ruleType"},
				},
			},
			want: &model.SpecObject{
				API:  "/v1/config/projects/{project}/eventing/rules/{id}",
				Type: "eventing-rule",
				Meta: map[string]string{
					"id":      "ruleType",
					"project": "project",
				},
				Spec: map[string]interface{}{
					"rule": "allow",
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

			got, err := generateEventingRule()
			if (err != nil) != tt.wantErr {
				t.Errorf("generateEventingRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateEventingRule() = %v, want %v", got, tt.want)
			}

			mockSurvey.AssertExpectations(t)
		})
	}
}

func Test_generateEventingSchema(t *testing.T) {
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
			name: "error surveying ruleType",
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
			name: "error surveying schema",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter rule type"}, &surveyReturnValue, mock.Anything},
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
			name: "eventing schema spec object created",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "project"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter rule type"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "ruleType"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{mock.Anything, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "type col {id : ID @primary}"},
				},
			},
			want: &model.SpecObject{
				API:  "/v1/config/projects/{project}/eventing/schema/{id}",
				Type: "eventing-schema",
				Meta: map[string]string{
					"id":      "ruleType",
					"project": "project",
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

			got, err := generateEventingSchema()
			if (err != nil) != tt.wantErr {
				t.Errorf("generateEventingSchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateEventingSchema() = %v, want %v", got, tt.want)
			}

			mockSurvey.AssertExpectations(t)
		})
	}
}

func Test_generateEventingConfig(t *testing.T) {
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
			name: "error surveying dbAlias",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter project"}, &surveyReturnValue, mock.Anything},
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
			name: "eventing config spec object created",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter project"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "project"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{mock.Anything, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "db"},
				},
			},
			want: &model.SpecObject{
				API:  "/v1/config/projects/{project}/eventing/config/{id}",
				Type: "eventing-config",
				Meta: map[string]string{
					"project": "project",
					"id":      "eventing-config",
				},
				Spec: map[string]interface{}{
					"dbAlias": "db",
					"enabled": true,
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

			got, err := generateEventingConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("generateEventingConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateEventingConfig() = %v, want %v", got, tt.want)
			}

			mockSurvey.AssertExpectations(t)
		})
	}
}

func Test_generateEventingTrigger(t *testing.T) {
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
			name: "error surveying trigger name",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter project"}, &surveyReturnValue, mock.Anything},
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
			name: "error surveying source",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter project"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "trigger name"}, &surveyReturnValue, mock.Anything},
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
			name: "error surveying trigger operation with file storage source",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter project"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "trigger name"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select source ", Options: []string{"Database", "File Storage", "Custom"}}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "File Storage"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select trigger operation", Options: []string{"FILE_CREATE", "FILE_DELETE"}}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "error surveying trigger operation with custom source",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter project"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "trigger name"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select source ", Options: []string{"Database", "File Storage", "Custom"}}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "Custom"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter trigger type"}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "error surveying trigger operation with database source",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter project"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "trigger name"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select source ", Options: []string{"Database", "File Storage", "Custom"}}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "Database"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select trigger operation", Options: []string{"DB_INSERT", "DB_UPDATE", "DB_DELETE"}}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "error surveying db alias with database source",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter project"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "trigger name"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select source ", Options: []string{"Database", "File Storage", "Custom"}}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "Database"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select trigger operation", Options: []string{"DB_INSERT", "DB_UPDATE", "DB_DELETE"}}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Database Alias "}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "error surveying collection name with database source",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter project"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "trigger name"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select source ", Options: []string{"Database", "File Storage", "Custom"}}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "Database"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select trigger operation", Options: []string{"DB_INSERT", "DB_UPDATE", "DB_DELETE"}}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Database Alias "}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter collection/table name"}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "error surveying webhook url with database source",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter project"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "trigger name"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select source ", Options: []string{"Database", "File Storage", "Custom"}}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "Database"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select trigger operation", Options: []string{"DB_INSERT", "DB_UPDATE", "DB_DELETE"}}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Database Alias "}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter collection/table name"}, mock.Anything, mock.Anything},
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
			name: "error surveying advanced settings with database source",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter project"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "trigger name"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select source ", Options: []string{"Database", "File Storage", "Custom"}}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "Database"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select trigger operation", Options: []string{"DB_INSERT", "DB_UPDATE", "DB_DELETE"}}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Database Alias "}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter collection/table name"}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "webhook url"}, mock.Anything, mock.Anything},
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
			name: "error surveying retries of advanced settings true with database source",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter project"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "trigger name"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select source ", Options: []string{"Database", "File Storage", "Custom"}}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "Database"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select trigger operation", Options: []string{"DB_INSERT", "DB_UPDATE", "DB_DELETE"}}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Database Alias "}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter collection/table name"}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "webhook url"}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Do you want advanced settings? (Y / n) ?", Default: "n"}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil, "y"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Retries count", Default: "3"}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "error surveying timeout of advanced settings true with database source",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter project"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "trigger name"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select source ", Options: []string{"Database", "File Storage", "Custom"}}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "Database"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select trigger operation", Options: []string{"DB_INSERT", "DB_UPDATE", "DB_DELETE"}}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Database Alias "}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter collection/table name"}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "webhook url"}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Do you want advanced settings? (Y / n) ?", Default: "n"}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil, "y"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Retries count", Default: "3"}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Timeout", Default: "5000"}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "spec object created when advanced settings true with database source",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter project"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "project"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "trigger name"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "trigger"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select source ", Options: []string{"Database", "File Storage", "Custom"}}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "Database"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select trigger operation", Options: []string{"DB_INSERT", "DB_UPDATE", "DB_DELETE"}}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil, "DB_INSERT"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Database Alias "}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil, "db"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter collection/table name"}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil, "col"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "webhook url"}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil, "url"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Do you want advanced settings? (Y / n) ?", Default: "n"}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil, "y"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Retries count", Default: "3"}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Timeout", Default: "5000"}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
			},
			want: &model.SpecObject{
				API:  "/v1/config/projects/{project}/eventing/triggers/{id}",
				Type: "eventing-triggers",
				Meta: map[string]string{
					"project": "project",
					"id":      "trigger",
				},
				Spec: map[string]interface{}{
					"type":    "DB_INSERT",
					"url":     "url",
					"retries": 3,
					"timeout": 5000,
					"options": map[string]interface{}{"db": "db", "col": "col"},
				},
			},
		},
		{
			name: "spec object created with file storage source",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter project"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "project"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "trigger name"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "trigger"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select source ", Options: []string{"Database", "File Storage", "Custom"}}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "File Storage"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select trigger operation", Options: []string{"FILE_CREATE", "FILE_DELETE"}}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil, "FILE_CREATE"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "webhook url"}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil, "url"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Do you want advanced settings? (Y / n) ?", Default: "n"}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil, "n"},
				},
			},
			want: &model.SpecObject{
				API:  "/v1/config/projects/{project}/eventing/triggers/{id}",
				Type: "eventing-triggers",
				Meta: map[string]string{
					"project": "project",
					"id":      "trigger",
				},
				Spec: map[string]interface{}{
					"type":    "FILE_CREATE",
					"url":     "url",
					"retries": 3,
					"timeout": 5000,
					"options": map[string]interface{}{},
				},
			},
		},
		{
			name: "spec object created with custom source",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter project"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "project"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "trigger name"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "trigger"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select source ", Options: []string{"Database", "File Storage", "Custom"}}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "Custom"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter trigger type"}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil, "custom"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "webhook url"}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil, "url"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Do you want advanced settings? (Y / n) ?", Default: "n"}, mock.Anything, mock.Anything},
					paramsReturned: []interface{}{nil, "n"},
				},
			},
			want: &model.SpecObject{
				API:  "/v1/config/projects/{project}/eventing/triggers/{id}",
				Type: "eventing-triggers",
				Meta: map[string]string{
					"project": "project",
					"id":      "trigger",
				},
				Spec: map[string]interface{}{
					"type":    "custom",
					"url":     "url",
					"retries": 3,
					"timeout": 5000,
					"options": map[string]interface{}{},
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

			got, err := generateEventingTrigger()
			if (err != nil) != tt.wantErr {
				t.Errorf("generateEventingTrigger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateEventingTrigger() = %v, want %v", got, tt.want)
			}

			mockSurvey.AssertExpectations(t)
		})
	}
}
