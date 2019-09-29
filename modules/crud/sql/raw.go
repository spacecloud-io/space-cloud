package sql

import (
	"context"
)

// RawBatch performs a batch operation for schema creation
// NOTE: not to be exposed externally
func (s *SQL) RawBatch(ctx context.Context, queries []string) error {

	// Skip if length of queries == 0
	if len(queries) == 0 {
		return nil
	}

	tx, err := s.client.Beginx()
	if err != nil {
		return err
	}
	for _, query := range queries {
		_, err := tx.Exec(query)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// RawExec performs an operation for schema creation
// NOTE: not to be exposed externally
func (s *SQL) RawExec(ctx context.Context, query string) error {
	_, err := s.client.ExecContext(ctx, query, []interface{}{}...)
	return err
}
