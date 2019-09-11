package mgo

import (
	"context"
	"errors"
)

// RawBatch performs a batch operation for schema creation
// NOTE: not to be exposed externally
func (s *Mongo) RawBatch(ctx context.Context, queries []string) error {
	return errors.New("Schema creation operation cannot be performed over mongo")
}
