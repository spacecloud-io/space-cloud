package remoteservices

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

func Test_generateService(t *testing.T) {
	// surveyReturnValue stores the values returned from the survey
	surveyReturnValue := ""
	// want stores the yes/no answer for adding another endpoint, initialized with y to add first endpoint
	want := "y"
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
			name: "error surveying service",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID: "}, &surveyReturnValue, mock.Anything},
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
			name: "error surveying url",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID: "}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service Name: "}, &surveyReturnValue, mock.Anything},
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
			name: "error surveying endpointName",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID: "}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service Name: "}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service URL: ", Help: "e.g -> http://localhost:8090"}, &surveyReturnValue, mock.Anything},
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
			name: "error surveying method",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID: "}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service Name: "}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service URL: ", Help: "e.g -> http://localhost:8090"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Endpoint Name: "}, &surveyReturnValue, mock.Anything},
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
			name: "error surveying path",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID: "}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service Name: "}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service URL: ", Help: "e.g -> http://localhost:8090"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Endpoint Name: "}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select Method: ", Options: []string{"POST", "PUT", "GET", "DELETE"}}, &surveyReturnValue, mock.Anything},
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
			name: "error surveying another endpoint yes/no",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID: "}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service Name: "}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service URL: ", Help: "e.g -> http://localhost:8090"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Endpoint Name: "}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select Method: ", Options: []string{"POST", "PUT", "GET", "DELETE"}}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter URL Path:", Default: "/"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{mock.Anything, &want, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), ""},
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
					args:           []interface{}{&survey.Input{Message: "Enter Service Name: "}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "service"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service URL: ", Help: "e.g -> http://localhost:8090"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "url"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Endpoint Name: "}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "endpointName"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select Method: ", Options: []string{"POST", "PUT", "GET", "DELETE"}}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "POST"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter URL Path:", Default: "/"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "/"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{mock.Anything, &want, mock.Anything},
					paramsReturned: []interface{}{nil, "n"},
				},
			},
			want: &model.SpecObject{
				API:  "/v1/config/projects/{project}/remote-service/service/{id}",
				Type: "remote-services",
				Meta: map[string]string{
					"id":      "service",
					"project": "project",
				},
				Spec: map[string]interface{}{
					"url":       "url",
					"endpoints": []interface{}{map[string]interface{}{"endpointName": map[string]interface{}{"method": "POST", "path": "/"}}},
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

			got, err := generateService()
			if (err != nil) != tt.wantErr {
				t.Errorf("generateService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateService() = %v, want %v", got, tt.want)
			}

			mockSurvey.AssertExpectations(t)
		})
	}
}
