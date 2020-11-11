package eventing

import (
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

func TestModule_SetConfig(t *testing.T) {
	type args struct {
		eventing *config.EventingConfig
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
			name:           "eventing is not enabled",
			m:              &Module{config: &config.Eventing{Enabled: true}},
			args:           args{eventing: &config.EventingConfig{Enabled: false, DBAlias: "mysql"}},
			schemaMockArgs: []mockArgs{},
		},
		{
			name:    "DBType not mentioned",
			m:       &Module{config: &config.Eventing{Enabled: true}},
			args:    args{eventing: &config.EventingConfig{Enabled: true, DBAlias: ""}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockSchema := &mockSchemaEventingInterface{}

			for _, m := range tt.schemaMockArgs {
				mockSchema.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			tt.m.schema = mockSchema

			if err := tt.m.SetConfig("projectID", tt.args.eventing); (err != nil) != tt.wantErr {
				t.Errorf("Module.SetConfig() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockSchema.AssertExpectations(t)
		})
	}
}

// TODO: New function && write test case for len(schemaType["dummyDBName"][eventType]) != 0
