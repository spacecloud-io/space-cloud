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

func (a *App) getDBWhereClause(project, db, tableName string, fieldSchemas model.FieldSchemas) *graphql.ArgumentConfig {
	whereClauseType := a.rootDBWhereTypes[project][getTableWhereClauseName(db, tableName)]

	// Add one field for each column
	for fieldName, schema := range fieldSchemas {
		// Ignore the linked field. We cannot put where clauses on that
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
		"_eq":  &graphql.InputObjectFieldConfig{Type: anyType},
		"_ne":  &graphql.InputObjectFieldConfig{Type: anyType},
		"_gt":  &graphql.InputObjectFieldConfig{Type: anyType},
		"_gte": &graphql.InputObjectFieldConfig{Type: anyType},
		"_lt":  &graphql.InputObjectFieldConfig{Type: anyType},
		"_lte": &graphql.InputObjectFieldConfig{Type: anyType},
		"_in":  &graphql.InputObjectFieldConfig{Type: graphql.NewList(anyType)},
		"_nin": &graphql.InputObjectFieldConfig{Type: graphql.NewList(anyType)},
	}

	// Specific operators for string type
	if graphqlType == graphql.String {
		whereClauseOperators["_regex"] = &graphql.InputObjectFieldConfig{Type: anyType}
		whereClauseOperators["_like"] = &graphql.InputObjectFieldConfig{Type: anyType}
	}

	// TODO: Specific operators for json type

	return graphql.NewInputObject(graphql.InputObjectConfig{
		Name:   fmt.Sprintf("DB_%sFilter", graphqlType.Name()),
		Fields: whereClauseOperators,
	})
}

func adjustWhereClause(tableName string, where map[string]interface{}, s *store, path *graphql.ResponsePath) map[string]interface{} {
	newWhereClause := make(map[string]interface{}, len(where))
	for fieldName, val := range where {
		if fieldName == "_or" {
			arr := val.([]interface{})
			for i, item := range arr {
				arr[i] = adjustWhereClause(tableName, item.(map[string]interface{}), s, path)
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
			if temp, p := s.load(v, path); p {
				v = temp
			}
			newOperatorMap["$"+k[1:]] = v
		}

		key := fieldName
		if tableName != "" {
			key = fmt.Sprintf("%s.%s", tableName, fieldName)
		}
		newWhereClause[key] = newOperatorMap
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

// func modifyTheAggregateField(tableName string, fieldAST *ast.Field, aggregate map[string][]string) bool {
// 	if fieldAST.Name.Value != "_aggregate" {
// 		return false
// 	}

// 	var operation, field string
// 	for _, arg := range fieldAST.Arguments {
// 		switch arg.Name.Value {
// 		case "field":
// 			t, _ := utils.ParseGraphqlValue(arg.Value, nil)
// 			field = t.(string)
// 		case "op":
// 			t, _ := utils.ParseGraphqlValue(arg.Value, nil)
// 			operation = t.(string)
// 		}
// 	}

// 	// Add aggregate for this operation
// 	alias := field
// 	if fieldAST.Alias != nil {
// 		alias = fieldAST.Alias.Value
// 	}
// 	aggregate[operation] = append(aggregate[operation], fmt.Sprintf("%s:%s.%s", alias, tableName, field))
// 	return true
// }

// var aggregateOperationType = graphql.NewEnum(graphql.EnumConfig{
// 	Name:        "DB_AggregateOpEnum",
// 	Description: "Enum to choose the operations allowed for aggregations",
// 	Values: graphql.EnumValueConfigMap{
// 		"sum":   &graphql.EnumValueConfig{Value: "sum"},
// 		"max":   &graphql.EnumValueConfig{Value: "max"},
// 		"min":   &graphql.EnumValueConfig{Value: "min"},
// 		"avg":   &graphql.EnumValueConfig{Value: "avg"},
// 		"count": &graphql.EnumValueConfig{Value: "count"},
// 	},
// })
