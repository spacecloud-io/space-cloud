package rootvalue

import (
	"fmt"

	"github.com/graph-gophers/dataloader"

	"github.com/spacecloud-io/space-cloud/modules/graphql/types"
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

// CreateDataloaderKey created a key for the dataloader
func CreateDataloaderKey(sourceType, sourceName string) string {
	return fmt.Sprintf("%s:%s", sourceType, sourceName)
}

// CreateOrStoreDataLoader creates a dataloader if it doesn't already exists
func (root *RootValue) CreateOrStoreDataLoader(key string, createLoader func() *dataloader.Loader) *dataloader.Loader {
	root.dlMutex.Lock()
	defer root.dlMutex.Unlock()

	if l, p := root.dataloaders[key]; p {
		return l
	}

	loader := createLoader()
	root.dataloaders[key] = loader
	return loader
}
