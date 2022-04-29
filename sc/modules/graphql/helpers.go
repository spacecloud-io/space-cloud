package graphql

import (
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/spacecloud-io/space-cloud/model"
	"github.com/spacecloud-io/space-cloud/utils"
)

func (a *App) createRootGraphQLTypes(project string) {
	a.rootGraphQLTypes[project] = map[string]*graphql.Object{}
	a.rootDBWhereTypes[project] = map[string]*graphql.InputObject{}

	for dbAlias, parsedSchema := range a.dbSchemas[project] {
		for tableName := range parsedSchema {
			a.rootGraphQLTypes[project][getTableFieldName(dbAlias, tableName)] = graphql.NewObject(graphql.ObjectConfig{
				Name:        getTableFieldName(dbAlias, tableName),
				Description: fmt.Sprintf("Record object from %s", tableName),
				Fields:      graphql.Fields{},
			})
			a.rootDBWhereTypes[project][getTableWhereClauseName(dbAlias, tableName)] = graphql.NewInputObject(graphql.InputObjectConfig{
				Name:        getTableWhereClauseName(dbAlias, tableName),
				Fields:      make(graphql.InputObjectConfigFieldMap),
				Description: fmt.Sprintf("Where clause type for %s", strings.Title(tableName)),
			})
		}
	}
}

func (a *App) getQueryType(project string) *graphql.Object {
	// Create root graphql types
	a.createRootGraphQLTypes(project)

	// Create the root query
	queryType := graphql.NewObject(graphql.ObjectConfig{Name: "Query", Fields: graphql.Fields{}})

	// Create the join object
	joinRetObj := map[string]string{}
	// Populate all the fields in the root query
	for dbAlias, parsedSchema := range a.dbSchemas[project] {
		for k, v := range a.getDatabaseFields(project, dbAlias, parsedSchema, queryType, joinRetObj) {
			queryType.AddFieldConfig(k, v)
		}
	}

	return queryType
}

func (a *App) getDatabaseFields(project, db string, schemas model.CollectionSchemas, queryType *graphql.Object, joinRetObj map[string]string) graphql.Fields {
	fields := make(graphql.Fields, len(schemas))

	for tableName, tableSchema := range schemas {
		// Add the fields for this table
		for fieldName, fieldSchema := range tableSchema {
			if !fieldSchema.IsLinked {
				a.rootGraphQLTypes[project][getTableFieldName(db, tableName)].AddFieldConfig(fieldName, &graphql.Field{
					Type:    scToGraphQLType(fieldSchema.Kind),
					Resolve: a.literalResolveFn,
				})
				continue
			}

			if fieldSchema.IsLinked {
				var graphqlType graphql.Output = a.rootGraphQLTypes[project][getTableFieldName(db, fieldSchema.LinkedTable.Table)]
				if fieldSchema.IsList {
					graphqlType = graphql.NewList(graphqlType)
				}
				a.rootGraphQLTypes[project][getTableFieldName(db, tableName)].AddFieldConfig(fieldName, &graphql.Field{
					Type:    graphqlType,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) { return a.literalResolveFn(p) },
					Args: graphql.FieldConfigArgument{
						"where": &graphql.ArgumentConfig{
							Type:         a.rootDBWhereTypes[project][getTableWhereClauseName(db, fieldSchema.LinkedTable.Table)],
							DefaultValue: map[string]interface{}{},
						},
					},
				})
				continue
			}
		}

		// Create a record object for the table
		graphqlObject := a.rootGraphQLTypes[project][getTableFieldName(db, tableName)]
		graphqlArguments := graphql.FieldConfigArgument{
			"where": a.getDBWhereClause(project, db, tableName, tableSchema),
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
			Resolve: a.dbReadResolveFn(project, db, tableName, utils.All, schemas),
		}
		fields[fieldKeySingle] = &graphql.Field{
			Type:    graphqlObject,
			Args:    graphqlArguments,
			Resolve: a.dbReadResolveFn(project, db, tableName, utils.One, schemas),
		}

		// Add to join return objects
		joinRetObj[fieldKeyMultiple] = ""
		joinRetObj[fieldKeySingle] = ""
		a.registeredQueries[fieldKeyMultiple] = struct{}{}
		a.registeredQueries[fieldKeySingle] = struct{}{}

		// Add the join field
		a.rootGraphQLTypes[project][getTableFieldName(db, tableName)].AddFieldConfig("_join", &graphql.Field{
			Type: queryType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return func() (interface{}, error) {
					return joinRetObj, nil
				}, nil
			},
		})

		// // Add the agg field
		// a.rootGraphQLTypes[project][getTableFieldName(db, tableName)].AddFieldConfig("_aggregate", &graphql.Field{
		// 	Type: graphql.Float,
		// 	Args: graphql.FieldConfigArgument{
		// 		"field": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
		// 		"op":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(aggregateOperationType)},
		// 	},
		// })
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
		return anyType

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

func getTableFieldName(db, table string) string {
	return fmt.Sprintf("%s_%s", strings.Title(db), strings.Title(table))
}

func getTableWhereClauseName(db, table string) string {
	return fmt.Sprintf("%s_%s_WhereClause", strings.Title(db), strings.Title(table))
}
