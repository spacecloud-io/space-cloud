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
	fields, foreignkeys, indexes, err := s.crud.DescribeTable(ctx, dbAlias, col)

	if err != nil {
		return nil, err
	}
	return generateInspection(dbType, col, fields, foreignkeys, indexes)
}

func generateInspection(dbType, col string, fields []model.InspectorFieldType, foreignkeys []model.ForeignKeysType, indexes []model.IndexType) (model.Collection, error) {
	inspectionCollection := model.Collection{}
	inspectionFields := model.Fields{}

	for _, field := range fields {
		fieldDetails := model.FieldType{FieldName: field.FieldName}

		// check if field nullable (!)
		if field.FieldNull == "NO" {
			fieldDetails.IsFieldTypeRequired = true
		}

		// field type
		switch model.DBType(dbType) {
		case model.Postgres:
			if err := inspectionPostgresCheckFieldType(field.VarcharSize, field.FieldType, &fieldDetails); err != nil {
				return nil, err
			}
		case model.MySQL:
			if err := inspectionMySQLCheckFieldType(field.VarcharSize, field.FieldType, &fieldDetails); err != nil {
				return nil, err
			}
		case model.SQLServer:
			if err := inspectionSQLServerCheckFieldType(field.VarcharSize, field.FieldType, &fieldDetails); err != nil {
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

		// check if list
		if field.FieldKey == "PRI" {
			fieldDetails.IsPrimary = true
			fieldDetails.PrimaryKeyInfo = &model.TableProperties{}
			// Set auto increment
			if field.AutoIncrement == "true" {
				fieldDetails.PrimaryKeyInfo.IsAutoIncrement = true
			}
			if model.DBType(dbType) == model.Postgres && strings.HasPrefix(field.FieldDefault, "nextval") {
				// override the default value, this is a special case if a postgres column has a auto increment value, the default value that database returns is -> ( nextval(auto_increment_test_auto_increment_test_seq )
				fieldDetails.Default = ""
				fieldDetails.IsDefault = false
				fieldDetails.PrimaryKeyInfo.IsAutoIncrement = true
			}
		}

		// check foreignKey & identify if relation exists
		for _, foreignValue := range foreignkeys {
			if foreignValue.ColumnName == field.FieldName && foreignValue.RefTableName != "" && foreignValue.RefColumnName != "" {
				fieldDetails.IsForeign = true
				fieldDetails.JointTable = &model.TableProperties{Table: foreignValue.RefTableName, To: foreignValue.RefColumnName, OnDelete: foreignValue.DeleteRule, ConstraintName: foreignValue.ConstraintName}
			}
		}
		for _, indexValue := range indexes {
			if indexValue.ColumnName == field.FieldName {
				fieldDetails.IsIndex = true
				fieldDetails.IsUnique = indexValue.IsUnique == "yes"
				fieldDetails.IndexInfo = &model.TableProperties{Order: indexValue.Order, Sort: indexValue.Sort, ConstraintName: indexValue.IndexName}
				if strings.HasPrefix(indexValue.IndexName, "index__") {
					// index is created through gateway, as it follows our naming convention
					indexValue.IndexName = getGroupNameFromIndexName(indexValue.IndexName)
				}
				fieldDetails.IndexInfo.Group = indexValue.IndexName
			}
		}
		// field name
		inspectionFields[field.FieldName] = &fieldDetails
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

func inspectionMySQLCheckFieldType(size int, typeName string, fieldDetails *model.FieldType) error {
	if typeName == "varchar(-1)" || typeName == "varchar(max)" {
		fieldDetails.Kind = model.TypeString
		return nil
	}
	if strings.HasPrefix(typeName, "varchar(") {
		fieldDetails.Kind = model.TypeID
		fieldDetails.TypeIDSize = size
		return nil
	}

	result := strings.Split(typeName, "(")

	switch result[0] {
	case "date":
		fieldDetails.Kind = model.TypeDate
	case "time":
		fieldDetails.Kind = model.TypeTime
	case "char", "tinytext", "text", "blob", "mediumtext", "mediumblob", "longtext", "longblob", "decimal":
		fieldDetails.Kind = model.TypeString
	case "smallint", "mediumint", "int", "bigint":
		fieldDetails.Kind = model.TypeInteger
	case "float", "double":
		fieldDetails.Kind = model.TypeFloat
	case "datetime", "timestamp", "datetimeoffset":
		fieldDetails.Kind = model.TypeDateTime
	case "tinyint", "boolean", "bit":
		fieldDetails.Kind = model.TypeBoolean
	case "json":
		fieldDetails.Kind = model.TypeJSON
	default:
		return errors.New("Inspection type check : no match found got " + result[0])
	}
	return nil
}

func inspectionSQLServerCheckFieldType(size int, typeName string, fieldDetails *model.FieldType) error {
	if typeName == "varchar(-1)" || typeName == "varchar(max)" {
		fieldDetails.Kind = model.TypeString
		return nil
	}
	if typeName == "nvarchar(-1)" || typeName == "nvarchar(max)" {
		fieldDetails.Kind = model.TypeString
		return nil
	}
	if strings.HasPrefix(typeName, "varchar(") {
		fieldDetails.Kind = model.TypeID
		fieldDetails.TypeIDSize = size
		return nil
	}

	result := strings.Split(typeName, "(")

	switch result[0] {
	case "date":
		fieldDetails.Kind = model.TypeDate
	case "time":
		fieldDetails.Kind = model.TypeTime
	case "char", "tinytext", "text", "blob", "mediumtext", "mediumblob", "longtext", "longblob", "decimal", "nvarchar":
		fieldDetails.Kind = model.TypeString
	case "smallint", "mediumint", "int", "bigint":
		fieldDetails.Kind = model.TypeInteger
	case "float", "double":
		fieldDetails.Kind = model.TypeFloat
	case "datetime", "timestamp", "datetimeoffset":
		fieldDetails.Kind = model.TypeDateTime
	case "tinyint", "boolean", "bit":
		fieldDetails.Kind = model.TypeBoolean
	case "json":
		fieldDetails.Kind = model.TypeJSON
	default:
		return errors.New("Inspection type check : no match found got " + result[0])
	}
	return nil
}

func inspectionPostgresCheckFieldType(size int, typeName string, fieldDetails *model.FieldType) error {
	if typeName == "character varying" {
		fieldDetails.Kind = model.TypeID
		fieldDetails.TypeIDSize = size
		return nil
	}

	result := strings.Split(typeName, " ")
	result = strings.Split(result[0], "(")

	switch result[0] {
	case "uuid":
		fieldDetails.Kind = model.TypeUUID
	case "date":
		fieldDetails.Kind = model.TypeDate
	case "time":
		fieldDetails.Kind = model.TypeTime
	case "character", "bit", "text":
		fieldDetails.Kind = model.TypeString
	case "bigint", "bigserial", "integer", "smallint", "smallserial", "serial":
		fieldDetails.Kind = model.TypeInteger
	case "float", "double", "real", "numeric":
		fieldDetails.Kind = model.TypeFloat
	case "datetime", "timestamp", "interval", "datetimeoffset":
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
