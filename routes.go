package main

import (
	"github.com/spaceuptech/space-cloud/config"
)

func (s *server) routes() {
	// Initialize the routes for config management
	s.router.Methods("POST").Path("/v1/api/config").HandlerFunc(config.HandleConfig(s.isProd, s.loadConfig))

	// Initialize the route for websocket
	s.router.HandleFunc("/v1/api/socket/json", s.handleWebsocket())

	// Initialize the routes for functions service
	s.router.Methods("POST").Path("/v1/api/functions/{service}/{func}").HandlerFunc(s.functions.HandleRequest(s.auth))

	// Initialize the routes for the crud operations
	s.router.Methods("POST").Path("/v1/api/{project}/crud/{dbType}/batch").HandlerFunc(s.handleBatch())

	crudRouter := s.router.Methods("POST").PathPrefix("/v1/api/{project}/crud/{dbType}/{col}").Subrouter()
	crudRouter.HandleFunc("/create", s.handleCreate())
	crudRouter.HandleFunc("/read", s.handleRead())
	crudRouter.HandleFunc("/update", s.handleUpdate())
	crudRouter.HandleFunc("/delete", s.handleDelete())
	crudRouter.HandleFunc("/aggr", s.handleAggregate())

	// Initialize the routes for the user management operations
	userRouter := s.router.PathPrefix("/v1/api/{project}/auth/{dbType}").Subrouter()
	userRouter.Methods("POST").Path("/email/signin").HandlerFunc(s.user.HandleEmailSignIn())
	userRouter.Methods("POST").Path("/email/signup").HandlerFunc(s.user.HandleEmailSignUp())
	userRouter.Methods("GET").Path("/profile/{id}").HandlerFunc(s.user.HandleProfile())
	userRouter.Methods("GET").Path("/profiles").HandlerFunc(s.user.HandleProfiles())

	// Initialize the routes for the file management operations
	s.router.Methods("POST").Path("/v1/api/{project}/files").HandlerFunc(s.file.HandleCreateFile(s.auth))
	s.router.Methods("GET").PathPrefix("/v1/api/{project}/files").HandlerFunc(s.file.HandleRead(s.auth))
	s.router.Methods("DELETE").PathPrefix("/v1/api/{project}/files").HandlerFunc(s.file.HandleDelete(s.auth))

	// Initialize the route for handling static files
	s.router.PathPrefix("/").HandlerFunc(s.static.HandleRequest)
}
