package eventing

import (
	"testing"
	"time"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

func TestModule_ProcessTransmittedEvents(t *testing.T) {
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	type args struct {
		eventDocs []*model.EventDocument
	}
	tests := []struct {
		name         string
		m            *Module
		args         args
		synvMockArgs []mockArgs
	}{
		{
			name: "token less than start",
			m:    &Module{},
			args: args{eventDocs: []*model.EventDocument{
				&model.EventDocument{
					ID:        "id",
					Token:     -1,
					Timestamp: time.Now().Format(time.RFC3339),
				},
			}},
			synvMockArgs: []mockArgs{
				mockArgs{
					method:         "GetAssignedTokens",
					args:           []interface{}{},
					paramsReturned: []interface{}{1, 100},
				},
			},
		},
		{
			name: "error parsing timestamp",
			m:    &Module{},
			args: args{eventDocs: []*model.EventDocument{
				&model.EventDocument{
					ID:        "id",
					Token:     50,
					Timestamp: "",
				},
			}},
			synvMockArgs: []mockArgs{
				mockArgs{
					method:         "GetAssignedTokens",
					args:           []interface{}{},
					paramsReturned: []interface{}{1, 100},
				},
			},
		},
		{
			name: "current timestamp not equal to or after timestamp",
			m:    &Module{},
			args: args{eventDocs: []*model.EventDocument{
				&model.EventDocument{
					ID:        "id",
					Token:     50,
					Timestamp: "5020-03-31T16:16:26+05:30",
				},
			}},
			synvMockArgs: []mockArgs{
				mockArgs{
					method:         "GetAssignedTokens",
					args:           []interface{}{},
					paramsReturned: []interface{}{1, 100},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockSync := mockSyncmanEventingInterface{}

			for _, m := range tt.synvMockArgs {
				mockSync.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			tt.m.syncMan = &mockSync

			tt.m.ProcessTransmittedEvents(tt.args.eventDocs)

			mockSync.AssertExpectations(t)
		})
	}
}

// TODO: cover the goroutine as well
