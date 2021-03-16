package db

import (
	"context"

	"github.com/spaceuptech/space-api-go/config"
	"github.com/spaceuptech/space-api-go/types"
)

// Delete contains the methods for the delete operation
type Delete struct {
	op     string
	find   types.M
	config *config.Config
	meta   *types.Meta
}

func initDelete(db, col, op string, config *config.Config) *Delete {
	return &Delete{
		meta:   &types.Meta{DbType: db, Col: col, Token: config.Token, Project: config.Project, Operation: types.Delete},
		op:     op,
		find:   make(types.M),
		config: config}
}

// Where sets the where clause for the request
func (d *Delete) Where(conds ...types.M) *Delete {
	if len(conds) == 1 {
		d.find = types.GenerateFind(conds[0])
	} else {
		d.find = types.GenerateFind(types.And(conds...))
	}
	return d
}

// Apply executes the operation and returns the result
func (d *Delete) Apply(ctx context.Context) (*types.Response, error) {
	return d.config.Transport.DoDBRequest(ctx, d.meta, d.createDeleteReq())
}

func (d *Delete) getProject() string {
	return d.config.Project
}
func (d *Delete) getDb() string {
	return d.meta.DbType
}
func (d *Delete) getToken() string {
	return d.config.Token
}
func (d *Delete) getCollection() string {
	return d.meta.Col
}
func (d *Delete) getOperation() string {
	return d.op
}
func (d *Delete) getFind() types.M {
	return d.find
}

func (d *Delete) createDeleteReq() *types.DeleteRequest {
	return &types.DeleteRequest{Find: d.find, Operation: d.op}
}
