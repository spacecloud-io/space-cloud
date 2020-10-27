package eventing

import (
	"errors"
	"net/http"
	"testing"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils/transport"
)

func Test_deleteEventingConfig(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		project string
	}
	tests := []struct {
		name              string
		args              args
		transportMockArgs []mockArgs
		wantErr           bool
	}{
		{
			name: "Unable to delete eventing config",
			args: args{project: "myproject"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodPost,
						"/v1/config/projects/myproject/eventing/config/eventing-config",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						errors.New("bad request"),
						map[string]interface{}{
							"statusCode": 400,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Eventing config deleted succesfully",
			args: args{project: "myproject"},
			transportMockArgs: []mockArgs{
				{
					method: "MakeHTTPRequest",
					args: []interface{}{
						http.MethodPost,
						"/v1/config/projects/myproject/eventing/config/eventing-config",
						map[string]string{},
						new(model.Response),
					},
					paramsReturned: []interface{}{
						nil,
						map[string]interface{}{
							"statusCode": 200,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockTransport := transport.MocketAuthProviders{}

			for _, m := range tt.transportMockArgs {
				mockTransport.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			transport.Client = &mockTransport

			if err := deleteEventingConfig(tt.args.project); (err != nil) != tt.wantErr {
				t.Errorf("deleteEventingConfig() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockTransport.AssertExpectations(t)
		})
	}
}
