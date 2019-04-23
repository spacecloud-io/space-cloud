package main

import "github.com/spaceuptech/space-cloud/config"

func (s *server) routes() {
	// Initialize the routes for config management
	s.router.Methods("POST").Path("/api/v1/config").HandlerFunc(config.HandleConfig(s.isProd, s.loadConfig))

	// Initialze the route for websocket
	s.router.HandleFunc("/api/v1/socket/json", handleWebsocket(s.realtime, s.auth, s.crud))

	// Initialize the routes for faas engine
	s.router.Methods("POST").Path("/api/v1/faas/{engine}/{func}").HandlerFunc(s.faas.HandleRequest(s.auth))

	// Initialize the routes for the crud operations
	s.router.Methods("POST").Path("/api/v1/{project}/crud/{dbType}/batch").HandlerFunc(s.handleBatch())

	crudRouter := s.router.Methods("POST").PathPrefix("/api/v1/{project}/crud/{dbType}/{col}").Subrouter()
	crudRouter.HandleFunc("/create", s.handleCreate())
	crudRouter.HandleFunc("/read", s.handleRead())
	crudRouter.HandleFunc("/update", s.handleUpdate())
	crudRouter.HandleFunc("/delete", s.handleDelete())
	crudRouter.HandleFunc("/aggr", s.handleAggregate())

	// Initialize the routes for the user management operations
	userRouter := s.router.PathPrefix("/api/v1/{project}/auth/{dbType}").Subrouter()
	userRouter.Methods("POST").Path("/email/signin").HandlerFunc(s.user.HandleEmailSignIn())
	userRouter.Methods("POST").Path("/email/signup").HandlerFunc(s.user.HandleEmailSignUp())
	userRouter.Methods("GET").Path("/profile/{id}").HandlerFunc(s.user.HandleProfile())
	userRouter.Methods("GET").Path("/profiles").HandlerFunc(s.user.HandleProfiles())

	// Initialize the routes for the file management operations
	s.router.Methods("POST").Path("/api/v1/{project}/files").HandlerFunc(s.file.HandleCreateFile(s.auth))
	s.router.Methods("GET").PathPrefix("/api/v1/{project}/files").HandlerFunc(s.file.HandleRead(s.auth))
	s.router.Methods("DELETE").PathPrefix("/api/v1/{project}/files").HandlerFunc(s.file.HandleDelete(s.auth))

}
