package accounts

import (
	"errors"
	"testing"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spaceuptech/space-cli/cmd/utils"
	"github.com/spaceuptech/space-cli/cmd/utils/file"
	"github.com/spaceuptech/space-cli/cmd/utils/input"
)

func Test_setAccount(t *testing.T) {
	// singleMatchCase stores the value returned from the survey when prefix matches one account ID
	// initialized with "lo" to be consistent with test case
	singleMatchCase := "lo"
	// multipleMatchCase stores the value returned from the survey when prefix matches multiple account IDs
	// initialized with "l" to be consistent with test case
	multipleMatchCase := "l"
	// noMatchCase stores the value returned from the survey when prefix matches no account ID
	// initialized with "a" to be consistent with test case
	noMatchCase := "a"
	// emptyPrefix stores the value returned from survey when prefix is empty
	// initialized with "" to be consistent with test case
	emptyPrefix := ""

	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		prefix string
	}
	tests := []struct {
		name           string
		args           args
		schemaMockArgs []mockArgs
		surveyMockArgs []mockArgs
		wantErr        bool
	}{
		{
			name: "could not fetch stored credentials",
			args: args{prefix: "T"},
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					paramsReturned: []interface{}{[]byte{}, errors.New("couldn't read accounts.yaml file")},
				},
			},
			wantErr: true,
		},
		{
			name: "prefix matches one account but could not survey account id",
			args: args{prefix: "lo"},
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					paramsReturned: []interface{}{[]byte("accounts:\n- id: local-admin\n  username: local-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\nselectedAccount: local-admin"), nil},
				},
			},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Choose the account ID: ", Options: []string{"local-admin"}, Default: "local-admin"}, &singleMatchCase},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "prefix matches one account ID and options surveyed successfully but could not update accouts.yaml file",
			args: args{prefix: "lo"},
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					paramsReturned: []interface{}{[]byte("accounts:\n- id: local-admin\n  username: local-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\nselectedAccount: local-admin"), nil},
				},
				{
					method:         "Stat",
					paramsReturned: []interface{}{nil, errors.New("could not get file info")},
				},
				{
					method:         "IsNotExist",
					paramsReturned: []interface{}{true},
				},
				{
					method:         "MkdirAll",
					paramsReturned: []interface{}{errors.New("could not make directory")},
				},
			},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Choose the account ID: ", Options: []string{"local-admin"}, Default: "local-admin"}, &singleMatchCase},
					paramsReturned: []interface{}{nil, "local-admin"},
				},
			},
			wantErr: true,
		},
		{
			name: "prefix matches one account and option surveyed successfully and yaml file updated",
			args: args{prefix: "lo"},
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					paramsReturned: []interface{}{[]byte("accounts:\n- id: local-admin\n  username: local-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\nselectedAccount: local-admin"), nil},
				},
				{
					method:         "Stat",
					paramsReturned: []interface{}{nil, nil},
				},
				{
					method:         "IsNotExist",
					paramsReturned: []interface{}{false},
				},
				{
					method:         "WriteFile",
					paramsReturned: []interface{}{nil},
				},
			},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Choose the account ID: ", Options: []string{"local-admin"}, Default: "local-admin"}, &singleMatchCase},
					paramsReturned: []interface{}{nil, "local-admin"},
				},
			},
		},
		{
			name: "prefix matches multiple accounts but could not survey accout id",
			args: args{prefix: "l"},
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					paramsReturned: []interface{}{[]byte("accounts:\n- id: local-admin\n  username: local-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\n- id: last-admin\n  username: last-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\nselectedAccount: local-admin"), nil},
				},
			},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Choose the account ID: ", Options: []string{"local-admin", "last-admin"}, Default: "local-admin"}, &multipleMatchCase},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), ""},
				},
			},
			wantErr: true,
		},
		{
			name: "prefix matches multiple accounts and option surveyed successfully but yaml file not updated",
			args: args{prefix: "l"},
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					paramsReturned: []interface{}{[]byte("accounts:\n- id: local-admin\n  username: local-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\n- id: last-admin\n  username: last-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\nselectedAccount: local-admin"), nil},
				},
				{
					method:         "Stat",
					paramsReturned: []interface{}{nil, errors.New("could not get file info")},
				},
				{
					method:         "IsNotExist",
					paramsReturned: []interface{}{true},
				},
				{
					method:         "MkdirAll",
					paramsReturned: []interface{}{errors.New("could not make directory")},
				},
			},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Choose the account ID: ", Options: []string{"local-admin", "last-admin"}, Default: "local-admin"}, &multipleMatchCase},
					paramsReturned: []interface{}{nil, "local-admin"},
				},
			},
			wantErr: true,
		},
		{
			name: "prefix matches multiple accounts and option surveyed successfully and yaml file updated",
			args: args{prefix: "l"},
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					paramsReturned: []interface{}{[]byte("accounts:\n- id: local-admin\n  username: local-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\n- id: last-admin\n  username: last-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\nselectedAccount: local-admin"), nil},
				},
				{
					method:         "Stat",
					paramsReturned: []interface{}{nil, nil},
				},
				{
					method:         "IsNotExist",
					paramsReturned: []interface{}{false},
				},
				{
					method:         "WriteFile",
					paramsReturned: []interface{}{nil},
				},
			},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Choose the account ID: ", Options: []string{"local-admin", "last-admin"}, Default: "local-admin"}, &multipleMatchCase},
					paramsReturned: []interface{}{nil, "local-admin"},
				},
			},
		},
		{
			name: "prefix matches no accounts and option not surveyed successfully",
			args: args{prefix: "a"},
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					paramsReturned: []interface{}{[]byte("accounts:\n- id: local-admin\n  username: local-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\n- id: last-admin\n  username: last-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\nselectedAccount: local-admin"), nil},
				},
			},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Choose the account ID: ", Options: []string{"local-admin", "last-admin"}, Default: "local-admin"}, &noMatchCase},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), "local-admin"},
				},
			},
			wantErr: true,
		},
		{
			name: "prefix matches no accounts and option surveyed successfully and yaml file not updated",
			args: args{prefix: "a"},
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					paramsReturned: []interface{}{[]byte("accounts:\n- id: local-admin\n  username: local-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\n- id: last-admin\n  username: last-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\nselectedAccount: local-admin"), nil},
				},
				{
					method:         "Stat",
					paramsReturned: []interface{}{nil, errors.New("could not get file info")},
				},
				{
					method:         "IsNotExist",
					paramsReturned: []interface{}{true},
				},
				{
					method:         "MkdirAll",
					paramsReturned: []interface{}{errors.New("could not make directory")},
				},
			},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Choose the account ID: ", Options: []string{"local-admin", "last-admin"}, Default: "local-admin"}, &noMatchCase},
					paramsReturned: []interface{}{nil, "local-admin"},
				},
			},
			wantErr: true,
		},
		{
			name: "prefix matches no accounts and option surveyed successfully and yaml file updated",
			args: args{prefix: "a"},
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					paramsReturned: []interface{}{[]byte("accounts:\n- id: local-admin\n  username: local-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\n- id: last-admin\n  username: last-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\nselectedAccount: local-admin"), nil},
				},
				{
					method:         "Stat",
					paramsReturned: []interface{}{nil, nil},
				},
				{
					method:         "IsNotExist",
					paramsReturned: []interface{}{false},
				},
				{
					method:         "WriteFile",
					paramsReturned: []interface{}{nil},
				},
			},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Choose the account ID: ", Options: []string{"local-admin", "last-admin"}, Default: "local-admin"}, &noMatchCase},
					paramsReturned: []interface{}{nil, "local-admin"},
				},
			},
		},
		{
			name: "empty prefix case and option not surveyed successfully",
			args: args{prefix: ""},
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					paramsReturned: []interface{}{[]byte("accounts:\n- id: local-admin\n  username: local-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\n- id: last-admin\n  username: last-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\nselectedAccount: local-admin"), nil},
				},
			},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Choose the account ID: ", Options: []string{"local-admin", "last-admin"}, Default: "local-admin"}, &emptyPrefix},
					paramsReturned: []interface{}{errors.New("unable to call AskOne"), "local-admin"},
				},
			},
			wantErr: true,
		},
		{
			name: "empty prefix case and option surveyed successfully and yaml file not updated",
			args: args{prefix: ""},
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					paramsReturned: []interface{}{[]byte("accounts:\n- id: local-admin\n  username: local-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\n- id: last-admin\n  username: last-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\nselectedAccount: local-admin"), nil},
				},
				{
					method:         "Stat",
					paramsReturned: []interface{}{nil, errors.New("could not get file info")},
				},
				{
					method:         "IsNotExist",
					paramsReturned: []interface{}{true},
				},
				{
					method:         "MkdirAll",
					paramsReturned: []interface{}{errors.New("could not make directory")},
				},
			},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Choose the account ID: ", Options: []string{"local-admin", "last-admin"}, Default: "local-admin"}, &emptyPrefix},
					paramsReturned: []interface{}{nil, "local-admin"},
				},
			},
			wantErr: true,
		},
		{
			name: "empty prefix case and option surveyed successfully and yaml file updated",
			args: args{prefix: ""},
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					paramsReturned: []interface{}{[]byte("accounts:\n- id: local-admin\n  username: local-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\n- id: last-admin\n  username: last-admin\n  key: 81WZUGRTtHbG\n  serverurl: http://localhost:4122\nselectedAccount: local-admin"), nil},
				},
				{
					method:         "Stat",
					paramsReturned: []interface{}{nil, nil},
				},
				{
					method:         "IsNotExist",
					paramsReturned: []interface{}{false},
				},
				{
					method:         "WriteFile",
					paramsReturned: []interface{}{nil},
				},
			},
			surveyMockArgs: []mockArgs{
				{
					method:         "AskOne",
					args:           []interface{}{&survey.Select{Message: "Choose the account ID: ", Options: []string{"local-admin", "last-admin"}, Default: "local-admin"}, &emptyPrefix},
					paramsReturned: []interface{}{nil, "local-admin"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockSchema := file.Mocket{}
			mockSurvey := utils.MockInputInterface{}

			for _, m := range tt.schemaMockArgs {
				mockSchema.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			for _, m := range tt.surveyMockArgs {
				mockSurvey.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			file.File = &mockSchema
			input.Survey = &mockSurvey

			if err := setAccount(tt.args.prefix); (err != nil) != tt.wantErr {
				t.Errorf("setAccount() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockSchema.AssertExpectations(t)
			mockSurvey.AssertExpectations(t)
		})
	}
}
