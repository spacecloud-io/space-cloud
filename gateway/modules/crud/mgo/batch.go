package mgo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// Batch performs the provided operations in a single Batch
func (m *Mongo) Batch(ctx context.Context, req *model.BatchRequest) ([]int64, error) {
	counts := make([]int64, len(req.Requests))
	err := m.client.UseSession(ctx, func(session mongo.SessionContext) error {
		err := session.StartTransaction()
		if err != nil {
			return err
		}
		for i, req := range req.Requests {
			col := req.Col

			switch req.Type {
			case string(utils.Create):
				doc := req.Document
				op := req.Operation

				counts[i], err = m.Create(session, col, &model.CreateRequest{Document: doc, Operation: op})
				if err != nil {
					_ = session.AbortTransaction(session)
					return err
				}
			case string(utils.Update):
				find := req.Find
				op := req.Operation
				update := req.Update

				counts[i], err = m.Update(session, col, &model.UpdateRequest{Find: find, Operation: op, Update: update})
				if err != nil {
					_ = session.AbortTransaction(session)
					return err
				}
			case string(utils.Delete):
				find := req.Find
				op := req.Operation

				counts[i], err = m.Delete(session, col, &model.DeleteRequest{Find: find, Operation: op})
				if err != nil {
					_ = session.AbortTransaction(session)
					return err
				}
			}
		}
		err = session.CommitTransaction(session)
		if err != nil {
			_ = session.AbortTransaction(session)
			return err
		}
		return nil
	})

	return counts, err
}
