package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"

	"github.com/spaceuptech/space-cloud/modules/crud"
)

func (m *Module) matchRule(ctx context.Context, project string, rule *config.Rule, args, auth map[string]interface{}) (*PostProcess, error) {
	if project != m.project {
		return &PostProcess{}, errors.New("invalid project details provided")
	}

	if rule.Rule == "allow" || rule.Rule == "authenticated" {
		return &PostProcess{}, nil
	}

	if idTemp, p := auth["id"]; p {
		if id, ok := idTemp.(string); ok && id == utils.InternalUserID {
			return &PostProcess{}, nil
		}
	}

	switch rule.Rule {
	case "deny":
		return &PostProcess{}, ErrIncorrectMatch

	case "match":
		return &PostProcess{}, match(rule, args)

	case "and":
		return m.matchAnd(ctx, project, rule, args, auth)

	case "or":
		return m.matchOr(ctx, project, rule, args, auth)

	case "webhook":
		return &PostProcess{}, matchFunc(ctx, rule, m.makeHttpRequest, args)

	case "query":
		return &PostProcess{}, matchQuery(ctx, project, rule, m.crud, args)

	case "force":
		return matchForce(rule, args)

	case "remove":
		return matchRemove(rule, args)
	default:
		return &PostProcess{}, ErrIncorrectMatch
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

func (m *Module) matchAnd(ctx context.Context, projectID string, rule *config.Rule, args, auth map[string]interface{}) (*PostProcess, error) {
	completeAction := &PostProcess{}
	for _, r := range rule.Clauses {
		postProcess, err := m.matchRule(ctx, projectID, r, args, auth)
		// if err is not nil then return error without checking the other clauses.
		if err != nil {
			return &PostProcess{}, err
		}
		completeAction.postProcessAction = append(completeAction.postProcessAction, postProcess.postProcessAction...)
	}
	return completeAction, nil
	// return nil
}

func (m *Module) matchOr(ctx context.Context, projectID string, rule *config.Rule, args, auth map[string]interface{}) (*PostProcess, error) {
	//append all parameters returned by all clauses! and then return mainStruct
	for _, r := range rule.Clauses {
		postProcess, err := m.matchRule(ctx, projectID, r, args, auth)
		if err == nil {
			//if condition is satisfied -> exit the function
			return postProcess, nil
		}
	}
	//if condition is not satisfied -> return empty PostProcess and error
	return &PostProcess{}, ErrIncorrectMatch
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

func matchForce(rule *config.Rule, args map[string]interface{}) (*PostProcess, error) {
	value := rule.Value
	if stringValue, ok := rule.Value.(string); ok {
		loadedValue, err := utils.LoadValue(stringValue, args)
		if err == nil {
			value = loadedValue
		}
	}
	//"res" - add to structure for post processing || "args" - store in args
	if strings.HasPrefix(rule.Field, "res") {
		addToStruct := PostProcessAction{Action: "force", Field: rule.Field, Value: value}
		return &PostProcess{postProcessAction: []PostProcessAction{addToStruct}}, nil
	} else if strings.HasPrefix(rule.Field, "args") {
		err := utils.StoreValue(rule.Field, value, args)
		return &PostProcess{}, err
	} else {
		return nil, ErrIncorrectRuleFieldType
	}
}

func matchRemove(rule *config.Rule, args map[string]interface{}) (*PostProcess, error) {
	actions := &PostProcess{}
	for _, field := range rule.Fields {
		//"res" - add field to structure for post processing || "args" - delete field from args
		if strings.HasPrefix(field, "res") {
			addToStruct := PostProcessAction{Action: "remove", Field: field, Value: nil}
			actions.postProcessAction = append(actions.postProcessAction, addToStruct)
		} else if strings.HasPrefix(field, "args") {
			// Since it depends on the request itself, delete the field from args
			if err := utils.DeleteValue(field, args); err != nil {
				return nil, err
			}
		} else {
			return nil, ErrIncorrectRuleFieldType
		}
	}
	return actions, nil
}
