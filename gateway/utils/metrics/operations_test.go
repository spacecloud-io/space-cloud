package metrics

import (
	"reflect"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func testFuncIsEqual(t *testing.T, x, y interface{}) {
	if !reflect.DeepEqual(x, y) {
		t.Errorf("isEqual() got = %v, want %v", x, y)
	}
}

func testFuncGenerateEventingSyncMapData(key string) *sync.Map {
	v := sync.Map{}
	value, _ := v.LoadOrStore(key, newMetrics())
	metrics := value.(*metrics)
	atomic.AddUint64(&metrics.eventing, uint64(1))
	return &v
}

func testFuncGenerateFunctionSyncMapData(key string) *sync.Map {
	v := sync.Map{}
	value, _ := v.LoadOrStore(key, newMetrics())
	metrics := value.(*metrics)
	atomic.AddUint64(&metrics.function, uint64(1))
	return &v
}

func testFuncGenerateDatabaseSyncMapData(key string, create, read, update, delete int) *sync.Map {
	v := sync.Map{}
	value, _ := v.LoadOrStore(key, newMetrics())
	metrics := value.(*metrics)
	atomic.AddUint64(&metrics.crud.create, uint64(create))
	atomic.AddUint64(&metrics.crud.read, uint64(read))
	atomic.AddUint64(&metrics.crud.update, uint64(update))
	atomic.AddUint64(&metrics.crud.delete, uint64(delete))
	return &v
}

func testFuncGenerateFileSyncMapData(key string) *sync.Map {
	v := sync.Map{}
	value, _ := v.LoadOrStore(key, newMetrics())
	metrics := value.(*metrics)
	atomic.AddUint64(&metrics.fileStore.create, uint64(1))
	atomic.AddUint64(&metrics.fileStore.read, uint64(1))
	atomic.AddUint64(&metrics.fileStore.list, uint64(1))
	atomic.AddUint64(&metrics.fileStore.delete, uint64(1))
	return &v
}

func testFuncIsMetricOpEqual(t *testing.T, gotValue, wantValue interface{}, op utils.OperationType, module string) {
	gotMetrics := gotValue.(*metrics)
	wantMetrics := wantValue.(*metrics)
	switch module {
	case eventingModule:
		testFuncIsEqual(t, gotMetrics.eventing, wantMetrics.eventing)
	case remoteServiceModule:
		testFuncIsEqual(t, gotMetrics.function, wantMetrics.function)
	case fileModule:
		switch op {
		case utils.Create:
			testFuncIsEqual(t, gotMetrics.fileStore.create, wantMetrics.fileStore.create)
		case utils.Read:
			testFuncIsEqual(t, gotMetrics.fileStore.read, wantMetrics.fileStore.read)
		case utils.List:
			testFuncIsEqual(t, gotMetrics.fileStore.list, wantMetrics.fileStore.list)
		case utils.Delete:
			testFuncIsEqual(t, gotMetrics.fileStore.delete, wantMetrics.fileStore.delete)
		}
	case databaseModule:
		switch op {
		case utils.Create:
			testFuncIsEqual(t, gotMetrics.crud.create, wantMetrics.crud.create)
		case utils.Read:
			testFuncIsEqual(t, gotMetrics.crud.read, wantMetrics.crud.read)
		case utils.Update:
			testFuncIsEqual(t, gotMetrics.crud.update, wantMetrics.crud.update)
		case utils.Delete:
			testFuncIsEqual(t, gotMetrics.crud.delete, wantMetrics.crud.delete)
		}
	}
}
func TestModule_AddDBOperation(t *testing.T) {
	type args struct {
		project string
		dbAlias string
		col     string
		count   int64
		op      utils.OperationType
	}
	tests := []struct {
		name   string
		args   args
		fields *Module
		want   *sync.Map
	}{
		{
			name: "valid create case",
			args: args{
				project: "projectID",
				dbAlias: "dbAlias",
				col:     "table",
				count:   100,
				op:      utils.Create,
			},
			fields: &Module{},
			want:   testFuncGenerateDatabaseSyncMapData(generateDatabaseKey("projectID", "dbAlias", "table"), 100, 0, 0, 0),
		},
		{
			name: "valid read case",
			args: args{
				project: "projectID",
				dbAlias: "dbAlias",
				col:     "table",
				count:   100,
				op:      utils.Read,
			},
			fields: &Module{},
			want:   testFuncGenerateDatabaseSyncMapData(generateDatabaseKey("projectID", "dbAlias", "table"), 0, 100, 0, 0),
		},
		{
			name: "valid update case",
			args: args{
				project: "projectID",
				dbAlias: "dbAlias",
				col:     "table",
				count:   100,
				op:      utils.Update,
			},
			fields: &Module{},
			want:   testFuncGenerateDatabaseSyncMapData(generateDatabaseKey("projectID", "dbAlias", "table"), 0, 0, 100, 0),
		},
		{
			name: "valid delete case",
			args: args{
				project: "projectID",
				dbAlias: "dbAlias",
				col:     "table",
				count:   100,
				op:      utils.Delete,
			},
			fields: &Module{},
			want:   testFuncGenerateDatabaseSyncMapData(generateDatabaseKey("projectID", "dbAlias", "table"), 0, 0, 0, 100),
		},
		{
			name: "valid case metric disabled",
			args: args{
				project: "projectID",
				dbAlias: "local",
			},
			fields: &Module{isMetricDisabled: true},
			want:   &sync.Map{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.AddDBOperation(tt.args.project, tt.args.dbAlias, tt.args.col, tt.args.count, tt.args.op)
			tt.want.Range(func(wantKey, wantValue interface{}) bool {
				gotValue, ok := tt.fields.projects.Load(wantKey)
				if !ok {
					t.Errorf("AddDBOperation() key doesn't exist in result want %v", wantKey)
				}
				testFuncIsMetricOpEqual(t, gotValue, wantValue, tt.args.op, databaseModule)
				return false
			})
		})
	}
}

func TestModule_AddEventingType(t *testing.T) {
	type args struct {
		project      string
		eventingType string
	}
	tests := []struct {
		name   string
		args   args
		fields *Module
		want   *sync.Map
	}{
		{
			name: "valid case",
			args: args{
				project:      "projectID",
				eventingType: "type",
			},
			fields: &Module{},
			want:   testFuncGenerateEventingSyncMapData(generateEventingKey("projectID", "type")),
		},
		{
			name: "valid case metric disabled",
			args: args{
				project:      "projectID",
				eventingType: "type",
			},
			fields: &Module{isMetricDisabled: true},
			want:   &sync.Map{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.AddEventingType(tt.args.project, tt.args.eventingType)
			tt.want.Range(func(wantKey, wantValue interface{}) bool {
				gotValue, ok := tt.fields.projects.Load(wantKey)
				if !ok {
					t.Errorf("AddEventingType() key doesn't exist in result want %v", wantKey)
				}
				if !reflect.DeepEqual(gotValue, wantValue) {
					t.Errorf("AddEventingType() got value = %v %T want = %v %T", gotValue, gotValue, wantValue, wantValue)
				}
				return false
			})
		})
	}
}

func TestModule_AddFileOperation(t *testing.T) {
	type args struct {
		project   string
		storeType string
		op        utils.OperationType
	}
	tests := []struct {
		name   string
		args   args
		fields *Module
		want   *sync.Map
	}{
		{
			name: "valid create case",
			args: args{
				project:   "projectID",
				storeType: "local",
				op:        utils.Create,
			},
			fields: &Module{},
			want:   testFuncGenerateFileSyncMapData(generateFileKey("projectID", "local")),
		},
		{
			name: "valid read case",
			args: args{
				project:   "projectID",
				storeType: "local",
				op:        utils.Read,
			},
			fields: &Module{},
			want:   testFuncGenerateFileSyncMapData(generateFileKey("projectID", "local")),
		},
		{
			name: "valid list case",
			args: args{
				project:   "projectID",
				storeType: "local",
				op:        utils.List,
			},
			fields: &Module{},
			want:   testFuncGenerateFileSyncMapData(generateFileKey("projectID", "local")),
		},
		{
			name: "valid delete case",
			args: args{
				project:   "projectID",
				storeType: "local",
				op:        utils.Delete,
			},
			fields: &Module{},
			want:   testFuncGenerateFileSyncMapData(generateFileKey("projectID", "local")),
		},
		{
			name: "valid case metric disabled",
			args: args{
				project:   "projectID",
				storeType: "local",
			},
			fields: &Module{isMetricDisabled: true},
			want:   &sync.Map{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.AddFileOperation(tt.args.project, tt.args.storeType, tt.args.op)
			tt.want.Range(func(wantKey, wantValue interface{}) bool {
				gotValue, ok := tt.fields.projects.Load(wantKey)
				if !ok {
					t.Errorf("AddFileOperation() key doesn't exist in result want %v", wantKey)
				}
				testFuncIsMetricOpEqual(t, gotValue, wantValue, tt.args.op, fileModule)
				return false
			})
		})
	}
}

func TestModule_AddFunctionOperation(t *testing.T) {
	type args struct {
		project  string
		service  string
		function string
	}
	tests := []struct {
		name   string
		args   args
		fields *Module
		want   *sync.Map
	}{
		{
			name: "valid case",
			args: args{
				project:  "projectID",
				service:  "service",
				function: "function",
			},
			fields: &Module{},
			want:   testFuncGenerateFunctionSyncMapData(generateFunctionKey("projectID", "service", "function")),
		},
		{
			name: "valid case metric disabled",
			args: args{
				project: "projectID",
			},
			fields: &Module{isMetricDisabled: true},
			want:   &sync.Map{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.AddFunctionOperation(tt.args.project, tt.args.service, tt.args.function)
			tt.want.Range(func(wantKey, wantValue interface{}) bool {
				gotValue, ok := tt.fields.projects.Load(wantKey)
				if !ok {
					t.Errorf("AddFunctionOperation() key doesn't exist in result want %v", wantKey)
				}
				if !reflect.DeepEqual(gotValue, wantValue) {
					t.Errorf("AddFunctionOperation() got value = %v %T want = %v %T", gotValue, gotValue, wantValue, wantValue)
				}
				return false
			})
		})
	}
}
