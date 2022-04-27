package graphql

import (
	"fmt"
	"reflect"

	"github.com/graphql-go/graphql"

	"github.com/spacecloud-io/space-cloud/model"
)

func (a *App) dbReadResolveFn(project, db, tableName, op string) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		s := p.Info.RootValue.(*store)

		// Prepare the database where clause
		where := adjustWhereClause(p.Args["where"].(map[string]interface{}), s, p.Info.Path)

		// Generate the options
		options := &model.ReadOptions{}
		options.Sort = adjustSortArgument(p.Args["sort"].(map[string]interface{}))

		// Get Skip and Limit
		options.Skip = extractIntegerFromArg("skip", p.Args)
		options.Limit = extractIntegerFromArg("limit", p.Args)

		// We return a thunk function since we want to execute this resolver concurrently
		return func() (interface{}, error) {
			r, _, err := a.database.Read(p.Context, project, db, tableName, &model.ReadRequest{Operation: op, Find: where, Options: options}, model.RequestParams{})
			return r, err
		}, nil
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
