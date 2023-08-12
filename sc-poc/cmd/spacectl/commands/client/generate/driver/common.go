package driver

import "github.com/getkin/kin-openapi/openapi3"

func IsOperationValidForTypeGen(operation *openapi3.Operation) bool {
	if operation == nil {
		return false
	}

	if _, p := operation.Extensions["x-client-gen"]; p {
		return true
	}

	return false
}

func AddPadding(depth int) string {
	s := ""
	for i := 0; i < depth; i++ {
		s += "	"
	}

	return s
}
