package types

type (
	// AuthResult describes whether or not the request is authenticated
	AuthResult struct {
		IsAuthenticated bool
		Claims          map[string]interface{}
	}
)
