package utils

import (
	"context"
	"net/http"

	"github.com/spacecloud-io/space-cloud/modules/auth/types"
)

// StoreAuthenticationResult returns a new request object with the auth result in the context object
func StoreAuthenticationResult(r *http.Request, result *types.AuthResult) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), authenticationResultKey, result))
}

// GetAuthenticationResult returns the result of the authentication middleware
func GetAuthenticationResult(r *http.Request) (*types.AuthResult, bool) {
	result := r.Context().Value(authenticationResultKey)
	if result == nil {
		return nil, false
	}

	return result.(*types.AuthResult), true
}

type (
	contextKey int
)

const (
	authenticationResultKey contextKey = iota
)
