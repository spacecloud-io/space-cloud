package utils

import "errors"

// ErrInvalidParams is thrown when the input parameters for an operation are invalid
var ErrInvalidParams = errors.New("Invalid parameter provided")

// ErrDatabaseDisabled is thrown when an operation is requested on a disabled database
var ErrDatabaseDisabled = errors.New("Database is disabled. Please enable it")

// ErrUnsupportedDatabase is thrown when an invalid db type is provided
var ErrUnsupportedDatabase = errors.New("Unsupported database. Make sure your database type is correct")

// ErrDatabaseConnection is thrown when SC was unable to connect to the requested database
var ErrDatabaseConnection = errors.New("Could not connect to database. Make sure it is up and connection string provided to SC is correct")

// ErrDatabaseConfigAbsent is thrown when database config is not present
var ErrDatabaseConfigAbsent = errors.New("No such database found in SC config file")
