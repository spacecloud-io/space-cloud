package openapi

// GetAllTodosRequest
type GetAllTodosRequest struct {
	Eq *GetAllTodos_Eq `json:"_eq,omitempty"`
}

// GetAllTodos_Eq
type GetAllTodos_Eq struct {
}

// GetAllTodosResponse
type GetAllTodosResponse struct {
	HasuraPersons []GetAllTodos_HasuraPersons `json:"hasura_persons"`
}

// GetAllTodos_HasuraPersons
type GetAllTodos_HasuraPersons struct {
	Addresses []HasuraPersons_Addresses `json:"addresses"`
	Age       int32                     `json:"age"`
	Id        int32                     `json:"id"`
	Name      []string                  `json:"name"`
	Number    *HasuraPersons_Number     `json:"number,omitempty"`
}

// HasuraPersons_Addresses
type HasuraPersons_Addresses struct {
	City string `json:"city"`
	Id   int32  `json:"id"`
}

// HasuraPersons_Number
type HasuraPersons_Number struct {
	Number string `json:"number"`
	Id     int32  `json:"id"`
}
