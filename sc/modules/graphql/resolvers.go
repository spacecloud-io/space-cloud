package graphql

import (
	"github.com/graphql-go/graphql"

	"github.com/spacecloud-io/space-cloud/model"
	"github.com/spacecloud-io/space-cloud/utils"
)

func (a *App) dbReadResolveFn(project, db, tableName string) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		// TODO: prepare the database read parameters
		where := adjustWhereClause(p.Args["where"].(map[string]interface{}))

		// We return a thunk function since we want to execute this resolver concurrently
		return func() (interface{}, error) {
			r, _, err := a.database.Read(p.Context, project, db, tableName, &model.ReadRequest{Operation: utils.All, Find: where}, model.RequestParams{})
			return r, err
		}, nil
	}
}
