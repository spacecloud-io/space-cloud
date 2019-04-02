package auth

import (
	"context"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"

	"github.com/spaceuptech/space-cloud/modules/crud"
)

func (m *Module) matchRule(rule *config.Rule, args map[string]interface{}) error {
	switch rule.Rule {
	case "allow", "authenticated":
		return nil

	case "deny":
		return ErrIncorrectMatch

	case "match":
		return match(rule, args)

	case "and":
		return matchAnd(rule, args)

	case "or":
		return matchOr(rule, args)

	case "query":
		return matchQuery(rule, m.crud, args)

	default:
		return ErrIncorrectMatch
	}
}

func matchQuery(rule *config.Rule, crud *crud.Module, args map[string]interface{}) error {
	// Adjust the find object to load any variables referenced from state
	rule.Find = utils.Adjust(rule.Find, args).(map[string]interface{})

	// Create a new read request
	req := &model.ReadRequest{Find: rule.Find, Operation: utils.One}

	// Execute the read request
	_, err := crud.Read(context.TODO(), rule.DB, args["project"].(string), rule.Col, req)
	return err
}

func matchAnd(rule *config.Rule, args map[string]interface{}) error {
	for _, r := range rule.Clauses {
		err := match(r, args)
		if err != nil {
			return err
		}
	}

	return nil
}

func matchOr(rule *config.Rule, args map[string]interface{}) error {
	for _, r := range rule.Clauses {
		err := match(r, args)
		if err == nil {
			return nil
		}
	}

	return ErrIncorrectMatch
}

func match(rule *config.Rule, args map[string]interface{}) error {
	switch rule.Type {
	case "string":
		return matchString(rule, args)

	case "number":
		return matchNumber(rule, args)

	case "bool":
		return matchBool(rule, args)
	}

	return ErrIncorrectMatch
}
