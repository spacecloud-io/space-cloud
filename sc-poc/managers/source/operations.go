package source

import (
	"fmt"
	"strings"
	"sync"

	"github.com/caddyserver/caddy/v2"
	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
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

// GetSources returns an array of sources applicable for that provider
func (a *App) GetSources(provider string) []Source {
	return a.sourceMap[provider]
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

// ResolveDependencies helpers providers resolve a sources dependency
func ResolveDependencies(ctx caddy.Context, callerAppName string, source Source) error {
	providers := source.GetProviders()

	// Resolve all providers except the caller
	for _, provider := range providers {
		if provider == callerAppName {
			return nil
		}

		if _, err := ctx.App(provider); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) GetPlugins() []v1alpha1.HTTPPlugin {
	return a.plugins
}
