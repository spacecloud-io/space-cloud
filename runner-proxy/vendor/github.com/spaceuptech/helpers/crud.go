package helpers

import (
	"context"
	"fmt"

	"github.com/fatih/structs"
	scAPI "github.com/spaceuptech/space-api-go"
	"github.com/spaceuptech/space-api-go/db"

	"net/http"
	"os"

	"github.com/spaceuptech/space-api-go/types"
)

type Crud struct {
	DbClient *db.DB
}

const (
	scEnvKey = "SPACE_CLOUD_ADDR"
	scToken  = "SPACE_CLOUD_TOKEN"
)

func InitCrud() (*Crud, error) {
	value := os.Getenv(scEnvKey)
	if value == "" {
		return nil, fmt.Errorf("couldn't find env variable %s", scEnvKey)
	}

	// Create space api object
	a := scAPI.New("spacecloud", value, false)
	value = os.Getenv(scToken)
	if value == "" {
		return nil, fmt.Errorf("couldn't find env variable %s", scToken)
	}
	a.SetToken(value)

	c := a.DB("db")
	if c == nil {
		return nil, fmt.Errorf("unable to initialize space api go, check if space cloud is running")
	}
	return &Crud{DbClient: c}, nil
}

func (m *Crud) Insert(ctx context.Context, tableName string, obj interface{}) error {
	_, err := m.CheckErrors(m.DbClient.Insert(tableName).Doc(obj).Apply(ctx))
	if err != nil {
		return err
	}
	return nil
}

func (m *Crud) Upsert(ctx context.Context, tableName string, whereClause types.M, obj interface{}) error {
	_, err := m.CheckErrors(m.DbClient.Upsert(tableName).Where(whereClause).Set(structs.Map(obj)).Apply(ctx))
	if err != nil {
		return err
	}
	return nil
}

func (m *Crud) Update(ctx context.Context, tableName string, whereClause types.M, obj interface{}) error {
	_, err := m.CheckErrors(m.DbClient.Update(tableName).Where(whereClause).Set(structs.Map(obj)).Apply(ctx))
	if err != nil {
		return err
	}
	return nil
}

func (m *Crud) GetOne(ctx context.Context, tableName string, whereClause types.M, response interface{}) error {
	result, err := m.CheckErrors(m.DbClient.GetOne(tableName).Where(whereClause).Apply(ctx))
	if err != nil {
		return err
	}
	if err := result.Unmarshal(response); err != nil {
		return fmt.Errorf("unable to un-marshal database response (%v)", err)
	}
	return nil
}

func (m *Crud) GetAll(ctx context.Context, tableName string, whereClause types.M, response interface{}) error {
	result, err := m.CheckErrors(m.DbClient.Get(tableName).Where(whereClause).Apply(ctx))
	if err != nil {
		return err
	}
	if err := result.Unmarshal(response); err != nil {
		return fmt.Errorf("unable to un-marshal database response (%v)", err)
	}
	return nil
}

func (m *Crud) Delete(ctx context.Context, tableName string, whereClause types.M) error {
	_, err := m.CheckErrors(m.DbClient.DeleteOne(tableName).Where(whereClause).Apply(ctx))
	if err != nil {
		return err
	}
	return nil
}

func (m *Crud) CheckErrors(result *types.Response, err error) (*types.Response, error) {
	if err != nil {
		// Network error
		return nil, fmt.Errorf("network error occurred while querying database (%v)", err)

	}
	if result.Status != http.StatusOK || result.Error != "" {
		// Query processing error
		return nil, fmt.Errorf("invalid status code (%d) received - (%v)", result.Status, result.Error)
	}
	return result, nil
}
