package eventing

import (
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

func TestModule_SetConfig(t *testing.T) {
	type args struct {
		eventing *config.EventingConfig
	}
	tests := []struct {
		name    string
		m       *Module
		args    args
		wantErr bool
	}{
		{
			name: "eventing is not enabled",
			m:    &Module{config: &config.Eventing{Enabled: true}},
			args: args{eventing: &config.EventingConfig{Enabled: false, DBAlias: "mysql"}},
		},
		{
			name:    "DBAlias not mentioned",
			m:       &Module{config: &config.Eventing{Enabled: true}},
			args:    args{eventing: &config.EventingConfig{Enabled: true, DBAlias: ""}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.SetConfig("projectID", tt.args.eventing); (err != nil) != tt.wantErr {
				t.Errorf("Module.SetConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TODO: New function && write test case for len(schemaType["dummyDBName"][eventType]) != 0
