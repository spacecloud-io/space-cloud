package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// LoadValue loads a value from the state
func LoadValue(key string, state map[string]interface{}) (interface{}, error) {
	if key == "" {
		return nil, errors.New("invalid key")
	}

	tempArray := splitVariable(key, '.')
	length := len(tempArray) - 1

	if tempArray[0] == "utils" {
		function := tempArray[1]
		pre := strings.IndexRune(function, '(')
		post := strings.IndexRune(function, ')')
		// params := splitVariable(function[pre+1:len(function)-1], ',')

		// TODO: Fix this if support for mongo is needed
		// // Method check - stringToObjectId
		// if strings.HasPrefix(function, "stringToObjectId") {
		// 	value, err := LoadValue(function[pre+1:post], state)
		// 	if err != nil {
		// 		return nil, err
		// 	}

		// 	if v, ok := value.(primitive.A); ok {
		// 		value = []interface{}(v)
		// 	}

		// 	// The value can be a string or an array of string
		// 	switch v := value.(type) {
		// 	case primitive.ObjectID:
		// 		return v, nil
		// 	case string:
		// 		return primitive.ObjectIDFromHex(v)
		// 	case []interface{}:
		// 		array := make([]interface{}, len(v))
		// 		for i, item := range v {
		// 			s, ok := item.(string)
		// 			if !ok {
		// 				return nil, fmt.Errorf("value (%v) cannot be converted to ObjectID", item)
		// 			}

		// 			objID, err := primitive.ObjectIDFromHex(s)
		// 			if err != nil {
		// 				return nil, err
		// 			}

		// 			array[i] = objID
		// 		}
		// 		return array, nil
		// 	default:
		// 		return nil, fmt.Errorf("invalid type provided (%s) for object id conversion", reflect.TypeOf(value))
		// 	}
		// }

		// // Method check - objectIdToString
		// if strings.HasPrefix(function, "objectIdToString") {
		// 	value, err := LoadValue(function[pre+1:post], state)
		// 	if err != nil {
		// 		return nil, err
		// 	}

		// 	switch val := value.(type) {
		// 	case string:
		// 		return val, nil

		// 	case primitive.ObjectID:
		// 		return val.Hex(), nil

		// 	default:
		// 		return nil, fmt.Errorf("invalid type provided (%s) for object id conversion", reflect.TypeOf(value))
		// 	}
		// }

		if strings.HasPrefix(function, "exists") {
			_, err := LoadValue(function[pre+1:post], state)
			return err == nil, nil
		}
		if strings.HasPrefix(function, "length") {
			value, err := LoadValue(function[pre+1:post], state)
			if err != nil {
				return nil, err
			}
			switch v := value.(type) {
			case []interface{}:
				return int64(len(v)), nil
			case map[string]interface{}:
				return int64(len(v)), nil
			case string:
				return int64(len(v)), nil
			default:
				return nil, fmt.Errorf("unable to calculate length of variable '%s' - got type (%s) want object or array", key, reflect.TypeOf(value))
			}
		}
		if strings.HasPrefix(function, "now") {
			return time.Now().UTC().Format(time.RFC3339Nano), nil
		}

		// TODO: Add this if needed
		// if strings.HasPrefix(function, "addDuration") {
		// 	params0 := strings.TrimSpace(params[0])
		// 	params0 = strings.Trim(params0, "'")
		// 	params1 := strings.TrimSpace(params[1])
		// 	params1 = strings.Trim(params1, "'")

		// 	params0, err := LoadStringIfExists(params0, state)
		// 	if err != nil {
		// 		return "", err
		// 	}

		// 	paresedtime, err := time.ParseDuration(params1)
		// 	if err != nil {
		// 		return "", err
		// 	}

		// 	param0, err := CheckParse(params0)
		// 	if err != nil {
		// 		return "", err
		// 	}
		// 	return param0.Add(paresedtime).Format(time.RFC3339Nano), nil
		// }
		// if strings.HasPrefix(function, "roundUpDate") {
		// 	params0 := strings.TrimSpace(params[0])
		// 	params0 = strings.Trim(params0, "'")
		// 	params1 := strings.TrimSpace(params[1])
		// 	params1 = strings.Trim(params1, "'")

		// 	params0, err := LoadStringIfExists(params0, state)
		// 	if err != nil {
		// 		return "", err
		// 	}

		// 	param0, err := CheckParse(params0)
		// 	if err != nil {
		// 		return "", err
		// 	}

		// 	var timeDate time.Time
		// 	switch params1 {
		// 	case "year":
		// 		timeDate = time.Date(param0.Year(), time.January, 1, 0, 0, 0, 0, param0.Location())
		// 	case "month":
		// 		timeDate = time.Date(param0.Year(), param0.Month(), 1, 0, 0, 0, 0, param0.Location())
		// 	case "day":
		// 		timeDate = time.Date(param0.Year(), param0.Month(), param0.Day(), 0, 0, 0, 0, param0.Location())
		// 	case "hour":
		// 		timeDate = time.Date(param0.Year(), param0.Month(), param0.Day(), param0.Hour(), 0, 0, 0, param0.Location())
		// 	case "minute":
		// 		timeDate = time.Date(param0.Year(), param0.Month(), param0.Day(), param0.Hour(), param0.Minute(), 0, 0, param0.Location())
		// 	case "second":
		// 		timeDate = time.Date(param0.Year(), param0.Month(), param0.Day(), param0.Hour(), param0.Minute(), param0.Second(), 0, param0.Location())
		// 	default:
		// 		return nil, fmt.Errorf("Invalid parameter (%s) provided for space cloud internal function (utils.roundUpDate)", params1)
		// 	}
		// 	return timeDate.Format(time.RFC3339Nano), nil
		// }

		return nil, errors.New("invalid utils operation")
	}

	var obj interface{}
	obj = state

	for index, k := range tempArray {
		if strings.Contains(k, "]") {
			pre := strings.IndexRune(k, '[')
			post := strings.IndexRune(k, ']')
			var err error
			obj, err = getValue(k[0:pre], obj)
			if err != nil {
				return nil, err
			}

			// Load the value within brackets
			subVal, err := LoadValue(k[pre+1:post], state)
			if err != nil {
				return nil, err
			}

			// Get the key value
			switch v := subVal.(type) {
			case int64, float64, int, float32:
				k = fmt.Sprintf("%v", v)
			case string:
				k = v
			default:
				return nil, fmt.Errorf("key (%s) is of unknown type", reflect.TypeOf(subVal))
			}
		}

		var err error
		obj, err = getValue(k, obj)
		if err != nil {
			return nil, err
		}

		// If we are at the final element, it means we need to return that value
		if index == length {
			return obj, nil
		}
	}

	return nil, errors.New("key not found")
}

// LoadStringIfExists loads a value if its present else returns the same
func LoadStringIfExists(value string, state map[string]interface{}) (string, error) {
	if !strings.HasPrefix(value, "args.") && !strings.HasPrefix(value, "utils.") {
		return value, nil
	}
	temp, err := LoadValue(value, state)
	if err != nil {
		return "", err
	}
	tempString, ok := temp.(string)
	if !ok {
		return "", fmt.Errorf("unable to load variable '%s' - got '%s'; want 'string'", value, reflect.TypeOf(temp))
	}
	value = tempString
	return value, nil
}

// LoadNumber loads a key as a float. Throws error
func LoadNumber(key interface{}, args map[string]interface{}) (float64, error) {
	// Create a temporary copy of key
	temp := key

	// Load value from argument if key was string i.e. it points to a variable in the argument
	if tempString, ok := key.(string); ok {
		val, err := LoadValue(tempString, args)
		if err != nil {
			return 0, err
		}
		temp = val
	}

	switch v := temp.(type) {
	case int64:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case float64:
		return v, nil
	}

	return 0, errors.New("store: Cloud not load value")
}

// LoadBool loads a key as a float. Throws error
func LoadBool(key interface{}, args map[string]interface{}) (bool, error) {
	// Create a temporary copy of key
	temp := key

	// Load value from argument if key was string i.e. it points to a variable in the argument
	if tempString, ok := key.(string); ok {
		val, err := LoadValue(tempString, args)
		if err != nil {
			return false, err
		}
		temp = val
	}

	if v, ok := temp.(bool); ok {
		return v, nil
	}

	return false, errors.New("store: Cloud not load value")
}

// StoreValue stores a value in the provided state
func StoreValue(key string, value interface{}, state map[string]interface{}) error {
	keyArray := splitVariable(key, '.')
	length := len(keyArray) - 1

	obj := state

	for i, k := range keyArray {
		if i == length {
			if strings.HasSuffix(k, "]") {
				pre := strings.IndexRune(k, '[')
				post := strings.IndexRune(k, ']')

				var err error
				obj, err = convertOrCreate(k[0:pre], obj)
				if err != nil {
					return err
				}

				subVal, err := LoadValue(k[pre+1:post], state)
				if err != nil {
					return err
				}
				subKey, ok := subVal.(string)
				if !ok {
					return errors.New("key not of type string")
				}

				obj[subKey] = value
				return nil
			}
			obj[k] = value
			return nil
		}
		if strings.HasSuffix(k, "]") {
			pre := strings.IndexRune(k, '[')
			post := strings.IndexRune(k, ']')

			var err error
			obj, err = convertOrCreate(k[0:pre], obj)
			if err != nil {
				return err
			}

			subVal, err := LoadValue(k[pre+1:post], state)
			if err != nil {
				return err
			}
			subKey, ok := subVal.(string)
			if !ok {
				return errors.New("key not of type string")
			}

			obj, err = convertOrCreate(subKey, obj)
			if err != nil {
				return err
			}
			continue
		}
		var err error
		obj, err = convertOrCreate(k, obj)
		if err != nil {
			return err
		}
	}

	return nil
}

func getValue(key string, obj interface{}) (interface{}, error) {
	switch val := obj.(type) {
	case []interface{}:
		// The key should be a number (index) if the object is an array
		index, err := strconv.Atoi(key)
		if err != nil {
			return nil, fmt.Errorf("key (%s) provided instead of index", key)
		}

		// Check if index is not out of bounds otherwise return value at that index
		if index >= len(val) {
			return nil, fmt.Errorf("index (%d) out of bounds", index)
		}
		return val[index], nil

	// TODO: Add this if support for mongo is needed
	// case primitive.A:
	// 	// The key should be a number (index) if the object is an array
	// 	index, err := strconv.Atoi(key)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("Key (%s) provided instead of index", key), err, nil)
	// 	}

	// 	// Check if index is not out of bounds otherwise return value at that index
	// 	if index >= len(val) {
	// 		return nil, fmt.Errorf("Index (%d) out of bounds", index)
	// 	}
	// 	return val[index], nil

	case map[string]interface{}:
		// Throw error if key is not present in state. Otherwise return value
		tempObj, p := val[key]
		if !p {
			return nil, fmt.Errorf("key (%s) not present in state", key)
		}
		return tempObj, nil

	default:
		return nil, fmt.Errorf("unsupported data type (%s)", reflect.TypeOf(obj))
	}
}

func splitVariable(key string, delimiter rune) []string {
	var inBracket1 int
	var inBracket2 int

	var lastIndex int
	array := []string{}
	for i, c := range key {
		if c == '[' {
			inBracket1++
		}
		if c == '(' {
			inBracket2++
		}
		if c == ']' {
			inBracket1--
		}
		if c == ')' {
			inBracket2--
		}
		if c == delimiter && inBracket1 == 0 && inBracket2 == 0 {
			sub := key[lastIndex:i]
			array = append(array, sub)
			lastIndex = i + 1
		}
		if i == len(key)-1 {
			sub := key[lastIndex : i+1]
			array = append(array, sub)
		}
	}
	return array
}

func convertOrCreate(k string, obj map[string]interface{}) (map[string]interface{}, error) {
	tempObj, present := obj[k]
	if !present {
		tempObj = make(map[string]interface{})
		obj[k] = tempObj
	}

	var ok bool
	obj2, ok := tempObj.(map[string]interface{})
	if !ok {
		return nil, errors.New("the variable cannot be mapped")
	}
	return obj2, nil
}
