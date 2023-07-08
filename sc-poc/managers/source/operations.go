package source

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/caddyserver/caddy/v2"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	sourcesGVR  []schema.GroupVersionResource
	sourcesLock sync.Mutex
)

// RegisterSource takes a source and registers it with caddy.
// Additionally, it populates the global GVR map for configman
// to provision its configurations
func RegisterSource(src caddy.Module, gvr schema.GroupVersionResource) {
	sourcesLock.Lock()
	defer sourcesLock.Unlock()
	caddy.RegisterModule(src)
	sourcesGVR = append(sourcesGVR, gvr)
}

// GetRegisteredSources returns the global variable sourcesGVR which stores
// a slice of all registered sources' GVR
func GetRegisteredSources() []schema.GroupVersionResource {
	return sourcesGVR
}

// GetWorkspaces returns a list of workspaces registered with this SC instance
func (a *App) GetWorkspaces() []string {
	return a.workspaces
}

// GetSources returns an array of sources applicable for that provider
func (a *App) GetSources(workspace, provider string) map[string]Source {
	sources := make(map[string]Source)

	// First load all the sources from the main workspace
	if workspace, p := a.sourceMap["main"]; p {
		for _, source := range workspace[provider] {
			sources[getUniqueSourceName(source)] = source
		}
	}

	// Now superimpose the sources from the requested workspace
	if workspace, p := a.sourceMap[workspace]; p {
		for _, source := range workspace[provider] {
			sources[getUniqueSourceName(source)] = source
		}
	}

	// TODO: Add support for marking certain sources as deleted so we can remove it from the requested workspace.
	// This can be done by making a new source called delete marker whose sole purpose is to mark other sources
	// as deleted

	return sources
}

// GetModuleName returns a caddy compatible module name
func GetModuleName(gvr schema.GroupVersionResource) string {
	moduleName := fmt.Sprintf("%s--%s--%s", gvr.Group, gvr.Version, gvr.Resource)

	// Replace the periods with `---`
	moduleName = strings.Join(strings.Split(moduleName, "."), "---")

	// Replace the `/` with `----`
	moduleName = strings.Join(strings.Split(moduleName, "/"), "----")

	return strings.ToLower(fmt.Sprintf("source.%s", moduleName))
}

// GetResourceGVR returns the api version and kind of the resource
func GetResourceGVR(moduleName string) schema.GroupVersionResource {
	moduleName = strings.TrimPrefix(moduleName, "source.")
	moduleName = strings.Join(strings.Split(moduleName, "----"), "/")
	moduleName = strings.Join(strings.Split(moduleName, "---"), ".")
	arr := strings.Split(moduleName, "--")
	return schema.GroupVersionResource{Group: arr[0], Version: arr[1], Resource: arr[2]}
}

// TODO: Remove this
// // ResolveDependencies helpers providers resolve a sources dependency
// func ResolveDependencies(ctx caddy.Context, callerAppName string, source Source) error {
// 	providers := source.GetProviders()

// 	// Resolve all providers except the caller
// 	for _, provider := range providers {
// 		if provider == callerAppName {
// 			return nil
// 		}

// 		if _, err := ctx.App(provider); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// GetWorkspaceNameFromSource gets the workspace name from a source
func GetWorkspaceNameFromSource(s Source) string {
	labels := s.GetLabels()
	if labels == nil {
		labels = map[string]string{}
	}
	workspace := labels[workspaceLabel]
	if workspace == "" || workspace == "root" {
		workspace = "main"
	}

	return workspace
}

// GetWorkspaceNameFromHeaders gets the workspace name from a headers
func GetWorkspaceNameFromHeaders(r *http.Request) string {
	ws := r.Header.Get("x-sc-workspace")
	if ws == "" {
		ws = "main"
	}
	return ws
}

func getUniqueSourceName(s Source) string {
	return fmt.Sprintf("%s-%s", s.GroupVersionKind().String(), s.GetName())
}
