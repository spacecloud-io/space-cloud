package utils

import (
	"regexp"
	"strings"
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
	plural := word
	if match, _ := regexp.MatchString("[sxz]$", word); !match {
		plural += "s"
	} else if match, _ := regexp.MatchString("[^aeioudgkprt]h$", word); match {
		plural += "es"
	} else if match, _ := regexp.MatchString("[^aeiou]y$", word); match {
		plural = strings.Replace(plural, "y", "ies", -1)
	}
	return strings.ToLower(plural)
}
