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
		if err := inspectionCheckFieldType(value.FieldType, &fieldDetails); err != nil {
			return nil, err
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

	inspectionCollection[col] = inspectionFields
	inspectionDb[dbType] = inspectionCollection
	return inspectionCollection, nil
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
