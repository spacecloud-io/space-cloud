package crud

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func TestModule_closeBatchOperation(t *testing.T) {
	type fields struct {
		RWMutex             sync.RWMutex
		block               Crud
		dbType              string
		alias               string
		project             string
		schema              model.SchemaCrudInterface
		queries             map[string]*config.PreparedQuery
		batchMapTableToChan batchMap
		dataLoader          loader
		hooks               *model.CrudHooks
		metricHook          model.MetricCrudHook
		getSecrets          utils.GetSecrets
	}
	tests := []struct {
		name    string
		project string
		dbAlias string
		fields  fields
	}{
		{
			name:    "Correct values",
			project: "myproject",
			dbAlias: "db",
			fields: fields{
				batchMapTableToChan: map[string]map[string]map[string]batchChannels{
					"myproject": {
						"db": {
							"orders": {
								closeC: make(chan struct{}),
							},
							"customers": {
								closeC: make(chan struct{}),
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			m := &Module{
				batchMapTableToChan: tt.fields.batchMapTableToChan,
			}
			var wg sync.WaitGroup
			wg.Add(len(tt.fields.batchMapTableToChan[tt.project][tt.dbAlias]))
			for _, info := range tt.fields.batchMapTableToChan[tt.project][tt.dbAlias] {
				info := info
				go func() {
					<-info.closeC
					wg.Done()
				}()
			}
			m.closeBatchOperation()
			wg.Wait()
		})
	}
}

func TestModule_sendResponses(t *testing.T) {
	type args struct {
		responseChannels []batchResponseChan
		response         batchResponse
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Correct value",
			args: args{
				responseChannels: []batchResponseChan{make(batchResponseChan), make(batchResponseChan)},
				response:         batchResponse{err: nil},
			},
		},
	}
	m := Init()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			var wg sync.WaitGroup
			wg.Add(len(tt.args.responseChannels))
			for _, channel := range tt.args.responseChannels {
				channel := channel
				go func() {
					<-channel
					wg.Done()
				}()
			}
			m.sendResponses(tt.args.responseChannels, tt.args.response)
			wg.Wait()
		})
	}
}
