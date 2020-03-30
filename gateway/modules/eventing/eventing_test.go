package eventing

import (
	"context"
	"errors"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/stretchr/testify/mock"
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
			args: args{project: "abc", eventing: &config.Eventing{Enabled: true, DBType: "mysql", Schemas: map[string]config.SchemaObject{"eventType": config.SchemaObject{ID: "id", Schema: "schema"}}}},
			schemaMockArgs: []mockArgs{
				mockArgs{
					method:         "Parser",
					args:           []interface{}{config.Crud{"dummyDBName": &config.CrudStub{Collections: map[string]*config.TableRule{"eventType": &config.TableRule{Schema: "schema"}}}}},
					paramsReturned: []interface{}{nil, errors.New("some error")},
				},
			},
			wantErr: true,
		},
		{
			name: "eventing is not enabled",
			m:    &Module{config: &config.Eventing{Enabled: true}},
			args: args{project: "abc", eventing: &config.Eventing{Enabled: false, DBType: "mysql", Schemas: map[string]config.SchemaObject{"eventType": config.SchemaObject{ID: "id", Schema: "schema"}}}},
			schemaMockArgs: []mockArgs{
				mockArgs{
					method:         "Parser",
					args:           []interface{}{config.Crud{"dummyDBName": &config.CrudStub{Collections: map[string]*config.TableRule{"eventType": &config.TableRule{Schema: "schema"}}}}},
					paramsReturned: []interface{}{model.Type{"dummyDBName": model.Collection{"eventType": model.Fields{}}}, nil},
				},
			},
		},
		{
			name: "DBType not mentioned",
			m:    &Module{config: &config.Eventing{Enabled: true}},
			args: args{project: "abc", eventing: &config.Eventing{Enabled: true, DBType: "", Schemas: map[string]config.SchemaObject{"eventType": config.SchemaObject{ID: "id", Schema: "schema"}}}},
			schemaMockArgs: []mockArgs{
				mockArgs{
					method:         "Parser",
					args:           []interface{}{config.Crud{"dummyDBName": &config.CrudStub{Collections: map[string]*config.TableRule{"eventType": &config.TableRule{Schema: "schema"}}}}},
					paramsReturned: []interface{}{model.Type{"dummyDBName": model.Collection{"eventType": model.Fields{}}}, nil},
				},
			},
			wantErr: true,
		},
		{
			name: "config is set",
			m:    &Module{},
			args: args{project: "abc", eventing: &config.Eventing{Enabled: true, DBType: "mysql", Schemas: map[string]config.SchemaObject{"eventType": config.SchemaObject{ID: "id", Schema: "schema"}}}},
			schemaMockArgs: []mockArgs{
				mockArgs{
					method:         "Parser",
					args:           []interface{}{config.Crud{"dummyDBName": &config.CrudStub{Collections: map[string]*config.TableRule{"eventType": &config.TableRule{Schema: "schema"}}}}},
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

type mockSchemaEventingInterface struct {
	mock.Mock
}

func (m *mockSchemaEventingInterface) CheckIfEventingIsPossible(dbAlias, col string, obj map[string]interface{}, isFind bool) (findForUpdate map[string]interface{}, present bool) {
	c := m.Called(dbAlias, col, obj, isFind)
	return nil, c.Bool(1)
}

func (m *mockSchemaEventingInterface) Parser(crud config.Crud) (model.Type, error) {
	c := m.Called(crud)
	return nil, c.Error(1)
}

func (m *mockSchemaEventingInterface) SchemaValidator(col string, collectionFields model.Fields, doc map[string]interface{}) (map[string]interface{}, error) {
	c := m.Called(col, collectionFields, doc)
	return nil, c.Error(1)
}

func (m *mockSchemaEventingInterface) SchemaModifyAll(ctx context.Context, dbAlias, project string, tables map[string]*config.TableRule) error {
	c := m.Called(ctx, dbAlias, project, tables)
	return c.Error(0)
}
