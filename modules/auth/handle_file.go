package auth

import (
	"strings"
	"os"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/utils"
)

// IsFileOpAuthorised checks if the caller is authorized to make the request
func (m *Module) IsFileOpAuthorised(project, token, path string, op utils.FileOpType, args map[string]interface{}) error {
	m.RLock()
	defer m.RUnlock()

	// Get the rules corresponding to the requested path
	params, rules, err := m.getFileRule(path)
	if err != nil {
		return err
	}
	rule := rules.Rule[string(op)]
	if rule.Rule == "allow" {
		return nil
	}

	auth, err := m.parseToken(token)
	if err != nil {
		return err
	}

	// Add the path params and auth object to the arguments list
	args["params"] = params
	args["auth"] = auth

	// Match the rule
	return m.matchRule(project, rule, map[string]interface{}{"args": args})
}

func (m *Module) getFileRule(path string) (map[string]interface{}, *config.FileRule, error) {
	pathParams := make(map[string]interface{})

	in1 := strings.Split(path, string(os.PathSeparator))
	// Remove last element if it is  Empty
	if in1[len(in1)-1] == "" {
		in1 = in1[:len(in1)-1]
	}

	for _, r := range m.fileRules {

		rulePath := strings.Split(r.Prefix, string(os.PathSeparator))

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
