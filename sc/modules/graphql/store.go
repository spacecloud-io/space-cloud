package graphql

import (
	"strconv"
	"strings"
	"sync"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/language/ast"

	"github.com/spacecloud-io/space-cloud/model"
	"github.com/spacecloud-io/space-cloud/utils"
)

type store struct {
	lock    sync.RWMutex
	m       map[string]interface{}
	exports map[string]int

	// Mutation related variables
	isMutation         bool
	mutationFieldCount int
	allReq             []*model.AllRequest
	err                error
	wg                 sync.WaitGroup
}

func newStore(graphqlAST *ast.Document) (*store, error) {
	s := &store{m: map[string]interface{}{}, exports: map[string]int{}, isMutation: false, mutationFieldCount: 0}

	dbName := ""

	for _, node := range graphqlAST.Definitions {
		opDef, ok := node.(*ast.OperationDefinition)
		if !ok {
			continue
		}

		if opDef.Operation == ast.OperationTypeMutation {
			s.isMutation = true

			for _, v := range opDef.SelectionSet.Selections {
				field := v.(*ast.Field)
				tempDBName := strings.Split(field.Name.Value, "_")[0]
				if dbName == "" || dbName == tempDBName {
					dbName = tempDBName
					s.mutationFieldCount++
					s.wg.Add(1)
					continue
				}

				err := gqlerrors.NewError("Cannot perform mutations on more than one database.", []ast.Node{field}, "", nil, nil, nil)
				return nil, gqlerrors.FormatError(err)
			}
		}
	}

	return s, nil
}

func (s *store) addMutationRequest(allReq []*model.AllRequest) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.allReq = append(s.allReq, allReq...)

	s.mutationFieldCount--
	return s.mutationFieldCount == 0
}

func (s *store) done() {
	s.wg.Done()
}

func (s *store) wait() {
	s.wg.Wait()
}

func (s *store) load(key interface{}, path *graphql.ResponsePath) (interface{}, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	keyString, ok := key.(string)
	if !ok {
		return nil, false
	}

	// Return if value is not refering to exports
	if !strings.HasPrefix(keyString, "exports.") {
		return nil, false
	}
	keyString = strings.TrimPrefix(keyString, "exports.")
	rootKey := strings.Split(keyString, ".")[0]

	exportsPrefixLen, p := s.exports[rootKey]
	if !p {
		return nil, false
	}

	pathPrefix := traversePath(path, false)
	pathPrefix = strings.Join(strings.Split(pathPrefix, ".")[:exportsPrefixLen], ".")

	v, p := s.m[rootKey+":"+pathPrefix]
	if !p {
		return nil, false
	}

	// We might have to do a proper load value if key string contained a nested structure
	if obj, ok := v.(map[string]interface{}); ok && rootKey != keyString {
		v, err := utils.LoadValue(keyString, obj)
		if err != nil {
			return nil, false
		}

		return v, true
	}
	return v, p
}

func (s *store) store(key string, value interface{}, path *graphql.ResponsePath) {
	s.lock.Lock()
	defer s.lock.Unlock()

	pathPrefix := traversePath(path, false)
	storeKey := key + ":" + pathPrefix
	s.m[storeKey] = value

	s.exports[key] = len(strings.Split(pathPrefix, "."))
}

func traversePath(path *graphql.ResponsePath, capture bool) string {
	// Exit if path is empty
	if path == nil {
		return ""
	}

	// Mark capture mode to true if key is an index
	if i, ok := path.Key.(int); ok {
		return strconv.Itoa(i) + "." + traversePath(path.Prev, true)
	}

	key := traversePath(path.Prev, capture)
	if capture {
		// We add the path key only if in capture mode
		key = path.Key.(string) + "." + key

		// Remove the trailing dot
		key = strings.TrimSuffix(key, ".")
	}

	return key
}
