package graphql

import (
	"fmt"
	"strings"

	"github.com/graph-gophers/dataloader"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"

	"github.com/spacecloud-io/space-cloud/modules/graphql/rootvalue"
	"github.com/spacecloud-io/space-cloud/modules/graphql/types"
	"github.com/spacecloud-io/space-cloud/utils"
)

func (s *GraphqlSource) resolveRemoteGraphqlQuery() graphql.FieldResolveFn {
	type channelMsg struct {
		data interface{}
		err  error
	}

	return func(p graphql.ResolveParams) (interface{}, error) {
		root := p.Info.RootValue.(*rootvalue.RootValue)
		loader := root.CreateOrStoreDataLoader(rootvalue.CreateDataloaderKey("graphql", s.Name), func() *dataloader.Loader {
			return dataloader.NewBatchedLoader(s.grapqhlLoaderBatchFn(p.Info.Operation.(*ast.OperationDefinition), p.Info.VariableValues))
		})

		fieldAst := p.Info.FieldASTs[0]

		// Get field level query
		vars := map[string]struct{}{}
		query, loaderKey, exportedVars := s.extractFieldQuery(root, fieldAst, vars, p.Info.Path, true)

		c := make(chan channelMsg)
		go func() {
			data, err := loader.Load(p.Context, &GraphqlLoaderKey{
				FieldName:    loaderKey,
				Query:        query,
				AllowedVars:  vars,
				ExportedVars: exportedVars,
			})()
			c <- channelMsg{data, err}
			close(c)
		}()

		return func() (interface{}, error) {
			msg := <-c

			// First handle the error
			if msg.err != nil {
				switch v := msg.err.(type) {
				case *types.GraphqlError:
					root.AddFormatedErrors(v.FormatedErrors)
					return nil, msg.err
				default:
					return nil, msg.err
				}
			}

			return msg.data, nil
		}, nil
	}
}

func (s *GraphqlSource) extractFieldQuery(root *rootvalue.RootValue, fieldAst *ast.Field, allowedVars map[string]struct{}, path *graphql.ResponsePath, allowExportingOfVars bool) (string, string, map[string]*types.StoreValue) {
	query := ""

	// Add the arguments if any
	queryArgs := ""
	if len(fieldAst.Arguments) > 0 {
		queryArgs += "("
		for _, arg := range fieldAst.Arguments {
			// Add the value to the query
			queryArgs += arg.Name.Value + ": "
			queryArgs += string(arg.Value.GetLoc().Source.Body[arg.Value.GetLoc().Start:arg.Value.GetLoc().End])
			queryArgs += ", "

			// See if the value is a variable. Need to store this mapping.
			checkForVariablesInValue(arg.Value, allowedVars)
		}
		queryArgs += ") "
	}

	// Get values of all variables which were exported

	exportedVars := root.GetExportedVarsWithValues(allowedVars, path)
	newExportedVars := make(map[string]*types.StoreValue, len(exportedVars))

	// Prepare a loader key which will be the same as the field name by default
	loaderKey := strings.TrimPrefix(fieldAst.Name.Value, fmt.Sprintf("%s_", s.Name))
	if len(exportedVars) > 0 && allowExportingOfVars {

		// Loader key will become the alias if exported vars are present. This is done to
		// make sure each field in the query is unique

		// Prepare a random suffix
		randomSuffix := ""
		for _, kv := range exportedVars {
			randomSuffix += kv.Key
			randomSuffix += fmt.Sprintf("%v", kv.Value.Value)
		}
		randomSuffix = utils.Hash(randomSuffix)

		// Add the suffix to the loader key
		loaderKey += randomSuffix

		for _, kv := range exportedVars {
			// Append key and value to the random suffix
			// We will replace the existing variable with a new one for each exported variable.
			// This allows us to have different variables for each query.
			ogKey := fmt.Sprintf("$%s", kv.Key)
			newKey := fmt.Sprintf("%s%s", kv.Key, randomSuffix)
			queryArgs = strings.ReplaceAll(queryArgs, ogKey, "$"+newKey)

			// Lets replace the old key with the new one from the allowed variables map
			delete(allowedVars, kv.Key)
			allowedVars[newKey] = struct{}{}

			// Add the new key to the new exported variables map. This will be used to populate the final
			// graphql variables sent to the remote graphql source.
			newExportedVars[newKey] = &kv.Value
		}

		query += loaderKey + ": "
	}

	// Add the field name. We need to add a custom alias if exported variables are used
	query += strings.TrimPrefix(fieldAst.Name.Value, fmt.Sprintf("%s_", s.Name)) + " "

	// Add the field args
	query += queryArgs

	// Add directives if any
	if len(fieldAst.Directives) > 0 {
		for _, d := range fieldAst.Directives {
			// TODO: Only remove those directives which are not allowed by that particular source or
			// we can remove all directives used by space cloud itself
			if utils.StringExists(d.Name.Value, "export", "injectClaim", "tag") {
				continue
			}

			query += string(d.Loc.Source.Body[d.Loc.Start:d.Loc.End])
			query += " "
		}
	}

	// Add the selection set if provided
	if fieldAst.SelectionSet != nil && len(fieldAst.SelectionSet.Selections) > 0 {
		query += "{\n"
		for _, selection := range fieldAst.SelectionSet.Selections {
			f := selection.(*ast.Field)

			// Skip if its the join field. That will not be part of the query sent to remote
			if f.Name.Value == "_join" {
				continue
			}

			// Add the selection set.
			// TODO: Account for the exported variables used in the selection set
			newQuery, _, _ := s.extractFieldQuery(root, f, allowedVars, path.WithKey(fieldAst.Name.Value), false)
			query += newQuery
		}
		query += "}"
	}

	query += "\n"

	return query, loaderKey, newExportedVars
}
