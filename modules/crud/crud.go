package crud

import (
	"context"
	"log"
	"sync"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"

	"github.com/spaceuptech/space-cloud/modules/crud/mgo"
	"github.com/spaceuptech/space-cloud/modules/crud/sql"
)

// Module is the root block providing convenient wrappers
type Module struct {
	sync.RWMutex
	blocks    map[string]Crud
	primaryDB string
	schema    []SchemaType
}

type (
	// SchemaType is the data structure for storing the parsed values of schema string
	SchemaType               map[string]schemaCollection
	schemaCollection         map[string]schemaField
	schemaField              map[string]*schemaFieldType
	schemaFieldDirectiveArgs map[string]string
	fieldType                int

	schemaFieldType struct {
		IsFieldTypeRequired bool
		IsList              bool
		DirectiveType       map[string]schemaFieldDirectiveArgs
		Kind                fieldType
	}
)

const (
	TypeString fieldType = iota
	TypeInteger
	TypeFloat
	TypeBoolean
	TypeDateTime
	TypeEnum
	TypeJSON
	TypeID
	TypeJoin
)

// Crud abstracts the implementation crud operations of databases
type Crud interface {
	Create(ctx context.Context, project, col string, req *model.CreateRequest) error
	Read(ctx context.Context, project, col string, req *model.ReadRequest) (interface{}, error)
	Update(ctx context.Context, project, col string, req *model.UpdateRequest) error
	Delete(ctx context.Context, project, col string, req *model.DeleteRequest) error
	Aggregate(ctx context.Context, project, col string, req *model.AggregateRequest) (interface{}, error)
	Batch(ctx context.Context, project string, req *model.BatchRequest) error
	GetDBType() utils.DBType
	IsClientSafe() error
	Close() error
}

// Init create a new instance of the Module object
func Init() *Module {
	return &Module{blocks: make(map[string]Crud)}
}

func initBlock(dbType utils.DBType, enabled bool, connection string) (Crud, error) {
	switch dbType {
	case utils.Mongo:
		return mgo.Init(enabled, connection)

	case utils.MySQL, utils.Postgres:
		return sql.Init(dbType, enabled, connection)

	default:
		return nil, utils.ErrInvalidParams
	}
}

func (m *Module) getCrudBlock(dbType string) (Crud, error) {
	if crud, p := m.blocks[dbType]; p {
		return crud, nil
	}

	return nil, utils.ErrDatabaseConfigAbsent
}

// SetConfig set the rules adn secret key required by the crud block
func (m *Module) SetConfig(crud config.Crud) {
	m.Lock()
	defer m.Unlock()

	// Close the previous database connections
	for _, v := range m.blocks {
		v.Close()
	}
	m.blocks = make(map[string]Crud, len(crud))

	// Create a new crud blocks
	for k, v := range crud {
		c, err := initBlock(utils.DBType(k), v.Enabled, v.Conn)
		m.blocks[k] = c

		if err != nil {
			log.Println("Error connecting to " + k + " : " + err.Error())
		} else {
			log.Println("Successfully connected to " + k)
		}
	}
}
