package json

import (
	"reflect"
	"strings"
	"unicode"
)

func getTag(field reflect.StructField) string {
	return field.Tag.Get("json")
}

func isIgnoredStructField(field reflect.StructField) bool {
	if field.PkgPath != "" {
		if field.Anonymous {
			if !(field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct) && field.Type.Kind() != reflect.Struct {
				return true
			}
		} else {
			// private field
			return true
		}
	}
	tag := getTag(field)
	return tag == "-"
}

type structTag struct {
	key         string
	isTaggedKey bool
	isOmitEmpty bool
	isString    bool
	field       reflect.StructField
}

type structTags []*structTag

func (t structTags) existsKey(key string) bool {
	for _, tt := range t {
		if tt.key == key {
			return true
		}
	}
	return false
}

func isValidTag(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		switch {
		case strings.ContainsRune("!#$%&()*+-./:<=>?@[]^_{|}~ ", c):
			// Backslash and quote chars are reserved, but
			// otherwise any punctuation chars are allowed
			// in a tag name.
		case !unicode.IsLetter(c) && !unicode.IsDigit(c):
			return false
		}
	}
	return true
}

func structTagFromField(field reflect.StructField) *structTag {
	keyName := field.Name
	tag := getTag(field)
	st := &structTag{field: field}
	opts := strings.Split(tag, ",")
	if len(opts) > 0 {
		if opts[0] != "" && isValidTag(opts[0]) {
			keyName = opts[0]
			st.isTaggedKey = true
		}
	}
	st.key = keyName
	if len(opts) > 1 {
		st.isOmitEmpty = opts[1] == "omitempty"
		st.isString = opts[1] == "string"
	}
	return st
}
