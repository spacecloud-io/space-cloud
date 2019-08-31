package schema

type (
	Schema struct {
		SchemaDoc SchemaType
	}

	// SchemaType is the data structure for storing the parsed values of schema string
	SchemaType       map[string]SchemaCollection // key is database name
	SchemaCollection map[string]SchemaField      // key is collection name
	SchemaField      map[string]*SchemaFieldType // key is field name
	DirectiveArgs    map[string]string           // key is directive's argument name
	fieldType        int

	SchemaFieldType struct {
		IsFieldTypeRequired bool
		IsList              bool
		Directive           DirectiveProperties
		Kind                string
		NestedObject        SchemaField
		TableJoin           string
	}
	DirectiveProperties struct {
		Kind  string
		Value DirectiveArgs
	}
)

const (
	TypeString   string = "String"
	TypeInteger  string = "Integer"
	TypeFloat    string = "Float"
	TypeBoolean  string = "Boolean"
	TypeDateTime string = "DateTime"
	TypeEnum     string = "Enum"
	TypeJSON     string = "Json"
	TypeID       string = "ID"
	TypeJoin     string = "Object"
	TypeRelation string = "Join"
)
