package crud

import (
	"errors"
	"reflect"
	"sync"
	"testing"

	"github.com/graph-gophers/dataloader"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func TestModule_getLoader(t *testing.T) {
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
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *dataloader.Loader
		want1  bool
	}{
		{
			name: "Get Loader For Specified key",
			fields: fields{
				dataLoader: loader{loaderMap: map[string]*dataloader.Loader{}},
			},
			args: args{
				key: "some-key",
			},
			want:  nil,
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Module{
				RWMutex:             tt.fields.RWMutex,
				block:               tt.fields.block,
				dbType:              tt.fields.dbType,
				alias:               tt.fields.alias,
				project:             tt.fields.project,
				schema:              tt.fields.schema,
				queries:             tt.fields.queries,
				batchMapTableToChan: tt.fields.batchMapTableToChan,
				dataLoader:          tt.fields.dataLoader,
				hooks:               tt.fields.hooks,
				metricHook:          tt.fields.metricHook,
				getSecrets:          tt.fields.getSecrets,
			}
			tt.want = m.createLoader(tt.args.key)
			got, got1 := m.getLoader(tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getLoader() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getLoader() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_resultsHolder_addResult(t *testing.T) {
	type fields struct {
		Mutex        sync.Mutex
		results      []*dataloader.Result
		whereClauses []interface{}
	}
	type args struct {
		i      int
		result *dataloader.Result
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Correct values",
			fields: fields{
				results: make([]*dataloader.Result, 1),
			},
			args: args{
				i: 0,
				result: &dataloader.Result{
					Data: map[string]interface{}{
						"id":   "1",
						"name": "John",
					},
					Error: nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			holder := &resultsHolder{
				results:      tt.fields.results,
				whereClauses: tt.fields.whereClauses,
			}
			holder.addResult(tt.args.i, tt.args.result)
			if holder.results[tt.args.i] != tt.args.result {
				t.Errorf("addResult() got %v want %v", holder.results[tt.args.i], tt.args.result)
			}
		})
	}
}

func Test_resultsHolder_addWhereClause(t *testing.T) {
	type fields struct {
		Mutex        sync.Mutex
		results      []*dataloader.Result
		whereClauses []interface{}
	}
	type args struct {
		whereClause map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Correct value",
			fields: fields{
				whereClauses: make([]interface{}, 1),
			},
			args: args{
				whereClause: map[string]interface{}{
					"id": map[string]interface{}{
						"$eq": "1",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			holder := &resultsHolder{whereClauses: tt.fields.whereClauses}
			holder.addWhereClause(tt.args.whereClause)
			if reflect.DeepEqual(tt.args.whereClause, []interface{}{tt.fields.whereClauses}) {
				t.Errorf("addResult() got %v want %v", tt.args.whereClause, []interface{}{tt.fields.whereClauses})
			}
		})
	}
}

func Test_resultsHolder_fillErrorMessage(t *testing.T) {
	type fields struct {
		Mutex        sync.Mutex
		results      []*dataloader.Result
		whereClauses []interface{}
	}
	type args struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Correct value",
			fields: fields{
				results: []*dataloader.Result{nil, nil},
			},
			args: args{err: errors.New("some database error")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			holder := &resultsHolder{results: tt.fields.results}
			holder.fillErrorMessage(tt.args.err)
			for _, result := range tt.fields.results {
				if result.Error != tt.args.err {
					t.Errorf("fillErrorMessage() got %v want %v", result.Error, tt.args.err)
				}
			}
		})
	}
}

func Test_resultsHolder_fillResults(t *testing.T) {
	type fields struct {
		Mutex        sync.Mutex
		results      []*dataloader.Result
		whereClauses []interface{}
	}
	type args struct {
		res []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*dataloader.Result
	}{
		{
			name: "Result already has a value",
			fields: fields{
				results: make([]*dataloader.Result, 2),
				whereClauses: []interface{}{
					map[string]interface{}{
						"id": map[string]interface{}{
							"$eq": "1",
						},
					},
					map[string]interface{}{
						"id": map[string]interface{}{
							"$eq": "2",
						},
					},
				},
			},
			args: args{
				res: []interface{}{
					map[string]interface{}{
						"id":   "1",
						"name": "John",
					},
					map[string]interface{}{
						"id":   "2",
						"name": "Sam",
					},
				},
			},
			want: []*dataloader.Result{
				{
					Data: map[string]interface{}{
						"id":   "1",
						"name": "John",
					},
					Error: nil,
				},
				{
					Data: map[string]interface{}{
						"id":   "2",
						"name": "Sam",
					},
					Error: nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			holder := &resultsHolder{
				results:      tt.fields.results,
				whereClauses: tt.fields.whereClauses,
			}
			holder.fillResults(tt.args.res)
			if len(tt.args.res) != len(holder.results) {
				t.Errorf("fillResults() length got %v want %v", len(tt.args.res), len(holder.results))
			}
			for i, result := range holder.results {
				if reflect.DeepEqual(tt.args.res[i], result.Data) {
					t.Errorf("fillResults() got %v want %v", tt.args.res[i], result.Data)
				}
			}
		})
	}
}

func Test_resultsHolder_getResults(t *testing.T) {
	type fields struct {
		Mutex        sync.Mutex
		results      []*dataloader.Result
		whereClauses []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   []*dataloader.Result
	}{
		{
			name: "Correct value",
			fields: fields{
				results: []*dataloader.Result{
					{
						Data: map[string]interface{}{
							"id":   "1",
							"name": "John",
						},
						Error: nil,
					},
					{
						Data: map[string]interface{}{
							"id":   "2",
							"name": "Sam",
						},
						Error: nil,
					},
				},
			},
			want: []*dataloader.Result{
				{
					Data: map[string]interface{}{
						"id":   "1",
						"name": "John",
					},
					Error: nil,
				},
				{
					Data: map[string]interface{}{
						"id":   "2",
						"name": "Sam",
					},
					Error: nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			holder := &resultsHolder{
				results: tt.fields.results,
			}
			if got := holder.getResults(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getResults() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_resultsHolder_getWhereClauses(t *testing.T) {
	type fields struct {
		Mutex        sync.Mutex
		results      []*dataloader.Result
		whereClauses []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   []interface{}
	}{
		{
			name: "Correct value",
			fields: fields{
				whereClauses: []interface{}{
					map[string]interface{}{
						"id": map[string]interface{}{
							"$eq": "1",
						},
					},
					map[string]interface{}{
						"id": map[string]interface{}{
							"$eq": "2",
						},
					},
				},
			},
			want: []interface{}{
				map[string]interface{}{
					"id": map[string]interface{}{
						"$eq": "1",
					},
				},
				map[string]interface{}{
					"id": map[string]interface{}{
						"$eq": "2",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			holder := &resultsHolder{
				whereClauses: tt.fields.whereClauses,
			}
			if got := holder.getWhereClauses(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getWhereClauses() = %v, want %v", got, tt.want)
			}
		})
	}
}
