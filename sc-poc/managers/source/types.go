package source

import (
	"sort"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	workspaceLabel string = "space-cloud.io/workspace"
)

type (
	// Workspace describes a workspace module
	Workspace interface {
		GetWorkspaceName() string
	}

	// Source describes a source module
	Source interface {
		GetPriority() int
		GetName() string
		GroupVersionKind() schema.GroupVersionKind
		GetProviders() []string
		GetLabels() map[string]string
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
