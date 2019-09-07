package schema

type (

	// schemaType is the data structure for storing the parsed values of schema string
	schemaType       map[string]schemaCollection // key is database name
	schemaCollection map[string]schemaField      // key is collection name
	schemaField      map[string]*schemaFieldType // key is field name
	directiveArgs    map[string]string           // key is directive's argument name
	fieldType        int

	schemaFieldType struct {
		isFieldTypeRequired bool
		isList              bool
		directive           directiveProperties
		Kind                string
		nestedObject        schemaField
		tableJoin           string
	}
	directiveProperties struct {
		Kind  string
		Value directiveArgs
	}
)

const (
	typeInteger  string = "Integer"
	typeString   string = "String"
	typeFloat    string = "Float"
	typeBoolean  string = "Boolean"
	typeDateTime string = "DateTime"
	typeEnum     string = "Enum"
	typeJSON     string = "Json"
	typeID       string = "ID"
	typeJoin     string = "Object"
	typeRelation string = "Join"
)
