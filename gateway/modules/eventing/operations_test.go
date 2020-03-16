package eventing

import (
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
	}{
		{
			name: "no rules with prefix 'realtime'",
			m:    &Module{config: &config.Eventing{InternalRules: map[string]config.EventingRule{"notrealtime": config.EventingRule{}}}},
			args: args{eventingRules: []config.EventingRule{config.EventingRule{}}},
		},
		{
			name: "rules with prefix 'realtime'",
			m:    &Module{config: &config.Eventing{InternalRules: map[string]config.EventingRule{"realtime-abc": config.EventingRule{}}}},
			args: args{eventingRules: []config.EventingRule{config.EventingRule{}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.SetRealtimeTriggers(tt.args.eventingRules)
		})
	}
}
