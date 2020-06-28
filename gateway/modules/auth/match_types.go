package auth

import (
	"errors"
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
		return ErrIncorrectRuleFieldType
	}

	f1, err := utils.LoadStringIfExists(f1String, args)
	if err != nil {
		return err
	}
	switch v := rule.F2.(type) {
	case string:
		if strings.HasPrefix(v, "args.") {
			temp, err := utils.LoadValue(v, args)
			if err != nil {
				return err
			}
			f2, ok = temp.(string)
			if !ok {
				f2String, ok = temp.([]interface{})
				if !ok {
					return fmt.Errorf("invalid second field (%v) provided - wanted array of string", temp)
				}
			}

		} else {
			f2, err = utils.LoadStringIfExists(v, args)
			if err != nil {
				return err
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
			return matchIn(f2String, f1, args)
		case "notin":
			return matchNotIn(f2String, f1, args)
		}
	case []interface{}:
		f2String = v
		switch rule.Eval {
		case "in":
			return matchIn(f2String, f1, args)
		case "notin":
			return matchNotIn(f2String, f1, args)
		}
	default:
		return ErrIncorrectRuleFieldType
	}
	return ErrIncorrectMatch
}

func matchIn(f2 []interface{}, f1 interface{}, state map[string]interface{}) error {
	for _, Field2 := range f2 {
		switch v := Field2.(type) {
		case string:
			f2, err := utils.LoadStringIfExists(v, state)
			if err != nil {
				return err
			}
			if f1 == f2 {
				return nil
			}
		case float64, int, int32, int64, float32:
			f2, err := utils.LoadNumber(v, state)
			if err != nil {
				return err
			}
			if f1 == f2 {
				return nil
			}
		}
	}
	return ErrIncorrectMatch
}

func matchNotIn(f2 []interface{}, f1 interface{}, state map[string]interface{}) error {
	for _, Field2 := range f2 {
		switch v := Field2.(type) {
		case string:
			f2, err := utils.LoadStringIfExists(v, state)
			if err != nil {
				return err
			}
			if f1 == f2 {
				return ErrIncorrectMatch
			}
		case float64, int, int32, int64, float32:
			f2, err := utils.LoadNumber(v, state)
			if err != nil {
				return err
			}
			if f1 == f2 {
				return ErrIncorrectMatch
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
		return err
	}
	f2, err = utils.LoadNumber(rule.F2, args)
	if err != nil {
		switch v := rule.F2.(type) {
		case string:
			if strings.HasPrefix(v, "args.") {
				temp, err := utils.LoadValue(rule.F2.(string), args)
				if err != nil {
					return err
				}
				t, ok := temp.(float64)
				if !ok {
					tArr, ok := temp.([]interface{})
					if !ok {
						return fmt.Errorf("invalid second field (%v) provided - wanted array of numbers", temp)
					}
					f2Number = tArr
				}
				f2 = t
			}

		case []interface{}:
			f2Number = v
		default:
			return ErrIncorrectRuleFieldType
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
		return matchIn(f2Number, f1, args)
	case "notin":
		return matchNotIn(f2Number, f1, args)
	}

	return ErrIncorrectMatch
}

func matchBool(rule *config.Rule, args map[string]interface{}) error {

	f1, err := utils.LoadBool(rule.F1, args)
	if err != nil {
		return err
	}

	f2, err := utils.LoadBool(rule.F2, args)
	if err != nil {
		return err
	}

	switch rule.Eval {
	case "==":
		if f1 == f2 {
			return nil
		}
		return ErrIncorrectMatch

	case "!=":
		if f1 != f2 {
			return nil
		}
		return ErrIncorrectMatch
	}
	return ErrIncorrectRuleFieldType
}

func matchdate(rule *config.Rule, args map[string]interface{}) error {
	f1String, ok := rule.F1.(string)
	if !ok {
		return ErrIncorrectRuleFieldType
	}

	f1String, err := utils.LoadStringIfExists(f1String, args)
	if err != nil {
		return err
	}

	f2String, ok := rule.F2.(string)
	if !ok {
		return ErrIncorrectRuleFieldType
	}

	f2String, err = utils.LoadStringIfExists(f2String, args)
	if err != nil {
		return err
	}

	f1, err := utils.CheckParse(f1String)
	if err != nil {
		return err
	}

	f2, err := utils.CheckParse(f2String)
	if err != nil {
		return err
	}
	switch rule.Eval {
	case "==":
		if f1.Equal(f2) {
			return nil
		}

	case "<=":
		if f1.Before(f2) || f1.Equal(f2) {
			return nil
		}

	case ">=":
		if f1.After(f2) || f1.Equal(f2) {
			return nil
		}

	case "<":
		if f1.Before(f2) {
			return nil
		}

	case ">":
		if f1.After(f2) {
			return nil
		}

	case "!=":
		if !f1.Equal(f2) {
			return nil
		}
	}
	return errors.New("date match failed")
}
