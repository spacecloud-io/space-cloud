package graphql

import (
	"fmt"
	"reflect"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"

	"github.com/spacecloud-io/space-cloud/model"
)

func (a *App) dbReadResolveFn(project, db, tableName, op string, dbSchema model.CollectionSchemas) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		s := p.Info.RootValue.(*store)

		// Prepare the database where clause
		where := adjustWhereClause(p.Args["where"].(map[string]interface{}), s, p.Info.Path)

		// Generate the options
		options := &model.ReadOptions{Select: make(map[string]int32)}
		options.Sort = adjustSortArgument(p.Args["sort"].(map[string]interface{}))

		// Get Skip and Limit
		options.Skip = extractIntegerFromArg("skip", p.Args)
		options.Limit = extractIntegerFromArg("limit", p.Args)

		// Get joins
		calculateJoins(tableName, p.Info.FieldASTs[0], options.Select, &options.Join, dbSchema)

		// We return a thunk function since we want to execute this resolver concurrently
		return func() (interface{}, error) {
			r, _, err := a.database.Read(p.Context, project, db, tableName, &model.ReadRequest{Operation: op, Find: where, Options: options}, model.RequestParams{})
			return r, err
		}, nil
	}
}

func calculateJoins(parentTable string, parentFieldAST *ast.Field, selectedField map[string]int32, join *[]*model.JoinOption, dbSchema model.CollectionSchemas) {
	for _, t := range parentFieldAST.SelectionSet.Selections {
		// Get the field name
		fieldAST := t.(*ast.Field)
		fieldName := fieldAST.Name.Value
		fieldSchema := dbSchema[parentTable][fieldName]

		// Skip the join field
		if fieldName == "_join" {
			continue
		}

		// First add the field to select
		if !fieldSchema.IsLinked {
			selectedField[fmt.Sprintf("%s.%s", parentTable, fieldName)] = 1
		}

		// Add a join clause if a linked field is requested for
		if fieldSchema.IsLinked {
			joinOption := &model.JoinOption{
				Op:    "one",
				Type:  "LEFT",
				Table: fieldSchema.LinkedTable.Table,
				As:    fieldName,
				On: map[string]interface{}{
					fmt.Sprintf("%s.%s", parentTable, fieldSchema.LinkedTable.From): map[string]interface{}{"$eq": fmt.Sprintf("%s.%s", fieldSchema.LinkedTable.Table, fieldSchema.LinkedTable.To)},
				},
				Join: []*model.JoinOption{},
			}
			if fieldSchema.IsList {
				joinOption.Op = "all"
			}
			*join = append(*join, joinOption)

			// Iterate object joint table to add them to select
			calculateJoins(fieldSchema.LinkedTable.Table, fieldAST, selectedField, &joinOption.Join, dbSchema)
		}
	}
}

func (a *App) literalResolveFn(p graphql.ResolveParams) (interface{}, error) {
	// Get the value from the source map
	srcMap, ok := p.Source.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid type '%s' received for field '%s'", reflect.TypeOf(p.Source).String(), p.Info.FieldName)
	}
	val := srcMap[p.Info.FieldName]

	// Return if value is nil
	if val == nil {
		return nil, nil
	}

	// Store the source in the main store if the value isn't nil
	fieldAST := p.Info.FieldASTs[0]
	if key, ok := containsExportDirective(fieldAST); ok {
		s := p.Info.RootValue.(*store)
		s.store(key, val, p.Info.Path)
	}
	return val, nil
}
