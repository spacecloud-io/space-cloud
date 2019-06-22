package server

import (
	"github.com/spaceuptech/space-cloud/utils/handlers"
)

// Routes initialises the http endpoints
func (s *Server) Routes() {
	// Initialize the routes for config management
	//s.router.Methods("POST").Path("/v1/api/config").HandlerFunc(config.HandleConfig(s.isProd, s.loadConfig))

	// Initialize the route for websocket
	s.router.HandleFunc("/v1/api/socket/json", s.handleWebsocket())

	// Initialize the routes for functions service
	s.router.Methods("POST").Path("/v1/api/{project}/functions/{service}/{func}").HandlerFunc(handlers.HandleFunctionCall(s.projects))

	// Initialize the routes for the crud operations
	s.router.Methods("POST").Path("/v1/api/{project}/crud/{dbType}/batch").HandlerFunc(handlers.HandleCrudBatch(s.isProd, s.projects))

	crudRouter := s.router.Methods("POST").PathPrefix("/v1/api/{project}/crud/{dbType}/{col}").Subrouter()
	crudRouter.HandleFunc("/create", handlers.HandleCrudCreate(s.isProd, s.projects))
	crudRouter.HandleFunc("/read", handlers.HandleCrudRead(s.projects))
	crudRouter.HandleFunc("/update", handlers.HandleCrudUpdate(s.isProd, s.projects))
	crudRouter.HandleFunc("/delete", handlers.HandleCrudDelete(s.isProd, s.projects))
	crudRouter.HandleFunc("/aggr", handlers.HandleCrudAggregate(s.projects))

	// Initialize the routes for the user management operations
	userRouter := s.router.PathPrefix("/v1/api/{project}/auth/{dbType}").Subrouter()
	userRouter.Methods("POST").Path("/email/signin").HandlerFunc(handlers.HandleEmailSignIn(s.projects))
	userRouter.Methods("POST").Path("/email/signup").HandlerFunc(handlers.HandleEmailSignUp(s.projects))
	userRouter.Methods("GET").Path("/profile/{id}").HandlerFunc(handlers.HandleProfile(s.projects))
	userRouter.Methods("GET").Path("/profiles").HandlerFunc(handlers.HandleProfiles(s.projects))
	userRouter.Methods("GET").Path("/edit_profile/{id}").HandlerFunc(handlers.HandleEmailEditProfile(s.projects))

	// Initialize the routes for the file management operations
	s.router.Methods("POST").Path("/v1/api/{project}/files").HandlerFunc(handlers.HandleCreateFile(s.projects))
	s.router.Methods("GET").PathPrefix("/v1/api/{project}/files").HandlerFunc(handlers.HandleRead(s.projects))
	s.router.Methods("DELETE").PathPrefix("/v1/api/{project}/files").HandlerFunc(handlers.HandleDelete(s.projects))

	// Initialize the route for handling static files
	s.router.PathPrefix("/").HandlerFunc(handlers.HandleStaticRequest(s.projects))
}
