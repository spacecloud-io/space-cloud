package graphql

import (
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/spacecloud-io/space-cloud/model"
)

func (a *App) getQueryType(project string) *graphql.Object {
	// Create the root query
	queryType := graphql.NewObject(graphql.ObjectConfig{Name: "Query", Fields: graphql.Fields{}})

	// Populate all the fields in the root query
	for dbAlias, parsedSchema := range a.dbSchemas[project] {
		for k, v := range a.getTableFields(project, dbAlias, parsedSchema, queryType) {
			queryType.AddFieldConfig(k, v)
		}
	}

	return queryType
}

func (a *App) getTableFields(project, db string, schemas model.CollectionSchemas, queryType *graphql.Object) graphql.Fields {
	fields := make(graphql.Fields, len(schemas))

	for tableName, tableSchema := range schemas {
		tableFields := make(graphql.Fields, len(tableSchema))

		for fieldName, fieldSchema := range tableSchema {
			tableFields[fieldName] = a.getArbitaryField(fieldSchema)
		}

		fields[fmt.Sprintf("%s_many%s", db, strings.Title(tableName))] = &graphql.Field{
			Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
				Name:        fmt.Sprintf("%s_%s", strings.Title(db), strings.Title(tableName)),
				Description: fmt.Sprintf("Get multiple records from %s", tableName),
				Fields:      tableFields,
			})),
			Resolve: a.dbReadResolveFn(project, db, tableName),
		}

		// tableFields["_join"] = &graphql.Field{
		// 	Type: queryType,
		// 	// TODO: Add a generic join resolver here
		// }
	}

	return fields
}

func (a *App) getArbitaryField(fieldType *model.FieldType) *graphql.Field {
	switch fieldType.Kind {
	case model.TypeDateTimeWithZone, model.TypeDateTime, model.TypeDate, model.TypeTime, model.TypeUUID, model.TypeChar, model.TypeVarChar, model.TypeString, model.TypeID, model.TypeEnum:
		return &graphql.Field{
			Type:    graphql.String,
			Resolve: graphql.DefaultResolveFn,
		}
	case model.TypeInteger, model.TypeSmallInteger, model.TypeBigInteger:
		return &graphql.Field{
			Type:    graphql.Int,
			Resolve: graphql.DefaultResolveFn,
		}
	case model.TypeDecimal, model.TypeFloat:
		return &graphql.Field{
			Type:    graphql.Float,
			Resolve: graphql.DefaultResolveFn,
		}
	case model.TypeBoolean:
		return &graphql.Field{
			Type:    graphql.Boolean,
			Resolve: graphql.DefaultResolveFn,
		}
	// TODO: Add a case for JSON types
	default:
		// TODO: Return an `any` scaler type
		return &graphql.Field{
			Type:    graphql.String,
			Resolve: graphql.DefaultResolveFn,
		}
	}
}
