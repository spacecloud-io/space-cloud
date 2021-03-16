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
	// For handling http related stuff
	router *mux.Router

	debounce *utils.Debounce
	cache    *utils.TTLMap
	auth     *utils.Auth
}

// New creates a new instance of the runner
func New(secret string) (*Server, error) {
	debounce := utils.NewDebounce()
	m := utils.Tick()
	a := utils.New(secret)
	// Return a new runner instance
	return &Server{
		router:   mux.NewRouter(),
		debounce: debounce,
		cache:    m,
		auth:     a,
	}, nil
}

// Start begins the runner
func (s *Server) Start(port string) error {
	// Initialise the various routes of the s
	s.routes()

	// Start http server
	corsObj := utils.CreateCorsObject()
	helpers.Logger.LogInfo(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Starting server proxy on port %s", port), nil)
	return http.ListenAndServe(":"+port, corsObj.Handler(loggerMiddleWare(s.router)))
}
