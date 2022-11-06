package tmpl

import (
	"encoding/json"
	"strings"
	"text/template"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/getlantern/deepcopy"
)

// ExecGoTemplate compiles and executes a go template
func ExecGoTemplate(tmpl string, params interface{}) (string, error) {
	t, err := template.New("test").Funcs(CreateGoFuncMaps()).Parse(tmpl)
	if err != nil {
		return "", err
	}

	var b strings.Builder
	if err := t.Execute(&b, params); err != nil {
		return "", err
	}

	return b.String(), nil
}

// CreateGoFuncMaps creates the helper functions that can be used in go templates
func CreateGoFuncMaps() template.FuncMap {
	m := sprig.TxtFuncMap()
	m["marshalJSON"] = func(a interface{}) (string, error) {
		data, err := json.Marshal(a)
		return string(data), err
	}
	m["copy"] = func(a interface{}) (interface{}, error) {
		var b interface{}
		err := deepcopy.Copy(&b, a)
		return b, err
	}
	m["parseTimeInMillis"] = func(a interface{}) time.Time {
		var t int64
		switch v := a.(type) {
		case float32:
			t = int64(v)
		case float64:
			t = int64(v)
		case int64:
			t = v
		case int32:
			t = int64(v)
		case int:
			t = int64(v)
		}
		return time.Unix(0, t*int64(time.Millisecond))
	}

	return m
}
