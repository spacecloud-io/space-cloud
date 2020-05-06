package auth

import (
	"fmt"
	"strings"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func matchString(rule *config.Rule, args map[string]interface{}) error {
	var f2String []interface{}
	var f2 string
	f1String, ok := rule.F1.(string)
	if !ok {
		return formaterror(rule, ErrIncorrectRuleFieldType)
	}

	f1, err := utils.LoadStringIfExists(f1String, args)
	if err != nil {
		return formaterror(rule, err)
	}
	switch v := rule.F2.(type) {
	case string:
		if strings.HasPrefix(v, "args.") {
			temp, err := utils.LoadValue(v, args)
			if err != nil {
				return formaterror(rule, err)
			}
			f2, ok = temp.(string)
			if !ok {
				f2String, ok = temp.([]interface{})
				if !ok {
					return formaterror(rule, fmt.Errorf("invalid second field (%v) provided - wanted array of string", temp))
				}
			}

		} else {
			f2, err = utils.LoadStringIfExists(v, args)
			if err != nil {
				return formaterror(rule, err)
			}
		}
		switch rule.Eval {
		case "==":
			if f1 == f2 {
				return nil
			}

		case "!=":
			if f1 != f2 {
				return nil
			}
		case "in":
			return matchIn(rule, f2String, f1, args)
		case "notin":
			return matchNotIn(rule, f2String, f1, args)
		}
	case []interface{}:
		f2String = v
		switch rule.Eval {
		case "in":
			return matchIn(rule, f2String, f1, args)
		case "notin":
			return matchNotIn(rule, f2String, f1, args)
		}
	default:
		return formaterror(rule, ErrIncorrectRuleFieldType)
	}
	return formaterror(rule, ErrIncorrectMatch)
}

func matchIn(rule *config.Rule, f2 []interface{}, f1 interface{}, state map[string]interface{}) error {
	for _, Field2 := range f2 {
		switch v := Field2.(type) {
		case string:
			f2, err := utils.LoadStringIfExists(v, state)
			if err != nil {
				return formaterror(rule, err)
			}
			if f1 == f2 {
				return nil
			}
		case float64, int, int32, int64, float32:
			f2, err := utils.LoadNumber(v, state)
			if err != nil {
				return formaterror(rule, err)
			}
			if f1 == f2 {
				return nil
			}
		}
	}
	return formaterror(rule, ErrIncorrectMatch)
}

func matchNotIn(rule *config.Rule, f2 []interface{}, f1 interface{}, state map[string]interface{}) error {
	for _, Field2 := range f2 {
		switch v := Field2.(type) {
		case string:
			f2, err := utils.LoadStringIfExists(v, state)
			if err != nil {
				return formaterror(rule, err)
			}
			if f1 == f2 {
				return formaterror(rule, ErrIncorrectMatch)
			}
		case float64, int, int32, int64, float32:
			f2, err := utils.LoadNumber(v, state)
			if err != nil {
				return formaterror(rule, err)
			}
			if f1 == f2 {
				return formaterror(rule, ErrIncorrectMatch)
			}
		}
	}
	return nil
}

func matchNumber(rule *config.Rule, args map[string]interface{}) error {
	var f2Number []interface{}
	var f2 float64
	f1, err := utils.LoadNumber(rule.F1, args)
	if err != nil {
		return formaterror(rule, err)
	}
	f2, err = utils.LoadNumber(rule.F2, args)
	if err != nil {
		switch v := rule.F2.(type) {
		case string:
			if strings.HasPrefix(v, "args.") {
				temp, err := utils.LoadValue(rule.F2.(string), args)
				if err != nil {
					return formaterror(rule, err)
				}
				t, ok := temp.(float64)
				if !ok {
					tArr, ok := temp.([]interface{})
					if !ok {
						return formaterror(rule, fmt.Errorf("invalid second field (%v) provided - wanted array of numbers", temp))
					}
					f2Number = tArr
				}
				f2 = t
			}

		case []interface{}:
			f2Number = v
		default:
			return formaterror(rule, ErrIncorrectRuleFieldType)
		}
	}
	switch rule.Eval {
	case "==":
		if f1 == f2 {
			return nil
		}

	case "<=":
		if f1 <= f2 {
			return nil
		}

	case ">=":
		if f1 >= f2 {
			return nil
		}

	case "<":
		if f1 < f2 {
			return nil
		}

	case ">":
		if f1 > f2 {
			return nil
		}

	case "!=":
		if f1 != f2 {
			return nil
		}
	case "in":
		return matchIn(rule, f2Number, f1, args)
	case "notin":
		return matchNotIn(rule, f2Number, f1, args)
	}

	return formaterror(rule, ErrIncorrectMatch)
}

func matchBool(rule *config.Rule, args map[string]interface{}) error {

	f1, err := utils.LoadBool(rule.F1, args)
	if err != nil {
		return formaterror(rule, err)
	}

	f2, err := utils.LoadBool(rule.F2, args)
	if err != nil {
		return formaterror(rule, err)
	}

	switch rule.Eval {
	case "==":
		if f1 == f2 {
			return nil
		}
		return formaterror(rule, ErrIncorrectMatch)

	case "!=":
		if f1 != f2 {
			return nil
		}
		return formaterror(rule, ErrIncorrectMatch)
	}
	return formaterror(rule, ErrIncorrectRuleFieldType)
}
