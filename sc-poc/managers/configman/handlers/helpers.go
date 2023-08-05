package handlers

import "strings"

func getName(path string) string {
	arr := strings.Split(path, "/")
	if len(arr) != 8 {
		return ""
	}

	return arr[len(arr)-1]
}
