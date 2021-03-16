package types

import (
	"errors"

	"github.com/mitchellh/mapstructure"
)

// Response is the object recieved from the server. Status is the status code received from the server.
// Data is either a map[string]interface{} or
type Response struct {
	Status int
	Data   M
	Error  string
}

// Unmarshal parses the response data and stores the result in the value pointed to by v.
// If v is nil or not a pointer, Unmarshal returns an error.
func (res *Response) Unmarshal(v interface{}) error {
	if res.Status < 200 || res.Status >= 300 {
		return errors.New("Result not present")
	}
	return mapstructure.Decode(res.Data["result"], v)
}

// Raw returns the raw map
func (res *Response) Raw() map[string]interface{} {
	return res.Data
}

// GetStatus returns the response status
func (res *Response) GetStatus() int {
	return res.Status
}

// GetError returns the error message receieved from the server
func (res *Response) GetError() error {
	return errors.New(res.Error)
}
