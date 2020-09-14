package services

import (
	"errors"
	"reflect"
	"testing"

	"github.com/AlecAivazis/survey/v2"
	"github.com/stretchr/testify/mock"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils/input"
)

func TestGenerateService(t *testing.T) {
	// surveyReturnValue stores the values returned from the survey
	surveyReturnValue := ""
	// initialized with N to to pass into test case where want is used again
	want := "N"
	notAutoDockerImage := "not-auto"
	var port int32
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		projectID   string
		dockerImage string
	}
	tests := []struct {
		name           string
		args           args
		surveyMockArgs []mockArgs
		want           *model.SpecObject
		wantErr        bool
	}{
		{
			name: "error surveying project id",
			args: args{projectID: "", dockerImage: "not-auto"},
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
			name: "error surveying service ID",
			args: args{projectID: "projectID", dockerImage: "not-auto"},
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
			name: "error surveying service version",
			args: args{projectID: "projectID", dockerImage: "not-auto"},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service ID"}, &surveyReturnValue, mock.Anything},
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
			name: "error surveying port",
			args: args{projectID: "projectID", dockerImage: "not-auto"},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service Version", Default: "v1"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{mock.Anything, &port, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "error getting project config",
			args: args{projectID: "projectID", dockerImage: "auto"},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service Version", Default: "v1"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service Port", Default: "8080"}, &port, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
			},
			wantErr: true,
		},
		{
			name: "error getting surveying about private docker registry",
			args: args{projectID: "projectID", dockerImage: notAutoDockerImage},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service Version", Default: "v1"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service Port", Default: "8080"}, &port, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Are you using private docker registry (Y / N) ?", Default: "N"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "error getting surveying docker secret",
			args: args{projectID: "projectID", dockerImage: notAutoDockerImage},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service Version", Default: "v1"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service Port", Default: "8080"}, &port, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Are you using private docker registry (Y / N) ?", Default: "N"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "Y"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Docker Secret"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "error surveying other secrets",
			args: args{projectID: "projectID", dockerImage: notAutoDockerImage},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service Version", Default: "v1"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service Port", Default: "8080"}, &port, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Are you using private docker registry (Y / N) ?", Default: "N"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "N"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Do you want to add other secrets (Y / N) ?", Default: "N"}, &want, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "error surveying File and Environment Secret (CSV)",
			args: args{projectID: "projectID", dockerImage: notAutoDockerImage},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service Version", Default: "v1"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service Port", Default: "8080"}, &port, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Are you using private docker registry (Y / N) ?", Default: "N"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "N"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Do you want to add other secrets (Y / N) ?", Default: "N"}, &want, mock.Anything},
					paramsReturned: []interface{}{nil, "Y"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter File & Environment Secret (CSV)"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "error surveying replica range",
			args: args{projectID: "projectID", dockerImage: notAutoDockerImage},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service Version", Default: "v1"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service Port", Default: "8080"}, &port, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Are you using private docker registry (Y / N) ?", Default: "N"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "N"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Do you want to add other secrets (Y / N) ?", Default: "N"}, &want, mock.Anything},
					paramsReturned: []interface{}{nil, "Y"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter File & Environment Secret (CSV)"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "filesecret1,filesecret2"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Replica Range", Default: "1-100"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "error converting replica min to int",
			args: args{projectID: "projectID", dockerImage: notAutoDockerImage},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service Version", Default: "v1"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service Port", Default: "8080"}, &port, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Are you using private docker registry (Y / N) ?", Default: "N"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "N"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Do you want to add other secrets (Y / N) ?", Default: "N"}, &want, mock.Anything},
					paramsReturned: []interface{}{nil, "Y"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter File & Environment Secret (CSV)"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "filesecret1,filesecret2"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Replica Range", Default: "1-100"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "notInt-90"},
				},
			},
			wantErr: true,
		},
		{
			name: "error converting replica max to int",
			args: args{projectID: "projectID", dockerImage: notAutoDockerImage},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service Version", Default: "v1"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service Port", Default: "8080"}, &port, mock.Anything},
					paramsReturned: []interface{}{nil, ""},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Are you using private docker registry (Y / N) ?", Default: "N"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "N"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Do you want to add other secrets (Y / N) ?", Default: "N"}, &want, mock.Anything},
					paramsReturned: []interface{}{nil, "Y"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter File & Environment Secret (CSV)"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "filesecret1,filesecret2"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Replica Range", Default: "1-100"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "10-notInt"},
				},
			},
			wantErr: true,
		},
		{
			name: "spec object created",
			args: args{projectID: "projectID", dockerImage: "image"},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service ID"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "service"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service Version", Default: "v1"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "v1"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Service Port", Default: "8080"}, &port, mock.Anything},
					paramsReturned: []interface{}{nil, "8080"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Are you using private docker registry (Y / N) ?", Default: "N"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "N"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Do you want to add other secrets (Y / N) ?", Default: "N"}, &want, mock.Anything},
					paramsReturned: []interface{}{nil, "N"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Replica Range", Default: "1-100"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "10-90"},
				},
			},
			want: &model.SpecObject{
				API:  "/v1/runner/{project}/services/{id}/{version}",
				Type: "service",
				Meta: map[string]string{
					"id":      "service",
					"project": "projectID",
					"version": "v1",
				},
				Spec: &model.Service{
					Labels: map[string]string{},
					Scale:  model.ScaleConfig{Replicas: int32(10), MinReplicas: int32(10), MaxReplicas: int32(90), Concurrency: 50, Mode: "parallel"},
					Tasks: []model.Task{
						{
							ID:        "service",
							Ports:     []model.Port{{Name: "http", Protocol: "http", Port: port}},
							Resources: model.Resources{CPU: 250, Memory: 512},
							Docker:    model.Docker{ImagePullPolicy: model.PullIfNotExists, Image: "image", Secret: "", Cmd: []string{}},
							Runtime:   model.Image,
							Secrets:   []string{},
							Env:       map[string]string{},
						},
					},
					Affinity:  []model.Affinity{},
					Whitelist: []model.Whitelist{{ProjectID: "projectID", Service: "*"}},
					Upstreams: []model.Upstream{{ProjectID: "projectID", Service: "*"}},
				},
			},
		},
		// TODO: write test cases where no error while getting project config
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockSurvey := utils.MockInputInterface{}

			for _, m := range tt.surveyMockArgs {
				mockSurvey.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			input.Survey = &mockSurvey

			got, err := GenerateService(tt.args.projectID, tt.args.dockerImage)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateService() = %v, want %v", got, tt.want)
			}

			mockSurvey.AssertExpectations(t)
		})
	}
}

func TestGenerateServiceRoute(t *testing.T) {
	// surveyReturnValue stores the values returned from the survey
	surveyReturnValue := ""
	// surveyReturnPortValue stores the values returned from the survey for port
	var surveyReturnPortValue int32
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		projectID string
	}
	tests := []struct {
		name           string
		args           args
		surveyMockArgs []mockArgs
		want           *model.SpecObject
		wantErr        bool
	}{
		{
			name: "unable to survey project id",
			args: args{projectID: ""},
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
			name: "unable to survey id",
			args: args{projectID: "project"},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Route ID:"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), "routeID"},
				},
			},
			wantErr: true,
		},
		{
			name: "unable to survey port",
			args: args{projectID: "project"},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Route ID:"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "routeID"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Port:", Default: "8080"}, &surveyReturnPortValue, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), "8080"},
				},
			},
			wantErr: true,
		},
		{
			name: "unable to survey version",
			args: args{projectID: "project"},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Route ID:"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "routeID"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Port:", Default: "8080"}, &surveyReturnPortValue, mock.Anything},
					paramsReturned: []interface{}{nil, "8080"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Version:"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), "version"},
				},
			},
			wantErr: true,
		},
		{
			name: "service route generated",
			args: args{projectID: "project"},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Route ID:"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "routeID"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Port:", Default: "8080"}, &surveyReturnPortValue, mock.Anything},
					paramsReturned: []interface{}{nil, int32(8080)},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Version:"}, &surveyReturnValue, mock.Anything},
					paramsReturned: []interface{}{nil, "version"},
				},
			},
			want: &model.SpecObject{
				API:  "/v1/runner/{project}/service-routes/{id}",
				Type: "service-route",
				Meta: map[string]string{
					"id":      "routeID",
					"project": "project",
				},
				Spec: map[string]interface{}{
					"routes": []interface{}{
						map[string]interface{}{
							"source": map[string]interface{}{
								"port": int32(8080),
							},
							"targets": []interface{}{
								map[string]interface{}{
									"type":    "internal",
									"version": "version",
									"weight":  100,
								},
							},
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

			got, err := GenerateServiceRoute(tt.args.projectID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateServiceRoute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateServiceRoute() = %v, want %v", got, tt.want)
			}

			mockSurvey.AssertExpectations(t)
		})
	}
}
