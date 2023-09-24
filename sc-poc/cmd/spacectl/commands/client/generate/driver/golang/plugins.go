package golang

import (
	"fmt"
	"strings"

	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
	"golang.org/x/tools/imports"
)

func (goDriver *Golang) GeneratePlugins(plugins []v1alpha1.HTTPPlugin) (string, string, error) {
	fileName := "plugins.gen.go"
	var b strings.Builder

	// package name and imports
	pkgOut := fmt.Sprintf("package %s\n\n", goDriver.pkgName)
	_, _ = b.WriteString(pkgOut)

	s := `
type Plugins struct {}

type PluginDetails struct {
	name string
	driver string
}

func (plugin PluginDetails) Name() string {
	return plugin.name
}

func (plugin PluginDetails) Driver() string {
	return plugin.driver
}

`

	for _, plugin := range plugins {
		driverName := getTypeName(plugin.Driver, false) + getTypeName(plugin.Name, false)
		s += fmt.Sprintf("func (plugin Plugins) %s() PluginDetails {\n", driverName)
		s += fmt.Sprintf("%s := PluginDetails{\n", driverName)
		s += fmt.Sprintf("name: %q,\n", plugin.Name)
		s += fmt.Sprintf("driver: %q,\n}\n", plugin.Driver)
		s += fmt.Sprintf("return %s\n}\n\n", driverName)
	}
	_, _ = b.WriteString(s)

	outBytes, err := imports.Process(goDriver.pkgName+".go", []byte(b.String()), nil)
	if err != nil {
		return "", "", fmt.Errorf("error formatting Go code: %w", err)
	}
	return string(outBytes), fileName, nil
}
