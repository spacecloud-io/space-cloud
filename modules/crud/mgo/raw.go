package mgo

import (
	"context"
	"errors"
)

// RawBatch performs a batch operation for schema creation
// NOTE: not to be exposed externally
func (m *Mongo) RawBatch(ctx context.Context, queries []string) error {
	return errors.New("raw batch operation cannot be performed on mongo")
}

// RawExec performs an operation for schema creation
// NOTE: not to be exposed externally
func (m *Mongo) RawExec(ctx context.Context, query string) error {
	return errors.New("raw exec operation cannot be performed on mongo")
}
