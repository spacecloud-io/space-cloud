package types

type (
	// StoreValue stores the value of the exported variable along with its type
	StoreValue struct {
		Value  interface{}
		TypeOf string
	}

	// StoreValue stores the value of the exported variable along with its key & type
	StoreKeyValue struct {
		Key   string
		Value StoreValue
	}
)
