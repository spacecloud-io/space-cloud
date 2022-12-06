package graphql

import (
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"

	"github.com/spacecloud-io/space-cloud/modules/graphql/rootvalue"
)

func (a *App) resolveJoin() graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		// We need to wait for our siblings to be done
		return func() (interface{}, error) {
			return a.rootJoinObj, nil
		}, nil
	}
}

func (a *App) resolveMiscField(srcName string) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		fieldAST := p.Info.FieldASTs[0]
		fieldValue := p.Source.(map[string]interface{})[p.Info.FieldName]

		root := p.Info.RootValue.(*rootvalue.RootValue)

		// Check if field value is to be exported
		for _, d := range fieldAST.Directives {
			if d.Name.Value == "export" {
				as := d.Arguments[0].Value.(*ast.StringValue).Value
				root.StoreExportedValue(as, fieldValue, strings.TrimPrefix(p.Info.ReturnType.Name(), srcName), p.Info.Path)
			}
		}

		return fieldValue, nil
	}
}
