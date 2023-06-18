package golang

import "github.com/getkin/kin-openapi/openapi3"

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
		s += "	"
	}

	return s
}
