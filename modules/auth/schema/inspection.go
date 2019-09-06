package schema

import (
	"context"
	"errors"
	"strings"
)

func (s *Schema) schemaInspection(ctx context.Context, dbType, project, col string) (string, error) {
	if dbType == "mongo" {
		return "", errors.New("Inspection cannot be performed over mongo")
	}
	// todo change project ot s.project
	fields, foreignkeys, err := s.crud.DescribeTable(ctx, dbType, project, col)
	if err != nil {
		return "", err
	}

	inspectionDb := SchemaType{}
	inspectionCollection := SchemaCollection{}
	inspectionFields := SchemaField{}
	for _, value := range fields {
		fieldDetails := SchemaFieldType{}
		// dirARgs := DirectiveArgs{}
		// check if field nullable (!)
		if value.FieldNull == "NO" {
			fieldDetails.IsFieldTypeRequired = true
		}
		// field type
		if err := inspectionCheckFieldType(value.FieldType, &fieldDetails); err != nil {
			return "", err
		}
		// check if list
		if value.FieldKey == "PRI" {
			fieldDetails.Directive.Kind = "id"
			fieldDetails.Kind = TypeID
		} else if value.FieldKey == "UNI" {
			fieldDetails.Directive.Kind = "unique"
		} else {
			fieldDetails.IsList = true
		}
		// check foreignKey & identify if relation exists
		for _, foreignValue := range foreignkeys {
			if foreignValue.ColumnName == value.FieldName {
				fieldDetails.TableJoin = foreignValue.RefTableName
				fieldDetails.Kind = TypeJoin
				fieldDetails.Directive.Kind = "relation"
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

func inspectionCheckFieldType(typeName string, fieldDetails *SchemaFieldType) error {
	// TODO: what about my-sql set type
	// TODO: sql types int(11), varchar(255) is there any thing else
	result := strings.Split(typeName, "(")
	switch result[0] {
	case "char", "varchar", "tinytext", "text", "blob", "mediumtext", "mediumblob", "longtext", "longblob", "decimal":
		fieldDetails.Kind = TypeString
	case "tinyint", "smallint", "mediumint", "int", "bigint":
		fieldDetails.Kind = TypeInteger
	case "float", "double":
		fieldDetails.Kind = TypeFloat
	case "date", "time", "datetime", "timestamp":
		fieldDetails.Kind = TypeDateTime
	case "boolean":
		fieldDetails.Kind = TypeBoolean
	default:
		return errors.New("Inspection type check : no match found")
	}
	return nil
}
