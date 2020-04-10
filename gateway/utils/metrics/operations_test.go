package metrics

import (
	"reflect"
	"sync"
	"testing"

	"github.com/spaceuptech/space-api-go/db"

	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
	"github.com/spaceuptech/space-cloud/gateway/utils/syncman"
)

func generateSyncMap(key string) sync.Map {
	v := sync.Map{}
	v.LoadOrStore(key, uint64(1))
	return v
}
func TestModule_AddDBOperation(t *testing.T) {
	type fields struct {
		lock             sync.RWMutex
		isProd           bool
		clusterID        string
		nodeID           string
		projects         sync.Map
		isMetricDisabled bool
		sink             *db.DB
		adminMan         *admin.Manager
		syncMan          *syncman.Manager
	}
	type args struct {
		project string
		dbType  string
		col     string
		count   int64
		op      utils.OperationType
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Module{
				lock:             tt.fields.lock,
				isProd:           tt.fields.isProd,
				clusterID:        tt.fields.clusterID,
				nodeID:           tt.fields.nodeID,
				projects:         tt.fields.projects,
				isMetricDisabled: tt.fields.isMetricDisabled,
				sink:             tt.fields.sink,
				adminMan:         tt.fields.adminMan,
				syncMan:          tt.fields.syncMan,
			}
			m.AddDBOperation(tt.args.project, tt.args.dbType, tt.args.col, tt.args.count, tt.args.op)
		})
	}
}

func TestModule_AddEventingType(t *testing.T) {
	type fields struct {
		lock             sync.RWMutex
		isProd           bool
		clusterID        string
		nodeID           string
		projects         sync.Map
		isMetricDisabled bool
		sink             *db.DB
		adminMan         *admin.Manager
		syncMan          *syncman.Manager
	}
	type args struct {
		project      string
		eventingType string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   sync.Map
	}{
		{
			name:   "valid case",
			fields: fields{},
			args: args{
				project:      "projectID",
				eventingType: "type",
			},
			want: generateSyncMap("eventing:projectID:type"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Module{
				lock:             tt.fields.lock,
				isProd:           tt.fields.isProd,
				clusterID:        tt.fields.clusterID,
				nodeID:           tt.fields.nodeID,
				projects:         tt.fields.projects,
				isMetricDisabled: tt.fields.isMetricDisabled,
				sink:             tt.fields.sink,
				adminMan:         tt.fields.adminMan,
				syncMan:          tt.fields.syncMan,
			}
			m.AddEventingType(tt.args.project, tt.args.eventingType)
			// if tt.want != m.projects {
			// 	t.Errorf("AddEventingType")
			// }
		})
	}
}

func TestModule_AddFileOperation(t *testing.T) {
	type fields struct {
		lock             sync.RWMutex
		isProd           bool
		clusterID        string
		nodeID           string
		projects         sync.Map
		isMetricDisabled bool
		sink             *db.DB
		adminMan         *admin.Manager
		syncMan          *syncman.Manager
	}
	type args struct {
		project   string
		storeType string
		op        utils.OperationType
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Module{
				lock:             tt.fields.lock,
				isProd:           tt.fields.isProd,
				clusterID:        tt.fields.clusterID,
				nodeID:           tt.fields.nodeID,
				projects:         tt.fields.projects,
				isMetricDisabled: tt.fields.isMetricDisabled,
				sink:             tt.fields.sink,
				adminMan:         tt.fields.adminMan,
				syncMan:          tt.fields.syncMan,
			}
			m.AddFileOperation(tt.args.project, tt.args.storeType, tt.args.op)
		})
	}
}

func TestModule_AddFunctionOperation(t *testing.T) {
	type fields struct {
		lock             sync.RWMutex
		isProd           bool
		clusterID        string
		nodeID           string
		projects         sync.Map
		isMetricDisabled bool
		sink             *db.DB
		adminMan         *admin.Manager
		syncMan          *syncman.Manager
	}
	type args struct {
		project  string
		service  string
		function string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Module{
				lock:             tt.fields.lock,
				isProd:           tt.fields.isProd,
				clusterID:        tt.fields.clusterID,
				nodeID:           tt.fields.nodeID,
				projects:         tt.fields.projects,
				isMetricDisabled: tt.fields.isMetricDisabled,
				sink:             tt.fields.sink,
				adminMan:         tt.fields.adminMan,
				syncMan:          tt.fields.syncMan,
			}
			m.AddFunctionOperation(tt.args.project, tt.args.service, tt.args.function)
		})
	}
}

func TestModule_LoadMetrics(t *testing.T) {
	type fields struct {
		lock             sync.RWMutex
		isProd           bool
		clusterID        string
		nodeID           string
		projects         sync.Map
		isMetricDisabled bool
		sink             *db.DB
		adminMan         *admin.Manager
		syncMan          *syncman.Manager
	}
	tests := []struct {
		name   string
		fields fields
		want   []interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Module{
				lock:             tt.fields.lock,
				isProd:           tt.fields.isProd,
				clusterID:        tt.fields.clusterID,
				nodeID:           tt.fields.nodeID,
				projects:         tt.fields.projects,
				isMetricDisabled: tt.fields.isMetricDisabled,
				sink:             tt.fields.sink,
				adminMan:         tt.fields.adminMan,
				syncMan:          tt.fields.syncMan,
			}
			if got := m.LoadMetrics(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}
