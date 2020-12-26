package schema

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

// SchemaInspection returns the schema in schema definition language (SDL)
func (s *Schema) SchemaInspection(ctx context.Context, dbAlias, project, col string) (string, error) {
	dbType, err := s.crud.GetDBType(dbAlias)
	if err != nil {
		return "", err
	}

	if dbType == "mongo" {
		return "", nil
	}

	inspectionCollection, err := s.Inspector(ctx, dbAlias, dbType, project, col)
	if err != nil {
		return "", err
	}

	return generateSDL(inspectionCollection)

}

// Inspector generates schema
func (s *Schema) Inspector(ctx context.Context, dbAlias, dbType, project, col string) (model.Collection, error) {
	fields, indexes, err := s.crud.DescribeTable(ctx, dbAlias, col)

	if err != nil {
		return nil, err
	}
	return generateInspection(dbType, col, fields, indexes)
}

func generateInspection(dbType, col string, fields []model.InspectorFieldType, indexes []model.IndexType) (model.Collection, error) {
	inspectionCollection := model.Collection{}
	inspectionFields := model.Fields{}

	for _, field := range fields {
		fieldDetails := model.FieldType{FieldName: field.ColumnName}

		// check if field nullable (!)
		if field.FieldNull == "NO" {
			fieldDetails.IsFieldTypeRequired = true
		}

		// field type
		if model.DBType(dbType) == model.Postgres {
			if err := inspectionPostgresCheckFieldType(field, &fieldDetails); err != nil {
				return nil, err
			}
		} else {
			if err := inspectionMySQLCheckFieldType(field, &fieldDetails); err != nil {
				return nil, err
			}
		}

		// default key
		if field.FieldDefault != "" {
			fieldDetails.IsDefault = true
			if model.DBType(dbType) == model.SQLServer {
				if fieldDetails.Kind == model.TypeBoolean {
					if field.FieldDefault == "1" {
						field.FieldDefault = "true"
					} else {
						field.FieldDefault = "false"
					}
				}
			}

			// add string between quotes
			if fieldDetails.Kind == model.TypeString || fieldDetails.Kind == model.TypeID || fieldDetails.Kind == model.TypeDateTime {
				field.FieldDefault = fmt.Sprintf("\"%s\"", field.FieldDefault)
			}
			fieldDetails.Default = field.FieldDefault
		}

		// Set auto increment
		if field.AutoIncrement == "true" {
			fieldDetails.IsAutoIncrement = true
		}
		if model.DBType(dbType) == model.Postgres && strings.HasPrefix(field.FieldDefault, "nextval") {
			// override the default value, this is a special case if a postgres column has a auto increment value, the default value that database returns is -> ( nextval(auto_increment_test_auto_increment_test_seq )
			fieldDetails.Default = ""
			fieldDetails.IsDefault = false
			fieldDetails.IsAutoIncrement = true
		}

		// check foreignKey & identify if relation exists
		if field.RefTableName != "" && field.RefColumnName != "" {
			fieldDetails.IsForeign = true
			fieldDetails.JointTable = &model.TableProperties{Table: field.RefTableName, To: field.RefColumnName, OnDelete: field.DeleteRule, ConstraintName: field.ConstraintName}
		}

		for _, indexValue := range indexes {
			if indexValue.ColumnName == field.ColumnName {
				if indexValue.IsUnique {
					fieldDetails.IsUnique = true
				} else if indexValue.IsPrimary {
					fieldDetails.IsUnique = false
					fieldDetails.IsPrimary = true
					continue
				} else {
					fieldDetails.IsIndex = true
				}

				fieldDetails.IndexInfo = &model.TableProperties{Order: indexValue.Order, Sort: indexValue.Sort, ConstraintName: indexValue.IndexName}
				if strings.HasPrefix(indexValue.IndexName, "index__") {
					// index is created through gateway, as it follows our naming convention
					indexValue.IndexName = getGroupNameFromIndexName(indexValue.IndexName)
				}
				fieldDetails.IndexInfo.Group = indexValue.IndexName
			}
		}
		// field name
		inspectionFields[field.ColumnName] = &fieldDetails
	}

	if len(inspectionFields) != 0 {
		inspectionCollection[col] = inspectionFields
	}
	return inspectionCollection, nil
}

func getGroupNameFromIndexName(indexName string) string {
	// ignoring the length check as the length is assured to be 3
	return strings.Split(indexName, "__")[2]
}

func inspectionMySQLCheckFieldType(field model.InspectorFieldType, fieldDetails *model.FieldType) error {
	if field.FieldType == "varchar(-1)" || field.FieldType == "varchar(max)" {
		fieldDetails.Kind = model.TypeString
		return nil
	}

	result := strings.Split(field.FieldType, "(")

	switch result[0] {
	case "varchar":
		fieldDetails.Kind = model.TypeID
		fieldDetails.TypeIDSize = field.VarcharSize
	case "date":
		fieldDetails.Kind = model.TypeDate
	case "time":
		fieldDetails.Kind = model.TypeTime
		if field.NumericScale > 0 {
			fieldDetails.Args = &model.FieldArgs{
				Scale: field.NumericScale,
			}
		}
	case "char", "tinytext", "text", "blob", "mediumtext", "mediumblob", "longtext", "longblob":
		fieldDetails.Kind = model.TypeString
	case "smallint", "mediumint", "int", "bigint":
		fieldDetails.Kind = model.TypeInteger
	case "float", "double", "decimal":
		fieldDetails.Kind = model.TypeFloat
		fieldDetails.Args = &model.FieldArgs{
			Precision: field.NumericPrecision,
			Scale:     field.NumericScale,
		}
	case "datetime", "timestamp", "datetimeoffset":
		fieldDetails.Kind = model.TypeDateTime
		if field.NumericScale > 0 {
			fieldDetails.Args = &model.FieldArgs{
				Scale: field.NumericScale,
			}
		}
	case "tinyint", "boolean", "bit":
		fieldDetails.Kind = model.TypeBoolean
	case "json":
		fieldDetails.Kind = model.TypeJSON
	default:
		return errors.New("Inspection type check : no match found got " + result[0])
	}
	return nil
}

func inspectionPostgresCheckFieldType(field model.InspectorFieldType, fieldDetails *model.FieldType) error {
	result := strings.Split(field.FieldType, "(")

	switch result[0] {
	case "character varying":
		fieldDetails.Kind = model.TypeID
		fieldDetails.TypeIDSize = field.VarcharSize
	case "uuid":
		fieldDetails.Kind = model.TypeUUID
	case "date":
		fieldDetails.Kind = model.TypeDate
	case "time", "time without time zone":
		fieldDetails.Kind = model.TypeTime
		if field.NumericScale > 0 {
			fieldDetails.Args = &model.FieldArgs{
				Scale: field.NumericScale,
			}
		}
	case "character", "bit", "text":
		fieldDetails.Kind = model.TypeString
	case "bigint", "bigserial", "integer", "smallint", "smallserial", "serial":
		fieldDetails.Kind = model.TypeInteger
	case "float", "double", "real", "numeric", "double precision":
		fieldDetails.Kind = model.TypeFloat
	case "datetime", "timestamp", "interval", "datetimeoffset", "timestamp without time zone":
		fieldDetails.Kind = model.TypeDateTime
	case "boolean":
		fieldDetails.Kind = model.TypeBoolean
	case "jsonb", "json":
		fieldDetails.Kind = model.TypeJSON
	default:
		return errors.New("Inspection type check : no match found got " + result[0])
	}
	return nil
}

// GetCollectionSchema returns schemas of collection aka tables for specified project & database
func (s *Schema) GetCollectionSchema(ctx context.Context, project, dbType string) (map[string]*config.TableRule, error) {

	collections := []string{}
	for _, dbSchema := range s.dbSchemas {
		if dbSchema.DbAlias == dbType {
			collections = append(collections, dbSchema.Table)
			break
		}
	}

	projectConfig := config.Crud{}
	projectConfig[dbType] = &config.CrudStub{}
	for _, colName := range collections {
		if colName == "default" {
			continue
		}
		schema, err := s.SchemaInspection(ctx, dbType, project, colName)
		if err != nil {
			return nil, err
		}

		if projectConfig[dbType].Collections == nil {
			projectConfig[dbType].Collections = map[string]*config.TableRule{}
		}
		projectConfig[dbType].Collections[colName] = &config.TableRule{Schema: schema}
	}
	return projectConfig[dbType].Collections, nil
}
