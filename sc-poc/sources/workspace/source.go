package workspace

import (
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/caddyserver/caddy/v2"
	"github.com/spacecloud-io/space-cloud/managers/source"
	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
)

var workspaceResource = schema.GroupVersionResource{Group: "core.space-cloud.io", Version: "v1alpha1", Resource: "workspaces"}

func init() {
	source.RegisterSource(WorkspaceSource{}, workspaceResource)
}

// WorkspaceSource describes a Workspace source
type WorkspaceSource struct {
	v1alpha1.Workspace
}

// CaddyModule returns the Caddy module information.
func (WorkspaceSource) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  caddy.ModuleID(source.GetModuleName(workspaceResource)),
		New: func() caddy.Module { return new(WorkspaceSource) },
	}
}
