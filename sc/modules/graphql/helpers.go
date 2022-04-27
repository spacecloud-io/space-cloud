package graphql

import (
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/spacecloud-io/space-cloud/model"
	"github.com/spacecloud-io/space-cloud/utils"
)

func (a *App) getQueryType(project string) *graphql.Object {
	// Create the root query
	queryType := graphql.NewObject(graphql.ObjectConfig{Name: "Query", Fields: graphql.Fields{}})

	// Create the join object
	joinRetObj := map[string]string{}
	// Populate all the fields in the root query
	for dbAlias, parsedSchema := range a.dbSchemas[project] {
		for k, v := range a.getTableFields(project, dbAlias, parsedSchema, queryType, joinRetObj) {
			queryType.AddFieldConfig(k, v)
		}
	}

	return queryType
}

func (a *App) getTableFields(project, db string, schemas model.CollectionSchemas, queryType *graphql.Object, joinRetObj map[string]string) graphql.Fields {
	fields := make(graphql.Fields, len(schemas))

	for tableName, tableSchema := range schemas {
		tableFields := make(graphql.Fields, len(tableSchema))

		for fieldName, fieldSchema := range tableSchema {
			tableFields[fieldName] = &graphql.Field{
				Type:    scToGraphQLType(fieldSchema.Kind),
				Resolve: a.literalResolveFn,
			}
		}

		// Create a record object for the table
		graphqlObject := graphql.NewObject(graphql.ObjectConfig{
			Name:        fmt.Sprintf("%s_%s", strings.Title(db), strings.Title(tableName)),
			Description: fmt.Sprintf("Record object from %s", tableName),
			Fields:      tableFields,
		})
		graphqlArguments := graphql.FieldConfigArgument{
			"where": getDBWhereClause(db, tableName, tableSchema),
			"sort":  getDBSort(db, tableName, tableSchema),
			"limit": &graphql.ArgumentConfig{Type: graphql.Int},
			"skip":  &graphql.ArgumentConfig{Type: graphql.Int},
		}

		// Create the queries that can be performed to read from table
		fieldKeyMultiple := fmt.Sprintf("%s_findMultiple%s", db, strings.Title(tableName))
		fieldKeySingle := fmt.Sprintf("%s_findOne%s", db, strings.Title(tableName))
		fields[fieldKeyMultiple] = &graphql.Field{
			Type:    graphql.NewList(graphqlObject),
			Args:    graphqlArguments,
			Resolve: a.dbReadResolveFn(project, db, tableName, utils.All),
		}
		fields[fieldKeySingle] = &graphql.Field{
			Type:    graphqlObject,
			Args:    graphqlArguments,
			Resolve: a.dbReadResolveFn(project, db, tableName, utils.One),
		}

		// Add to join return objects
		joinRetObj[fieldKeyMultiple] = ""
		joinRetObj[fieldKeySingle] = ""
		a.registeredQueries[fieldKeyMultiple] = struct{}{}
		a.registeredQueries[fieldKeySingle] = struct{}{}

		tableFields["_join"] = &graphql.Field{
			Type: queryType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return func() (interface{}, error) {
					return joinRetObj, nil
				}, nil
			},
		}
	}

	return fields
}

func scToGraphQLType(kind string) graphql.Output {
	switch kind {
	case model.TypeDateTimeWithZone, model.TypeDateTime, model.TypeDate, model.TypeTime, model.TypeUUID, model.TypeChar, model.TypeVarChar, model.TypeString, model.TypeID, model.TypeEnum:
		return graphql.String
	case model.TypeInteger, model.TypeSmallInteger, model.TypeBigInteger:
		return graphql.Int
	case model.TypeDecimal, model.TypeFloat:
		return graphql.Float
	case model.TypeBoolean:
		return graphql.Boolean

	// TODO: Add a case for JSON types
	default:
		// TODO: Return an `any` scaler type
		return graphql.String

	}
}

func containsExportDirective(field *ast.Field) (string, bool) {
	if len(field.Directives) == 0 {
		return "", false
	}

	for _, dir := range field.Directives {
		if dir.Name.Value == "export" {
			key := field.Name.Value
			if field.Alias != nil {
				key = field.Alias.Value
			}

			for _, arg := range dir.Arguments {
				if arg.Name.Value == "as" {
					t, _ := utils.ParseGraphqlValue(arg.Value, map[string]interface{}{})
					key = t.(string)
					break
				}
			}
			return key, true
		}
	}

	return "", false
}
