package types

import "github.com/graphql-go/graphql/gqlerrors"

// GraphqlType is our version
type GraphqlType interface {
	GetName() string
	GetKind() string
}

type (
	// GraphqlError stores the errors encountered while processing the graphql request
	GraphqlError struct {
		FormatedErrors gqlerrors.FormattedErrors
	}
)

// Error returns a string representation of the error
func (e *GraphqlError) Error() string {
	return "graphql query returned with an error"
}
