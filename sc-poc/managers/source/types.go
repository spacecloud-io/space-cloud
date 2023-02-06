package source

import (
	"sort"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

type (
	// Source describes a source module
	Source interface {
		GetPriority() int
		GetName() string
		GroupVersionKind() schema.GroupVersionKind
		GetProviders() []string
	}

	// Sources is an array of Source
	Sources []Source
)

// Sort sorts the array of sources based on their priority
func (s Sources) Sort() {
	sort.SliceStable(s, func(i, j int) bool {
		return s[i].GetPriority() > s[j].GetPriority()
	})
}
