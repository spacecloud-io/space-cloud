package rtk

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func isOperationValidForTypeGen(operation *openapi3.Operation) bool {
	if operation == nil {
		return false
	}

	if _, p := operation.Extensions["x-client-gen"]; p {
		return true
	}

	return false
}

func addPadding(depth int) string {
	s := ""
	for i := 0; i < depth; i++ {
		s += "  "
	}

	return s
}

func getTypeName(name string, skipFirst bool) string {
	arr := strings.Split(name, "-")
	for i, item := range arr {
		if i == 0 && skipFirst {
			arr[i] = item
			continue
		}

		arr[i] = strings.Title(item)
	}

	return strings.Join(arr, "")
}
