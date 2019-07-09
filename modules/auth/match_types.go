package auth

import (
	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/utils"
)

func matchString(rule *config.Rule, args map[string]interface{}) error {
	f1String, ok := rule.F1.(string)
	if !ok {
		return ErrIncorrectRuleFieldType
	}
	f2String, ok := rule.F2.(string)
	if !ok {
		return ErrIncorrectRuleFieldType
	}

	f1 := utils.LoadStringIfExists(f1String, args)
	f2 := utils.LoadStringIfExists(f2String, args)
	switch rule.Eval {
	case "==":
		if f1 == f2 {
			return nil
		}

	case "!=":
		if f1 != f2 {
			return nil
		}
	}

	return ErrIncorrectMatch
}

func matchNumber(rule *config.Rule, args map[string]interface{}) error {

	f1, err := utils.LoadNumber(rule.F1, args)
	if err != nil {
		return err
	}

	f2, err := utils.LoadNumber(rule.F2, args)
	if err != nil {
		return err
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
	}

	return ErrIncorrectRuleFieldType
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
		} else {
			return ErrIncorrectMatch
		}

	case "!=":
		if f1 != f2 {
			return nil
		} else {
			return ErrIncorrectMatch
		}
	}
	return ErrIncorrectRuleFieldType
}
