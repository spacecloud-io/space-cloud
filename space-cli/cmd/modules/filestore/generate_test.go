package filestore

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

func Test_generateFilestoreRule(t *testing.T) {
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
					args:           []interface{}{&survey.Input{Message: "Enter Project ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "error surveying id/rule name",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Rule Name"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "error surveying prefix",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Rule Name"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Prefix", Default: "/"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "file store rule spec object created",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "project"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Rule Name"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "rule"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Prefix", Default: "/"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "/"},
				},
			},
			want: &model.SpecObject{
				API:  "/v1/config/projects/{project}/file-storage/rules/{id}",
				Type: "filestore-rule",
				Meta: map[string]string{
					"id":      "rule",
					"project": "project",
				},
				Spec: map[string]interface{}{
					"prefix": "/",
					"rule": map[string]interface{}{
						"create": map[string]interface{}{
							"rule": "allow",
						},
						"delete": map[string]interface{}{
							"rule": "allow",
						},
						"read": map[string]interface{}{
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

			got, err := generateFilestoreRule()
			if (err != nil) != tt.wantErr {
				t.Errorf("generateFilestoreRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateFilestoreRule() = %v, want %v", got, tt.want)
			}

			mockSurvey.AssertExpectations(t)
		})
	}
}

func Test_generateFilestoreConfig(t *testing.T) {
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
					args:           []interface{}{&survey.Input{Message: "Enter Project ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "error surveying store type",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Enter Storetype", Options: []string{"local", "amazon-s3", "gcp-storage"}}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "store type default case",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Enter Storetype", Options: []string{"local", "amazon-s3", "gcp-storage"}}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "default"},
				},
			},
			wantErr: true,
		},
		{
			name: "error surveying connection with storetype local",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Enter Storetype", Options: []string{"local", "amazon-s3", "gcp-storage"}}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "local"},
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
			name: "error surveying connection with storetype amazon-s3",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Enter Storetype", Options: []string{"local", "amazon-s3", "gcp-storage"}}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "amazon-s3"},
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
			name: "error surveying endpoint with storetype amazon-s3",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Enter Storetype", Options: []string{"local", "amazon-s3", "gcp-storage"}}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "amazon-s3"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter connection"}, &surveyReturnValue, mock.Anything},
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
			name: "error surveying bucket with storetype amazon-s3",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Enter Storetype", Options: []string{"local", "amazon-s3", "gcp-storage"}}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "amazon-s3"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter connection"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter endpoint"}, &surveyReturnValue, mock.Anything},
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
			name: "file store config spec object created with storetype amazon-s3",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "project"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Enter Storetype", Options: []string{"local", "amazon-s3", "gcp-storage"}}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "amazon-s3"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter connection"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "connection-string"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter endpoint"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "www.some-endpoint/"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{mock.Anything, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "bucket"},
				},
			},
			want: &model.SpecObject{
				API:  "/v1/config/projects/{project}/file-storage/config/{id}",
				Type: "filestore-config",
				Meta: map[string]string{
					"project": "project",
					"id":      "filestore-config",
				},
				Spec: map[string]interface{}{
					"bucket":    "bucket",
					"conn":      "connection-string",
					"enabled":   true,
					"endpoint":  "www.some-endpoint/",
					"storeType": "amazon-s3",
				},
			},
		},
		{
			name: "error surveying bucket with storetype gcp-storage",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Enter Storetype", Options: []string{"local", "amazon-s3", "gcp-storage"}}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "gcp-storage"},
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
			name: "file store config spec object created with store type gcp",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "project"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Enter Storetype", Options: []string{"local", "amazon-s3", "gcp-storage"}}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "gcp-storage"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{mock.Anything, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "bucket"},
				},
			},
			want: &model.SpecObject{
				API:  "/v1/config/projects/{project}/file-storage/config/{id}",
				Type: "filestore-config",
				Meta: map[string]string{
					"project": "project",
					"id":      "filestore-config",
				},
				Spec: map[string]interface{}{
					"bucket":    "bucket",
					"conn":      "",
					"enabled":   true,
					"endpoint":  "",
					"storeType": "gcp-storage",
				},
			},
		},
		{
			name: "file store config spec object created with store type local",
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Project ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "project"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Enter Storetype", Options: []string{"local", "amazon-s3", "gcp-storage"}}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "local"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{mock.Anything, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "connection-string"},
				},
			},
			want: &model.SpecObject{
				API:  "/v1/config/projects/{project}/file-storage/config/{id}",
				Type: "filestore-config",
				Meta: map[string]string{
					"project": "project",
					"id":      "filestore-config",
				},
				Spec: map[string]interface{}{
					"bucket":    "",
					"conn":      "connection-string",
					"enabled":   true,
					"endpoint":  "",
					"storeType": "local",
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

			got, err := generateFilestoreConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("generateFilestoreConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateFilestoreConfig() = %v, want %v", got, tt.want)
			}

			mockSurvey.AssertExpectations(t)
		})
	}
}
