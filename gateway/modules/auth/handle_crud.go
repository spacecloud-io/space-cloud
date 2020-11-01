package auth

import (
	"context"
	"fmt"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// IsCreateOpAuthorised checks if the crud operation is authorised
func (m *Module) IsCreateOpAuthorised(ctx context.Context, project, dbAlias, col, token string, req *model.CreateRequest) (model.RequestParams, error) {
	m.RLock()
	defer m.RUnlock()

	rule, auth, err := m.authenticateCrudRequest(ctx, project, dbAlias, col, token, model.Create)
	if err != nil {
		return model.RequestParams{}, err
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
		_, err := m.matchRule(ctx, project, rule, map[string]interface{}{"args": args}, auth, model.ReturnWhereStub{})
		if err != nil {
			return model.RequestParams{}, err
		}
	}

	attr := map[string]string{"project": project, "db": dbAlias, "col": col}
	return model.RequestParams{Claims: auth, Resource: "db-create", Op: "access", Attributes: attr}, nil
}

// IsReadOpAuthorised checks if the crud operation is authorised
func (m *Module) IsReadOpAuthorised(ctx context.Context, project, dbAlias, col, token string, req *model.ReadRequest, stub model.ReturnWhereStub) (*model.PostProcess, model.RequestParams, error) {
	m.RLock()
	defer m.RUnlock()

	rule, auth, err := m.authenticateCrudRequest(ctx, project, dbAlias, col, token, model.Read)
	if err != nil {
		return nil, model.RequestParams{}, err
	}

	opts := map[string]interface{}{}
	if req.Options != nil {
		if req.Options.Limit != nil {
			opts["limit"] = *req.Options.Limit
		}
		if req.Options.Skip != nil {
			opts["skip"] = *req.Options.Skip
		}
	}
	args := map[string]interface{}{"op": req.Operation, "auth": auth, "find": req.Find, "token": token, "opts": opts}
	actions, err := m.matchRule(ctx, project, rule, map[string]interface{}{"args": args}, auth, stub)
	if err != nil {
		return nil, model.RequestParams{}, err
	}

	attr := map[string]string{"project": project, "db": dbAlias, "col": col}
	return actions, model.RequestParams{Claims: auth, Resource: "db-read", Op: "access", Attributes: attr}, nil
}

// IsUpdateOpAuthorised checks if the crud operation is authorised
func (m *Module) IsUpdateOpAuthorised(ctx context.Context, project, dbAlias, col, token string, req *model.UpdateRequest) (model.RequestParams, error) {
	m.RLock()
	defer m.RUnlock()

	rule, auth, err := m.authenticateCrudRequest(ctx, project, dbAlias, col, token, model.Update)
	if err != nil {
		return model.RequestParams{}, err
	}

	args := map[string]interface{}{"op": req.Operation, "auth": auth, "find": req.Find, "update": req.Update, "token": token}
	_, err = m.matchRule(ctx, project, rule, map[string]interface{}{"args": args}, auth, model.ReturnWhereStub{})
	if err != nil {
		return model.RequestParams{}, err
	}

	attr := map[string]string{"project": project, "db": dbAlias, "col": col}
	return model.RequestParams{Claims: auth, Resource: "db-update", Op: "access", Attributes: attr}, nil
}

// IsDeleteOpAuthorised checks if the crud operation is authorised
func (m *Module) IsDeleteOpAuthorised(ctx context.Context, project, dbAlias, col, token string, req *model.DeleteRequest) (model.RequestParams, error) {
	m.RLock()
	defer m.RUnlock()

	rule, auth, err := m.authenticateCrudRequest(ctx, project, dbAlias, col, token, model.Delete)
	if err != nil {
		return model.RequestParams{}, err
	}

	args := map[string]interface{}{"op": req.Operation, "auth": auth, "find": req.Find, "token": token}
	_, err = m.matchRule(ctx, project, rule, map[string]interface{}{"args": args}, auth, model.ReturnWhereStub{})
	if err != nil {
		return model.RequestParams{}, err
	}

	attr := map[string]string{"project": project, "db": dbAlias, "col": col}
	return model.RequestParams{Claims: auth, Resource: "db-delete", Op: "access", Attributes: attr}, nil
}

// IsAggregateOpAuthorised checks if the crud operation is authorised
func (m *Module) IsAggregateOpAuthorised(ctx context.Context, project, dbAlias, col, token string, req *model.AggregateRequest) (model.RequestParams, error) {
	m.RLock()
	defer m.RUnlock()

	rule, auth, err := m.authenticateCrudRequest(ctx, project, dbAlias, col, token, model.Aggregation)
	if err != nil {
		return model.RequestParams{}, err
	}

	args := map[string]interface{}{"op": req.Operation, "auth": auth, "pipeline": req.Pipeline, "token": token}
	_, err = m.matchRule(ctx, project, rule, map[string]interface{}{"args": args}, auth, model.ReturnWhereStub{})
	if err != nil {
		return model.RequestParams{}, err
	}

	attr := map[string]string{"project": project, "db": dbAlias, "col": col}
	return model.RequestParams{Claims: auth, Resource: "db-aggregate", Op: "access", Attributes: attr}, nil
}

// IsPreparedQueryAuthorised checks if the crud operation is authorised
func (m *Module) IsPreparedQueryAuthorised(ctx context.Context, project, dbAlias, id, token string, req *model.PreparedQueryRequest) (*model.PostProcess, model.RequestParams, error) {
	m.RLock()
	defer m.RUnlock()

	rule, auth, err := m.authenticatePreparedQueryRequest(ctx, project, dbAlias, id, token)
	if err != nil {
		return nil, model.RequestParams{}, err
	}

	args := map[string]interface{}{"auth": auth, "params": req.Params, "token": token}
	actions, err := m.matchRule(ctx, project, rule, map[string]interface{}{"args": args}, auth, model.ReturnWhereStub{})
	if err != nil {
		return nil, model.RequestParams{}, err
	}

	attr := map[string]string{"project": project, "db": dbAlias}
	return actions, model.RequestParams{Claims: auth, Resource: "db-prepared-query", Op: "access", Attributes: attr}, nil
}

func (m *Module) authenticateCrudRequest(ctx context.Context, projectID, dbAlias, col, token string, op model.OperationType) (rule *config.Rule, auth map[string]interface{}, err error) {
	// Get rule
	rule, err = m.getCrudRule(ctx, projectID, dbAlias, col, op)
	if err != nil {
		return
	}

	// Return if rule is allow
	if rule.Rule == "allow" {
		return
	}

	// Parse token
	auth, err = m.jwt.ParseToken(ctx, token)
	return
}

func (m *Module) authenticatePreparedQueryRequest(ctx context.Context, projectID, dbAlias, id, token string) (rule *config.Rule, auth map[string]interface{}, err error) {
	// Get rule
	rule, err = m.getPrepareQueryRule(ctx, projectID, dbAlias, id)
	if err != nil {
		return
	}

	// Return if rule is allow
	if rule.Rule == "allow" {
		return
	}

	// Parse token
	auth, err = m.jwt.ParseToken(ctx, token)
	return
}

func (m *Module) getCrudRule(ctx context.Context, projectID, dbAlias, col string, query model.OperationType) (*config.Rule, error) {
	resourceIDs := []string{
		config.GenerateResourceID(m.clusterID, projectID, config.ResourceDatabaseRule, dbAlias, col, "rule"),
		config.GenerateResourceID(m.clusterID, projectID, config.ResourceDatabaseRule, dbAlias, "default", "rule"),
	}
	for _, resourceID := range resourceIDs {
		rule, ok := m.dbRules[resourceID]
		if ok {
			if r, p3 := rule.Rules[string(query)]; p3 {
				return r, nil
			}
		}
	}
	return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Security rule not defined for collection/table (%s) in database (%s). Ensure your table has correct access rights", col, dbAlias), nil, nil)
}

func (m *Module) getPrepareQueryRule(ctx context.Context, projectID, dbAlias, id string) (*config.Rule, error) {
	resourceIDs := []string{
		config.GenerateResourceID(m.clusterID, projectID, config.ResourceDatabasePreparedQuery, dbAlias, id),
		config.GenerateResourceID(m.clusterID, projectID, config.ResourceDatabasePreparedQuery, dbAlias, "default"),
	}
	for _, resourceID := range resourceIDs {
		rule, ok := m.dbPrepQueryRules[resourceID]
		if ok && rule.Rule != nil {
			return rule.Rule, nil
		}
	}
	return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("No security rule found for prepared Query (%s) in database with alias (%s)", id, dbAlias), nil, nil)
}
