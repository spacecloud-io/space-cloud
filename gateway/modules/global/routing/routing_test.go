package routing

import (
	"reflect"
	"sync"
	"testing"
	"text/template"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *Routing
	}{
		// TODO: Add test cases.
		{
			name: "New Routing instance",
			want: &Routing{
				lock:        sync.RWMutex{},
				routes:      routeMapping{},
				goTemplates: map[string]*template.Template{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
