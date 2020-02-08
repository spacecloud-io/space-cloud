package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// IsCreateOpAuthorised checks if the crud operation is authorised
func (m *Module) IsCreateOpAuthorised(ctx context.Context, project, dbType, col, token string, req *model.CreateRequest) (int, error) {
	m.RLock()
	defer m.RUnlock()

	rule, auth, err := m.authenticateCrudRequest(dbType, col, token, utils.Create)
	if err != nil {
		return http.StatusUnauthorized, err
	}

	args := map[string]interface{}{"op": req.Operation, "auth": auth, "token": token}

	var rows []interface{}
	switch req.Operation {
	case utils.One:
		rows = []interface{}{req.Document}
	case utils.All:
		rows = req.Document.([]interface{})
	default:
		rows = []interface{}{}
	}

	for _, row := range rows {
		args["doc"] = row
		_, err := m.matchRule(ctx, project, rule, map[string]interface{}{"args": args}, auth)
		if err != nil {
			return http.StatusForbidden, err
		}
	}

	if err := m.schema.ValidateCreateOperation(dbType, col, req); err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}

// IsReadOpAuthorised checks if the crud operation is authorised
func (m *Module) IsReadOpAuthorised(ctx context.Context, project, dbType, col, token string, req *model.ReadRequest) (*PostProcess, int, error) {
	m.RLock()
	defer m.RUnlock()

	rule, auth, err := m.authenticateCrudRequest(dbType, col, token, utils.Read)
	if err != nil {
		return &PostProcess{}, http.StatusUnauthorized, err
	}

	args := map[string]interface{}{"op": req.Operation, "auth": auth, "find": req.Find, "token": token}
	actions, err := m.matchRule(ctx, project, rule, map[string]interface{}{"args": args}, auth)
	if err != nil {
		return &PostProcess{}, http.StatusForbidden, err
	}

	return actions, http.StatusOK, nil
}

// IsUpdateOpAuthorised checks if the crud operation is authorised
func (m *Module) IsUpdateOpAuthorised(ctx context.Context, project, dbType, col, token string, req *model.UpdateRequest) (int, error) {
	m.RLock()
	defer m.RUnlock()

	rule, auth, err := m.authenticateCrudRequest(dbType, col, token, utils.Update)
	if err != nil {
		return http.StatusUnauthorized, err
	}

	args := map[string]interface{}{"op": req.Operation, "auth": auth, "find": req.Find, "update": req.Update, "token": token}
	_, err = m.matchRule(ctx, project, rule, map[string]interface{}{"args": args}, auth)
	if err != nil {
		return http.StatusForbidden, err
	}

	if err := m.schema.ValidateUpdateOperation(dbType, col, req.Operation, req.Update, req.Find); err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}

// IsDeleteOpAuthorised checks if the crud operation is authorised
func (m *Module) IsDeleteOpAuthorised(ctx context.Context, project, dbType, col, token string, req *model.DeleteRequest) (int, error) {
	m.RLock()
	defer m.RUnlock()

	rule, auth, err := m.authenticateCrudRequest(dbType, col, token, utils.Delete)
	if err != nil {
		return http.StatusUnauthorized, err
	}

	args := map[string]interface{}{"op": req.Operation, "auth": auth, "find": req.Find, "token": token}
	_, err = m.matchRule(ctx, project, rule, map[string]interface{}{"args": args}, auth)
	if err != nil {
		return http.StatusForbidden, err
	}

	return http.StatusOK, nil
}

// IsAggregateOpAuthorised checks if the crud operation is authorised
func (m *Module) IsAggregateOpAuthorised(ctx context.Context, project, dbType, col, token string, req *model.AggregateRequest) (int, error) {
	m.RLock()
	defer m.RUnlock()

	rule, auth, err := m.authenticateCrudRequest(dbType, col, token, utils.Aggregation)
	if err != nil {
		return http.StatusUnauthorized, err
	}

	args := map[string]interface{}{"op": req.Operation, "auth": auth, "pipeline": req.Pipeline, "token": token}
	_, err = m.matchRule(ctx, project, rule, map[string]interface{}{"args": args}, auth)
	if err != nil {
		return http.StatusForbidden, err
	}

	return http.StatusOK, nil
}

// PostProcessMethod to do processing on result
func (m *Module) PostProcessMethod(postProcess *PostProcess, result interface{}) error {
	// Gracefully exist if the result is nil
	if result == nil {
		return nil
	}

	// convert to array of interfaces
	var resultArr []interface{}
	switch val := result.(type) {
	case map[string]interface{}:
		resultArr = []interface{}{val} //make an array of interface with val element
	case []interface{}:
		resultArr = val
	default:
		return errors.New("result is of invalid type")
	}

	for _, doc := range resultArr {
		for _, field := range postProcess.postProcessAction {
			// apply Action on all elements
			switch field.Action {
			case "force":
				if err := utils.StoreValue(field.Field, field.Value, map[string]interface{}{"res": doc}); err != nil {
					return err
				}

			case "remove":
				if err := utils.DeleteValue(field.Field, map[string]interface{}{"res": doc}); err != nil {
					return err
				}
			default:
				err := fmt.Errorf("invalid action (%s) received in post processing read op", field.Action)
				return err
			}
		}
	}
	return nil
}

func (m *Module) authenticateCrudRequest(dbType, col, token string, op utils.OperationType) (rule *config.Rule, auth map[string]interface{}, err error) {
	// Get rule
	rule, err = m.getCrudRule(dbType, col, op)
	if err != nil {
		return
	}

	// Return if rule is allow
	if rule.Rule == "allow" {
		return
	}

	// Parse token
	auth, err = m.parseToken(token)
	return
}

func (m *Module) getCrudRule(dbType, col string, query utils.OperationType) (*config.Rule, error) {
	if dbRules, p1 := m.rules[dbType]; p1 {
		if collection, p2 := dbRules.Collections[col]; p2 {
			if rule, p3 := collection.Rules[string(query)]; p3 {
				return rule, nil
			}
		}
		if defaultCol, p2 := dbRules.Collections["default"]; p2 {
			if rule, p3 := defaultCol.Rules[string(query)]; p3 {
				return rule, nil
			}
		}
	}
	return nil, fmt.Errorf("no rule found for collection %s in database %s", col, dbType)
}
