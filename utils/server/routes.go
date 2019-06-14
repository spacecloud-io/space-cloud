package server

import (
	"github.com/spaceuptech/space-cloud/config"
)

func (s *Server) Routes() {
	// Initialize the routes for config management
	s.Router.Methods("POST").Path("/v1/api/config").HandlerFunc(config.HandleConfig(s.IsProd, s.LoadConfig))

	// Initialize the route for websocket
	s.Router.HandleFunc("/v1/api/socket/json", HandleWebsocket(s.Realtime, s.Auth, s.Crud))

	// Initialize the routes for functions engine
	s.Router.Methods("POST").Path("/v1/api/functions/{engine}/{func}").HandlerFunc(s.Functions.HandleRequest(s.Auth))

	// Initialize the routes for the crud operations
	s.Router.Methods("POST").Path("/v1/api/{project}/crud/{dbType}/batch").HandlerFunc(s.handleBatch())

	// Initialize the routes for the CRUD operations
	crudRouter := s.Router.Methods("POST").PathPrefix("/v1/api/{project}/crud/{dbType}/{col}").Subrouter()
	crudRouter.HandleFunc("/create", s.handleCreate())
	crudRouter.HandleFunc("/read", s.handleRead())
	crudRouter.HandleFunc("/update", s.handleUpdate())
	crudRouter.HandleFunc("/delete", s.handleDelete())
	crudRouter.HandleFunc("/aggr", s.handleAggregate())

	// Initialize the routes for the user management operations
	userRouter := s.Router.PathPrefix("/v1/api/{project}/auth/{dbType}").Subrouter()
	userRouter.Methods("POST").Path("/email/signin").HandlerFunc(s.User.HandleEmailSignIn())
	userRouter.Methods("POST").Path("/email/signup").HandlerFunc(s.User.HandleEmailSignUp())
	userRouter.Methods("GET").Path("/profile/{id}").HandlerFunc(s.User.HandleProfile())
	userRouter.Methods("GET").Path("/profiles").HandlerFunc(s.User.HandleProfiles())

	// Initialize the routes for the file management operations
	s.Router.Methods("POST").Path("/v1/api/{project}/files").HandlerFunc(s.File.HandleCreateFile(s.Auth))
	s.Router.Methods("GET").PathPrefix("/v1/api/{project}/files").HandlerFunc(s.File.HandleRead(s.Auth))
	s.Router.Methods("DELETE").PathPrefix("/v1/api/{project}/files").HandlerFunc(s.File.HandleDelete(s.Auth))

	// Initialize the route for handling static files
	s.Static.HandleRequest(s.Router)
}
