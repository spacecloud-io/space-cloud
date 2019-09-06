package schema

import (
	"bytes"
	"html/template"
)

func generateSDL(schemaCol SchemaCollection) (string, error) {
	schema := `type {{range $k,$v := .}} {{$k}} { {{range $fieldName, $fieldValue := $v}}
		{{$fieldName}} :  {{if eq $fieldValue.Kind "Object"}}{{$fieldValue.TableJoin}}{{else}}{{$fieldValue.Kind}}{{end}}{{if $fieldValue.IsFieldTypeRequired}}!{{end}} {{if $fieldValue.Directive.Kind}} @{{$fieldValue.Directive.Kind}}{{end}} {{end}} {{end}}
	}`

	buf := &bytes.Buffer{}
	t := template.Must(template.New("greet").Parse(schema))
	if err := t.Execute(buf, schemaCol); err != nil {
		return "", err
	}
	return buf.String(), nil
}
