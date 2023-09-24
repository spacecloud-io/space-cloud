package typescript

import (
	"fmt"
	"strings"

	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
)

func (t *Typescript) GeneratePlugins(plugins []v1alpha1.HTTPPlugin) (string, string, error) {
	fileName := "plugins.ts"
	var b strings.Builder

	s :=
		`interface PluginDetails {
    name: string,
    driver: string
}

export class Plugins {
`

	for _, plugin := range plugins {
		driverName := getTypeName(plugin.Driver, false) + getTypeName(plugin.Name, false)
		s += fmt.Sprintf("    %s = (): PluginDetails => {\n", driverName)
		s += fmt.Sprintf("        const %s: PluginDetails = {\n", driverName)
		s += fmt.Sprintf("            name: %q,\n", plugin.Name)
		s += fmt.Sprintf("            driver: %q,\n", plugin.Driver)
		s += "        }\n"
		s += fmt.Sprintf("        return %s\n", driverName)
		s += "    }\n\n"
	}

	s += "}\n"
	_, _ = b.WriteString(s)
	return b.String(), fileName, nil
}
