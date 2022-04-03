package utils

import "net/url"

// GetQueryParams generates query params map
func GetQueryParams(queryParams url.Values) map[string]string {
	params := make(map[string]string)

	for key, val := range queryParams {
		params[key] = val[0]
	}

	return params
}
