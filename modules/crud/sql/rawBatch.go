package sql

import (
	"context"
	"fmt"
)

// RawBatch performs a batch operation for schema creation
// NOTE: not to be exposed externally
func (s *SQL) RawBatch(ctx context.Context, queries []string) error {
	fmt.Println(s, "\nsql-----------------", s.client, s.client == nil)
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
