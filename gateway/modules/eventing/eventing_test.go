package eventing

import (
	"errors"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

func TestModule_SetConfig(t *testing.T) {
	type args struct {
		project  string
		eventing *config.Eventing
	}
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	tests := []struct {
		name           string
		m              *Module
		args           args
		schemaMockArgs []mockArgs
		wantErr        bool
	}{
		{
			name: "unable to parse schema",
			m:    &Module{},
			args: args{project: "abc", eventing: &config.Eventing{Enabled: true, DBAlias: "mysql", Schemas: map[string]config.SchemaObject{"eventType": {ID: "id", Schema: "schema"}}}},
			schemaMockArgs: []mockArgs{
				{
					method:         "Parser",
					args:           []interface{}{config.Crud{"dummyDBName": &config.CrudStub{Collections: map[string]*config.TableRule{"eventType": {Schema: "schema"}}}}},
					paramsReturned: []interface{}{nil, errors.New("some error")},
				},
			},
			wantErr: true,
		},
		{
			name:           "eventing is not enabled",
			m:              &Module{config: &config.Eventing{Enabled: true}},
			args:           args{project: "abc", eventing: &config.Eventing{Enabled: false, DBAlias: "mysql", Schemas: map[string]config.SchemaObject{"eventType": {ID: "id", Schema: "schema"}}}},
			schemaMockArgs: []mockArgs{},
		},
		{
			name: "DBType not mentioned",
			m:    &Module{config: &config.Eventing{Enabled: true}},
			args: args{project: "abc", eventing: &config.Eventing{Enabled: true, DBAlias: "", Schemas: map[string]config.SchemaObject{"eventType": {ID: "id", Schema: "schema"}}}},
			schemaMockArgs: []mockArgs{
				{
					method:         "Parser",
					args:           []interface{}{config.Crud{"dummyDBName": &config.CrudStub{Collections: map[string]*config.TableRule{"eventType": {Schema: "schema"}}}}},
					paramsReturned: []interface{}{model.Type{"dummyDBName": model.Collection{"eventType": model.Fields{}}}, nil},
				},
			},
			wantErr: true,
		},
		{
			name: "config is set",
			m:    &Module{config: &config.Eventing{}},
			args: args{project: "abc", eventing: &config.Eventing{Enabled: true, DBAlias: "mysql", Schemas: map[string]config.SchemaObject{"eventType": {ID: "id", Schema: "schema"}}}},
			schemaMockArgs: []mockArgs{
				{
					method:         "Parser",
					args:           []interface{}{config.Crud{"dummyDBName": &config.CrudStub{Collections: map[string]*config.TableRule{"eventType": {Schema: "schema"}}}}},
					paramsReturned: []interface{}{model.Type{"dummyDBName": model.Collection{"eventType": model.Fields{}}}, nil},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockSchema := &mockSchemaEventingInterface{}

			for _, m := range tt.schemaMockArgs {
				mockSchema.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			tt.m.schema = mockSchema

			if err := tt.m.SetConfig(tt.args.project, tt.args.eventing); (err != nil) != tt.wantErr {
				t.Errorf("Module.SetConfig() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockSchema.AssertExpectations(t)
		})
	}
}

// TODO: New function && write test case for len(schemaType["dummyDBName"][eventType]) != 0
