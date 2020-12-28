package schema

import (
	"bytes"
	"text/template"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

func generateSDL(schemaCol model.Collection) (string, error) {
	schema := "type {{range $k,$v := .}} {{$k}} {\n {{range $fieldName, $fieldValue := $v}}" +
		// column name
		"\t{{$fieldName}}:" +

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

		// @unique directive
		"{{if $fieldValue.IsUnique}}" +
		"@unique(group: \"{{$fieldValue.IndexInfo.Group}}\", order: {{$fieldValue.IndexInfo.Order}}) " +
		"{{else}}" +
		"{{if $fieldValue.IsIndex}}" +
		"@index(group: \"{{$fieldValue.IndexInfo.Group}}\", sort: \"{{$fieldValue.IndexInfo.Sort}}\", order: {{$fieldValue.IndexInfo.Order}}) " +
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
