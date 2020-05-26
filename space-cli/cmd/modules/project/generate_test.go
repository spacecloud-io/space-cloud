package project

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

func Test_generateProject(t *testing.T) {
	// surveyReturnValue stores the values returned from the survey
	surveyReturnValue := ""
	// contextTime stores graphql query timeout, initialized as 0 in consistency with code
	contextTime := 0
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
			name: "empty project id error",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{mock.Anything, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
			},
		},
		{
			name: "error surveying project name",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID: "}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "project"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project Name: ", Default: "project"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), "projectName"},
				},
			},
			wantErr: true,
		},
		{
			name: "error surveying aes key",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID: "}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "project"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project Name: ", Default: "project"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "projectName"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter AES Key: "}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), "key"},
				},
			},
			wantErr: true,
		},
		{
			name: "error surveying context time",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID: "}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "project"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project Name: ", Default: "project"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "projectName"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter AES Key: "}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "key"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Graphql Query Timeout: ", Default: "10"}, &contextTime, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), "key"},
				},
			},
			wantErr: true,
		},
		{
			name: "spec object created",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID: "}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "project"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project Name: ", Default: "project"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "projectName"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter AES Key: "}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "key"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Graphql Query Timeout: ", Default: "10"}, &contextTime, mock.Anything},
					paramsReturned: []interface{}{nil, 15},
				},
			},
			want: &model.SpecObject{
				API:  "/v1/config/projects/{project}",
				Type: "project",
				Meta: map[string]string{
					"project": "project",
				},
				Spec: map[string]interface{}{
					"id":                 "project",
					"aesKey":             "key",
					"name":               "projectName",
					"secrets":            []map[string]interface{}{{"isPrimary": true, "secret": "ksuid string"}},
					"contextTimeGraphQL": 15,
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

			got, err := generateProject()
			if (err != nil) != tt.wantErr {
				t.Errorf("generateProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != nil {
				got.Spec.(map[string]interface{})["secrets"] = []map[string]interface{}{{"isPrimary": true, "secret": "ksuid string"}}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateProject() = %v, want %v", got, tt.want)
			}

			mockSurvey.AssertExpectations(t)
		})
	}
}
