package schema

import (
	"context"
	"errors"
	"strings"

	"github.com/spaceuptech/space-cloud/utils"
)

// SchemaInspection resturn schema in schema definition language (SDL)
func (s *Schema) SchemaInspection(ctx context.Context, dbType, project, col string) (string, error) {
	if dbType == "mongo" {
		return "", errors.New("Inspection cannot be performed over mongo")
	}

	fields, foreignkeys, err := s.crud.DescribeTable(ctx, dbType, project, col)
	if err != nil {
		return "", err
	}
	inspectionDb := schemaType{}
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
				return "", err
			}
		} else {
			if err := inspectionCheckFieldType(value.FieldType, &fieldDetails); err != nil {
				return "", err
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
				fieldDetails.Kind = typeJoin
				fieldDetails.Directive = "relation"
			}
		}

		// field name
		inspectionFields[value.FieldName] = &fieldDetails

	}

	inspectionCollection[strings.Title(col)] = inspectionFields
	inspectionDb[dbType] = inspectionCollection

	schemaInSDL, err := generateSDL(inspectionCollection)
	if err != nil {
		return "", nil
	}
	return schemaInSDL, nil

}

func inspectionCheckFieldType(typeName string, fieldDetails *schemaFieldType) error {

	result := strings.Split(typeName, "(")

	switch result[0] {
	case "char", "varchar", "tinytext", "text", "blob", "mediumtext", "mediumblob", "longtext", "longblob", "decimal":
		fieldDetails.Kind = typeString
	case "tinyint", "smallint", "mediumint", "int", "bigint":
		fieldDetails.Kind = typeInteger
	case "float", "double":
		fieldDetails.Kind = typeFloat
	case "date", "time", "datetime", "timestamp":
		fieldDetails.Kind = typeDateTime
	case "boolean":
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
	case "character", "bit":
		fieldDetails.Kind = typeString
	case "bigint", "bigserial", "integer", "numeric", "smallint", "smallserial", "serial", "text":
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
