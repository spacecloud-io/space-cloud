package sql

import (
	"context"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

// Batch performs the provided operations in a single Batch
func (s *SQL) Batch(ctx context.Context, req *model.BatchRequest) ([]int64, error) {

	// Create an array to hold the counts
	counts := make([]int64, len(req.Requests))

	// Create a transaction object
	tx, err := s.getClient().BeginTxx(ctx, nil) // TODO - Write *sqlx.TxOption instead of nil
	if err != nil {
		return counts, err
	}

	for i, req := range req.Requests {
		switch req.Type {
		case string(model.Create):
			sqlQuery, args, err := s.generateCreateQuery(req.Col, &model.CreateRequest{Document: req.Document, Operation: req.Operation})
			if err != nil {
				return counts, err
			}
			res, err := doExecContext(ctx, sqlQuery, args, tx)
			if err != nil {
				return counts, err
			}
			counts[i], _ = res.RowsAffected()

		case string(model.Delete):
			sqlQuery, args, err := s.generateDeleteQuery(ctx, &model.DeleteRequest{Find: req.Find, Operation: req.Operation}, req.Col)
			if err != nil {
				return counts, err
			}
			res, err := doExecContext(ctx, sqlQuery, args, tx)
			if err != nil {
				return counts, err
			}
			counts[i], _ = res.RowsAffected()

		case string(model.Update):
			n, err := s.update(ctx, req.Col, &model.UpdateRequest{Find: req.Find, Operation: req.Operation, Update: req.Update}, tx)
			if err != nil {
				return counts, err
			}
			counts[i] = n

		}
	}
	return counts, tx.Commit() // commit the Batch
}
