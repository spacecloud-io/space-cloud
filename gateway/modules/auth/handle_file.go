package auth

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// IsFileOpAuthorised checks if the caller is authorized to make the request
func (m *Module) IsFileOpAuthorised(ctx context.Context, project, token, path string, op model.FileOpType, args map[string]interface{}) (*model.PostProcess, error) {
	m.RLock()
	defer m.RUnlock()

	// Get the rules corresponding to the requested path
	params, rules, err := m.getFileRule(path)
	if err != nil {
		return nil, err
	}
	rule := rules.Rule[string(op)]
	if rule.Rule == "allow" {
		if m.project == project {
			return &model.PostProcess{}, nil
		}
		return nil, errors.New("invalid project details provided")
	}

	var auth map[string]interface{}
	auth, err = m.jwt.ParseToken(ctx, token)
	if err != nil {
		return nil, err
	}

	// Add the path params and auth object to the arguments list
	args["params"] = params
	args["auth"] = auth
	args["token"] = token

	// Match the rule
	return m.matchRule(ctx, project, rule, map[string]interface{}{"args": args}, auth, model.ReturnWhereStub{})
}

func (m *Module) getFileRule(path string) (map[string]interface{}, *config.FileRule, error) {
	pathParams := make(map[string]interface{})
	ps := "/"
	if m.fileStoreType == string(utils.Local) {
		ps = string(os.PathSeparator)
	}

	// Check if its a valid absolute path
	if strings.Contains(path, "..") {
		return nil, nil, errors.New("Local: Provided path should be absolute")
	}

	in1 := strings.Split(path, ps)
	// Remove last element if it is empty
	if in1[len(in1)-1] == "" {
		in1 = in1[:len(in1)-1]
	}

	for _, r := range m.fileRules {

		rulePath := strings.Split(r.Prefix, ps)

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
