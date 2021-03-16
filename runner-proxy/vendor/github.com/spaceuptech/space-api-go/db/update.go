package db

import (
	"context"

	"github.com/spaceuptech/space-api-go/config"
	"github.com/spaceuptech/space-api-go/types"
)

// Update contains the methods for the update operation
type Update struct {
	op           string
	find, update types.M
	config       *config.Config
	meta         *types.Meta
}

func initUpdate(db, col, op string, config *config.Config) *Update {
	meta := &types.Meta{Col: col, DbType: db, Project: config.Project, Token: config.Token, Operation: types.Update}
	f := make(types.M)
	u := make(types.M)
	return &Update{op, f, u, config, meta}
}

// Where sets the where clause for the request
func (u *Update) Where(conds ...types.M) *Update {
	if len(conds) == 1 {
		u.find = types.GenerateFind(conds[0])
	} else {
		u.find = types.GenerateFind(types.And(conds...))
	}
	return u
}

// Set the value of the provided fields in the document
func (u *Update) Set(obj types.M) *Update {
	u.update["$set"] = obj
	return u
}

// Push adds an item to an array
func (u *Update) Push(obj types.M) *Update {
	u.update["$push"] = obj
	return u
}

// Remove removes the specified field from a document
func (u *Update) Remove(fields ...string) *Update {
	obj := make(types.M, len(fields))
	for _, field := range fields {
		obj[field] = 1
	}
	u.update["$unset"] = obj
	return u
}

// Rename renames the specified field
func (u *Update) Rename(obj types.M) *Update {
	u.update["$rename"] = obj
	return u
}

// Inc increments the value of the field by the specified amount
func (u *Update) Inc(obj types.M) *Update {
	u.update["$inc"] = obj
	return u
}

// Mul multiplies the value of the field by the specified amount
func (u *Update) Mul(obj types.M) *Update {
	u.update["$mul"] = obj
	return u
}

// Max updates the field if the specified value is greater than the existing field value
func (u *Update) Max(obj types.M) *Update {
	u.update["$max"] = obj
	return u
}

// Min updates the field if the specified value is lesser than the existing field value
func (u *Update) Min(obj types.M) *Update {
	u.update["$min"] = obj
	return u
}

// CurrentTimestamp sets the value of a field to current timestamp
func (u *Update) CurrentTimestamp(fields ...string) *Update {
	objTemp, p := u.update["$currentDate"]
	if !p {
		objTemp = types.M{}
	}

	obj := objTemp.(types.M)
	for _, field := range fields {
		obj[field] = types.M{"$type": "timestamp"}
	}

	u.update["$currentDate"] = obj
	return u
}

// CurrentDate sets the value of a field to current date
func (u *Update) CurrentDate(fields ...string) *Update {
	objTemp, p := u.update["$currentDate"]
	if !p {
		objTemp = types.M{}
	}

	obj := objTemp.(types.M)
	for _, field := range fields {
		obj[field] = types.M{"$type": "date"}
	}

	u.update["$currentDate"] = obj
	return u
}

// Apply executes the operation and returns the result
func (u *Update) Apply(ctx context.Context) (*types.Response, error) {
	return u.config.Transport.DoDBRequest(ctx, u.meta, u.createUpdateReq())
}

func (u *Update) getProject() string {
	return u.config.Project
}
func (u *Update) getDb() string {
	return u.meta.DbType
}
func (u *Update) getToken() string {
	return u.config.Token
}
func (u *Update) getCollection() string {
	return u.meta.Col
}
func (u *Update) getOperation() string {
	return u.op
}
func (u *Update) getUpdate() interface{} {
	return u.update
}
func (u *Update) getFind() interface{} {
	return u.find
}

func (u *Update) createUpdateReq() *types.UpdateRequest {
	return &types.UpdateRequest{Find: u.find, Operation: u.op, Update: u.update}
}
