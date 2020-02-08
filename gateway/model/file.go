package model

import "io"

// File is the struct returned for file reads
type File struct {
	File  io.Reader
	Close func() error
}

// CreateFileRequest is the request received to create a new file or directory
type CreateFileRequest struct {
	Meta    map[string]interface{} `json:"meta"`
	Path    string                 `json:"path"`
	Name    string                 `json:"name"`
	Type    string                 `json:"type"`    // Either file or dir
	MakeAll bool                   `json:"makeAll"` // This option is only available for creating directories
}

// DeleteFileRequest is the request received to delete a new file or directory
type DeleteFileRequest struct {
	Meta map[string]interface{} `json:"meta"`
	Path string                 `json:"path"`
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

// FilePayload is body of request to file module
type FilePayload struct {
	Meta map[string]interface{} `json:"meta"`
	Path string                 `json:"path"`
}

// FileReader is a function type used for file streaming
type FileReader func(io.Reader) (int, error)
