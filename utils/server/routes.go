package server

import (
	"net/http/pprof"

	"github.com/gorilla/mux"
	"github.com/spaceuptech/space-cloud/utils/handlers"
)

// InitRoutes initialises the http routes
func (s *Server) InitRoutes(profiler bool, staticPath string) {
	s.routes(s.router, profiler, staticPath)
}

// InitSecureRoutes initialises the http routes
func (s *Server) InitSecureRoutes(profiler bool, staticPath string) {
	s.routes(s.routerSecure, profiler, staticPath)
}

func (s *Server) routes(router *mux.Router, profiler bool, staticPath string) {
	// Initialize the routes for config management
	router.Methods("GET").Path("/v1/config/env").HandlerFunc(handlers.HandleLoadEnv(s.adminMan))
	router.Methods("POST").Path("/v1/config/login").HandlerFunc(handlers.HandleAdminLogin(s.adminMan, s.syncMan))
	router.Methods("POST").Path("/v1/config/projects").HandlerFunc(handlers.HandleCreateProject(s.adminMan, s.syncMan))
	router.Methods("POST").Path("/v1/config/projects/{project}/config").HandlerFunc(handlers.HandleGlobalConfig(s.adminMan, s.syncMan))
	router.Methods("GET").Path("/v1/config/projects").HandlerFunc(handlers.HandleLoadProjects(s.adminMan, s.syncMan))
	router.Methods("PUT").Path("/v1/config/projects/{project}").HandlerFunc(handlers.HandleStoreProjectConfig(s.adminMan, s.syncMan))
	router.Methods("DELETE").Path("/v1/config/projects/{project}").HandlerFunc(handlers.HandleDeleteProjectConfig(s.adminMan, s.syncMan))
	// added endpoints for service
	router.Methods("POST").Path("/v1/config/projects/{project}/services/{service}").HandlerFunc(handlers.HandleAddService(s.adminMan, s.syncMan))
	router.Methods("DELETE").Path("/v1/config/projects/{project}/services/{service}").HandlerFunc(handlers.HandleDeleteService(s.adminMan, s.syncMan))
	//Initialize route for graphql schema inspection
	//Initialize route for user management config
	router.Methods("POST").Path("/v1/config/projects/{project}/user-management/{provider}").HandlerFunc(handlers.HandleUserManagement(s.adminMan, s.syncMan))
	//Initialize route for eventing config
	router.Methods("POST").Path("/v1/config/projects/{project}/event-triggers/rules/{ruleName}").HandlerFunc(handlers.HandleAddEventingRule(s.adminMan, s.syncMan))
	router.Methods("DELETE").Path("/v1/config/projects/{project}/event-triggers/rules/{ruleName}").HandlerFunc(handlers.HandleDeleteEventingRule(s.adminMan, s.syncMan))
	router.Methods("POST").Path("/v1/config/projects/{project}/event-triggers/config").HandlerFunc(handlers.HandleSetEventingConfig(s.adminMan, s.syncMan))
	//Initialize route for file storage config
	router.Methods("POST").Path("/v1/config/projects/{project}/file-storage/config").HandlerFunc(handlers.HandleSetFileStore(s.adminMan, s.syncMan))
	router.Methods("GET").Path("/v1/config/projects/{project}/file-storage/connection-state").HandlerFunc(handlers.HandleGetFileState(s.adminMan, s.syncMan))
	router.Methods("POST").Path("/v1/config/projects/{project}/file-storage/rules/{ruleName}").HandlerFunc(handlers.HandleSetFileRule(s.adminMan, s.syncMan))
	router.Methods("DELETE").Path("/v1/config/projects/{project}/file-storage/rules/{ruleName}").HandlerFunc(handlers.HandleDeleteFileRule(s.adminMan, s.syncMan))

	//Initialize route for getting database config
	router.Methods("GET").Path("/v1/config/projects/{project}/database/{dbType}/connection-state").HandlerFunc(handlers.HandleGetConnectionState(s.adminMan, s.projects))
	router.Methods("GET").Path("/v1/config/projects/{project}/database/{dbType}/list-collections").HandlerFunc(handlers.HandleGetCollections(s.adminMan, s.projects)) // TODO: Check response type
	router.Methods("POST").Path("/v1/config/projects/{project}/database/{dbType}/collections/{col}/rules").HandlerFunc(handlers.HandleCollectionRules(s.adminMan, s.syncMan))
	router.Methods("DELETE").Path("/v1/config/projects/{project}/database/{dbType}/collections/{col}").HandlerFunc(handlers.HandleDeleteCollection(s.adminMan, s.projects, s.syncMan))
	router.Methods("POST").Path("/v1/config/projects/{project}/database/{dbType}/config").HandlerFunc(handlers.HandleDatabaseConnection(s.adminMan, s.syncMan))
	router.Methods("POST").Path("/v1/config/projects/{project}/database/{dbType}/modify-schema").HandlerFunc(handlers.HandleModifyAllSchema(s.adminMan, s.projects, s.syncMan))
	router.Methods("POST").Path("/v1/config/projects/{project}/database/{dbType}/collections/{col}/modify-schema").HandlerFunc(handlers.HandleModifySchema(s.adminMan, s.projects, s.syncMan))
	router.Methods("POST").Path("/v1/config/projects/{project}/database/{dbType}/reload-schema").HandlerFunc(handlers.HandleReloadSchema(s.adminMan, s.projects, s.syncMan))
	router.Methods("GET").Path("/v1/config/projects/{project}/database/{dbType}/collections/{col}/inspect-schema").HandlerFunc(handlers.HandleSchemaInspection(s.adminMan, s.projects, s.syncMan))

	//Initialize route for getting all schemas for all the collections present in config.crud
	router.Methods("GET").Path("/v1/config/inspect/{project}/{dbType}").HandlerFunc(handlers.HandleGetCollectionSchemas(s.adminMan, s.projects))

	// Initialize route for graphql
	router.Path("/v1/api/{project}/graphql").HandlerFunc(handlers.HandleGraphQLRequest(s.projects))

	// Initialize the route for websocket
	router.HandleFunc("/v1/api/{project}/socket/json", s.handleWebsocket())

	// Initialize the route for graphql websocket
	router.HandleFunc("/v1/api/{project}/graphql/socket", s.handleGraphqlSocket())

	// Initialize the routes for functions service
	router.Methods("POST").Path("/v1/api/{project}/functions/{service}/{func}").HandlerFunc(handlers.HandleFunctionCall(s.projects))

	// Initialize the routes for realtime service
	router.Methods("POST").Path("/v1/api/{project}/realtime/handle").HandlerFunc(handlers.HandleRealtimeEvent(s.projects))
	router.Methods("POST").Path("/v1/api/{project}/realtime/process").HandlerFunc(handlers.HandleRealtimeProcessRequest(s.projects))

	// Initialize the routes for eventing service
	router.Methods("POST").Path("/v1/api/{project}/event-triggers/queue").HandlerFunc(handlers.HandleQueueEvent(s.adminMan, s.projects))
	router.Methods("POST").Path("/v1/api/{project}/eventing/process").HandlerFunc(handlers.HandleProcessEvent(s.adminMan, s.projects))

	// Initialize the routes for the crud operations
	router.Methods("POST").Path("/v1/api/{project}/crud/{dbType}/batch").HandlerFunc(handlers.HandleCrudBatch(s.projects))

	crudRouter := router.Methods("POST").PathPrefix("/v1/api/{project}/crud/{dbType}/{col}").Subrouter()
	crudRouter.HandleFunc("/create", handlers.HandleCrudCreate(s.projects))
	crudRouter.HandleFunc("/read", handlers.HandleCrudRead(s.projects))
	crudRouter.HandleFunc("/update", handlers.HandleCrudUpdate(s.projects))
	crudRouter.HandleFunc("/delete", handlers.HandleCrudDelete(s.projects))
	crudRouter.HandleFunc("/aggr", handlers.HandleCrudAggregate(s.projects))

	// Initialize the routes for the user management operations
	userRouter := router.PathPrefix("/v1/api/{project}/auth/{dbType}").Subrouter()
	userRouter.Methods("POST").Path("/email/signin").HandlerFunc(handlers.HandleEmailSignIn(s.projects))
	userRouter.Methods("POST").Path("/email/signup").HandlerFunc(handlers.HandleEmailSignUp(s.projects))
	userRouter.Methods("GET").Path("/profile/{id}").HandlerFunc(handlers.HandleProfile(s.projects))
	userRouter.Methods("GET").Path("/profiles").HandlerFunc(handlers.HandleProfiles(s.projects))
	userRouter.Methods("GET").Path("/edit_profile/{id}").HandlerFunc(handlers.HandleEmailEditProfile(s.projects))

	// Initialize the routes for the file management operations
	router.Methods("POST").Path("/v1/api/{project}/files").HandlerFunc(handlers.HandleCreateFile(s.projects))
	router.Methods("GET").PathPrefix("/v1/api/{project}/files").HandlerFunc(handlers.HandleRead(s.projects))
	router.Methods("DELETE").PathPrefix("/v1/api/{project}/files").HandlerFunc(handlers.HandleDelete(s.projects))

	// Register pprof handlers if profiler set to true
	if profiler {
		router.HandleFunc("/debug/pprof/", pprof.Index)
		router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		router.HandleFunc("/debug/pprof/profile", pprof.Profile)
		router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
		router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
		router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
		router.Handle("/debug/pprof/block", pprof.Handler("block"))
	}

	router.PathPrefix("/mission-control").HandlerFunc(handlers.HandleMissionControl(staticPath))

	// Initialize the route for handling static files
	// router.PathPrefix("/").HandlerFunc(handlers.HandleStaticRequest(s.static))
}
