package rootvalue

import (
	"sync"

	"github.com/graph-gophers/dataloader"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/spacecloud-io/space-cloud/modules/graphql/types"
)

type (
	// RootValue is a global store used in every graphql query execution.
	// Note: Make sure that a new root value is used for the invocation of each graphql request.
	RootValue struct {
		dlMutex    sync.Mutex
		errMutex   sync.Mutex
		storeMutex sync.RWMutex

		// Query generation vars
		operationAST *ast.OperationDefinition

		// Data loaders
		graphqlLoaders map[string]*dataloader.Loader

		// Error handling
		formatedErrors gqlerrors.FormattedErrors

		// Store for all exported variables
		exportsKeyLength map[string]int
		store            map[string]*types.StoreValue
		exportedVars     map[string]struct{}
	}
)

func New(graphqlDoc *ast.Document) *RootValue {
	root := new(RootValue)
	root.operationAST = graphqlDoc.Definitions[0].(*ast.OperationDefinition)

	// Create a new map for data loader
	root.graphqlLoaders = make(map[string]*dataloader.Loader)

	// Create new maps for store management
	root.store = make(map[string]*types.StoreValue)
	root.exportsKeyLength = make(map[string]int)
	root.exportedVars = make(map[string]struct{})

	return root
}
