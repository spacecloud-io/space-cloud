package source

import (
	"testing"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

type mockSource struct {
	name     string
	priority int
}

func (s mockSource) GetName() string {
	return s.name
}

func (s mockSource) GetPriority() int {
	return s.priority
}

func (s mockSource) GetProviders() []string {
	return nil
}

func (s mockSource) GroupVersionKind() schema.GroupVersionKind {
	return schema.GroupVersionKind{}
}

func TestSources_Sort(t *testing.T) {

	tests := []struct {
		name string
		s    Sources
		want []string
	}{
		{
			name: "sorting test",
			s: Sources{
				mockSource{name: "app1", priority: 70},
				mockSource{name: "app2", priority: 50},
				mockSource{name: "app3", priority: 100},
			},
			want: []string{"app3", "app1", "app2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Sort()
			for i, s := range tt.s {
				if s.GetName() != tt.want[i] {
					t.Errorf("Invalid sort order: wanted - %v; got - %v", tt.want, tt.s)
					return
				}
			}
		})
	}
}
