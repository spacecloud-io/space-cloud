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
		"{{else if and $fieldValue.IsLinked $fieldValue.IsList}}" +
		"[{{$fieldValue.Kind}}]" +
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
		"{{if and $fieldValue.Args.Precision $fieldValue.Args.Scale}}" +
		"precision: {{$fieldValue.Args.Precision}},  scale: {{$fieldValue.Args.Scale}}" +
		"{{else if $fieldValue.Args.Precision}}" +
		"precision: {{$fieldValue.Args.Precision}}" +
		"{{else if $fieldValue.Args.Scale}}" +
		"scale: {{$fieldValue.Args.Scale}}" +
		"{{end}}" +
		") " +
		"{{end}}" +

		// @primary directive
		"{{if $fieldValue.IsPrimary}}" +
		"@primary" +
		"{{end}}" +
		"{{if $fieldValue.PrimaryKeyInfo}}" +
		"{{if $fieldValue.PrimaryKeyInfo.Order}}" +
		"(order: {{$fieldValue.PrimaryKeyInfo.Order}}) " +
		"{{end}}" +
		"{{end}}" +

		// @autoIncrement directive
		"{{if $fieldValue.IsAutoIncrement}}" +
		"@autoIncrement " +
		"{{end}}" +

		// @size directive for type ID
		"{{if (or (eq $fieldValue.Kind \"Char\") (eq $fieldValue.Kind \"ID\") (eq $fieldValue.Kind \"Varchar\")) }}" +
		"{{if eq $fieldValue.TypeIDSize -1 }}" +
		"@size(value: \"max\") " +
		"{{else}}" +
		"@size(value: {{$fieldValue.TypeIDSize}}) " +
		"{{end}}" +
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

		// @link directive
		"{{if $fieldValue.IsLinked}}" +
		"{{if and $fieldValue.LinkedTable.Table $fieldValue.LinkedTable.From $fieldValue.LinkedTable.To $fieldValue.LinkedTable.To $fieldValue.LinkedTable.DBType $fieldValue.LinkedTable.Field}}" +
		"@link(table: \"{{$fieldValue.LinkedTable.Table}}\", from: \"{{$fieldValue.LinkedTable.From}}\", to: \"{{$fieldValue.LinkedTable.To}}\", db: \"{{$fieldValue.LinkedTable.DBType}}\", field: \"{{$fieldValue.LinkedTable.Field}}\" ) " +
		"{{else if and $fieldValue.LinkedTable.Table $fieldValue.LinkedTable.From $fieldValue.LinkedTable.To $fieldValue.LinkedTable.To $fieldValue.LinkedTable.DBType}}" +
		"@link(table: \"{{$fieldValue.LinkedTable.Table}}\", from: \"{{$fieldValue.LinkedTable.From}}\", to: \"{{$fieldValue.LinkedTable.To}}\", db: \"{{$fieldValue.LinkedTable.DBType}}\" ) " +
		"{{else if and $fieldValue.LinkedTable.Table $fieldValue.LinkedTable.From $fieldValue.LinkedTable.To $fieldValue.LinkedTable.To $fieldValue.LinkedTable.Field}}" +
		"@link(table: \"{{$fieldValue.LinkedTable.Table}}\", from: \"{{$fieldValue.LinkedTable.From}}\", to: \"{{$fieldValue.LinkedTable.To}}\", field: \"{{$fieldValue.LinkedTable.Field}}\" ) " +
		"{{else if and $fieldValue.LinkedTable.Table $fieldValue.LinkedTable.From $fieldValue.LinkedTable.To}}" +
		"@link(table: \"{{$fieldValue.LinkedTable.Table}}\", from: \"{{$fieldValue.LinkedTable.From}}\", to: \"{{$fieldValue.LinkedTable.To}}\" ) " +
		"{{end}}" +
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
