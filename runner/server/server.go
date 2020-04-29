package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/spaceuptech/space-cloud/runner/metrics"

	"github.com/dgraph-io/badger"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/runner/model"
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

	// For internal use
	auth     *auth.Module
	driver   driver.Interface
	debounce *utils.Debounce

	// For autoscaler
	db       *badger.DB
	chAppend chan *model.ProxyMessage
}

// New creates a new instance of the runner
func New(c *Config) (*Server, error) {
	// Add the proxy port to the driver config
	proxyPort, err := strconv.Atoi(c.ProxyPort)
	if err != nil {
		return nil, fmt.Errorf("invalid proxy port (%s) provided", c.ProxyPort)
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

	opts := badger.DefaultOptions("/tmp/runner.db")
	opts.Logger = &logrus.Logger{Out: ioutil.Discard}
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	// Periodically run the garbage collector
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
		again:
			err := db.RunValueLogGC(0.7)
			if err == nil {
				goto again
			}
		}
	}()

	// Return a new runner instance
	return &Server{
		config: c,
		router: mux.NewRouter(),

		// For internal use
		auth:     a,
		driver:   d,
		debounce: debounce,

		// For autoscaler
		db:       db,
		chAppend: make(chan *model.ProxyMessage, 10),
	}, nil
}

// Start begins the runner
func (s *Server) Start() error {
	// Initialise the various routes of the s
	s.routes()

	// Start necessary routines for autoscaler
	go s.routineAdjustScale()
	for i := 0; i < 10; i++ {
		go s.routineDumpDetails()
	}

	// Start proxy server
	go func() {
		// Create a new router
		router := mux.NewRouter()
		router.PathPrefix("/").HandlerFunc(s.handleProxy())

		// Start http server
		corsObj := utils.CreateCorsObject()
		logrus.Infof("Starting server proxy on port %s", s.config.ProxyPort)
		if err := http.ListenAndServe(":"+s.config.ProxyPort, corsObj.Handler(router)); err != nil {
			logrus.Fatalln("Proxy server failed:", err)
		}
	}()

	// Start the http server
	corsObj := utils.CreateCorsObject()
	logrus.Infof("Starting server on port %s", s.config.Port)
	return http.ListenAndServe(":"+s.config.Port, corsObj.Handler(s.router))
}
