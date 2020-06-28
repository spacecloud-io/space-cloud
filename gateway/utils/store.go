package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
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
		return "", fmt.Errorf("variable (%s) is of incorrect type (%s)", value, reflect.TypeOf(temp))
	}
	value = tempString
	return value, nil
}

// LoadValue loads a value from the state
func LoadValue(key string, state map[string]interface{}) (interface{}, error) {
	if key == "" {
		return nil, errors.New("Invalid key")
	}

	tempArray := splitVariable(key, '.')
	length := len(tempArray) - 1

	if length == 0 {
		return nil, errors.New("The variable does not map to internal state")
	}

	if tempArray[0] == "utils" {
		function := tempArray[1]
		pre := strings.IndexRune(function, '(')
		post := strings.IndexRune(function, ')')
		params := splitVariable(function[pre+1:len(function)-1], ',')
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
			default:
				return nil, fmt.Errorf("invalid type found for length")
			}
		}
		if strings.HasPrefix(function, "now") {
			return time.Now().UTC().Format(time.RFC3339), nil
		}
		if strings.HasPrefix(function, "addDuration") {
			params0 := strings.TrimSpace(params[0])
			params0 = strings.Trim(params0, "'")
			params1 := strings.TrimSpace(params[1])
			params1 = strings.Trim(params1, "'")

			params0, err := LoadStringIfExists(params0, state)
			if err != nil {
				return "", err
			}

			paresedtime, err := time.ParseDuration(params1)
			if err != nil {
				return "", err
			}

			param0, err := CheckParse(params0)
			if err != nil {
				return "", err
			}
			return param0.Add(paresedtime).Format(time.RFC3339), nil
		}
		if strings.HasPrefix(function, "roundUpDate") {
			params0 := strings.TrimSpace(params[0])
			params0 = strings.Trim(params0, "'")
			params1 := strings.TrimSpace(params[1])
			params1 = strings.Trim(params1, "'")

			params0, err := LoadStringIfExists(params0, state)
			if err != nil {
				return "", err
			}

			param0, err := CheckParse(params0)
			if err != nil {
				return "", err
			}

			var timeDate time.Time
			switch params1 {
			case "year":
				timeDate = time.Date(param0.Year(), time.January, 1, 0, 0, 0, 0, param0.Location())
			case "month":
				timeDate = time.Date(param0.Year(), param0.Month(), 1, 0, 0, 0, 0, param0.Location())
			case "day":
				timeDate = time.Date(param0.Year(), param0.Month(), param0.Day(), 0, 0, 0, 0, param0.Location())
			case "hour":
				timeDate = time.Date(param0.Year(), param0.Month(), param0.Day(), param0.Hour(), 0, 0, 0, param0.Location())
			case "minute":
				timeDate = time.Date(param0.Year(), param0.Month(), param0.Day(), param0.Hour(), param0.Minute(), 0, 0, param0.Location())
			case "second":
				timeDate = time.Date(param0.Year(), param0.Month(), param0.Day(), param0.Hour(), param0.Minute(), param0.Second(), 0, param0.Location())
			default:
				return nil, fmt.Errorf("invalid parameter (%s) provided in `utils.roundUp`", params1)
			}
			return timeDate.Format(time.RFC3339), nil
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
	case int32:
		return float64(v), nil
	case int:
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
		return nil, fmt.Errorf("could not convert value at key (%s) of type (%s) to object", key, reflect.TypeOf(tempObj))
	}
	return conv, nil
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

// StoreValue  -- stores a value in the provided state
func StoreValue(key string, value interface{}, state map[string]interface{}) error {
	keyArray := splitVariable(key, '.')
	length := len(keyArray) - 1
	if length == 0 {
		// return errors.New(ErrorInvalidVariable)
		return errors.New("Invalid Variable Error")
	}

	scope, present := state[keyArray[0]]
	if !present {
		return errors.New("Scope not present for given variable")
	}

	obj, ok := scope.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid type (%s) received for state", reflect.TypeOf(scope))
	}

	for i, k := range keyArray {
		if i == 0 {
			continue
		}

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
					return errors.New("Key not of type string")
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
				return errors.New("Key not of type string")
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

// DeleteValue  -- deletes a value in the provided state
func DeleteValue(key string, state map[string]interface{}) error {
	keyArray := strings.Split(key, ".")

	length := len(keyArray) - 1
	if length == 0 {
		return errors.New("invalid variable provided")
	}

	scope, present := state[keyArray[0]]
	if !present {
		return errors.New("Scope not present for given variable")
	}

	obj, ok := scope.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid type (%s) received for state", reflect.TypeOf(scope))
	}

	for i, k := range keyArray {
		if i == 0 {
			continue
		}

		if i == length {
			delete(obj, k)
			break
		}

		tempObj, present := obj[k]
		if !present {
			return errors.New("Cannot find property " + k + "of undefined")
		}

		var ok bool
		obj, ok = tempObj.(map[string]interface{})
		if !ok {
			return errors.New("The variable cannot be mapped")
		}
	}

	return nil
}

// StoreValueInObject -- stores a value in provided object
func StoreValueInObject(key string, value interface{}, obj map[string]interface{}) error {
	keyArray := strings.Split(key, ".")

	length := len(keyArray) - 1

	for i, k := range keyArray {
		if i == length {
			obj[k] = value
			break
		}

		tempObj, present := obj[k]
		if !present {
			obj[k] = make(map[string]interface{})
			obj, _ = obj[k].(map[string]interface{})
			continue
		}

		var ok bool
		obj, ok = tempObj.(map[string]interface{})
		if !ok {
			return errors.New("The variable cannot be mapped")
		}
	}

	return nil
}
