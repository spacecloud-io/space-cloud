package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/AlecAivazis/survey/v2"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/mock"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils/input"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils/transport"
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
		name              string
		args              args
		surveyMockArgs    []mockArgs
		transportMockArgs []mockArgs
		want              *model.SpecObject
		wantErr           bool
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
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodGet,
						"/v1/config/projects/projectID",
						mock.Anything,
						new(model.Response),
					},
					paramsReturned: []interface{}{
						errors.New("bad request"),
						model.Response{
							Result: []interface{}{
								map[string]interface{}{
									"statusCode": 400,
								},
							},
						},
					},
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
					Labels:                 map[string]string{},
					StatsInclusionPrefixes: "http.inbound,cluster_manager,listener_manager",
					AutoScale: &model.AutoScaleConfig{
						PollingInterval:  int32(15),
						CoolDownInterval: int32(120),
						MinReplicas:      int32(10),
						MaxReplicas:      int32(90),
						Triggers: []model.AutoScaleTrigger{
							{
								Name:             "Request per second",
								Type:             "requests-per-second",
								MetaData:         map[string]string{"target": "50"},
								AuthenticatedRef: nil,
							},
						},
					},
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
			mockTransport := transport.MocketAuthProviders{}

			for _, m := range tt.surveyMockArgs {
				mockSurvey.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			for _, m := range tt.transportMockArgs {
				mockTransport.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			input.Survey = &mockSurvey
			transport.Client = &mockTransport

			got, err := GenerateService(tt.args.projectID, tt.args.dockerImage)
			if (err != nil) != tt.wantErr {

				t.Errorf("GenerateService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if arr := deep.Equal(got, tt.want); len(arr) != 0 {
				t.Errorf("GenerateService() = %v", arr)
			}

			mockSurvey.AssertExpectations(t)
			mockTransport.AssertExpectations(t)
		})
	}
}

func TestGenerateServiceRoute(t *testing.T) {
	// surveyReturnValue stores the values returned from the survey
	surveyReturnValue := ""
	// surveyReturnIntValue stores the values returned from the survey for port
	var surveyReturnIntValue int
	var surveyReturnBoolValue bool
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
		// {
		// 	name: "unable to survey port",
		// 	args: args{projectID: "project"},
		// 	surveyMockArgs: []mockArgs{
		// 		{
		// 			method:         "AskOne",
		// 			args:           []interface{}{&survey.Input{Message: "Enter Route ID:"}, &surveyReturnValue, mock.Anything},
		// 			paramsReturned: []interface{}{nil, "routeID"},
		// 		},
		// 		{
		// 			method:         "AskOne",
		// 			args:           []interface{}{&survey.Input{Message: "Enter Port:", Default: "8080"}, &surveyReturnIntValue, mock.Anything},
		// 			paramsReturned: []interface{}{errors.New("unable to call AskOne"), "8080"},
		// 		},
		// 	},
		// 	wantErr: true,
		// },
		// {
		// 	name: "unable to survey version",
		// 	args: args{projectID: "project"},
		// 	surveyMockArgs: []mockArgs{
		// 		{
		// 			method:         "AskOne",
		// 			args:           []interface{}{&survey.Input{Message: "Enter Route ID:"}, &surveyReturnValue, mock.Anything},
		// 			paramsReturned: []interface{}{nil, "routeID"},
		// 		},
		// 		{
		// 			method:         "AskOne",
		// 			args:           []interface{}{&survey.Input{Message: "Enter Port:", Default: "8080"}, &surveyReturnIntValue, mock.Anything},
		// 			paramsReturned: []interface{}{nil, "8080"},
		// 		},
		// 		{
		// 			method:         "AskOne",
		// 			args:           []interface{}{&survey.Input{Message: "Enter Version:"}, &surveyReturnValue, mock.Anything},
		// 			paramsReturned: []interface{}{errors.New("unable to call AskOne"), "version"},
		// 		},
		// 	},
		// 	wantErr: true,
		// },
		{
			name: "service route generated versioned target with url matcher",
			args: args{projectID: "project"},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Route ID:"}, &surveyReturnValue},
					paramsReturned: []interface{}{nil, "new-service"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select request protocol:", Options: []string{"http", "tcp"}}, &surveyReturnValue},
					paramsReturned: []interface{}{nil, "http"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Port:", Default: "8080"}, &surveyReturnIntValue},
					paramsReturned: []interface{}{nil, 8080},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Retries:", Default: "3"}, &surveyReturnIntValue},
					paramsReturned: []interface{}{nil, 3},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Timeout in seconds:", Default: "180"}, &surveyReturnIntValue},
					paramsReturned: []interface{}{nil, 180},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select target type:", Options: []string{"version", "external"}}, &surveyReturnValue},
					paramsReturned: []interface{}{nil, "version"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter service port:"}, &surveyReturnIntValue},
					paramsReturned: []interface{}{nil, 8080},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Version:"}, &surveyReturnValue},
					paramsReturned: []interface{}{nil, "v1"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter weight"}, &surveyReturnIntValue},
					paramsReturned: []interface{}{nil, 100},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select matcher type:", Options: []string{"url", "header"}}, &surveyReturnValue},
					paramsReturned: []interface{}{nil, "url"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select match condition:", Options: []string{"exact", "prefix", "regex"}}, &surveyReturnValue},
					paramsReturned: []interface{}{nil, "exact"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter URL:"}, &surveyReturnValue},
					paramsReturned: []interface{}{nil, "/v2"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Confirm{Message: "Do you want to ignore case?"}, &surveyReturnBoolValue},
					paramsReturned: []interface{}{nil, true},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Confirm{Message: "Do you want to add another matcher?"}, &surveyReturnBoolValue},
					paramsReturned: []interface{}{nil, false},
				},
			},
			want: &model.SpecObject{
				API:  "/v1/runner/{project}/service-routes/{id}",
				Type: "service-route",
				Meta: map[string]string{
					"id":      "new-service",
					"project": "project",
				},
				Spec: map[string]interface{}{
					"routes": []interface{}{
						map[string]interface{}{
							"requestRetries": 3,
							"requestTimeout": 180,
							"source": map[string]interface{}{
								"port":     8080,
								"protocol": "http",
							},
							"targets": []interface{}{
								map[string]interface{}{
									"type":    "version",
									"version": "v1",
									"port":    8080,
									"weight":  100,
								},
							},
							"matchers": []interface{}{
								map[string]interface{}{
									"url": map[string]interface{}{
										"value":        "/v2",
										"type":         "exact",
										"isIgnoreCase": true,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "service route generated external target with header matcher",
			args: args{projectID: "project"},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Route ID:"}, &surveyReturnValue},
					paramsReturned: []interface{}{nil, "new-service"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select request protocol:", Options: []string{"http", "tcp"}}, &surveyReturnValue},
					paramsReturned: []interface{}{nil, "http"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Port:", Default: "8080"}, &surveyReturnIntValue},
					paramsReturned: []interface{}{nil, 8080},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Retries:", Default: "3"}, &surveyReturnIntValue},
					paramsReturned: []interface{}{nil, 3},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter Timeout in seconds:", Default: "180"}, &surveyReturnIntValue},
					paramsReturned: []interface{}{nil, 180},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select target type:", Options: []string{"version", "external"}}, &surveyReturnValue},
					paramsReturned: []interface{}{nil, "external"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter service port:"}, &surveyReturnIntValue},
					paramsReturned: []interface{}{nil, 8080},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter host address:"}, &surveyReturnValue},
					paramsReturned: []interface{}{nil, "project.svc.cluster.local"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter weight"}, &surveyReturnIntValue},
					paramsReturned: []interface{}{nil, 100},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select matcher type:", Options: []string{"url", "header"}}, &surveyReturnValue},
					paramsReturned: []interface{}{nil, "header"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Select match condition:", Options: []string{"exact", "prefix", "regex", "check-presence"}}, &surveyReturnValue},
					paramsReturned: []interface{}{nil, "exact"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter header key:"}, &surveyReturnValue},
					paramsReturned: []interface{}{nil, "key"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Input{Message: "Enter header value:"}, &surveyReturnValue},
					paramsReturned: []interface{}{nil, "value"},
				},
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Confirm{Message: "Do you want to add another matcher?"}, &surveyReturnBoolValue},
					paramsReturned: []interface{}{nil, false},
				},
			},
			want: &model.SpecObject{
				API:  "/v1/runner/{project}/service-routes/{id}",
				Type: "service-route",
				Meta: map[string]string{
					"id":      "new-service",
					"project": "project",
				},
				Spec: map[string]interface{}{
					"routes": []interface{}{
						map[string]interface{}{
							"requestRetries": 3,
							"requestTimeout": 180,
							"source": map[string]interface{}{
								"port":     8080,
								"protocol": "http",
							},
							"targets": []interface{}{
								map[string]interface{}{
									"type":   "external",
									"host":   "project.svc.cluster.local",
									"port":   8080,
									"weight": 100,
								},
							},
							"matchers": []interface{}{
								map[string]interface{}{
									"headers": []interface{}{
										map[string]interface{}{
											"key":   "key",
											"value": "value",
											"type":  "exact",
										},
									},
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

			if arr := deep.Equal(got, tt.want); len(arr) > 0 {
				a, _ := json.MarshalIndent(arr, "", " ")
				t.Errorf("GetServiceRoutes() diff = %v", string(a))
			}

			mockSurvey.AssertExpectations(t)
		})
	}
}
