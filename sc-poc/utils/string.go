package utils

// StringExists returns true if the given string exists in the array
func StringExists(value string, elements ...string) bool {
	for _, e := range elements {
		if e == value {
			return true
		}
	}
	return false
}
