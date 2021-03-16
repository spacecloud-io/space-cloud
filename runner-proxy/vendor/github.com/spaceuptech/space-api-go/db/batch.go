package db

import (
	"context"
	"errors"

	"github.com/spaceuptech/space-api-go/config"
	"github.com/spaceuptech/space-api-go/types"
)

// Delete contains the methods for the delete operation
type Batch struct {
	db     string
	config *config.Config
	reqs   []types.AllRequest
	meta   *types.Meta
}

type Request interface {
	getProject() string
	getDb() string
	getToken() string
	getCollection() string
	getOperation() string
}

func initBatch(db string, config *config.Config) *Batch {
	meta := &types.Meta{DbType: db, Project: config.Project, Token: config.Token, Operation: types.Batch}
	return &Batch{db, config, []types.AllRequest{}, meta}
}

// Add adds a delete request to batch
func (b *Batch) Add(request Request) error {
	req := request
	if b.config.Project != req.getProject() {
		return errors.New("Cannot Batch Requests of Different Projects")
	}
	if b.db != req.getDb() {
		return errors.New("Cannot Batch Requests of Different Database Types")
	}
	if b.config.Token != req.getToken() {
		return errors.New("Cannot Batch Requests using Different Tokens")
	}
	allReq := types.AllRequest{}
	allReq.Col = req.getCollection()
	allReq.Operation = req.getOperation()

	switch req := req.(type) {
	case *Insert:
		allReq.Document = req.getDoc()
		allReq.Type = types.Create
	case *Update:
		allReq.Find = req.getFind()
		allReq.Update = req.getUpdate()
		allReq.Type = types.Update
	case *Delete:
		allReq.Find = req.getFind()
		allReq.Type = types.Delete
	}
	b.reqs = append(b.reqs, allReq)
	return nil
}

// Apply executes the operation and returns the result
func (b *Batch) Apply(ctx context.Context) (*types.Response, error) {
	return b.config.Transport.DoDBRequest(ctx, b.meta, b.createBatchReq())
}

func (b *Batch) createBatchReq() *types.BatchRequest {
	return &types.BatchRequest{Requests: b.reqs}
}
