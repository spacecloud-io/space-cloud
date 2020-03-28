package eventing

import (
	"testing"
	"time"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
	"github.com/spaceuptech/space-cloud/gateway/utils/syncman"
)

func TestModule_ProcessTransmittedEvents(t *testing.T) {

	adminMod := admin.New()
	syncmanMod, _ := syncman.New("nodeID", "clusterID", "advertiseAddr", "storeType", "runnerAddr", "artifactAddr", adminMod)

	type args struct {
		eventDocs []*model.EventDocument
	}
	tests := []struct {
		name string
		m    *Module
		args args
	}{
		{
			name: "token less than start",
			m:    &Module{syncMan: syncmanMod},
			args: args{eventDocs: []*model.EventDocument{&model.EventDocument{ID: "id", BatchID: "batchid", Type: "DB_INSERT", RuleName: "some-rule", EventTimestamp: int64(100), Payload: "payload", Remark: "remark", Retries: 3, Status: "ok", Timestamp: time.Now().UTC().UnixNano()/int64(time.Millisecond) + int64(10000), Token: -1, URL: "url"}}},
		},
		{
			name: "current timestamp < timestamp",
			m:    &Module{syncMan: syncmanMod},
			args: args{eventDocs: []*model.EventDocument{&model.EventDocument{ID: "id", BatchID: "batchid", Type: "DB_INSERT", RuleName: "some-rule", EventTimestamp: int64(100), Payload: "payload", Remark: "remark", Retries: 3, Status: "ok", Timestamp: time.Now().UTC().UnixNano()/int64(time.Millisecond) + int64(10000), Token: 50, URL: "url"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.ProcessTransmittedEvents(tt.args.eventDocs)
		})
	}
}

// TODO: cover the goroutine as well
