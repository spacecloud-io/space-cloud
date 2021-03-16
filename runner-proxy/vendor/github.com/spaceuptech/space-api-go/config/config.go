package config

import (
	"github.com/spaceuptech/space-api-go/transport"
)

// Config holds the config of the API object
type Config struct {
	Project   string
	URL       string
	Token     string
	IsSecure  bool
	Transport *transport.Transport
}
