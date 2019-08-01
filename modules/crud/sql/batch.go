package sql

import (
	"context"
	"fmt"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

// Batch performs the provided operations in a single Batch
func (s *SQL) Batch(ctx context.Context, project string, txRequest *model.BatchRequest) error {

	tx, err := s.client.BeginTxx(ctx, nil) //TODO - Write *sqlx.TxOption instead of nil
	if err != nil {
		fmt.Println("Error in initiating Batch")
		return err
	}
	for _, req := range txRequest.Requests {
		switch req.Type {
		case string(utils.Create):
			sqlQuery, args, err := s.generateCreateQuery(ctx, project, req.Col, &model.CreateRequest{Document: req.Document, Operation: req.Operation})
			if err != nil {
				return err
			}
			_, err = doExecContext(ctx, sqlQuery, args, tx)
			if err != nil {
				return err
			}

		case string(utils.Delete):
			sqlQuery, args, err := s.generateDeleteQuery(ctx, project, req.Col, &model.DeleteRequest{Find: req.Find, Operation: req.Operation})
			if err != nil {
				return err
			}
			_, err = doExecContext(ctx, sqlQuery, args, tx)
			if err != nil {
				return err
			}

		case string(utils.Update):
			err = s.update(ctx, project, req.Col, &model.UpdateRequest{Find: req.Find, Operation: req.Operation, Update: req.Update}, tx)
			if err != nil {
				return err
			}

		}
	}
	return tx.Commit() // commit the Batch
}
