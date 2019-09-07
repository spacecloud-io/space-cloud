package schema

import (
	"bytes"
	"html/template"
)

func generateSDL(schemaCol schemaCollection) (string, error) {
	schema := `type {{range $k,$v := .}} {{$k}} { {{range $fieldName, $fieldValue := $v}}
		{{$fieldName}} :  {{if eq $fieldValue.Kind "Object"}}{{$fieldValue.tableJoin}}{{else}}{{$fieldValue.Kind}}{{end}}{{if $fieldValue.isFieldTypeRequired}}!{{end}} {{if $fieldValue.directive.Kind}} @{{$fieldValue.directive.Kind}}{{end}} {{end}} {{end}}
	}`

	buf := &bytes.Buffer{}
	t := template.Must(template.New("greet").Parse(schema))
	if err := t.Execute(buf, schemaCol); err != nil {
		return "", err
	}
	return buf.String(), nil
}
