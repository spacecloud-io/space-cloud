package rootvalue

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/graph-gophers/dataloader"
	"github.com/graphql-go/graphql"

	"github.com/spacecloud-io/space-cloud/modules/graphql/types"
	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
)

// GraphqlLoaderKey describes a key used by the graphql dataloader
type GraphqlLoaderKey struct {
	FieldName    string
	Query        string
	AllowedVars  map[string]struct{}
	ExportedVars map[string]*types.StoreValue
}

// String returns the string representation of the key
func (k *GraphqlLoaderKey) String() string {
	return k.FieldName
}

// Raw returns the raw query value
func (k *GraphqlLoaderKey) Raw() interface{} {
	return k.Query
}

// CreateRemoteGraphqlLoader creates a new dataloader for graphql sources
func (root *RootValue) CreateRemoteGraphqlLoader(sources []*v1alpha1.GraphqlSource, source *v1alpha1.GraphqlSource, vars map[string]interface{}) *dataloader.Loader {
	root.dlMutex.Lock()
	defer root.dlMutex.Unlock()

	if l, p := root.graphqlLoaders[source.Name]; p {
		return l
	}

	loader := dataloader.NewBatchedLoader(root.grapqhlLoaderBatchFn(sources, source, vars))
	root.graphqlLoaders[source.Name] = loader
	return loader
}

func (root *RootValue) grapqhlLoaderBatchFn(sources []*v1alpha1.GraphqlSource, source *v1alpha1.GraphqlSource, graphqlVars map[string]interface{}) dataloader.BatchFunc {
	return func(ctx context.Context, keys dataloader.Keys) (results []*dataloader.Result) {
		// Make a result object
		results = make([]*dataloader.Result, len(keys))

		// Prepare list of allowed & exported vars
		allowedVars := map[string]struct{}{}
		exportedVars := map[string]*types.StoreValue{}
		for _, temp := range keys {
			key := temp.(*GraphqlLoaderKey)
			for k, v := range key.AllowedVars {
				allowedVars[k] = v
			}
			for k, v := range key.ExportedVars {
				exportedVars[k] = v
			}
		}

		// Remove unneeded keys from GraphQL variable
		newGraphqlVars := make(map[string]interface{}, len(allowedVars))
		for k, v := range graphqlVars {
			if _, p := allowedVars[k]; p {
				newGraphqlVars[k] = v
			}
		}

		// Inject the vars that are exported
		for k, v := range exportedVars {
			newGraphqlVars[k] = v.Value
		}

		// For the graphql query for this source
		prefix, suffix := extractQueryPrefixSuffix(sources, root.operationAST, allowedVars, exportedVars)
		query := prefix

		// Add all the queries now
		for _, temp := range keys {
			key := temp.(*GraphqlLoaderKey)
			query += key.Query
		}

		query += suffix

		fmt.Println("=============================")
		fmt.Println(query)
		fmt.Println("=============================")

		reqBody, _ := json.Marshal(map[string]interface{}{
			"query":     query,
			"variables": newGraphqlVars,
		})
		resp, err := http.Post(source.Spec.Source.URL, "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			for i := range keys {
				results[i] = &dataloader.Result{Error: err}
			}
			return
		}
		defer resp.Body.Close()

		graphqlResult := new(graphql.Result)
		if err := json.NewDecoder(resp.Body).Decode(graphqlResult); err != nil {
			for i := range keys {
				results[i] = &dataloader.Result{Error: err}
			}
			return
		}

		if graphqlResult.HasErrors() {
			err := &types.GraphqlError{FormatedErrors: graphqlResult.Errors}
			for i := range keys {
				results[i] = &dataloader.Result{Error: err}
			}
			return
		}

		// Sort all the graphql results
		data := graphqlResult.Data.(map[string]interface{})
		for i, key := range keys {
			results[i] = &dataloader.Result{Data: data[key.String()]}
		}

		return
	}
}
