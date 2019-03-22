package sql

import (
	"context"
	"fmt"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

// Transaction performs the provided operations in a single transaction
func (s *SQL) Transaction(ctx context.Context, project string, txRequest *model.TransactionRequest) error {

	tx, err := s.client.BeginTxx(ctx, nil) //TODO - Wirte *sqlx.TxOption instead of nil
	if err != nil {
		fmt.Println("Error in initiating transactions")
		return err
	}
	for _, req := range txRequest.Requests {
		switch req.Type {
		case string(utils.Create):
			sqlQuery, args, err := s.generateCreateQuery(ctx, project, req.Col, &model.CreateRequest{Document: req.Document, Operation: req.Operation})
			if err != nil {
				return err
			}
			err = doTransactionExecContext(ctx, sqlQuery, args, tx)
			if err != nil {
				return err
			}

		case string(utils.Delete):
			sqlQuery, args, err := s.generateDeleteQuery(ctx, project, req.Col, &model.DeleteRequest{Find: req.Find, Operation: req.Operation})
			if err != nil {
				return err
			}
			err = doTransactionExecContext(ctx, sqlQuery, args, tx)
			if err != nil {
				return err
			}

		case string(utils.Update):
			sqlQuery, args, err := s.generateUpdateQuery(ctx, project, req.Col, &model.UpdateRequest{Find: req.Find, Operation: req.Operation, Update: req.Update})
			if err != nil {
				return err
			}
			err = doTransactionExecContext(ctx, sqlQuery, args, tx)
			if err != nil {
				return err
			}

		}
	}
	err = tx.Commit() // commit the transaction
	if err != nil {
		return err
	}
	return nil
}
