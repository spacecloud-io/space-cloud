package auth

import (
	"context"

	"github.com/spaceuptech/space-cloud/crud"
	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

func matchQuery(rule *Rule, crud *crud.Module, args map[string]interface{}) error {
	// Adjust the find object to load any variables referenced from state
	rule.Find = utils.Adjust(rule.Find, args).(map[string]interface{})

	// Create a new read request
	req := &model.ReadRequest{Find: rule.Find, Operation: utils.One}

	// Execute the read request
	_, err := crud.Read(context.TODO(), rule.DbType, args["project"].(string), rule.Col, req)
	return err
}

func matchAnd(rule *Rule, args map[string]interface{}) error {
	for _, r := range rule.Clauses {
		err := match(r, args)
		if err != nil {
			return err
		}
	}

	return nil
}

func matchOr(rule *Rule, args map[string]interface{}) error {
	for _, r := range rule.Clauses {
		err := match(r, args)
		if err == nil {
			return nil
		}
	}

	return ErrIncorrectMatch
}

func match(rule *Rule, args map[string]interface{}) error {
	switch rule.FieldType {
	case "string":
		return matchString(rule, args)

	case "number":
		return matchNumber(rule, args)

	case "bool":
		return matchNumber(rule, args)
	}

	return ErrIncorrectMatch
}
