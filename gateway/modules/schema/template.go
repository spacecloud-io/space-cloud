package schema

import (
	"bytes"
	"text/template"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

func generateSDL(schemaCol model.Collection) (string, error) {
	schema := "type {{range $k,$v := .}} {{$k}} {\n {{range $fieldName, $fieldValue := $v}}" +
		// column name
		"\t{{$fieldName}}: " +

		// column type
		"{{if eq $fieldValue.Kind \"Object\"}}" +
		"{{$fieldValue.JointTable.Table}}" +
		"{{else}}" +
		"{{$fieldValue.Kind}}" +
		"{{end}}" +

		// (!) is required
		"{{if $fieldValue.IsFieldTypeRequired}}" +
		"!" +
		"{{end}} " +

		// @args directive
		"{{if $fieldValue.Args}}" +
		"@args(" +
		"{{if $fieldValue.Args.Precision}}" +
		"precision: {{$fieldValue.Args.Precision}}," +
		"{{end}} " +
		"{{if $fieldValue.Args.Scale}}" +
		" scale: {{$fieldValue.Args.Scale}}" +
		"{{end}}" +
		")" +
		"{{end}}" +

		// @primary directive
		"{{if $fieldValue.IsPrimary}}" +
		"@primary" +
		"{{end}}" +
		"{{if $fieldValue.PrimaryKeyInfo}}" +
		"{{if and $fieldValue.PrimaryKeyInfo.IsAutoIncrement $fieldValue.PrimaryKeyInfo.Order}}" +
		"(autoIncrement: {{$fieldValue.PrimaryKeyInfo.IsAutoIncrement}},order: {{$fieldValue.PrimaryKeyInfo.Order}} ) " +
		"{{else if $fieldValue.PrimaryKeyInfo.IsAutoIncrement}}" +
		"(autoIncrement: {{$fieldValue.PrimaryKeyInfo.IsAutoIncrement}}) " +
		"{{else if $fieldValue.PrimaryKeyInfo.Order}}" +
		"(order: {{$fieldValue.PrimaryKeyInfo.Order}}) " +
		"{{end}}" +
		"{{end}}" +

		// @size directive for type ID
		"{{if eq $fieldValue.Kind \"ID\"}}" +
		"@size(value: {{$fieldValue.TypeIDSize}}) " +
		"{{end}}" +
		"{{if $fieldValue.IsCreatedAt}}" +
		"@createdAt " +
		"{{end}}" +
		"{{if $fieldValue.IsUpdatedAt}}" +
		"@updatedAt " +
		"{{end}}" +

		// @unique or @index directive
		"{{range $k,$v := $fieldValue.IndexInfo }}" +
		"{{if $v.IsUnique}}" +
		"@unique(group: \"{{$v.Group}}\", order: {{$v.Order}}) " +
		"{{else}}" +
		"{{if $v.IsIndex}}" +
		"@index(group: \"{{$v.Group}}\", sort: \"{{$v.Sort}}\", order: {{$v.Order}}) " +
		"{{end}}" +
		"{{end}}" +
		"{{end}}" +

		// @default directive
		"{{if $fieldValue.IsDefault}}" +
		"@default(value: {{$fieldValue.Default}}) " +
		"{{end}}" +
		"{{if $fieldValue.IsLinked}}" +
		"@link(table: {{$fieldValue.LinkedTable.Table}}, from: {{$fieldValue.LinkedTable.From}}, to: {{$fieldValue.LinkedTable.To}}, field: {{$fieldValue.LinkedTable.Field}}) " +
		"{{end}}" +

		// @foreign directive
		"{{if $fieldValue.IsForeign}}" +
		"@foreign(table: {{$fieldValue.JointTable.Table}}, field: {{$fieldValue.JointTable.To}}" +
		"{{if eq $fieldValue.JointTable.OnDelete \"CASCADE\"}}" +
		",onDelete: cascade" +
		"{{end}}" +
		") " +
		"{{end}}" +
		"\n" +
		"{{end}}" +
		"{{end}}" +
		"}"

	buf := &bytes.Buffer{}
	t := template.Must(template.New("greet").Parse(schema))
	if err := t.Execute(buf, schemaCol); err != nil {
		return "", err
	}
	return buf.String(), nil
}
