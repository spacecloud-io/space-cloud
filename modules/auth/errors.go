package auth

import "errors"

// ErrRuleNotFound is thrown when an error is not present in the Auth object
var ErrRuleNotFound = errors.New("Auth: No rule has been provided")

// ErrIncorrectRuleFieldType is thrown when the field type of a rule is of incorrect type
var ErrIncorrectRuleFieldType = errors.New("Auth: Incorrect rule field type")

// ErrIncorrectMatch is thrown when the field type of a rule is of incorrect type
var ErrIncorrectMatch = errors.New("Auth: The two fields do not match")
