package schema

import (
	"context"
	"errors"
	"strings"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/utils"
)

// SchemaInspection resturn schema in schema definition language (SDL)
func (s *Schema) SchemaInspection(ctx context.Context, dbType, project, col string) (string, error) {
	if dbType == "mongo" {
		return "", errors.New("Inspection cannot be performed over mongo")
	}

	inspectionCollection, err := s.Inspector(ctx, dbType, project, col)
	if err != nil {
		return "", err
	}

	return generateSDL(inspectionCollection)

}

// Inspector does something
func (s *Schema) Inspector(ctx context.Context, dbType, project, col string) (schemaCollection, error) {
	fields, foreignkeys, err := s.crud.DescribeTable(ctx, dbType, project, col)
	if err != nil {
		return nil, err
	}
	inspectionCollection := schemaCollection{}
	inspectionFields := schemaField{}

	for _, value := range fields {
		fieldDetails := schemaFieldType{}

		// check if field nullable (!)
		if value.FieldNull == "NO" {
			fieldDetails.IsFieldTypeRequired = true
		}

		// field type
		if utils.DBType(dbType) == utils.Postgres {
			if err := inspectionPostgresCheckFieldType(value.FieldType, &fieldDetails); err != nil {
				return nil, err
			}
		} else {
			if err := inspectionMySQLCheckFieldType(value.FieldType, &fieldDetails); err != nil {
				return nil, err
			}
		}
		// check if list
		if value.FieldKey == "PRI" {
			fieldDetails.Directive = "id"
			fieldDetails.Kind = typeID
		} else if value.FieldKey == "UNI" {
			fieldDetails.Directive = "unique"
		}

		// check foreignKey & identify if relation exists
		for _, foreignValue := range foreignkeys {
			if foreignValue.ColumnName == value.FieldName {
				fieldDetails.JointTable.TableName = foreignValue.RefTableName
				fieldDetails.JointTable.TableField = foreignValue.RefColumnName
				fieldDetails.Kind = foreignValue.RefTableName
				fieldDetails.Directive = "relation"
			}
		}

		// field name
		inspectionFields[value.FieldName] = &fieldDetails
	}
	if len(inspectionFields) != 0 {
		inspectionCollection[col] = inspectionFields
	}
	return inspectionCollection, nil
}

func inspectionMySQLCheckFieldType(typeName string, fieldDetails *schemaFieldType) error {

	result := strings.Split(typeName, "(")

	switch result[0] {
	case "char", "varchar", "tinytext", "text", "blob", "mediumtext", "mediumblob", "longtext", "longblob", "decimal":
		fieldDetails.Kind = typeString
	case "smallint", "mediumint", "int", "bigint":
		fieldDetails.Kind = typeInteger
	case "float", "double":
		fieldDetails.Kind = typeFloat
	case "date", "time", "datetime", "timestamp":
		fieldDetails.Kind = typeDateTime
	case "tinyint", "boolean": // sql stores boolean valuse as tinyint(1), TODO: what if tinyint(28) then it should come under integer
		fieldDetails.Kind = typeBoolean
	default:
		return errors.New("Inspection type check : no match found got " + result[0])
	}
	return nil
}

func inspectionPostgresCheckFieldType(typeName string, fieldDetails *schemaFieldType) error {

	result := strings.Split(typeName, " ")
	result = strings.Split(result[0], "(")

	switch result[0] {
	case "character", "bit", "text":
		fieldDetails.Kind = typeString
	case "bigint", "bigserial", "integer", "numeric", "smallint", "smallserial", "serial":
		fieldDetails.Kind = typeInteger
	case "float", "double", "real":
		fieldDetails.Kind = typeFloat
	case "date", "time", "datetime", "timestamp", "interval":
		fieldDetails.Kind = typeDateTime
	case "boolean":
		fieldDetails.Kind = typeBoolean
	case "json":
		fieldDetails.Kind = typeJSON

	default:
		return errors.New("Inspection type check : no match found got " + result[0])
	}
	return nil
}

// GetCollectionSchema returns schemas of collection aka tables for specified project & database
func (s *Schema) GetCollectionSchema(ctx context.Context, project, dbType string) (map[string]*config.TableRule, error) {

	collections := []string{}
	for dbName, crudValue := range s.config {
		if dbName == dbType {
			for colName := range crudValue.Collections {
				collections = append(collections, colName)
			}
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
