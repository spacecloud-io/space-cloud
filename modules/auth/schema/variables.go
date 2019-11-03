package schema

type (

	// schemaType is the data structure for storing the parsed values of schema string
	schemaType       map[string]schemaCollection // key is database name
	schemaCollection map[string]schemaField      // key is collection name
	schemaField      map[string]*schemaFieldType // key is field name
	directiveArgs    map[string]string           // key is Directive's argument name
	fieldType        int

	schemaFieldType struct {
		IsFieldTypeRequired bool
		IsList              bool
		Kind                string
		Directive           string
		nestedObject        schemaField
		JointTable          tableProperties
	}
	tableProperties struct {
		TableName, TableField string
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
	typeObject   string = "Object"
	typeJoin     string = "Join"

	directiveUnique    string = "unique"
	directiveRelation  string = "relation"
	directivePrimary   string = "primary"
	directiveCreatedAt string = "createdAt"
	directiveUpdatedAt string = "updatedAt"
)
