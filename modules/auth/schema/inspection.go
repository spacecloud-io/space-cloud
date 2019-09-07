package schema

import (
	"context"
	"errors"
	"strings"
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
			fieldDetails.isFieldTypeRequired = true
		}

		// field type
		if err := inspectionCheckFieldType(value.FieldType, &fieldDetails); err != nil {
			return "", err
		}

		// check if list
		if value.FieldKey == "PRI" {
			fieldDetails.directive.Kind = "id"
			fieldDetails.Kind = typeID
		} else if value.FieldKey == "UNI" {
			fieldDetails.directive.Kind = "unique"
		} else {
			fieldDetails.isList = true
		}

		// check foreignKey & identify if relation exists
		for _, foreignValue := range foreignkeys {
			if foreignValue.ColumnName == value.FieldName {
				fieldDetails.tableJoin = foreignValue.RefTableName
				fieldDetails.Kind = typeJoin
				fieldDetails.directive.Kind = "relation"
			}
		}

		// field name
		inspectionFields[value.FieldName] = &fieldDetails

	}

	inspectionCollection[col] = inspectionFields
	inspectionDb[dbType] = inspectionCollection

	schemaInSDL, err := generateSDL(inspectionCollection)
	if err != nil {
		return "", nil
	}
	return schemaInSDL, nil

}

func inspectionCheckFieldType(typeName string, fieldDetails *schemaFieldType) error {
	// TODO: what about my-sql set type

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
		return errors.New("Inspection type check : no match found")
	}
	return nil
}
