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
	someString := ""
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
					args:           []interface{}{mock.Anything, &someString, mock.Anything},
					paramsReturned: []interface{}{errors.New("some error"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "empty project id error",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{mock.Anything, &someString, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
			},
		},
		{
			name: "error surveying project name",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID: "}, &someString, mock.Anything},
					paramsReturned: []interface{}{nil, "project"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project Name: ", Default: "project"}, &someString, mock.Anything},
					paramsReturned: []interface{}{errors.New("some error"), "projectName"},
				},
			},
			wantErr: true,
		},
		{
			name: "error surveying aes key",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID: "}, &someString, mock.Anything},
					paramsReturned: []interface{}{nil, "project"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project Name: ", Default: "project"}, &someString, mock.Anything},
					paramsReturned: []interface{}{nil, "projectName"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter AES Key: "}, &someString, mock.Anything},
					paramsReturned: []interface{}{errors.New("some error"), "key"},
				},
			},
			wantErr: true,
		},
		{
			name: "error surveying context time",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID: "}, &someString, mock.Anything},
					paramsReturned: []interface{}{nil, "project"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project Name: ", Default: "project"}, &someString, mock.Anything},
					paramsReturned: []interface{}{nil, "projectName"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter AES Key: "}, &someString, mock.Anything},
					paramsReturned: []interface{}{nil, "key"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Graphql Query Timeout: ", Default: "10"}, &contextTime, mock.Anything},
					paramsReturned: []interface{}{errors.New("some error"), "key"},
				},
			},
			wantErr: true,
		},
		// {
		// 	name: "no error surveying anyting",
		// 	surveyMockArgs: []mockArgs{
		// 		{
		// 			method:         "AskOne",
		// 			args:           []interface{}{&survey.Input{Message: "Enter Project ID: "}, &someString, mock.Anything},
		// 			paramsReturned: []interface{}{nil, "project"},
		// 		},
		// 		{
		// 			method:         "AskOne",
		// 			args:           []interface{}{&survey.Input{Message: "Enter Project Name: ", Default: "project"}, &someString, mock.Anything},
		// 			paramsReturned: []interface{}{nil, "projectName"},
		// 		},
		// 		{
		// 			method:         "AskOne",
		// 			args:           []interface{}{&survey.Input{Message: "Enter AES Key: "}, &someString, mock.Anything},
		// 			paramsReturned: []interface{}{nil, "key"},
		// 		},
		// 		{
		// 			method:         "AskOne",
		// 			args:           []interface{}{&survey.Input{Message: "Enter Graphql Query Timeout: ", Default: "10"}, &contextTime, mock.Anything},
		// 			paramsReturned: []interface{}{nil, 15},
		// 		},
		// 	},
		// 	want: &model.SpecObject{
		// 		API:  "/v1/config/projects/{project}",
		// 		Type: "project",
		// 		Meta: map[string]string{
		// 			"project": "project",
		// 		},
		// 		Spec: map[string]interface{}{
		// 			"id":                 "project",
		// 			"aesKey":             "key",
		// 			"name":               "projectName",
		// 			"secrets":            []map[string]interface{}{{"isPrimary": true, "secret": ksuid.New().String()}},
		// 			"contextTimeGraphQL": 15,
		// 		},
		// 	},
		// },
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateProject() = %v, want %v", got, tt.want)
			}

			mockSurvey.AssertExpectations(t)
		})
	}
}
