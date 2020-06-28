package metrics

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
	"github.com/spaceuptech/space-cloud/gateway/utils/syncman"
)

func TestModule_createCrudDocuments(t *testing.T) {
	type args struct {
		key   string
		value *metricOperations
		t     string
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{
			name: "valid test case",
			args: args{
				key:   generateDatabaseKey("project", "dbAlias", "tableName"),
				value: &metricOperations{create: 100, update: 100, read: 100, delete: 100},
				t:     mock.Anything,
			},
			want: []interface{}{
				map[string]interface{}{"project_id": "project", "module": databaseModule, "type": utils.Create, "sub_type": "tableName", "ts": mock.Anything, "count": uint64(100), "driver": "dbAlias", "node_id": "nodeID", "cluster_id": "clusterID"},
				map[string]interface{}{"project_id": "project", "module": databaseModule, "type": utils.Read, "sub_type": "tableName", "ts": mock.Anything, "count": uint64(100), "driver": "dbAlias", "node_id": "nodeID", "cluster_id": "clusterID"},
				map[string]interface{}{"project_id": "project", "module": databaseModule, "type": utils.Update, "sub_type": "tableName", "ts": mock.Anything, "count": uint64(100), "driver": "dbAlias", "node_id": "nodeID", "cluster_id": "clusterID"},
				map[string]interface{}{"project_id": "project", "module": databaseModule, "type": utils.Delete, "sub_type": "tableName", "ts": mock.Anything, "count": uint64(100), "driver": "dbAlias", "node_id": "nodeID", "cluster_id": "clusterID"},
			},
		},
		{
			name: "valid test case read & update are zero",
			args: args{
				key:   generateDatabaseKey("project", "dbAlias", "tableName"),
				value: &metricOperations{create: 100, delete: 100},
				t:     mock.Anything,
			},
			want: []interface{}{
				map[string]interface{}{"project_id": "project", "module": databaseModule, "type": utils.Create, "sub_type": "tableName", "ts": mock.Anything, "count": uint64(100), "driver": "dbAlias", "node_id": "nodeID", "cluster_id": "clusterID"},
				map[string]interface{}{"project_id": "project", "module": databaseModule, "type": utils.Delete, "sub_type": "tableName", "ts": mock.Anything, "count": uint64(100), "driver": "dbAlias", "node_id": "nodeID", "cluster_id": "clusterID"},
			},
		},
	}
	m, _ := New("clusterID", "nodeID", false, admin.New("clusterID", &config.AdminUser{}), &syncman.Manager{}, false)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
		name string
		args args
		want interface{}
	}{
		// TODO: Add test cases.
	}
	m, _ := New("clusterID", "nodeID", false, admin.New("clusterID", &config.AdminUser{}), &syncman.Manager{}, false)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := m.createDocument(tt.args.project, tt.args.driver, tt.args.subType, tt.args.module, tt.args.op, tt.args.count, tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createDocument() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModule_createEventDocument(t *testing.T) {
	type args struct {
		key   string
		value uint64
		t     string
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{
			name: "valid test case",
			args: args{
				key:   generateEventingKey("project", "event-name"),
				value: 100,
				t:     mock.Anything,
			},
			want: []interface{}{
				map[string]interface{}{"project_id": "project", "module": eventingModule, "type": utils.OperationType("event-name"), "sub_type": notApplicable, "ts": mock.Anything, "count": uint64(100), "driver": notApplicable, "node_id": "nodeID", "cluster_id": "clusterID"},
			},
		},
		{
			name: "valid test case read & list are zero",
			args: args{
				key:   generateEventingKey("project", "event-name"),
				value: 0,
				t:     mock.Anything,
			},
			want: []interface{}{},
		},
	}

	m, _ := New("clusterID", "nodeID", false, admin.New("clusterID", &config.AdminUser{}), &syncman.Manager{}, false)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
	type args struct {
		key   string
		value *metricOperations
		t     string
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{
			name: "valid test case",
			args: args{
				key:   generateFileKey("project", "local"),
				value: &metricOperations{create: 100, list: 100, read: 100, delete: 100},
				t:     mock.Anything,
			},
			want: []interface{}{
				map[string]interface{}{"project_id": "project", "module": fileModule, "type": utils.Create, "sub_type": notApplicable, "ts": mock.Anything, "count": uint64(100), "driver": "local", "node_id": "nodeID", "cluster_id": "clusterID"},
				map[string]interface{}{"project_id": "project", "module": fileModule, "type": utils.Read, "sub_type": notApplicable, "ts": mock.Anything, "count": uint64(100), "driver": "local", "node_id": "nodeID", "cluster_id": "clusterID"},
				map[string]interface{}{"project_id": "project", "module": fileModule, "type": utils.Delete, "sub_type": notApplicable, "ts": mock.Anything, "count": uint64(100), "driver": "local", "node_id": "nodeID", "cluster_id": "clusterID"},
				map[string]interface{}{"project_id": "project", "module": fileModule, "type": utils.List, "sub_type": notApplicable, "ts": mock.Anything, "count": uint64(100), "driver": "local", "node_id": "nodeID", "cluster_id": "clusterID"},
			},
		},
		{
			name: "valid test case read & list are zero",
			args: args{
				key:   generateFileKey("project", "local"),
				value: &metricOperations{create: 100, delete: 100},
				t:     mock.Anything,
			},
			want: []interface{}{
				map[string]interface{}{"project_id": "project", "module": fileModule, "type": utils.Create, "sub_type": notApplicable, "ts": mock.Anything, "count": uint64(100), "driver": "local", "node_id": "nodeID", "cluster_id": "clusterID"},
				map[string]interface{}{"project_id": "project", "module": fileModule, "type": utils.Delete, "sub_type": notApplicable, "ts": mock.Anything, "count": uint64(100), "driver": "local", "node_id": "nodeID", "cluster_id": "clusterID"},
			},
		},
	}
	m, _ := New("clusterID", "nodeID", false, admin.New("clusterID", &config.AdminUser{}), &syncman.Manager{}, false)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
	type args struct {
		key   string
		value uint64
		t     string
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{
			name: "valid test case",
			args: args{
				key:   generateFunctionKey("project", "service", "function"),
				value: 100,
				t:     mock.Anything,
			},
			want: []interface{}{
				map[string]interface{}{"project_id": "project", "module": remoteServiceModule, "type": utils.OperationType("calls"), "sub_type": "function", "ts": mock.Anything, "count": uint64(100), "driver": "service", "node_id": "nodeID", "cluster_id": "clusterID"},
			},
		},
		{
			name: "valid test case read & list are zero",
			args: args{
				key:   generateFunctionKey("project", "service", "function"),
				value: 0,
				t:     mock.Anything,
			},
			want: []interface{}{},
		},
	}
	m, _ := New("clusterID", "nodeID", false, admin.New("clusterID", &config.AdminUser{}), &syncman.Manager{}, false)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
						t.Errorf("createFunctionDocument() key %v got value = %v %T want = %v %T", key, gotValue, gotValue, wantValue, wantValue)
					}
				}
			}
		})
	}
}
