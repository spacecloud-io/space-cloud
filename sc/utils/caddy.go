package utils

import (
	"encoding/json"
	"fmt"

	"github.com/caddyserver/caddy/v2"
)

// GetCaddyMatcherSet returns a caddy matcher set
func GetCaddyMatcherSet(path []string, methods []string) []caddy.ModuleMap {
	// We will always need to match based on the path
	set := map[string]json.RawMessage{
		"path": GetByteStringArray(path...),
	}

	// Match on method if provided
	if len(methods) > 0 {
		set["method"] = GetByteStringArray(methods...)
	}

	// Return the match set
	return []caddy.ModuleMap{set}
}

// GetByteStringArray returns an array of string in json form
func GetByteStringArray(val ...string) []byte {
	data, _ := json.Marshal(val)
	return data
}

// GetCaddyHandler returns a marshaled caddy handler config
func GetCaddyHandler(handlerName string, params map[string]interface{}) []json.RawMessage {
	handler := make(map[string]interface{})

	// Add the handler name / identifier
	handler["handler"] = fmt.Sprintf("sc_%s_handler", handlerName)

	// Add the params the handler needs
	for k, v := range params {
		handler[k] = v
	}

	data, _ := json.Marshal(handler)
	return []json.RawMessage{data}
}
