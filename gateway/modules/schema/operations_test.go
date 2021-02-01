package schema

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/go-test/deep"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	. "github.com/spaceuptech/space-cloud/gateway/modules/schema/helpers"
)

func TestManager_GetSchemas(t *testing.T) {
	type args struct {
		ctx     context.Context
		dbAlias string
		col     string
		format  string
	}
	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}
	tests := []struct {
		name         string
		crud         config.DatabaseSchemas
		args         args
		crudMockArgs []mockArgs
		want1        []interface{}
		want2        []interface{}
		wantErr      bool
	}{
		{
			name: "Invalid dbAlias provided while fetching specific table from specific database",
			crud: config.DatabaseSchemas{
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db1", "authors"):     &config.DatabaseSchema{DbAlias: "db1", Table: "authors", Schema: `type authors {id: ID! name: String genre_id: ID! @foreign(table: "genres",to: "id")}`},
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db1", "subscribers"): &config.DatabaseSchema{DbAlias: "db1", Table: "subscribers", Schema: `type subscribers {id: ID! name: String author_id: ID! @foreign(table: "authors",to: "id")}`},
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db1", "genres"):      &config.DatabaseSchema{DbAlias: "db1", Table: "genres", Schema: `type genres {id: ID! name: String }`},
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db2", "authors"):     &config.DatabaseSchema{DbAlias: "db2", Table: "authors", Schema: `type authors {id: ID! name: String genre_id: ID! @foreign(table: "genres",to: "id")}`},
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db2", "subscribers"): &config.DatabaseSchema{DbAlias: "db2", Table: "subscribers", Schema: `type subscribers {id: ID! name: String author_id: ID! @foreign(table: "authors",to: "id")}`},
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db2", "genres"):      &config.DatabaseSchema{DbAlias: "db2", Table: "genres", Schema: `type genres {id: ID! name: String }`},
			},
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"db1"},
					paramsReturned: []interface{}{"postgres"},
				},
				{
					method:         "GetDBType",
					args:           []interface{}{"db2"},
					paramsReturned: []interface{}{"postgres"},
				},
			},
			args:    args{ctx: context.Background(), col: "genres", dbAlias: "db3"},
			wantErr: true,
		},
		{
			name: "Get schema of specified database & table in json format",
			crud: config.DatabaseSchemas{
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db1", "authors"):     &config.DatabaseSchema{DbAlias: "db1", Table: "authors", Schema: `type authors {id: ID! name: String genre_id: ID! @foreign(table: "genres",to: "id")}`},
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db1", "subscribers"): &config.DatabaseSchema{DbAlias: "db1", Table: "subscribers", Schema: `type subscribers {id: ID! name: String author_id: ID! @foreign(table: "authors",to: "id")}`},
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db1", "genres"):      &config.DatabaseSchema{DbAlias: "db1", Table: "genres", Schema: `type genres {id: ID! name: String }`},
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db2", "authors"):     &config.DatabaseSchema{DbAlias: "db2", Table: "authors", Schema: `type authors {id: ID! name: String genre_id: ID! @foreign(table: "genres",to: "id")}`},
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db2", "subscribers"): &config.DatabaseSchema{DbAlias: "db2", Table: "subscribers", Schema: `type subscribers {id: ID! name: String author_id: ID! @foreign(table: "authors",to: "id")}`},
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db2", "genres"):      &config.DatabaseSchema{DbAlias: "db2", Table: "genres", Schema: `type genres {id: ID! name: String }`},
			},
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"db1"},
					paramsReturned: []interface{}{"postgres"},
				},
				{
					method:         "GetDBType",
					args:           []interface{}{"db2"},
					paramsReturned: []interface{}{"postgres"},
				},
			},
			args: args{ctx: context.Background(), col: "authors", dbAlias: "db1", format: "json"},
			want1: []interface{}{
				dbSchemaResponse{
					DbAlias: "db1",
					Col:     "authors",
					SchemaObj: model.Fields{
						"id": &model.FieldType{
							FieldName:           "id",
							IsFieldTypeRequired: true,
							Kind:                model.TypeID,
							TypeIDSize:          model.DefaultCharacterSize,
						},
						"name": &model.FieldType{
							FieldName: "name",
							Kind:      model.TypeString,
						},
						"genre_id": &model.FieldType{
							FieldName:           "genre_id",
							IsFieldTypeRequired: true,
							Kind:                model.TypeID,
							IsForeign:           true,
							TypeIDSize:          model.DefaultCharacterSize,
							JointTable: &model.TableProperties{
								Table:          "genres",
								To:             "id",
								OnDelete:       "NO ACTION",
								ConstraintName: GetConstraintName("authors", "genre_id"),
							},
						},
					}},
			},
		},
		{
			name: "Get schema of specified database in json format",
			crud: config.DatabaseSchemas{
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db1", "authors"):     &config.DatabaseSchema{DbAlias: "db1", Table: "authors", Schema: `type authors {id: ID! name: String genre_id: ID! @foreign(table: "genres",to: "id")}`},
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db1", "subscribers"): &config.DatabaseSchema{DbAlias: "db1", Table: "subscribers", Schema: `type subscribers {id: ID! name: String author_id: ID! @foreign(table: "authors",to: "id")}`},
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db1", "genres"):      &config.DatabaseSchema{DbAlias: "db1", Table: "genres", Schema: `type genres {id: ID! name: String }`},
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db2", "authors"):     &config.DatabaseSchema{DbAlias: "db2", Table: "authors", Schema: `type authors {id: ID! name: String genre_id: ID! @foreign(table: "genres",to: "id")}`},
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db2", "subscribers"): &config.DatabaseSchema{DbAlias: "db2", Table: "subscribers", Schema: `type subscribers {id: ID! name: String author_id: ID! @foreign(table: "authors",to: "id")}`},
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db2", "genres"):      &config.DatabaseSchema{DbAlias: "db2", Table: "genres", Schema: `type genres {id: ID! name: String }`},
			},
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"db1"},
					paramsReturned: []interface{}{"postgres"},
				},
				{
					method:         "GetDBType",
					args:           []interface{}{"db2"},
					paramsReturned: []interface{}{"postgres"},
				},
			},
			args: args{ctx: context.Background(), col: "*", dbAlias: "db1", format: "json"},
			want1: []interface{}{
				dbSchemaResponse{
					DbAlias: "db1",
					Col:     "genres",
					SchemaObj: model.Fields{
						"id": &model.FieldType{
							FieldName:           "id",
							IsFieldTypeRequired: true,
							Kind:                model.TypeID,
							TypeIDSize:          model.DefaultCharacterSize,
						},
						"name": &model.FieldType{
							FieldName: "name",
							Kind:      model.TypeString,
						},
					},
				},
				dbSchemaResponse{
					DbAlias: "db1",
					Col:     "authors",
					SchemaObj: model.Fields{
						"id": &model.FieldType{
							FieldName:           "id",
							IsFieldTypeRequired: true,
							Kind:                model.TypeID,
							TypeIDSize:          model.DefaultCharacterSize,
						},
						"name": &model.FieldType{
							FieldName: "name",
							Kind:      model.TypeString,
						},
						"genre_id": &model.FieldType{
							FieldName:           "genre_id",
							IsFieldTypeRequired: true,
							Kind:                model.TypeID,
							IsForeign:           true,
							TypeIDSize:          model.DefaultCharacterSize,
							JointTable: &model.TableProperties{
								Table:          "genres",
								To:             "id",
								OnDelete:       "NO ACTION",
								ConstraintName: GetConstraintName("authors", "genre_id"),
							},
						},
					},
				},
				dbSchemaResponse{
					DbAlias: "db1",
					Col:     "subscribers",
					SchemaObj: model.Fields{
						"id": &model.FieldType{
							FieldName:           "id",
							IsFieldTypeRequired: true,
							Kind:                model.TypeID,
							TypeIDSize:          model.DefaultCharacterSize,
						},
						"name": &model.FieldType{
							FieldName: "name",
							Kind:      model.TypeString,
						},
						"author_id": &model.FieldType{
							FieldName:           "author_id",
							IsFieldTypeRequired: true,
							Kind:                model.TypeID,
							IsForeign:           true,
							TypeIDSize:          model.DefaultCharacterSize,
							JointTable: &model.TableProperties{
								Table:          "authors",
								To:             "id",
								OnDelete:       "NO ACTION",
								ConstraintName: GetConstraintName("subscribers", "author_id"),
							},
						},
					},
				},
			},
		},
		// The output of this test case is not predictable
		// {
		// 	name: "Get schema of all database in json format",
		// 	crud: config.DatabaseSchemas{
		// 		config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db1", "authors"):     &config.DatabaseSchema{DbAlias: "db1", Table: "authors", Schema: `type authors {id: ID! name: String genre_id: ID! @foreign(table: "genres",to: "id")}`},
		// 		config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db1", "subscribers"): &config.DatabaseSchema{DbAlias: "db1", Table: "subscribers", Schema: `type subscribers {id: ID! name: String author_id: ID! @foreign(table: "authors",to: "id")}`},
		// 		config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db1", "genres"):      &config.DatabaseSchema{DbAlias: "db1", Table: "genres", Schema: `type genres {id: ID! name: String }`},
		// 		config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db2", "authors"):     &config.DatabaseSchema{DbAlias: "db2", Table: "authors", Schema: `type authors {id: ID! name: String genre_id: ID! @foreign(table: "genres",to: "id")}`},
		// 		config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db2", "subscribers"): &config.DatabaseSchema{DbAlias: "db2", Table: "subscribers", Schema: `type subscribers {id: ID! name: String author_id: ID! @foreign(table: "authors",to: "id")}`},
		// 		config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db2", "genres"):      &config.DatabaseSchema{DbAlias: "db2", Table: "genres", Schema: `type genres {id: ID! name: String }`},
		// 	},
		// 	args: args{ctx: context.Background(), col: "*", dbAlias: "*", format: "json"},
		// 	want1: []interface{}{
		// 		dbSchemaResponse{
		// 			DbAlias: "db1",
		// 			Col:     "genres",
		// 			SchemaObj: model.Fields{
		// 				"id": &model.FieldType{
		// 					FieldName:           "id",
		// 					IsFieldTypeRequired: true,
		// 					Kind:                model.TypeID,
		// 					TypeIDSize:          model.DefaultCharacterSize,
		// 				},
		// 				"name": &model.FieldType{
		// 					FieldName:  "name",
		// AutoIncrementInfo:new(model.AutoIncrementInfo),
		// 					Kind:       model.TypeString,
		// 					TypeIDSize: model.DefaultCharacterSize,
		// 				},
		// 			},
		// 		},
		// 		dbSchemaResponse{
		// 			DbAlias: "db1",
		// 			Col:     "authors",
		// 			SchemaObj: model.Fields{
		// 				"id": &model.FieldType{
		// 					FieldName:           "id",
		// 					IsFieldTypeRequired: true,
		// 					Kind:                model.TypeID,
		// 					TypeIDSize:          model.DefaultCharacterSize,
		// 				},
		// 				"name": &model.FieldType{
		// 					FieldName:  "name",
		// AutoIncrementInfo:new(model.AutoIncrementInfo),
		// 					Kind:       model.TypeString,
		// 					TypeIDSize: model.DefaultCharacterSize,
		// 				},
		// 				"genre_id": &model.FieldType{
		// 					FieldName:           "genre_id",
		// 					IsFieldTypeRequired: true,
		// 					Kind:                model.TypeID,
		// 					IsForeign:           true,
		// 					TypeIDSize:          model.DefaultCharacterSize,
		// 					JointTable: &model.TableProperties{
		// 						Table:          "genres",
		// 						To:             "id",
		// 						OnDelete:       "NO ACTION",
		// 						ConstraintName: GetConstraintName("authors", "genre_id"),
		// 					},
		// 				},
		// 			},
		// 		},
		// 		dbSchemaResponse{
		// 			DbAlias: "db1",
		// 			Col:     "subscribers",
		// 			SchemaObj: model.Fields{
		// 				"id": &model.FieldType{
		// 					FieldName:           "id",
		// AutoIncrementInfo:new(model.AutoIncrementInfo),
		// 					IsFieldTypeRequired: true,
		// 					Kind:                model.TypeID,
		// 					TypeIDSize:          model.DefaultCharacterSize,
		// 				},
		// 				"name": &model.FieldType{
		// 					FieldName:  "name",
		// 					Kind:       model.TypeString,
		// 					TypeIDSize: model.DefaultCharacterSize,
		// 				},
		// 				"author_id": &model.FieldType{
		// 					FieldName:           "author_id",
		// AutoIncrementInfo:new(model.AutoIncrementInfo),
		// 					IsFieldTypeRequired: true,
		// 					Kind:                model.TypeID,
		// 					IsForeign:           true,
		// 					TypeIDSize:          model.DefaultCharacterSize,
		// 					JointTable: &model.TableProperties{
		// 						Table:          "authors",
		// 						To:             "id",
		// 						OnDelete:       "NO ACTION",
		// 						ConstraintName: GetConstraintName("subscribers", "author_id"),
		// 					},
		// 				},
		// 			},
		// 		},
		// 		dbSchemaResponse{
		// 			DbAlias: "db2",
		// 			Col:     "genres",
		// 			SchemaObj: model.Fields{
		// 				"id": &model.FieldType{
		// 					FieldName:           "id",
		// 					IsFieldTypeRequired: true,
		// 					Kind:                model.TypeID,
		// 					TypeIDSize:          model.DefaultCharacterSize,
		// 				},
		// 				"name": &model.FieldType{
		// 					FieldName:  "name",
		// AutoIncrementInfo:new(model.AutoIncrementInfo),
		// 					Kind:       model.TypeString,
		// 					TypeIDSize: model.DefaultCharacterSize,
		// 				},
		// 			},
		// 		},
		// 		dbSchemaResponse{
		// 			DbAlias: "db2",
		// 			Col:     "authors",
		// 			SchemaObj: model.Fields{
		// 				"id": &model.FieldType{
		// 					FieldName:           "id",
		// 					IsFieldTypeRequired: true,
		// 					Kind:                model.TypeID,
		// 					TypeIDSize:          model.DefaultCharacterSize,
		// 				},
		// 				"name": &model.FieldType{
		// 					FieldName:  "name",
		// AutoIncrementInfo:new(model.AutoIncrementInfo),
		// 					Kind:       model.TypeString,
		// 					TypeIDSize: model.DefaultCharacterSize,
		// 				},
		// 				"genre_id": &model.FieldType{
		// 					FieldName:           "genre_id",
		// 					IsFieldTypeRequired: true,
		// 					Kind:                model.TypeID,
		// 					IsForeign:           true,
		// 					TypeIDSize:          model.DefaultCharacterSize,
		// 					JointTable: &model.TableProperties{
		// 						Table:          "genres",
		// 						To:             "id",
		// 						OnDelete:       "NO ACTION",
		// 						ConstraintName: GetConstraintName("authors", "genre_id"),
		// 					},
		// 				},
		// 			},
		// 		},
		// 		dbSchemaResponse{
		// 			DbAlias: "db2",
		// 			Col:     "subscribers",
		// 			SchemaObj: model.Fields{
		// 				"id": &model.FieldType{
		// 					FieldName:           "id",
		// AutoIncrementInfo:new(model.AutoIncrementInfo),
		// 					IsFieldTypeRequired: true,
		// 					Kind:                model.TypeID,
		// 					TypeIDSize:          model.DefaultCharacterSize,
		// 				},
		// 				"name": &model.FieldType{
		// 					FieldName:  "name",
		// 					Kind:       model.TypeString,
		// 					TypeIDSize: model.DefaultCharacterSize,
		// 				},
		// 				"author_id": &model.FieldType{
		// 					FieldName:           "author_id",
		// AutoIncrementInfo:new(model.AutoIncrementInfo),
		// 					IsFieldTypeRequired: true,
		// 					Kind:                model.TypeID,
		// 					IsForeign:           true,
		// 					TypeIDSize:          model.DefaultCharacterSize,
		// 					JointTable: &model.TableProperties{
		// 						Table:          "authors",
		// 						To:             "id",
		// 						OnDelete:       "NO ACTION",
		// 						ConstraintName: GetConstraintName("subscribers", "author_id"),
		// 					},
		// 				},
		// 			},
		// 		},
		// 	},
		// 	want2: []interface{}{
		// 		dbSchemaResponse{
		// 			DbAlias: "db2",
		// 			Col:     "genres",
		// 			SchemaObj: model.Fields{
		// 				"id": &model.FieldType{
		// 					FieldName:           "id",
		// 					IsFieldTypeRequired: true,
		// 					Kind:                model.TypeID,
		// 					TypeIDSize:          model.DefaultCharacterSize,
		// 				},
		// 				"name": &model.FieldType{
		// 					FieldName:  "name",
		// AutoIncrementInfo:new(model.AutoIncrementInfo),
		// 					Kind:       model.TypeString,
		// 					TypeIDSize: model.DefaultCharacterSize,
		// 				},
		// 			},
		// 		},
		// 		dbSchemaResponse{
		// 			DbAlias: "db2",
		// 			Col:     "authors",
		// 			SchemaObj: model.Fields{
		// 				"id": &model.FieldType{
		// 					FieldName:           "id",
		// 					IsFieldTypeRequired: true,
		// 					Kind:                model.TypeID,
		// 					TypeIDSize:          model.DefaultCharacterSize,
		// 				},
		// 				"name": &model.FieldType{
		// 					FieldName:  "name",
		// AutoIncrementInfo:new(model.AutoIncrementInfo),
		// 					Kind:       model.TypeString,
		// 					TypeIDSize: model.DefaultCharacterSize,
		// 				},
		// 				"genre_id": &model.FieldType{
		// 					FieldName:           "genre_id",
		// 					IsFieldTypeRequired: true,
		// 					Kind:                model.TypeID,
		// 					IsForeign:           true,
		// 					TypeIDSize:          model.DefaultCharacterSize,
		// 					JointTable: &model.TableProperties{
		// 						Table:          "genres",
		// 						To:             "id",
		// 						OnDelete:       "NO ACTION",
		// 						ConstraintName: GetConstraintName("authors", "genre_id"),
		// 					},
		// 				},
		// 			},
		// 		},
		// 		dbSchemaResponse{
		// 			DbAlias: "db2",
		// 			Col:     "subscribers",
		// 			SchemaObj: model.Fields{
		// 				"id": &model.FieldType{
		// 					FieldName:           "id",
		// AutoIncrementInfo:new(model.AutoIncrementInfo),
		// 					IsFieldTypeRequired: true,
		// 					Kind:                model.TypeID,
		// 					TypeIDSize:          model.DefaultCharacterSize,
		// 				},
		// 				"name": &model.FieldType{
		// 					FieldName:  "name",
		// 					Kind:       model.TypeString,
		// 					TypeIDSize: model.DefaultCharacterSize,
		// 				},
		// 				"author_id": &model.FieldType{
		// 					FieldName:           "author_id",
		// AutoIncrementInfo:new(model.AutoIncrementInfo),
		// 					IsFieldTypeRequired: true,
		// 					Kind:                model.TypeID,
		// 					IsForeign:           true,
		// 					TypeIDSize:          model.DefaultCharacterSize,
		// 					JointTable: &model.TableProperties{
		// 						Table:          "authors",
		// 						To:             "id",
		// 						OnDelete:       "NO ACTION",
		// 						ConstraintName: GetConstraintName("subscribers", "author_id"),
		// 					},
		// 				},
		// 			},
		// 		},
		// 		dbSchemaResponse{
		// 			DbAlias: "db1",
		// 			Col:     "genres",
		// 			SchemaObj: model.Fields{
		// 				"id": &model.FieldType{
		// 					FieldName:           "id",
		// 					IsFieldTypeRequired: true,
		// 					Kind:                model.TypeID,
		// 					TypeIDSize:          model.DefaultCharacterSize,
		// 				},
		// 				"name": &model.FieldType{
		// 					FieldName:  "name",
		// AutoIncrementInfo:new(model.AutoIncrementInfo),
		// 					Kind:       model.TypeString,
		// 					TypeIDSize: model.DefaultCharacterSize,
		// 				},
		// 			},
		// 		},
		// 		dbSchemaResponse{
		// 			DbAlias: "db1",
		// 			Col:     "authors",
		// 			SchemaObj: model.Fields{
		// 				"id": &model.FieldType{
		// 					FieldName:           "id",
		// 					IsFieldTypeRequired: true,
		// 					Kind:                model.TypeID,
		// 					TypeIDSize:          model.DefaultCharacterSize,
		// 				},
		// 				"name": &model.FieldType{
		// 					FieldName:  "name",
		// AutoIncrementInfo:new(model.AutoIncrementInfo),
		// 					Kind:       model.TypeString,
		// 					TypeIDSize: model.DefaultCharacterSize,
		// 				},
		// 				"genre_id": &model.FieldType{
		// 					FieldName:           "genre_id",
		// 					IsFieldTypeRequired: true,
		// 					Kind:                model.TypeID,
		// 					IsForeign:           true,
		// 					TypeIDSize:          model.DefaultCharacterSize,
		// 					JointTable: &model.TableProperties{
		// 						Table:          "genres",
		// 						To:             "id",
		// 						OnDelete:       "NO ACTION",
		// 						ConstraintName: GetConstraintName("authors", "genre_id"),
		// 					},
		// 				},
		// 			},
		// 		},
		// 		dbSchemaResponse{
		// 			DbAlias: "db1",
		// 			Col:     "subscribers",
		// 			SchemaObj: model.Fields{
		// 				"id": &model.FieldType{
		// 					FieldName:           "id",
		// AutoIncrementInfo:new(model.AutoIncrementInfo),
		// 					IsFieldTypeRequired: true,
		// 					Kind:                model.TypeID,
		// 					TypeIDSize:          model.DefaultCharacterSize,
		// 				},
		// 				"name": &model.FieldType{
		// 					FieldName:  "name",
		// 					Kind:       model.TypeString,
		// 					TypeIDSize: model.DefaultCharacterSize,
		// 				},
		// 				"author_id": &model.FieldType{
		// 					FieldName:           "author_id",
		// AutoIncrementInfo:new(model.AutoIncrementInfo),
		// 					IsFieldTypeRequired: true,
		// 					Kind:                model.TypeID,
		// 					IsForeign:           true,
		// 					TypeIDSize:          model.DefaultCharacterSize,
		// 					JointTable: &model.TableProperties{
		// 						Table:          "authors",
		// 						To:             "id",
		// 						OnDelete:       "NO ACTION",
		// 						ConstraintName: GetConstraintName("subscribers", "author_id"),
		// 					},
		// 				},
		// 			},
		// 		},
		// 	},
		// },
		{
			name: "Get schema of specified database & table",
			crud: config.DatabaseSchemas{
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db1", "authors"):     &config.DatabaseSchema{DbAlias: "db1", Table: "authors", Schema: `type authors {id: ID! name: String genre_id: ID! @foreign(table: "genres",to: "id")}`},
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db1", "subscribers"): &config.DatabaseSchema{DbAlias: "db1", Table: "subscribers", Schema: `type subscribers {id: ID! name: String author_id: ID! @foreign(table: "authors",to: "id")}`},
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db1", "genres"):      &config.DatabaseSchema{DbAlias: "db1", Table: "genres", Schema: `type genres {id: ID! name: String }`},
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db2", "authors"):     &config.DatabaseSchema{DbAlias: "db2", Table: "authors", Schema: `type authors {id: ID! name: String genre_id: ID! @foreign(table: "genres",to: "id")}`},
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db2", "subscribers"): &config.DatabaseSchema{DbAlias: "db2", Table: "subscribers", Schema: `type subscribers {id: ID! name: String author_id: ID! @foreign(table: "authors",to: "id")}`},
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db2", "genres"):      &config.DatabaseSchema{DbAlias: "db2", Table: "genres", Schema: `type genres {id: ID! name: String }`},
			},
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"db1"},
					paramsReturned: []interface{}{"postgres"},
				},
				{
					method:         "GetDBType",
					args:           []interface{}{"db2"},
					paramsReturned: []interface{}{"postgres"},
				},
			},
			args: args{ctx: context.Background(), col: "authors", dbAlias: "db1", format: ""},
			want1: []interface{}{
				dbSchemaResponse{DbAlias: "db1", Col: "authors", Schema: `type authors {id: ID! name: String genre_id: ID! @foreign(table: "genres",to: "id")}`},
			},
		},
		{
			name: "Get schema of specified database",
			crud: config.DatabaseSchemas{
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db1", "authors"):     &config.DatabaseSchema{DbAlias: "db1", Table: "authors", Schema: `type authors {id: ID! name: String genre_id: ID! @foreign(table: "genres",to: "id")}`},
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db1", "subscribers"): &config.DatabaseSchema{DbAlias: "db1", Table: "subscribers", Schema: `type subscribers {id: ID! name: String author_id: ID! @foreign(table: "authors",to: "id")}`},
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db1", "genres"):      &config.DatabaseSchema{DbAlias: "db1", Table: "genres", Schema: `type genres {id: ID! name: String }`},
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db2", "authors"):     &config.DatabaseSchema{DbAlias: "db2", Table: "authors", Schema: `type authors {id: ID! name: String genre_id: ID! @foreign(table: "genres",to: "id")}`},
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db2", "subscribers"): &config.DatabaseSchema{DbAlias: "db2", Table: "subscribers", Schema: `type subscribers {id: ID! name: String author_id: ID! @foreign(table: "authors",to: "id")}`},
				config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db2", "genres"):      &config.DatabaseSchema{DbAlias: "db2", Table: "genres", Schema: `type genres {id: ID! name: String }`},
			},
			crudMockArgs: []mockArgs{
				{
					method:         "GetDBType",
					args:           []interface{}{"db1"},
					paramsReturned: []interface{}{"postgres"},
				},
				{
					method:         "GetDBType",
					args:           []interface{}{"db2"},
					paramsReturned: []interface{}{"postgres"},
				},
			},
			args: args{ctx: context.Background(), col: "*", dbAlias: "db1", format: ""},
			want1: []interface{}{
				dbSchemaResponse{DbAlias: "db1", Col: "genres", Schema: `type genres {id: ID! name: String }`},
				dbSchemaResponse{DbAlias: "db1", Col: "authors", Schema: `type authors {id: ID! name: String genre_id: ID! @foreign(table: "genres",to: "id")}`},
				dbSchemaResponse{DbAlias: "db1", Col: "subscribers", Schema: `type subscribers {id: ID! name: String author_id: ID! @foreign(table: "authors",to: "id")}`},
			},
		},
		// The output of this case is not predictable
		// {
		// 	name: "Get schema of all databases",
		// 	crud: config.DatabaseSchemas{
		// 		config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db1", "authors"):     &config.DatabaseSchema{DbAlias: "db1", Table: "authors", Schema: `type authors {id: ID! name: String genre_id: ID! @foreign(table: "genres",to: "id")}`},
		// 		config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db1", "subscribers"): &config.DatabaseSchema{DbAlias: "db1", Table: "subscribers", Schema: `type subscribers {id: ID! name: String author_id: ID! @foreign(table: "authors",to: "id")}`},
		// 		config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db1", "genres"):      &config.DatabaseSchema{DbAlias: "db1", Table: "genres", Schema: `type genres {id: ID! name: String }`},
		// 		config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db2", "authors"):     &config.DatabaseSchema{DbAlias: "db2", Table: "authors", Schema: `type authors {id: ID! name: String genre_id: ID! @foreign(table: "genres",to: "id")}`},
		// 		config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db2", "subscribers"): &config.DatabaseSchema{DbAlias: "db2", Table: "subscribers", Schema: `type subscribers {id: ID! name: String author_id: ID! @foreign(table: "authors",to: "id")}`},
		// 		config.GenerateResourceID("chicago", "myproject", config.ResourceDatabaseSchema, "db2", "genres"):      &config.DatabaseSchema{DbAlias: "db2", Table: "genres", Schema: `type genres {id: ID! name: String }`},
		// 	},
		// 	args: args{ctx: context.Background(), col: "*", dbAlias: "*", format: ""},
		// 	want1: []interface{}{
		// 		dbSchemaResponse{DbAlias: "db1", Col: "genres", Schema: `type genres {id: ID! name: String }`},
		// 		dbSchemaResponse{DbAlias: "db1", Col: "authors", Schema: `type authors {id: ID! name: String genre_id: ID! @foreign(table: "genres",to: "id")}`},
		// 		dbSchemaResponse{DbAlias: "db1", Col: "subscribers", Schema: `type subscribers {id: ID! name: String author_id: ID! @foreign(table: "authors",to: "id")}`},
		// 		dbSchemaResponse{DbAlias: "db2", Col: "genres", Schema: `type genres {id: ID! name: String }`},
		// 		dbSchemaResponse{DbAlias: "db2", Col: "authors", Schema: `type authors {id: ID! name: String genre_id: ID! @foreign(table: "genres",to: "id")}`},
		// 		dbSchemaResponse{DbAlias: "db2", Col: "subscribers", Schema: `type subscribers {id: ID! name: String author_id: ID! @foreign(table: "authors",to: "id")}`},
		// 	},
		// 	want2: []interface{}{
		// 		dbSchemaResponse{DbAlias: "db2", Col: "genres", Schema: `type genres {id: ID! name: String }`},
		// 		dbSchemaResponse{DbAlias: "db2", Col: "authors", Schema: `type authors {id: ID! name: String genre_id: ID! @foreign(table: "genres",to: "id")}`},
		// 		dbSchemaResponse{DbAlias: "db2", Col: "subscribers", Schema: `type subscribers {id: ID! name: String author_id: ID! @foreign(table: "authors",to: "id")}`},
		// 		dbSchemaResponse{DbAlias: "db1", Col: "genres", Schema: `type genres {id: ID! name: String }`},
		// 		dbSchemaResponse{DbAlias: "db1", Col: "authors", Schema: `type authors {id: ID! name: String genre_id: ID! @foreign(table: "genres",to: "id")}`},
		// 		dbSchemaResponse{DbAlias: "db1", Col: "subscribers", Schema: `type subscribers {id: ID! name: String author_id: ID! @foreign(table: "authors",to: "id")}`},
		// 	},
		// },
	}

	schemaMod := &Schema{clusterID: "chicago"}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCrud := mockCrudSchemaInterface{}

			for _, m := range tt.crudMockArgs {
				mockCrud.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			schemaMod.crud = &mockCrud

			if err := schemaMod.SetDatabaseSchema(tt.crud, "myproject"); err != nil {
				t.Errorf("Manager.GetSchemas() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, err := schemaMod.GetSchemaForDB(context.Background(), tt.args.dbAlias, tt.args.col, tt.args.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.GetSchemas() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(got, tt.want1); diff != nil {
				if diff2 := deep.Equal(got, tt.want2); diff2 != nil {
					a, _ := json.MarshalIndent(diff2, "", " ")
					t.Error("Manager.GetSchemas() error2 \n", string(a))
				} else {
					return
				}
				a, _ := json.MarshalIndent(diff, "", " ")
				t.Error("Manager.GetSchemas() error1 \n", string(a))
			}
		})
	}
}
