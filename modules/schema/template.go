package schema

import (
	"bytes"
	"html/template"
)

func generateSDL(schemaCol schemaCollection) (string, error) {
	schema := `type {{range $k,$v := .}} {{$k}} { {{range $fieldName, $fieldValue := $v}}
		{{$fieldName}} :  {{if eq $fieldValue.Kind "Object"}}{{$fieldValue.JointTable.TableName}}{{else}}{{$fieldValue.Kind}}{{end}}{{if $fieldValue.IsFieldTypeRequired}}!{{end}} {{if $fieldValue.Directive}} @{{$fieldValue.Directive}}{{if eq $fieldValue.Directive "relation"}}(field : {{$fieldValue.JointTable.TableField}}) {{end}} {{end}} {{end}} {{end}}
	}`

	buf := &bytes.Buffer{}
	t := template.Must(template.New("greet").Parse(schema))
	if err := t.Execute(buf, schemaCol); err != nil {
		return "", err
	}
	return buf.String(), nil
}
