package mgo

import (
	"context"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

// Batch performs the provided operations in a single Batch
func (m *Mongo) Batch(ctx context.Context, project string, txRequest *model.BatchRequest) ([]int64, error) {
	counts := make([]int64, len(txRequest.Requests))
	err := m.client.UseSession(ctx, func(session mongo.SessionContext) error {
		err := session.StartTransaction()
		if err != nil {
			return err
		}
		for i, req := range txRequest.Requests {
			col := req.Col

			switch req.Type {
			case string(utils.Create):
				doc := req.Document
				op := req.Operation

				counts[i], err = m.Create(session, project, col, &model.CreateRequest{Document: doc, Operation: op})
				if err != nil {
					session.AbortTransaction(session)
					return err
				}
			case string(utils.Update):
				find := req.Find
				op := req.Operation
				update := req.Update

				counts[i], err = m.Update(session, project, col, &model.UpdateRequest{Find: find, Operation: op, Update: update})
				if err != nil {
					session.AbortTransaction(session)
					return err
				}
			case string(utils.Delete):
				find := req.Find
				op := req.Operation

				counts[i], err = m.Delete(session, project, col, &model.DeleteRequest{Find: find, Operation: op})
				if err != nil {
					session.AbortTransaction(session)
					return err
				}
			}
		}
		err = session.CommitTransaction(session)
		if err != nil {
			session.AbortTransaction(session)
			return err
		}
		return nil
	})

	return counts, err
}
