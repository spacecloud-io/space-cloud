package utils

import (
	"fmt"
	"log"
	"reflect"
	"regexp"

	"github.com/google/go-cmp/cmp"
)

func attemptConvertBoolToInt64(val interface{}) interface{} {
	if tempBool, ok := val.(bool); ok {
		val = int64(0)
		if tempBool {
			val = int64(1)
		}
	}
	return val
}

func attemptConvertIntToInt64(val interface{}) interface{} {
	if tempInt, ok := val.(int); ok {
		val = int64(tempInt)
	} else if tempInt, ok := val.(int32); ok {
		val = int64(tempInt)
	}
	return val
}

func attemptConvertInt64ToFloat(val interface{}) interface{} {
	if tempInt, ok := val.(int64); ok {
		val = float64(tempInt)
	}
	return val
}

func compare(v1, v2 interface{}) bool {
	if reflect.TypeOf(v1).String() == reflect.Int64.String() {
		return fmt.Sprintf("%v", v1) == fmt.Sprintf("%v", v2)
	}

	return cmp.Equal(v1, v2)
}

func adjustValTypes(v1, v2 interface{}) (interface{}, interface{}) {
	v1 = attemptConvertBoolToInt64(v1)
	v2 = attemptConvertBoolToInt64(v2)

	v1 = attemptConvertIntToInt64(v1)
	v2 = attemptConvertIntToInt64(v2)

	v1 = attemptConvertInt64ToFloat(v1)
	v2 = attemptConvertInt64ToFloat(v2)
	return v1, v2
}

// Validate checks if the provided document matches with the where clause
func Validate(where map[string]interface{}, obj interface{}) bool {
	if res, ok := obj.(map[string]interface{}); ok {
		for k, temp := range where {
			if k == "$or" {
				array, ok := temp.([]interface{})
				if !ok {
					return false
				}
				for _, val := range array {
					value := val.(map[string]interface{})
					if Validate(value, res) {
						return true
					}
				}
				return false
			}

			val, p := res[k]
			if !p {
				return false
			}

			// And clause
			cond, ok := temp.(map[string]interface{})
			if !ok {
				temp, val = adjustValTypes(temp, val)
				if !compare(temp, val) {
					return false
				}
				continue
			}

			// match condition
			for k2, v2 := range cond {
				v2, val = adjustValTypes(v2, val)
				if reflect.TypeOf(val) != reflect.TypeOf(v2) {
					return false
				}
				switch k2 {
				case "$eq":
					if !compare(val, v2) {
						return false
					}
				case "$ne":
					if compare(val, v2) {
						return false
					}
				case "$gt":
					switch val2 := val.(type) {
					case string:
						v2String := v2.(string)
						if val2 <= v2String {
							return false
						}
					case int64:
						v2Int := v2.(int64)
						if val2 <= v2Int {
							return false
						}
					case float64:
						v2Float := v2.(float64)
						if val2 <= v2Float {
							return false
						}
					default:
						return false
					}
				case "$gte":
					switch val2 := val.(type) {
					case string:
						v2String := v2.(string)
						if val2 < v2String {
							return false
						}
					case int64:
						v2Int := v2.(int64)
						if val2 < v2Int {
							return false
						}
					case float64:
						v2Float := v2.(float64)
						if val2 < v2Float {
							return false
						}
					default:
						return false
					}

				case "$lt":
					switch val2 := val.(type) {
					case string:
						v2String := v2.(string)
						if val2 >= v2String {
							return false
						}
					case int64:
						v2Int := v2.(int64)
						if val2 >= v2Int {
							return false
						}
					case float64:
						v2Float := v2.(float64)
						if val2 >= v2Float {
							return false
						}
					default:
						return false
					}

				case "$lte":
					switch val2 := val.(type) {
					case string:
						v2String := v2.(string)
						if val2 > v2String {
							return false
						}
					case int64:
						v2Int := v2.(int64)
						if val2 > v2Int {
							return false
						}
					case float64:
						v2Float := v2.(float64)
						if val2 > v2Float {
							return false
						}
					default:
						return false
					}
				case "$contains":
					switch v := v2.(type) {
					case map[string]interface{}:
						// check if result contains specified contains field
						result, ok := res[k]
						if !ok {
							return false
						}
						resultObj, ok := result.(map[string]interface{})
						if !ok {
							return false
						}
						// check if the top level key values are same
						for key, value := range v {
							objValue, ok := resultObj[key]
							if ok {
								if objValue == value {
									return true
								}
							}
						}
						// no top level match found
						return false
					default:
						return false
					}
				case "$regex":
					regex := v2.(string)
					vString := val.(string)
					r, err := regexp.Compile(regex)
					if err != nil {
						log.Println("Couldn't compile regex")
						return false
					}
					return r.MatchString(vString)
				default:
					log.Printf("Invalid operator (%s) provided\n", k2)
					return false
				}
			}
		}
		return true
	}
	if resArray, ok := obj.([]interface{}); ok {
		for _, res := range resArray {
			tempObj, ok := res.(map[string]interface{})
			if !ok {
				return false
			}
			if !Validate(where, tempObj) {
				return false
			}
		}
		return true
	}
	return false
}
