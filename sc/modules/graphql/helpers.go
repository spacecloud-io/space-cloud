package graphql

import (
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/spacecloud-io/space-cloud/model"
	"github.com/spacecloud-io/space-cloud/utils"
)

func (a *App) createRootMutationGraphQLTypes(project string) {
	a.rootGraphQLMutationTypes[project] = map[string]*graphql.Object{}
	a.rootDBDocsTypes[project] = map[string]*graphql.InputObject{}

	for dbAlias, parsedSchema := range a.dbSchemas[project] {
		for tableName := range parsedSchema {
			a.rootGraphQLMutationTypes[project][getInsertResponseTypeName(dbAlias, tableName)] = graphql.NewObject(graphql.ObjectConfig{
				Name:        fmt.Sprintf("%s_%s_Returning", strings.Title(dbAlias), strings.Title(tableName)),
				Description: fmt.Sprintf("Mutation record object for %s", tableName),
				Fields:      graphql.Fields{},
			})
			a.rootDBDocsTypes[project][getInsertDocsName(dbAlias, tableName)] = graphql.NewInputObject(graphql.InputObjectConfig{
				Name:        getInsertDocsName(dbAlias, tableName),
				Description: fmt.Sprintf("Record to be inserted for %s", tableName),
				Fields:      make(graphql.InputObjectConfigFieldMap),
			})
		}
	}
}

func (a *App) createRootQueryGraphQLTypes(project string) {
	a.rootGraphQLQueryTypes[project] = map[string]*graphql.Object{}
	a.rootDBWhereTypes[project] = map[string]*graphql.InputObject{}

	for dbAlias, parsedSchema := range a.dbSchemas[project] {
		for tableName := range parsedSchema {
			a.rootGraphQLQueryTypes[project][getTableFieldName(dbAlias, tableName)] = graphql.NewObject(graphql.ObjectConfig{
				Name:        getTableFieldName(dbAlias, tableName),
				Description: fmt.Sprintf("Record object for %s", tableName),
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

func (a *App) getMutationType(project string) *graphql.Object {
	// Create root mutation graphql types
	a.createRootMutationGraphQLTypes(project)

	// Create the root mutation type
	mutationType := graphql.NewObject(graphql.ObjectConfig{Name: "Mutation", Fields: graphql.Fields{}})

	// Populate all the fields in the root query
	for dbAlias, parsedSchema := range a.dbSchemas[project] {
		for tableName, tableSchema := range parsedSchema {
			tableReturningType := a.rootGraphQLMutationTypes[project][getInsertResponseTypeName(dbAlias, tableName)]
			tableDocsType := a.rootDBDocsTypes[project][getInsertDocsName(dbAlias, tableName)]
			updateSetObjectType := graphql.NewInputObject(graphql.InputObjectConfig{Name: fmt.Sprintf("%s_%sSetObject", strings.Title(dbAlias), strings.Title(tableName)), Fields: make(graphql.InputObjectConfigFieldMap)})
			updateNumericOpType := graphql.NewInputObject(graphql.InputObjectConfig{Name: fmt.Sprintf("%s_%sNumericOpObject", strings.Title(dbAlias), strings.Title(tableName)), Fields: make(graphql.InputObjectConfigFieldMap)})
			updateCurrentDateOpType := graphql.NewInputObject(graphql.InputObjectConfig{Name: fmt.Sprintf("%s_%sCurrentDateOpObject", strings.Title(dbAlias), strings.Title(tableName)), Fields: make(graphql.InputObjectConfigFieldMap)})
			for fieldName, fieldSchema := range tableSchema {
				if !fieldSchema.IsLinked {
					graphqlType := scToGraphQLType(fieldSchema.Kind)
					tableReturningType.AddFieldConfig(fieldName, &graphql.Field{Type: graphqlType, Resolve: a.literalResolveFn})
					tableDocsType.AddFieldConfig(fieldName, &graphql.InputObjectFieldConfig{Type: graphqlType})
					updateSetObjectType.AddFieldConfig(fieldName, &graphql.InputObjectFieldConfig{Type: graphqlType})
					if graphqlType == graphql.Int || graphqlType == graphql.Float {
						updateNumericOpType.AddFieldConfig(fieldName, &graphql.InputObjectFieldConfig{Type: graphqlType})
					}
					if graphqlType == graphql.DateTime {
						updateCurrentDateOpType.AddFieldConfig(fieldName, &graphql.InputObjectFieldConfig{Type: graphqlType})
					}
					continue
				}
				if fieldSchema.IsLinked {
					var graphqlOutputType graphql.Output = a.rootGraphQLMutationTypes[project][getInsertResponseTypeName(dbAlias, fieldSchema.LinkedTable.Table)]
					if fieldSchema.IsList {
						graphqlOutputType = graphql.NewList(graphqlOutputType)
					}
					tableReturningType.AddFieldConfig(fieldName, &graphql.Field{Type: graphqlOutputType, Resolve: a.literalResolveFn})

					var graphqlInputType graphql.Input = a.rootDBDocsTypes[project][getInsertDocsName(dbAlias, fieldSchema.LinkedTable.Table)]
					if fieldSchema.IsList {
						graphqlInputType = graphql.NewList(a.rootDBDocsTypes[project][getInsertDocsName(dbAlias, fieldSchema.LinkedTable.Table)])
					}
					tableDocsType.AddFieldConfig(fieldName, &graphql.InputObjectFieldConfig{Type: graphqlInputType})
					continue
				}
			}

			// Add the insert query to the root mutation type
			mutationType.AddFieldConfig(fmt.Sprintf("%s_insert%s", dbAlias, strings.Title(tableName)), &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name:   getInsertResponseTypeName(dbAlias, tableName),
					Fields: graphql.Fields{"returning": &graphql.Field{Name: fmt.Sprintf("%s_%s_ReturningList", strings.Title(dbAlias), strings.Title(tableName)), Type: graphql.NewList(tableReturningType)}},
				}),
				Args: graphql.FieldConfigArgument{
					"docs": &graphql.ArgumentConfig{Type: graphql.NewList(tableDocsType), DefaultValue: []interface{}{}},
				},
				Resolve: a.dbInsertResolveFn(project, dbAlias, tableName, parsedSchema),
			})

			// Add the delete query to the root mutation type
			mutationType.AddFieldConfig(fmt.Sprintf("%s_delete%s", dbAlias, strings.Title(tableName)), &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name:   getDeleteResponseTypeName(dbAlias, tableName),
					Fields: graphql.Fields{"status": &graphql.Field{Type: graphql.Int}},
				}),
				Args: graphql.FieldConfigArgument{
					"where": &graphql.ArgumentConfig{
						Type:         a.rootDBWhereTypes[project][getTableWhereClauseName(dbAlias, tableName)],
						DefaultValue: map[string]interface{}{},
					},
				},
				Resolve: a.dbDeleteResolveFn(project, dbAlias, tableName),
			})

			fmt.Println("Locs:", len(updateNumericOpType.Fields()), len(updateSetObjectType.Fields()), len(updateCurrentDateOpType.Fields()))

			// Prepare the update operation args
			updateArgs := graphql.FieldConfigArgument{
				"set": &graphql.ArgumentConfig{Type: updateSetObjectType, DefaultValue: map[string]interface{}{}},
			}
			if len(updateNumericOpType.Fields()) > 0 {
				updateArgs["inc"] = &graphql.ArgumentConfig{Type: updateNumericOpType, DefaultValue: map[string]interface{}{}}
				updateArgs["mul"] = &graphql.ArgumentConfig{Type: updateNumericOpType, DefaultValue: map[string]interface{}{}}
				updateArgs["min"] = &graphql.ArgumentConfig{Type: updateNumericOpType, DefaultValue: map[string]interface{}{}}
				updateArgs["max"] = &graphql.ArgumentConfig{Type: updateNumericOpType, DefaultValue: map[string]interface{}{}}
			}
			if len(updateCurrentDateOpType.Fields()) > 0 {
				updateArgs["curentDate"] = &graphql.ArgumentConfig{Type: updateCurrentDateOpType, DefaultValue: map[string]interface{}{}}
			}

			// Add update query to the root mutation type
			mutationType.AddFieldConfig(fmt.Sprintf("%s_delete%s", dbAlias, strings.Title(tableName)), &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name:   getUpdateResponseTypeName(dbAlias, tableName),
					Fields: graphql.Fields{"status": &graphql.Field{Type: graphql.Int}},
				}),
				Args: updateArgs,
			})
		}
	}

	return mutationType
}

func (a *App) getQueryType(project string) *graphql.Object {
	// Create root query graphql types
	a.createRootQueryGraphQLTypes(project)

	// Create the root query type
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
				a.rootGraphQLQueryTypes[project][getTableFieldName(db, tableName)].AddFieldConfig(fieldName, &graphql.Field{
					Type:    scToGraphQLType(fieldSchema.Kind),
					Resolve: a.literalResolveFn,
				})
				continue
			}

			if fieldSchema.IsLinked {
				var graphqlType graphql.Output = a.rootGraphQLQueryTypes[project][getTableFieldName(db, fieldSchema.LinkedTable.Table)]
				if fieldSchema.IsList {
					graphqlType = graphql.NewList(graphqlType)
				}
				a.rootGraphQLQueryTypes[project][getTableFieldName(db, tableName)].AddFieldConfig(fieldName, &graphql.Field{
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
		graphqlObject := a.rootGraphQLQueryTypes[project][getTableFieldName(db, tableName)]
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
		a.rootGraphQLQueryTypes[project][getTableFieldName(db, tableName)].AddFieldConfig("_join", &graphql.Field{
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
	case model.TypeDateTimeWithZone, model.TypeDateTime, model.TypeDate, model.TypeTime:
		return graphql.DateTime
	case model.TypeUUID, model.TypeChar, model.TypeVarChar, model.TypeString, model.TypeID, model.TypeEnum:
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
	return fmt.Sprintf("%s_%s_Query", strings.Title(db), strings.Title(table))
}

func getInsertResponseTypeName(db, table string) string {
	return fmt.Sprintf("%s_%s_InsertResponse", strings.Title(db), strings.Title(table))
}

func getDeleteResponseTypeName(db, table string) string {
	return fmt.Sprintf("%s_%s_DeleteResponse", strings.Title(db), strings.Title(table))
}

func getUpdateResponseTypeName(db, table string) string {
	return fmt.Sprintf("%s_%s_UpdateResponse", strings.Title(db), strings.Title(table))
}

func getInsertDocsName(db, table string) string {
	return fmt.Sprintf("%s_%s_InsertDocs", strings.Title(db), strings.Title(table))
}

func getTableWhereClauseName(db, table string) string {
	return fmt.Sprintf("%s_%s_WhereClause", strings.Title(db), strings.Title(table))
}
