package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/runner/metrics"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/runner/utils"
	"github.com/spaceuptech/space-cloud/runner/utils/auth"
	"github.com/spaceuptech/space-cloud/runner/utils/driver"
)

// Server is the module responsible to manage the runner
type Server struct {
	// For storing config
	config *Config

	// For handling http related stuff
	router *mux.Router

	// For sending metrics to runner
	metrics *metrics.Module

	// For internal use
	auth     *auth.Module
	driver   driver.Interface
	debounce *utils.Debounce
}

// New creates a new instance of the runner
func New(c *Config) (*Server, error) {
	// Add the proxy port to the driver config
	proxyPort, err := strconv.Atoi(c.ProxyPort)
	if err != nil {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("invalid proxy port (%s) provided", c.ProxyPort), err, nil)
	}
	c.Driver.ProxyPort = uint32(proxyPort)

	metric := metrics.New(c.IsMetricDisabled, c.Driver.DriverType)

	// Initialise all modules
	a, err := auth.New(c.Auth)
	if err != nil {
		return nil, err
	}

	d, err := driver.New(a, c.Driver, metric.AddServiceCall)
	if err != nil {
		return nil, err
	}

	debounce := utils.NewDebounce()

	// Return a new runner instance
	return &Server{
		config: c,
		router: mux.NewRouter(),

		metrics: metric,

		// For internal use
		auth:     a,
		driver:   d,
		debounce: debounce,
	}, nil
}

// Start begins the runner
func (s *Server) Start() error {
	// Initialise the various routes of the s
	s.routes()

	// Start proxy server
	go func() {
		// Create a new router
		router := mux.NewRouter()
		router.PathPrefix("/").HandlerFunc(s.handleProxy())

		// Start http server
		corsObj := utils.CreateCorsObject()
		helpers.Logger.LogInfo(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Starting server proxy on port %s", s.config.ProxyPort), nil)
		if err := http.ListenAndServe(":"+s.config.ProxyPort, corsObj.Handler(router)); err != nil {
			helpers.Logger.LogFatal(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Proxy server failed: - %v", err), nil)
		}
	}()

	// Start the http server
	corsObj := utils.CreateCorsObject()
	helpers.Logger.LogInfo(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Starting server on port %s", s.config.Port), nil)
	return http.ListenAndServe(":"+s.config.Port, corsObj.Handler(loggerMiddleWare(s.router)))
}
