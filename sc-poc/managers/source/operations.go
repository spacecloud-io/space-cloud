package source

import (
	"fmt"
	"strings"

	"github.com/caddyserver/caddy/v2"
)

// GetSources returns an array of sources applicable for that provider
func (a *App) GetSources(provider string) []Source {
	return a.sourceMap[provider]
}

// GetModuleName returns a caddy compatible module name
func GetModuleName(apiVersion, kind string) string {
	moduleName := fmt.Sprintf("%s--%s", apiVersion, kind)

	// Replace the periods with `---`
	moduleName = strings.Join(strings.Split(moduleName, "."), "---")

	// Replace the `/` with `----`
	moduleName = strings.Join(strings.Split(moduleName, "/"), "----")

	return strings.ToLower(fmt.Sprintf("source.%s", moduleName))
}

// GetResourceGVK returns the api version and kind of the resource
func GetResourceGVK(moduleName string) (apiVersion, kind string) {
	moduleName = strings.Join(strings.Split(moduleName, "----"), "/")
	moduleName = strings.Join(strings.Split(moduleName, "---"), ".")
	arr := strings.Split(moduleName, "--")
	return arr[0], arr[1]
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
