package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/spaceuptech/helpers"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/runner-proxy/utils"
)

// Server is the module responsible to manage the runner
type Server struct {
	// For storing config
	config *Config

	// For handling http related stuff
	router *mux.Router

	debounce *utils.Debounce
	m        *TTLMap
}

// New creates a new instance of the runner
func New(c *Config) (*Server, error) {
	debounce := utils.NewDebounce()
	m := tick()
	// Return a new runner instance
	return &Server{
		config:   c,
		router:   mux.NewRouter(),
		debounce: debounce,
		m:        m,
	}, nil
}

// Start begins the runner
func (s *Server) Start() error {
	// Initialise the various routes of the s
	s.routes()

	// Start http server
	corsObj := utils.CreateCorsObject()
	helpers.Logger.LogInfo(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Starting server proxy on port %s", s.config.Port), nil)
	return http.ListenAndServe(":"+s.config.Port, corsObj.Handler(loggerMiddleWare(s.router)))
}
