package utils

import "github.com/invopop/jsonschema"

// GetJSONSchemaReflector returns a reflector to use for generating jsonschema from struct
func GetJSONSchemaReflector() *jsonschema.Reflector {
	return &jsonschema.Reflector{
		ExpandedStruct: true,
		DoNotReference: true,
	}
}
