package auth

import "errors"

// ErrRuleNotFound is thrown when an error is not present in the auth object
var ErrRuleNotFound = errors.New("auth: No rule has been provided")

// ErrIncorrectRuleFieldType is thrown when the field type of a rule is of incorrect type
var ErrIncorrectRuleFieldType = errors.New("auth: Incorrect rule field type")

// ErrIncorrectMatch is thrown when the field type of a rule is of incorrect type
var ErrIncorrectMatch = errors.New("auth: The two fields do not match")
