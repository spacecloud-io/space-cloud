package graphql

import (
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/spacecloud-io/space-cloud/model"
)

func getDBSort(db, tableName string, fieldSchemas model.FieldSchemas) *graphql.ArgumentConfig {
	fields := make(graphql.InputObjectConfigFieldMap)
	for fieldName, schema := range fieldSchemas {
		if schema.IsLinked {
			continue
		}

		fields[fieldName] = &graphql.InputObjectFieldConfig{Type: enumSort}
	}

	return &graphql.ArgumentConfig{
		Type: graphql.NewInputObject(graphql.InputObjectConfig{
			Name:        fmt.Sprintf("%s_%s_Sort", strings.Title(db), strings.Title(tableName)),
			Description: fmt.Sprintf("Sort type for %s", strings.Title(tableName)),
			Fields:      fields,
		}),
		Description:  fmt.Sprintf("Sort argument for %s", strings.Title(tableName)),
		DefaultValue: map[string]interface{}{},
	}
}

// adjustSortArgument converts the sort map to an array the database driver understands
func adjustSortArgument(sort map[string]interface{}) []string {
	arr := make([]string, 0, len(sort))
	for fieldName, sortType := range sort {
		val := fieldName
		if sortType == "desc" {
			val = "-" + fieldName
		}
		arr = append(arr, val)
	}

	return arr
}

func getDBWhereClause(db, tableName string, fieldSchemas model.FieldSchemas) *graphql.ArgumentConfig {
	whereClauseType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name:        fmt.Sprintf("%s_%s_WhereClause", strings.Title(db), strings.Title(tableName)),
		Fields:      make(graphql.InputObjectConfigFieldMap),
		Description: fmt.Sprintf("Where clause type for %s", strings.Title(tableName)),
	})

	// Add one field for each column
	for fieldName, schema := range fieldSchemas {
		// Ignore the linked field. We cannot put where clauses on that
		// TODO: Allow where clause to be in the argument of linked field collection set
		if schema.IsLinked {
			continue
		}
		fieldType := dbWhereClauseLiteralFilters[scToGraphQLType(schema.Kind).String()]
		whereClauseType.AddFieldConfig(fieldName, &graphql.InputObjectFieldConfig{Type: fieldType})
	}

	// Add field for or clause
	whereClauseType.AddFieldConfig("_or", &graphql.InputObjectFieldConfig{Type: graphql.NewList(whereClauseType)})

	return &graphql.ArgumentConfig{
		Type:         whereClauseType,
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
	newWhereClause := make(map[string]interface{}, len(where))
	for fieldName, val := range where {
		if fieldName == "_or" {
			arr := val.([]interface{})
			for i, item := range arr {
				arr[i] = adjustWhereClause(item.(map[string]interface{}))
			}
			newWhereClause["$or"] = arr
			continue
		}

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

		newWhereClause[fieldName] = newOperatorMap
	}

	return newWhereClause
}

func extractIntegerFromArg(key string, args map[string]interface{}) *int64 {
	val, p := args[key]
	if !p {
		return nil
	}

	num, ok := val.(int)
	if !ok {
		return nil
	}

	v := int64(num)
	return &v
}

var dbWhereClauseLiteralFilters = map[string]*graphql.InputObject{
	graphql.String.String():  getDBFilter(graphql.String),
	graphql.Int.String():     getDBFilter(graphql.Int),
	graphql.Float.String():   getDBFilter(graphql.Float),
	graphql.Boolean.String(): getDBFilter(graphql.Boolean),
}

var enumSort = graphql.NewEnum(graphql.EnumConfig{
	Name:        "DB_SortEnum",
	Description: "Enum to choose between ascending and descending sort",
	Values: graphql.EnumValueConfigMap{
		"asc":  &graphql.EnumValueConfig{Value: "asc"},
		"desc": &graphql.EnumValueConfig{Value: "desc"},
	},
})
