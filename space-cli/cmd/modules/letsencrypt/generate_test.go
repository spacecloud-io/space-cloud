package letsencrypt

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

func Test_generateLetsEncryptDomain(t *testing.T) {
	someString := ""
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
			name: "error surveying white listed domains",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter White Listed Domain by comma seperated value: "}, &someString, mock.Anything},
					paramsReturned: []interface{}{errors.New("some error"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "error surveying project",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter White Listed Domain by comma seperated value: "}, &someString, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{mock.Anything, &someString, mock.Anything},
					paramsReturned: []interface{}{errors.New("some error"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "no error surveying anything",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter White Listed Domain by comma seperated value: "}, &someString, mock.Anything},
					paramsReturned: []interface{}{nil, "domain1,domain2"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{mock.Anything, &someString, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
			},
			want: &model.SpecObject{
				API:  "/v1/config/projects/{project}/letsencrypt/config/{id}",
				Type: "letsencrypt",
				Meta: map[string]string{
					"project": "",
					"id":      "letsencrypt-config",
				},
				Spec: map[string]interface{}{
					"domains": []string{"domain1", "domain2"},
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

			got, err := generateLetsEncryptDomain()
			if (err != nil) != tt.wantErr {
				t.Errorf("generateLetsEncryptDomain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateLetsEncryptDomain() = %v, want %v", got, tt.want)
			}

			mockSurvey.AssertExpectations(t)
		})
	}
}
