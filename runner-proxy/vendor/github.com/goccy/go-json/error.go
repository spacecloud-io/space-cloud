package json

import (
	"fmt"
	"reflect"
	"strconv"
)

// Before Go 1.2, an InvalidUTF8Error was returned by Marshal when
// attempting to encode a string value with invalid UTF-8 sequences.
// As of Go 1.2, Marshal instead coerces the string to valid UTF-8 by
// replacing invalid bytes with the Unicode replacement rune U+FFFD.
//
// Deprecated: No longer used; kept for compatibility.
type InvalidUTF8Error struct {
	S string // the whole string value that caused the error
}

func (e *InvalidUTF8Error) Error() string {
	return fmt.Sprintf("json: invalid UTF-8 in string: %s", strconv.Quote(e.S))
}

// An InvalidUnmarshalError describes an invalid argument passed to Unmarshal.
// (The argument to Unmarshal must be a non-nil pointer.)
type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "json: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return fmt.Sprintf("json: Unmarshal(non-pointer %s)", e.Type)
	}
	return fmt.Sprintf("json: Unmarshal(nil %s)", e.Type)
}

// A MarshalerError represents an error from calling a MarshalJSON or MarshalText method.
type MarshalerError struct {
	Type       reflect.Type
	Err        error
	sourceFunc string
}

func (e *MarshalerError) Error() string {
	srcFunc := e.sourceFunc
	if srcFunc == "" {
		srcFunc = "MarshalJSON"
	}
	return fmt.Sprintf("json: error calling %s for type %s: %s", srcFunc, e.Type, e.Err.Error())
}

// Unwrap returns the underlying error.
func (e *MarshalerError) Unwrap() error { return e.Err }

// A SyntaxError is a description of a JSON syntax error.
type SyntaxError struct {
	msg    string // description of error
	Offset int64  // error occurred after reading Offset bytes
}

func (e *SyntaxError) Error() string { return e.msg }

// An UnmarshalFieldError describes a JSON object key that
// led to an unexported (and therefore unwritable) struct field.
//
// Deprecated: No longer used; kept for compatibility.
type UnmarshalFieldError struct {
	Key   string
	Type  reflect.Type
	Field reflect.StructField
}

func (e *UnmarshalFieldError) Error() string {
	return fmt.Sprintf("json: cannot unmarshal object key %s into unexported field %s of type %s",
		strconv.Quote(e.Key), e.Field.Name, e.Type.String(),
	)
}

// An UnmarshalTypeError describes a JSON value that was
// not appropriate for a value of a specific Go type.
type UnmarshalTypeError struct {
	Value  string       // description of JSON value - "bool", "array", "number -5"
	Type   reflect.Type // type of Go value it could not be assigned to
	Offset int64        // error occurred after reading Offset bytes
	Struct string       // name of the struct type containing the field
	Field  string       // the full path from root node to the field
}

func (e *UnmarshalTypeError) Error() string {
	if e.Struct != "" || e.Field != "" {
		return fmt.Sprintf("json: cannot unmarshal %s into Go struct field %s.%s of type %s",
			e.Value, e.Struct, e.Field, e.Type,
		)
	}
	return fmt.Sprintf("json: cannot unmarshal %s into Go value of type %s", e.Value, e.Type)
}

// An UnsupportedTypeError is returned by Marshal when attempting
// to encode an unsupported value type.
type UnsupportedTypeError struct {
	Type reflect.Type
}

func (e *UnsupportedTypeError) Error() string {
	return fmt.Sprintf("json: unsupported type: %s", e.Type)
}

type UnsupportedValueError struct {
	Value reflect.Value
	Str   string
}

func (e *UnsupportedValueError) Error() string {
	return fmt.Sprintf("json: unsupported value: %s", e.Str)
}

func errExceededMaxDepth(c byte, cursor int64) *SyntaxError {
	return &SyntaxError{
		msg:    fmt.Sprintf(`invalid character "%c" exceeded max depth`, c),
		Offset: cursor,
	}
}

func errNotAtBeginningOfValue(cursor int64) *SyntaxError {
	return &SyntaxError{msg: "not at beginning of value", Offset: cursor}
}

func errUnexpectedEndOfJSON(msg string, cursor int64) *SyntaxError {
	return &SyntaxError{
		msg:    fmt.Sprintf("json: %s unexpected end of JSON input", msg),
		Offset: cursor,
	}
}

func errExpected(msg string, cursor int64) *SyntaxError {
	return &SyntaxError{msg: fmt.Sprintf("expected %s", msg), Offset: cursor}
}

func errInvalidCharacter(c byte, context string, cursor int64) *SyntaxError {
	if c == 0 {
		return &SyntaxError{
			msg:    fmt.Sprintf("json: invalid character as %s", context),
			Offset: cursor,
		}
	}
	return &SyntaxError{
		msg:    fmt.Sprintf("json: invalid character %c as %s", c, context),
		Offset: cursor,
	}
}
