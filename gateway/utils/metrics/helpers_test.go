package metrics

import (
	"reflect"
	"sync"
	"testing"

	"github.com/spaceuptech/space-api-go/db"
	"github.com/stretchr/testify/mock"

	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
	"github.com/spaceuptech/space-cloud/gateway/utils/syncman"
)

func TestModule_createCrudDocuments(t *testing.T) {
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
		key   string
		value *metricOperations
		t     string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []interface{}
	}{
		{
			name:   "valid test case",
			fields: fields{nodeID: "nodeID", clusterID: "clusterID"},
			args: args{
				key:   "project:dbAlias:tableName",
				value: &metricOperations{create: 100, update: 100, read: 100, delete: 100},
				t:     mock.Anything,
			},
			want: []interface{}{
				map[string]interface{}{"project_id": "project", "module": "db", "type": utils.Create, "sub_type": "tableName", "ts": mock.Anything, "count": uint64(100), "driver": "dbAlias", "node_id": "nodeID", "cluster_id": "clusterID"},
				map[string]interface{}{"project_id": "project", "module": "db", "type": utils.Read, "sub_type": "tableName", "ts": mock.Anything, "count": uint64(100), "driver": "dbAlias", "node_id": "nodeID", "cluster_id": "clusterID"},
				map[string]interface{}{"project_id": "project", "module": "db", "type": utils.Update, "sub_type": "tableName", "ts": mock.Anything, "count": uint64(100), "driver": "dbAlias", "node_id": "nodeID", "cluster_id": "clusterID"},
				map[string]interface{}{"project_id": "project", "module": "db", "type": utils.Delete, "sub_type": "tableName", "ts": mock.Anything, "count": uint64(100), "driver": "dbAlias", "node_id": "nodeID", "cluster_id": "clusterID"},
			},
		},
		{
			name:   "valid test case read & update are zero",
			fields: fields{nodeID: "nodeID", clusterID: "clusterID"},
			args: args{
				key:   "project:dbAlias:tableName",
				value: &metricOperations{create: 100, delete: 100},
				t:     mock.Anything,
			},
			want: []interface{}{
				map[string]interface{}{"project_id": "project", "module": "db", "type": utils.Create, "sub_type": "tableName", "ts": mock.Anything, "count": uint64(100), "driver": "dbAlias", "node_id": "nodeID", "cluster_id": "clusterID"},
				map[string]interface{}{"project_id": "project", "module": "db", "type": utils.Delete, "sub_type": "tableName", "ts": mock.Anything, "count": uint64(100), "driver": "dbAlias", "node_id": "nodeID", "cluster_id": "clusterID"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Module{
				isProd:           tt.fields.isProd,
				clusterID:        tt.fields.clusterID,
				nodeID:           tt.fields.nodeID,
				projects:         tt.fields.projects,
				isMetricDisabled: tt.fields.isMetricDisabled,
				sink:             tt.fields.sink,
				adminMan:         tt.fields.adminMan,
				syncMan:          tt.fields.syncMan,
			}
			got := m.createCrudDocuments(tt.args.key, tt.args.value, tt.args.t)
			if len(got) != len(tt.want) {
				t.Errorf("createCrudDocuments() want & got length mismatch got = %v want = %v", len(got), len(tt.want))
			}
			for index, value := range tt.want {
				for key, wantValue := range value.(map[string]interface{}) {
					if key == "id" {
						continue
					}
					gotValue, ok := got[index].(map[string]interface{})[key]
					if !ok {
						t.Errorf("createCrudDocuments() key = %s doesn't exist in result", key)
						continue
					}
					if !reflect.DeepEqual(gotValue, wantValue) {
						t.Errorf("createCrudDocuments() got value = %v %T want = %v %T", gotValue, gotValue, wantValue, wantValue)
					}
				}
			}
		})
	}
}

func TestModule_createDocument(t *testing.T) {
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
		driver  string
		subType string
		module  string
		op      utils.OperationType
		count   uint64
		t       string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Module{
				isProd:           tt.fields.isProd,
				clusterID:        tt.fields.clusterID,
				nodeID:           tt.fields.nodeID,
				projects:         tt.fields.projects,
				isMetricDisabled: tt.fields.isMetricDisabled,
				sink:             tt.fields.sink,
				adminMan:         tt.fields.adminMan,
				syncMan:          tt.fields.syncMan,
			}
			if got := m.createDocument(tt.args.project, tt.args.driver, tt.args.subType, tt.args.module, tt.args.op, tt.args.count, tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createDocument() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModule_createEventDocument(t *testing.T) {
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
		key   string
		value uint64
		t     string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []interface{}
	}{
		{
			name:   "valid test case",
			fields: fields{nodeID: "nodeID", clusterID: "clusterID"},
			args: args{
				key:   "project:event-name",
				value: 100,
				t:     mock.Anything,
			},
			want: []interface{}{
				map[string]interface{}{"project_id": "project", "module": "eventing", "type": utils.OperationType("event-name"), "sub_type": "na", "ts": mock.Anything, "count": uint64(100), "driver": "na", "node_id": "nodeID", "cluster_id": "clusterID"},
			},
		},
		{
			name:   "valid test case read & list are zero",
			fields: fields{nodeID: "nodeID", clusterID: "clusterID"},
			args: args{
				key:   "project:local",
				value: 0,
				t:     mock.Anything,
			},
			want: []interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Module{
				isProd:           tt.fields.isProd,
				clusterID:        tt.fields.clusterID,
				nodeID:           tt.fields.nodeID,
				projects:         tt.fields.projects,
				isMetricDisabled: tt.fields.isMetricDisabled,
				sink:             tt.fields.sink,
				adminMan:         tt.fields.adminMan,
				syncMan:          tt.fields.syncMan,
			}
			got := m.createEventDocument(tt.args.key, tt.args.value, tt.args.t)
			if len(got) != len(tt.want) {
				t.Errorf("createEventDocument() want & got length mismatch got = %v want = %v", len(got), len(tt.want))
			}
			for index, value := range tt.want {
				for key, wantValue := range value.(map[string]interface{}) {
					if key == "id" {
						continue
					}
					gotValue, ok := got[index].(map[string]interface{})[key]
					if !ok {
						t.Errorf("createEventDocument() key = %s doesn't exist in result", key)
						continue
					}
					if !reflect.DeepEqual(gotValue, wantValue) {
						t.Errorf("createEventDocument() got value = %v %T want = %v %T", gotValue, gotValue, wantValue, wantValue)
					}
				}
			}
		})
	}
}

func TestModule_createFileDocuments(t *testing.T) {
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
		key   string
		value *metricOperations
		t     string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []interface{}
	}{
		{
			name:   "valid test case",
			fields: fields{nodeID: "nodeID", clusterID: "clusterID"},
			args: args{
				key:   "project:local:tableName",
				value: &metricOperations{create: 100, list: 100, read: 100, delete: 100},
				t:     mock.Anything,
			},
			want: []interface{}{
				map[string]interface{}{"project_id": "project", "module": "file", "type": utils.Create, "sub_type": "na", "ts": mock.Anything, "count": uint64(100), "driver": "local", "node_id": "nodeID", "cluster_id": "clusterID"},
				map[string]interface{}{"project_id": "project", "module": "file", "type": utils.Read, "sub_type": "na", "ts": mock.Anything, "count": uint64(100), "driver": "local", "node_id": "nodeID", "cluster_id": "clusterID"},
				map[string]interface{}{"project_id": "project", "module": "file", "type": utils.Delete, "sub_type": "na", "ts": mock.Anything, "count": uint64(100), "driver": "local", "node_id": "nodeID", "cluster_id": "clusterID"},
				map[string]interface{}{"project_id": "project", "module": "file", "type": utils.List, "sub_type": "na", "ts": mock.Anything, "count": uint64(100), "driver": "local", "node_id": "nodeID", "cluster_id": "clusterID"},
			},
		},
		{
			name:   "valid test case read & list are zero",
			fields: fields{nodeID: "nodeID", clusterID: "clusterID"},
			args: args{
				key:   "project:local",
				value: &metricOperations{create: 100, delete: 100},
				t:     mock.Anything,
			},
			want: []interface{}{
				map[string]interface{}{"project_id": "project", "module": "file", "type": utils.Create, "sub_type": "na", "ts": mock.Anything, "count": uint64(100), "driver": "local", "node_id": "nodeID", "cluster_id": "clusterID"},
				map[string]interface{}{"project_id": "project", "module": "file", "type": utils.Delete, "sub_type": "na", "ts": mock.Anything, "count": uint64(100), "driver": "local", "node_id": "nodeID", "cluster_id": "clusterID"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Module{
				isProd:           tt.fields.isProd,
				clusterID:        tt.fields.clusterID,
				nodeID:           tt.fields.nodeID,
				projects:         tt.fields.projects,
				isMetricDisabled: tt.fields.isMetricDisabled,
				sink:             tt.fields.sink,
				adminMan:         tt.fields.adminMan,
				syncMan:          tt.fields.syncMan,
			}
			got := m.createFileDocuments(tt.args.key, tt.args.value, tt.args.t)
			if len(got) != len(tt.want) {
				t.Errorf("createFileDocuments() want & got length mismatch got = %v want = %v", len(got), len(tt.want))
			}
			for index, value := range tt.want {
				for key, wantValue := range value.(map[string]interface{}) {
					if key == "id" {
						continue
					}
					gotValue, ok := got[index].(map[string]interface{})[key]
					if !ok {
						t.Errorf("createFileDocuments() key = %s doesn't exist in result", key)
						continue
					}
					if !reflect.DeepEqual(gotValue, wantValue) {
						t.Errorf("createFileDocuments() got value = %v %T want = %v %T", gotValue, gotValue, wantValue, wantValue)
					}
				}
			}
		})
	}
}

func TestModule_createFunctionDocument(t *testing.T) {
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
		key   string
		value uint64
		t     string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []interface{}
	}{
		{
			name:   "valid test case",
			fields: fields{nodeID: "nodeID", clusterID: "clusterID"},
			args: args{
				key:   "project:service:endpoint",
				value: 100,
				t:     mock.Anything,
			},
			want: []interface{}{
				map[string]interface{}{"project_id": "project", "module": "function", "type": utils.OperationType("calls"), "sub_type": "endpoint", "ts": mock.Anything, "count": uint64(100), "driver": "service", "node_id": "nodeID", "cluster_id": "clusterID"},
			},
		},
		{
			name:   "valid test case read & list are zero",
			fields: fields{nodeID: "nodeID", clusterID: "clusterID"},
			args: args{
				key:   "project:local",
				value: 0,
				t:     mock.Anything,
			},
			want: []interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Module{
				isProd:           tt.fields.isProd,
				clusterID:        tt.fields.clusterID,
				nodeID:           tt.fields.nodeID,
				projects:         tt.fields.projects,
				isMetricDisabled: tt.fields.isMetricDisabled,
				sink:             tt.fields.sink,
				adminMan:         tt.fields.adminMan,
				syncMan:          tt.fields.syncMan,
			}
			got := m.createFunctionDocument(tt.args.key, tt.args.value, tt.args.t)
			if len(got) != len(tt.want) {
				t.Errorf("createFunctionDocument() want & got length mismatch got = %v want = %v", len(got), len(tt.want))
			}
			for index, value := range tt.want {
				for key, wantValue := range value.(map[string]interface{}) {
					if key == "id" {
						continue
					}
					gotValue, ok := got[index].(map[string]interface{})[key]
					if !ok {
						t.Errorf("createFunctionDocument() key = %s doesn't exist in result", key)
						continue
					}
					if !reflect.DeepEqual(gotValue, wantValue) {
						t.Errorf("createFunctionDocument() got value = %v %T want = %v %T", gotValue, gotValue, wantValue, wantValue)
					}
				}
			}
		})
	}
}
