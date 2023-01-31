package utils

import (
	"strings"

	pluralize "github.com/gertd/go-pluralize"
)

// StringExists returns true if the given string exists in the array
func StringExists(value string, elements ...string) bool {
	for _, e := range elements {
		if e == value {
			return true
		}
	}
	return false
}

func Pluralize(word string) string {
	pluralize := pluralize.NewClient()
	plural := pluralize.Plural(word)
	return strings.ToLower(plural)
}
