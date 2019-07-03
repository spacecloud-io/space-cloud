package model

import "io"

// File is the struct returned for file reads
type File struct {
	File  io.Reader
	Close func() error
}

// CreateFileRequest is the request received to create a new file or directory
type CreateFileRequest struct {
	Path    string `json:"path"`
	Name    string `json:"name"`
	Type    string `json:"type"`    // Either file or dir
	MakeAll bool   `json:"makeAll"` // This option is only available for creating directories
}

// ListFilesRequest is the request made to browse the contents inside a directory
type ListFilesRequest struct {
	Path string `json:"path"`
	Type string `json:"type"` // Type could be dir, file or all
}

// ListFilesResponse is the response given to browse the contents inside a directory
type ListFilesResponse struct {
	Name string `json:"name"`
	Type string `json:"type"` // Type could be dir or file
}

// FileReader is a function type used for file streaming
type FileReader func(io.Reader) (int, error)
