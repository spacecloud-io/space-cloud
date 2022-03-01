package common

import (
	"encoding/json"
	"fmt"

	"github.com/caddyserver/caddy/v2"
)

func getMatcherSet(method, path string) []caddy.ModuleMap {
	return []caddy.ModuleMap{
		{
			"path":   getByteStringArray(path),
			"method": getByteStringArray(method),
		},
	}
}

func getByteStringArray(val string) []byte {
	vals := []string{val}
	data, _ := json.Marshal(vals)
	return data
}

func getHandler(handlerName, operation string) []byte {
	handler := make(map[string]string)

	handler["handler"] = fmt.Sprintf("sc_%s_handler", handlerName)
	handler["operation"] = operation

	data, _ := json.Marshal(handler)
	return data
}
