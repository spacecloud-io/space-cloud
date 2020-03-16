package eventing

import (
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

func TestModule_SetRealtimeTriggers(t *testing.T) {
	type args struct {
		eventingRules []config.EventingRule
	}
	tests := []struct {
		name string
		m    *Module
		args args
		want map[string]config.EventingRule
	}{
		{
			name: "no rules with prefix 'realtime'",
			m:    &Module{config: &config.Eventing{InternalRules: map[string]config.EventingRule{"notrealtime": config.EventingRule{}}}},
			args: args{eventingRules: []config.EventingRule{config.EventingRule{Type: "type", Options: map[string]string{"db": "db", "col": "col"}}}},
			want: map[string]config.EventingRule{"notrealtime": config.EventingRule{}, "realtime-db-col-type": config.EventingRule{Type: "type", Options: map[string]string{"db": "db", "col": "col"}}},
		},
		{
			name: "rules with prefix 'realtime'",
			m:    &Module{config: &config.Eventing{InternalRules: map[string]config.EventingRule{"realtime-abc": config.EventingRule{}}}},
			args: args{eventingRules: []config.EventingRule{config.EventingRule{}}},
			want: map[string]config.EventingRule{"realtime---": config.EventingRule{}},
		},
		{
			name: "add eventing rules when no internal rules exist",
			m:    &Module{config: &config.Eventing{InternalRules: make(map[string]config.EventingRule)}},
			args: args{eventingRules: []config.EventingRule{config.EventingRule{Type: "type", Options: map[string]string{"db": "db", "col": "col"}}, config.EventingRule{}}},
			want: map[string]config.EventingRule{"realtime-db-col-type": config.EventingRule{Type: "type", Options: map[string]string{"db": "db", "col": "col"}}, "realtime---": config.EventingRule{}},
		},
		{
			name: "add eventing rules when no realtime internal rules exist",
			m:    &Module{config: &config.Eventing{InternalRules: map[string]config.EventingRule{"notrealtime": config.EventingRule{}}}},
			args: args{eventingRules: []config.EventingRule{config.EventingRule{Type: "type", Options: map[string]string{"db": "db", "col": "col"}}, config.EventingRule{}}},
			want: map[string]config.EventingRule{"notrealtime": config.EventingRule{}, "realtime-db-col-type": config.EventingRule{Type: "type", Options: map[string]string{"db": "db", "col": "col"}}, "realtime---": config.EventingRule{}},
		},
		{
			name: "add eventing rules when realtime internal rules exist",
			m:    &Module{config: &config.Eventing{InternalRules: map[string]config.EventingRule{"realtime-abc": config.EventingRule{}, "realtime-def": config.EventingRule{}}}},
			args: args{eventingRules: []config.EventingRule{config.EventingRule{Type: "type", Options: map[string]string{"db": "db", "col": "col"}}, config.EventingRule{}}},
			want: map[string]config.EventingRule{"realtime-db-col-type": config.EventingRule{Type: "type", Options: map[string]string{"db": "db", "col": "col"}}, "realtime---": config.EventingRule{}},
		},
		{
			name: "add eventing rules when realtime and non-realtime internal rules exist",
			m:    &Module{config: &config.Eventing{InternalRules: map[string]config.EventingRule{"realtime-abc": config.EventingRule{}, "nonrealtime-def": config.EventingRule{}}}},
			args: args{eventingRules: []config.EventingRule{config.EventingRule{Type: "type", Options: map[string]string{"db": "db", "col": "col"}}, config.EventingRule{}}},
			want: map[string]config.EventingRule{"nonrealtime-def": config.EventingRule{}, "realtime-db-col-type": config.EventingRule{Type: "type", Options: map[string]string{"db": "db", "col": "col"}}, "realtime---": config.EventingRule{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.SetRealtimeTriggers(tt.args.eventingRules)

			if !reflect.DeepEqual(tt.m.config.InternalRules, tt.want) {
				t.Errorf("Error: got %v; wanted %v", tt.m.config.InternalRules, tt.want)
				return
			}
		})
	}
}
