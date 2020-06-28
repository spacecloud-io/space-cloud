package userman

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

func Test_generateUserManagement(t *testing.T) {
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
			name: "error surveying provider",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project"}, &surveyReturnValue, mock.Anything},
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
			name: "spec object created",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "project"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{mock.Anything, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "provider"},
				},
			},
			want: &model.SpecObject{
				API:  "/v1/config/projects/{project}/user-management/provider/{id}",
				Type: "auth-providers",
				Meta: map[string]string{
					"project": "project",
					"id":      "provider",
				},
				Spec: map[string]interface{}{
					"enabled": true,
					"secret":  "",
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

			got, err := generateUserManagement()
			if (err != nil) != tt.wantErr {
				t.Errorf("generateUserManagement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateUserManagement() = %v, want %v", got, tt.want)
			}

			mockSurvey.AssertExpectations(t)
		})
	}
}
