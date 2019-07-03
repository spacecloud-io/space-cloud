package utils

import (
	"errors"
	"strings"
)

// Adjust loads value from state if referenced
func Adjust(obj interface{}, state map[string]interface{}) interface{} {
	switch v := obj.(type) {
	case map[string]interface{}:
		newObj := map[string]interface{}{}
		for key, valTemp := range v {
			newObj[key] = Adjust(valTemp, state)
		}
		return newObj

	case []interface{}:
		newArray := []interface{}{}
		for _, valTemp := range v {
			newArray = append(newArray, Adjust(valTemp, state))
		}
		return newArray

	case string:
		val, err := LoadValue(v, state)
		if err == nil {
			return Adjust(val, state)
		}

		return v

	default:
		return v
	}
}

// LoadStringIfExists loads a value if its present else returns the same
func LoadStringIfExists(value string, state map[string]interface{}) string {
	if temp, err := LoadValue(value, state); err == nil {
		if tempString, ok := temp.(string); ok {
			value = tempString
		}
	}
	return value
}

// LoadValue loads a value from the state
func LoadValue(key string, state map[string]interface{}) (interface{}, error) {
	if key == "" {
		return nil, errors.New("Invalid key")
	}

	tempArray := splitVariable(key)
	length := len(tempArray) - 1

	if length == 0 {
		return nil, errors.New("The variable does not map to internal state")
	}

	if tempArray[0] == "utils" {
		function := tempArray[1]
		pre := strings.IndexRune(function, '(')
		post := strings.IndexRune(function, ')')
		if strings.HasPrefix(function, "exists") {
			_, err := LoadValue(function[pre+1:post], state)
			return err == nil, nil
		}

		return nil, errors.New("Invalid utils operation")
	}

	scope, present := state[tempArray[0]]
	if !present {
		return nil, errors.New("Scope not present")
	}

	obj, ok := scope.(map[string]interface{})
	if !ok {
		return nil, errors.New("Invalid state object")
	}

	for index, k := range tempArray {
		if index < 1 {
			continue
		}
		if index == length {
			if strings.HasSuffix(k, "]") {
				pre := strings.IndexRune(k, '[')
				post := strings.IndexRune(k, ']')
				var err error
				obj, err = convert(k[0:pre], obj)
				if err != nil {
					return nil, err
				}
				subVal, err := LoadValue(k[pre+1:post], state)
				if err != nil {
					return nil, err
				}
				subKey, ok := subVal.(string)
				if !ok {
					return nil, errors.New("Key not of type string")
				}
				value, present := obj[subKey]
				if !present {
					return nil, errors.New("Key not present in state - " + key)
				}
				return value, nil
			}

			value, present := obj[k]
			if !present {
				return nil, errors.New("Key not present in state - " + key)
			}
			return value, nil
		}
		if strings.Contains(k, "]") {
			pre := strings.IndexRune(k, '[')
			post := strings.IndexRune(k, ']')
			var err error
			obj, err = convert(k[0:pre], obj)
			if err != nil {
				return nil, err
			}

			subVal, err := LoadValue(k[pre+1:post], state)
			if err != nil {
				return nil, err
			}
			subKey, ok := subVal.(string)
			if !ok {
				return nil, errors.New("Key not of type string")
			}

			obj, err = convert(subKey, obj)
			if err != nil {
				return nil, err
			}
			continue
		}
		var err error
		obj, err = convert(k, obj)
		if err != nil {
			return nil, err
		}
	}

	return nil, errors.New("Key not found")
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
	case float64:
		return v, nil
	}

	return 0, errors.New("Store: Cloud not load value")
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

	return false, errors.New("Store: Cloud not load value")
}

func convert(key string, obj map[string]interface{}) (map[string]interface{}, error) {
	tempObj, present := obj[key]
	if !present {
		return nil, errors.New("Key not present in state (convert) - " + key)
	}
	conv, ok := tempObj.(map[string]interface{})
	if !ok {
		return nil, errors.New("Incorrect type")
	}
	return conv, nil
}

func splitVariable(key string) []string {
	var inBracket bool
	var lastIndex int
	array := []string{}
	for i, c := range key {
		if c == '[' || c == '(' {
			inBracket = true
		} else if c == ']' || c == ')' {
			inBracket = false
		} else if c == '.' && !inBracket {
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
