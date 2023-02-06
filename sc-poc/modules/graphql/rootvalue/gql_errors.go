package rootvalue

import "github.com/graphql-go/graphql/gqlerrors"

// AddFormatedErrors stores an array of formated graphql errors in the root value
func (root *RootValue) AddFormatedErrors(errs gqlerrors.FormattedErrors) {
	root.errMutex.Lock()
	defer root.errMutex.Unlock()

	root.formatedErrors = append(root.formatedErrors, errs...)
}

// GetFormatedErrors returns a list of formated errors
func (root *RootValue) GetFormatedErrors() gqlerrors.FormattedErrors {
	root.errMutex.Lock()
	defer root.errMutex.Unlock()

	return root.formatedErrors
}

// HasErrors checks if any errors had occured while processing the graphql request
func (root *RootValue) HasErrors() bool {
	root.errMutex.Lock()
	defer root.errMutex.Unlock()

	return len(root.formatedErrors) > 0
}
