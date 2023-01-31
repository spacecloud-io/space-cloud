package source

import (
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestGetModuleName(t *testing.T) {

	tests := []struct {
		name string
		gvr  schema.GroupVersionResource
		want string
	}{
		{
			name: "test module name",
			gvr:  schema.GroupVersionResource{Group: "core.space-cloud.io", Version: "v1alpha1", Resource: "compiledgraphqlsources"},
			want: "source.core---space-cloud---io--v1alpha1--compiledgraphqlsources",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetModuleName(tt.gvr)
			if got != tt.want {
				t.Errorf("GetModuleName() got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func TestGetResourceGVR(t *testing.T) {

	tests := []struct {
		name       string
		moduleName string
		want       schema.GroupVersionResource
	}{
		{
			name:       "test module name to GVR",
			moduleName: "source.core---space-cloud---io--v1alpha1--compiledgraphqlsources",
			want:       schema.GroupVersionResource{Group: "core.space-cloud.io", Version: "v1alpha1", Resource: "compiledgraphqlsources"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetResourceGVR(tt.moduleName)
			if reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetResourceGVR() got = %v, want = %v", got, tt.want)
			}
		})
	}
}
