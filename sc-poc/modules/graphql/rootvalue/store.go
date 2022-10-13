package rootvalue

import (
	"sort"
	"strconv"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/spacecloud-io/space-cloud/modules/graphql/types"
)

// GetExportedVarsWithValues returns a list of exported vars used along with their corresponding values
func (root *RootValue) GetExportedVarsWithValues(usedVars map[string]struct{}, path *graphql.ResponsePath) []*types.StoreKeyValue {
	root.storeMutex.RLock()
	usedExportedVarsTemp := make(map[string]struct{})
	for k := range usedVars {
		if _, p := root.exportedVars[k]; p {
			usedExportedVarsTemp[k] = struct{}{}
		}
	}
	root.storeMutex.RUnlock()

	usedExportedVars := make([]*types.StoreKeyValue, 0, len(usedExportedVarsTemp))
	for k := range usedExportedVarsTemp {
		v, t, p := root.loadValue(k, path)
		if p {
			usedExportedVars = append(usedExportedVars, &types.StoreKeyValue{Key: k, Value: types.StoreValue{Value: v, TypeOf: t}})
		}
	}

	sort.SliceStable(usedExportedVars, func(i, j int) bool {
		return usedExportedVars[i].Key < usedExportedVars[j].Key
	})

	return usedExportedVars
}

func (root *RootValue) loadValue(key interface{}, path *graphql.ResponsePath) (interface{}, string, bool) {
	root.storeMutex.RLock()
	defer root.storeMutex.RUnlock()

	keyString, ok := key.(string)
	if !ok {
		return nil, "", false
	}

	// Return if value is not refering to exports
	// if !strings.HasPrefix(keyString, "exports.") {
	// 	return nil, false
	// }
	// keyString = strings.TrimPrefix(keyString, "exports.")
	rootKey := strings.Split(keyString, ".")[0]

	exportsPrefixLen, p := root.exportsKeyLength[rootKey]
	if !p {
		return nil, "", false
	}

	pathPrefix := traversePath(path, false)
	pathPrefix = strings.Join(strings.Split(pathPrefix, ".")[:exportsPrefixLen], ".")

	v, p := root.store[rootKey+":"+pathPrefix]
	if !p {
		return nil, "", false
	}

	// // We might have to do a proper load value if key string contained a nested structure
	// if obj, ok := v.(map[string]interface{}); ok && rootKey != keyString {
	// 	v, err := utils.LoadValue(keyString, obj)
	// 	if err != nil {
	// 		return nil, false
	// 	}

	// 	return v, true
	// }
	return v.Value, v.TypeOf, p
}

// StoreExportedValue stores the value against the exported variable
func (root *RootValue) StoreExportedValue(key string, value interface{}, typeOf string, path *graphql.ResponsePath) {
	root.storeMutex.Lock()
	defer root.storeMutex.Unlock()

	// First we store the variable keys that are exported globally
	root.exportedVars[key] = struct{}{}

	pathPrefix := traversePath(path, false)
	storeKey := key + ":" + pathPrefix
	root.store[storeKey] = &types.StoreValue{Value: value, TypeOf: typeOf}

	root.exportsKeyLength[key] = len(strings.Split(pathPrefix, "."))
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
