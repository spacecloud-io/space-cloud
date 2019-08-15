package auth

import (
	"errors"
	"strings"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/utils"
)

// isPubsubAuthorised checks if the caller is authorized to make the pubsub request
func (m *Module) isPubsubAuthorised(project, token, subject, queue string, op utils.PubsubType, args map[string]interface{}) error {
	m.RLock()
	defer m.RUnlock()

	// Get the rules corresponding to the requested path
	params, rules, err := m.getPubsubRule(subject)
	if err != nil {
		return err
	}
	rule := rules.Rule[string(op)]
	if rule.Rule == "allow" {
		if m.project == project {
			return nil
		}
		return errors.New("invalid project details provided")
	}

	auth, err := m.parseToken(token)
	if err != nil {
		return err
	}

	// Add the path params and auth object to the arguments list
	args["params"] = params
	args["auth"] = auth

	// Match the rule
	return m.matchRule(project, rule, map[string]interface{}{"args": args}, auth)
}

// IsSubscribeAuthorised checks if the caller is authorized to make the subscribe request
func (m *Module) IsSubscribeAuthorised(project, token, subject, queue string, args map[string]interface{}) error {
	return m.isPubsubAuthorised(project, token, subject, queue, utils.Subscribe, args)
}

// IsPublishAuthorised checks if the caller is authorized to make the publish request
func (m *Module) IsPublishAuthorised(project, token, subject string, args map[string]interface{}) error {
	return m.isPubsubAuthorised(project, token, subject, "", utils.Publish, args)
}

// getPubsubRule gets the pubsub rule of a particular subject
func (m *Module) getPubsubRule(subject string) (map[string]interface{}, *config.PubsubRule, error) {
	pathParams := make(map[string]interface{})
	
	in1 := strings.Split(subject, "/")
	// Remove last element if it is empty
	if in1[len(in1)-1] == "" {
		in1 = in1[:len(in1)-1]
	}

	for _, r := range m.pubsubRules {

		rulePath := strings.Split(r.Subject, "/")

		if rulePath[len(rulePath)-1] == "" {
			rulePath = rulePath[:len(rulePath)-1]
		}

		if len(in1) < len(rulePath) {
			continue
		}
		// Create a match flag
		validMatch := true

		for i, p := range rulePath {
			// Check if path segment is a variable
			if !strings.HasPrefix(p, ":") {

				// Break the current rule since its an invalid match
				if p != in1[i] {
					validMatch = false
					break
				}
				continue
			}

			// Store the path variable
			varName := strings.TrimPrefix(p, ":")
			pathParams[varName] = in1[i]
		}

		if validMatch {
			return pathParams, r, nil
		}
	}

	return nil, nil, ErrRuleNotFound
}
