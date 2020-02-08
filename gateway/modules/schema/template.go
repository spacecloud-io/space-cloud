package schema

import (
	"bytes"
	"text/template"
)

func generateSDL(schemaCol Collection) (string, error) {
	schema := `type {{range $k,$v := .}} {{$k}} { {{range $fieldName, $fieldValue := $v}}
	{{$fieldName}}: {{if eq $fieldValue.Kind "Object"}}{{$fieldValue.JointTable.Table}}{{else}}{{$fieldValue.Kind}}{{end}}{{if $fieldValue.IsFieldTypeRequired}}!{{end}} {{if $fieldValue.IsPrimary}}@primary{{end}} {{if $fieldValue.IsUnique}}@unique(group: "{{$fieldValue.IndexInfo.Group}}", order: {{$fieldValue.IndexInfo.Order}})  {{else}} {{if $fieldValue.IsIndex}}@index(group: "{{$fieldValue.IndexInfo.Group}}", sort: "{{$fieldValue.IndexInfo.Sort}}", order: {{$fieldValue.IndexInfo.Order}}){{end}}{{end}} {{if $fieldValue.IsDefault}}@default(value: {{$fieldValue.Default}}){{end}} {{if $fieldValue.IsForeign}}@foreign(table: {{$fieldValue.JointTable.Table}}, field: {{$fieldValue.JointTable.To}}){{end}}{{end}}{{end}}
}`

	buf := &bytes.Buffer{}
	t := template.Must(template.New("greet").Parse(schema))
	if err := t.Execute(buf, schemaCol); err != nil {
		return "", err
	}
	return buf.String(), nil
}
