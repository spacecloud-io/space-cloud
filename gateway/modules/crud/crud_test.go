package crud

import (
	"github.com/spaceuptech/space-cloud/gateway/modules/crud/sql"
	"github.com/spaceuptech/space-cloud/gateway/modules/schema"
	"reflect"
	"sync"
	"testing"

	"github.com/graph-gophers/dataloader"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name string
		want *Module
	}{
		{
			name: "Correct value",
			want: &Module{
				batchMapTableToChan: make(batchMap),
				dataLoader:          loader{loaderMap: map[string]*dataloader.Loader{}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Init(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModule_GetDBType(t *testing.T) {
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
		dbAlias string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "DB type found",
			fields: fields{
				dbType: string(utils.MySQL),
				alias:  "db",
			},
			args:    args{dbAlias: "db"},
			want:    string(utils.MySQL),
			wantErr: false,
		},
		{
			name: "DB type found for sql- alias for previous version compatibility",
			fields: fields{
				dbType: string(utils.MySQL),
				alias:  "db",
			},
			args:    args{dbAlias: "sql-db"},
			want:    string(utils.MySQL),
			wantErr: false,
		},
		{
			name: "DB type not found",
			fields: fields{
				dbType: string(utils.MySQL),
				alias:  "db",
			},
			args:    args{dbAlias: "john"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Module{
				dbType: tt.fields.dbType,
				alias:  tt.fields.alias,
			}
			got, err := m.GetDBType(tt.args.dbAlias)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDBType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetDBType() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModule_SetSchema(t *testing.T) {
	var v model.SchemaCrudInterface = &schema.Schema{}
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
		s model.SchemaCrudInterface
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "Set schema",
			fields: fields{},
			args: args{
				s: v,
			},
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
			m.SetSchema(v)
			if !reflect.DeepEqual(m.schema, v) {
				t.Errorf("getCrudBlock() got = %v, want %v", m.schema, v)
			}
		})
	}
}

func TestModule_getCrudBlock(t *testing.T) {
	var v Crud = &sql.SQL{}
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
		dbAlias string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Crud
		wantErr bool
	}{
		{
			name:    "Error block nil",
			fields:  fields{},
			args:    args{},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Error wrong db alias",
			fields: fields{
				alias: "db1",
			},
			args: args{
				dbAlias: "db",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Error wrong db alias",
			fields: fields{
				alias: "db",
				block: v,
			},
			args: args{
				dbAlias: "db",
			},
			want:    v,
			wantErr: false,
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
			got, err := m.getCrudBlock(tt.args.dbAlias)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCrudBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCrudBlock() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_splitConnectionString(t *testing.T) {
	type args struct {
		connection string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
		want2 bool
	}{
		{
			name:  "It is a secret string",
			args:  args{connection: "secrets.name.key"},
			want:  "name",
			want1: "key",
			want2: true,
		},
		{
			name:  "It is not secret string",
			args:  args{connection: "root:1234@tcp(localhost:3306)/"},
			want:  "",
			want1: "",
			want2: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := splitConnectionString(tt.args.connection)
			if got != tt.want {
				t.Errorf("splitConnectionString() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("splitConnectionString() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("splitConnectionString() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
