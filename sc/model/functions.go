package model

import "github.com/spacecloud-io/space-cloud/config"

// FunctionsRequest is the api call request
type FunctionsRequest struct {
	Params  interface{}              `json:"params"`
	Timeout int                      `json:"timeout"`
	Cache   *config.ReadCacheOptions `json:"cache"`
}
