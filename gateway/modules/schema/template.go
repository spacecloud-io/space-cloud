package schema

import (
	"bytes"
	"text/template"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

func generateSDL(schemaCol model.Collection) (string, error) {
	schema := `type {{range $k,$v := .}} {{$k}} { {{range $fieldName, $fieldValue := $v}}
	{{$fieldName}}: {{if eq $fieldValue.Kind "Object"}}{{$fieldValue.JointTable.Table}}{{else}}{{$fieldValue.Kind}}{{end}}{{if $fieldValue.IsFieldTypeRequired}}!{{end}} {{if $fieldValue.IsPrimary}}@primary{{end}} {{if $fieldValue.IsCreatedAt}}@createdAt{{end}} {{if $fieldValue.IsUpdatedAt}}@updatedAt{{end}} {{if $fieldValue.IsUnique}}@unique(group: "{{$fieldValue.IndexInfo.Group}}", order: {{$fieldValue.IndexInfo.Order}})  {{else}} {{if $fieldValue.IsIndex}}@index(group: "{{$fieldValue.IndexInfo.Group}}", sort: "{{$fieldValue.IndexInfo.Sort}}", order: {{$fieldValue.IndexInfo.Order}}){{end}}{{end}} {{if $fieldValue.IsDefault}}@default(value: {{$fieldValue.Default}}){{end}} {{if $fieldValue.IsLinked}}@link(table: {{$fieldValue.LinkedTable.Table}}, from: {{$fieldValue.LinkedTable.From}}, to: {{$fieldValue.LinkedTable.To}}, field: {{$fieldValue.LinkedTable.Field}}){{end}} {{if $fieldValue.IsForeign}}@foreign(table: {{$fieldValue.JointTable.Table}}, field: {{$fieldValue.JointTable.To}}{{if eq $fieldValue.JointTable.OnDelete "CASCADE"}} ,onDelete: cascade{{end}}){{end}}{{end}}{{end}}
}`

	buf := &bytes.Buffer{}
	t := template.Must(template.New("greet").Parse(schema))
	if err := t.Execute(buf, schemaCol); err != nil {
		return "", err
	}
	return buf.String(), nil
}
