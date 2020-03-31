package eventing

import (
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

func TestModule_ProcessTransmittedEvents(t *testing.T) {
	type args struct {
		eventDocs []*model.EventDocument
	}
	tests := []struct {
		name string
		m    *Module
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.ProcessTransmittedEvents(tt.args.eventDocs)
		})
	}
}

// TODO: cover the goroutine as well
