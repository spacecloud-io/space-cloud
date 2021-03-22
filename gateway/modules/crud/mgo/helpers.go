package mgo

import (
	"context"
	"strings"
)

func sanitizeWhereClause(ctx context.Context, col string, find map[string]interface{}) map[string]interface{} {
	for key, value := range find {
		arr := strings.Split(key, ".")
		if len(arr) > 1 && arr[0] == col {
			delete(find, key)
			find[strings.Join(arr[1:], ".")] = value
		}
		switch key {
		case "$or":
			objArr, ok := value.([]interface{})
			if ok {
				for _, obj := range objArr {
					t, ok := obj.(map[string]interface{})
					if ok {
						sanitizeWhereClause(ctx, col, t)
					}
				}
			}
		default:
			obj, ok := value.(map[string]interface{})
			if ok {
				sanitizeWhereClause(ctx, col, obj)
			}
		}
	}
	return find
}
