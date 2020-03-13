package utils

import "strconv"

// AcceptableIDType converts a provied id to string
func AcceptableIDType(id interface{}) (string, bool) {
	switch v := id.(type) {
	case string:
		return v, true
	case int:
		return strconv.Itoa(v), true
	case int32:
		return strconv.FormatInt(int64(v), 10), true
	case int64:
		return strconv.FormatInt(v, 10), true
	case float64:
		// json.Unmarshal converts all numbers to float64
		if float64(int64(v)) == v { // v is actually an int
			return strconv.FormatInt(int64(v), 10), true
		}
		return "", false
	default:
		return "", false
	}
}

// GetIDVariable gets the id variable for the provided db type
func GetIDVariable(dbAlias string) string {
	idVar := "id"
	if DBType(dbAlias) == Mongo {
		idVar = "_id"
	}

	return idVar
}

// ArrayContains checks if the array contains the value provided
func ArrayContains(array []interface{}, value interface{}) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}
