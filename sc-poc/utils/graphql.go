package utils

import "github.com/graphql-go/graphql"

// GraphqlFieldDefinitionToField converts the field definition to a field object
func GraphqlFieldDefinitionToField(fieldName string, fieldDef *graphql.FieldDefinition) *graphql.Field {
	// Prepare args
	args := make(graphql.FieldConfigArgument, len(fieldDef.Args))
	for _, arg := range fieldDef.Args {
		args[arg.Name()] = &graphql.ArgumentConfig{
			Type:         arg.Type,
			Description:  arg.Description(),
			DefaultValue: arg.DefaultValue,
		}
	}

	return &graphql.Field{
		Type:              fieldDef.Type,
		Args:              args,
		Description:       fieldDef.Description,
		DeprecationReason: fieldDef.DeprecationReason,
		Resolve:           fieldDef.Resolve,
	}
}
