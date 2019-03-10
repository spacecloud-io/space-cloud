package mgo

import (
	"context"

	"github.com/spaceuptech/space-cloud/model"
	"go.mongodb.org/mongo-driver/mongo"
)

// Transaction performs the provided operations in a single transaction
func (m *Mongo) Transaction(ctx context.Context, project string, reqs []map[string]interface{}) error {
	return m.client.UseSession(ctx, func(session mongo.SessionContext) error {
		err := session.StartTransaction()
		if err != nil {
			return err
		}
		for _, req := range reqs {
			col := req["col"].(string)

			switch req["type"].(string) {
			case "create":
				doc := req["doc"]
				op := req["op"].(string)

				err = m.Create(session, project, col, &model.CreateRequest{Document: doc, Operation: op})
				if err != nil {
					session.AbortTransaction(session)
					return err
				}
			case "update":
				find := req["find"].(map[string]interface{})
				op := req["op"].(string)
				update := req["update"].(map[string]interface{})

				err = m.Update(session, project, col, &model.UpdateRequest{Find: find, Operation: op, Update: update})
				if err != nil {
					session.AbortTransaction(session)
					return err
				}
			case "delete":
				find := req["find"].(map[string]interface{})
				op := req["op"].(string)

				err = m.Delete(session, project, col, &model.DeleteRequest{Find: find, Operation: op})
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
}
