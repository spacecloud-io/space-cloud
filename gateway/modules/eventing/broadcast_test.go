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
		args         args
		syncMockArgs []mockArgs
	}{
		{
			name: "token less than start",
			args: args{eventDocs: []*model.EventDocument{
				&model.EventDocument{
					ID:        "id",
					Token:     -1,
					Timestamp: time.Now().Format(time.RFC3339),
				},
			}},
			syncMockArgs: []mockArgs{
				mockArgs{
					method:         "GetAssignedTokens",
					args:           []interface{}{},
					paramsReturned: []interface{}{1, 100},
				},
			},
		},
		{
			name: "error parsing timestamp",
			args: args{eventDocs: []*model.EventDocument{
				&model.EventDocument{
					ID:        "id",
					Token:     50,
					Timestamp: "",
				},
			}},
			syncMockArgs: []mockArgs{
				mockArgs{
					method:         "GetAssignedTokens",
					args:           []interface{}{},
					paramsReturned: []interface{}{1, 100},
				},
			},
		},
		{
			name: "current timestamp not equal to or after timestamp",
			args: args{eventDocs: []*model.EventDocument{
				&model.EventDocument{
					ID:        "id",
					Token:     50,
					Timestamp: "5020-03-31T16:16:26+05:30",
				},
			}},
			syncMockArgs: []mockArgs{
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

			m := &Module{}

			mockSync := mockSyncmanEventingInterface{}

			for _, m := range tt.syncMockArgs {
				mockSync.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			m.syncMan = &mockSync

			m.ProcessTransmittedEvents(tt.args.eventDocs)

			mockSync.AssertExpectations(t)
		})
	}
}

// TODO: cover the goroutine as well
