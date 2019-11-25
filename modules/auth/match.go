package auth

import (
	"context"
	"errors"
	"time"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"

	"github.com/spaceuptech/space-cloud/modules/crud"
)

func (m *Module) matchRule(ctx context.Context, project string, rule *config.Rule, args map[string]interface{}, auth map[string]interface{}) error {
	if project != m.project {
		return errors.New("invalid project details provided")
	}

	if rule.Rule == "allow" || rule.Rule == "authenticated" {
		return nil
	}

	if idTemp, p := auth["id"]; p {
		if id, ok := idTemp.(string); ok && id == utils.InternalUserID {
			return nil
		}
	}

	switch rule.Rule {
	case "deny":
		return ErrIncorrectMatch

	case "match":
		return match(rule, args)

	case "and":
		return matchAnd(rule, args)

	case "or":
		return matchOr(rule, args)

	case "webhook":
		return matchFunc(ctx, rule, m.makeHttpRequest, args)

	case "query":
		return matchQuery(ctx, project, rule, m.crud, args)

	default:
		return ErrIncorrectMatch
	}
}

func matchFunc(ctx context.Context, rule *config.Rule, MakeHttpRequest utils.MakeHttpRequest, args map[string]interface{}) error {
	obj := args["args"].(map[string]interface{})
	token := obj["token"].(string)
	delete(obj, "token")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result interface{}
	return MakeHttpRequest(ctx, "POST", rule.Url, token, obj, &result)
}

func matchQuery(ctx context.Context, project string, rule *config.Rule, crud *crud.Module, args map[string]interface{}) error {
	// Adjust the find object to load any variables referenced from state
	rule.Find = utils.Adjust(rule.Find, args).(map[string]interface{})

	// Create a new read request
	req := &model.ReadRequest{Find: rule.Find, Operation: utils.One}

	// Execute the read request
	_, err := crud.Read(ctx, rule.DB, project, rule.Col, req)
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
