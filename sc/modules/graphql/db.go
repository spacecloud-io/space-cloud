package graphql

import (
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/spacecloud-io/space-cloud/model"
)

func getDBWhereClause(db, tableName string, fieldSchemas model.FieldSchemas) *graphql.ArgumentConfig {
	fields := make(graphql.InputObjectConfigFieldMap)
	for fieldName, schema := range fieldSchemas {
		fieldType := dbWhereClauseLiteralFilters[scToGraphQLType(schema.Kind).String()]
		fields[fieldName] = &graphql.InputObjectFieldConfig{Type: fieldType}
	}

	return &graphql.ArgumentConfig{
		Type: graphql.NewInputObject(graphql.InputObjectConfig{
			Name:        fmt.Sprintf("%s_WhereClause", strings.Title(tableName)),
			Fields:      fields,
			Description: fmt.Sprintf("Where clause type for %s", strings.Title(tableName)),
		}),
		DefaultValue: map[string]interface{}{},
		Description:  fmt.Sprintf("Where clause argument for %s", strings.Title(tableName)),
	}
}

func getDBFilter(graphqlType graphql.Output) *graphql.InputObject {
	// Add the default operators for all types first
	whereClauseOperators := graphql.InputObjectConfigFieldMap{
		"_eq":  &graphql.InputObjectFieldConfig{Type: graphqlType},
		"_ne":  &graphql.InputObjectFieldConfig{Type: graphqlType},
		"_gt":  &graphql.InputObjectFieldConfig{Type: graphqlType},
		"_gte": &graphql.InputObjectFieldConfig{Type: graphqlType},
		"_lt":  &graphql.InputObjectFieldConfig{Type: graphqlType},
		"_lte": &graphql.InputObjectFieldConfig{Type: graphqlType},
		"_in":  &graphql.InputObjectFieldConfig{Type: graphql.NewList(graphqlType)},
		"_nin": &graphql.InputObjectFieldConfig{Type: graphql.NewList(graphqlType)},
	}

	// Specific operators for string type
	if graphqlType == graphql.String {
		whereClauseOperators["_regex"] = &graphql.InputObjectFieldConfig{Type: graphqlType}
		whereClauseOperators["_like"] = &graphql.InputObjectFieldConfig{Type: graphqlType}
	}

	// TODO: Specific operators for json type

	return graphql.NewInputObject(graphql.InputObjectConfig{
		Name:   fmt.Sprintf("DB_%sFilter", graphqlType.Name()),
		Fields: whereClauseOperators,
	})
}

func adjustWhereClause(where map[string]interface{}) map[string]interface{} {
	for fieldName, val := range where {
		// Check if value is provided directly
		operatorMap, ok := val.(map[string]interface{})
		if !ok {
			continue
		}

		// Replace the '_' with a '$'
		newOperatorMap := make(map[string]interface{}, len(operatorMap))
		for k, v := range operatorMap {
			newOperatorMap["$"+k[1:]] = v
		}

		where[fieldName] = newOperatorMap
	}

	return where
}

var dbWhereClauseLiteralFilters = map[string]*graphql.InputObject{
	graphql.String.String():  getDBFilter(graphql.String),
	graphql.Int.String():     getDBFilter(graphql.Int),
	graphql.Float.String():   getDBFilter(graphql.Float),
	graphql.Boolean.String(): getDBFilter(graphql.Boolean),
}
