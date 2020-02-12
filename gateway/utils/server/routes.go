package server

import (
	"net/http"
	"net/http/pprof"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/gateway/utils/handlers"
)

func (s *Server) routes(profiler bool, staticPath, configDomain string) *mux.Router {
	router := mux.NewRouter()
	// TODO: Only limit the host of config and runner apis
	if configDomain != "" {
		router.Host(configDomain)
	}
	// Initialize the routes for config management
	router.Methods(http.MethodGet).Path("/v1/config/env").HandlerFunc(handlers.HandleLoadEnv(s.adminMan))
	router.Methods(http.MethodPost).Path("/v1/config/login").HandlerFunc(handlers.HandleAdminLogin(s.adminMan, s.syncMan))
	router.Methods(http.MethodGet).Path("/v1/config/refresh-token").HandlerFunc(handlers.HandleRefreshToken(s.adminMan, s.syncMan))
	router.Methods(http.MethodPost).Path("/v1/config/projects").HandlerFunc(handlers.HandleCreateProject(s.adminMan, s.syncMan))
	router.Methods(http.MethodPut).Path("/v1/config/projects/{project}/config").HandlerFunc(handlers.HandleGlobalConfig(s.adminMan, s.syncMan))
	router.Methods(http.MethodGet).Path("/v1/config/projects").HandlerFunc(handlers.HandleLoadProjects(s.adminMan, s.syncMan, s.configFilePath))
	router.Methods(http.MethodPut).Path("/v1/config/projects/{project}").HandlerFunc(handlers.HandleStoreProjectConfig(s.adminMan, s.syncMan, s.configFilePath))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}").HandlerFunc(handlers.HandleDeleteProjectConfig(s.adminMan, s.syncMan, s.configFilePath))

	// added endpoints for service
	router.Methods(http.MethodPut).Path("/v1/config/projects/{project}/services/{service}").HandlerFunc(handlers.HandleAddService(s.adminMan, s.syncMan))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/services/{service}").HandlerFunc(handlers.HandleDeleteService(s.adminMan, s.syncMan))
	// Initialize route for graphql schema inspection
	// Initialize route for user management config
	router.Methods(http.MethodPut).Path("/v1/config/projects/{project}/user-management/{provider}").HandlerFunc(handlers.HandleUserManagement(s.adminMan, s.syncMan))
	// Initialize route for eventing config
	router.Methods(http.MethodPut).Path("/v1/config/projects/{project}/event-triggers/rules/{ruleName}").HandlerFunc(handlers.HandleAddEventingRule(s.adminMan, s.syncMan))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/event-triggers/rules/{ruleName}").HandlerFunc(handlers.HandleDeleteEventingRule(s.adminMan, s.syncMan))
	router.Methods(http.MethodPut).Path("/v1/config/projects/{project}/event-triggers/config").HandlerFunc(handlers.HandleSetEventingConfig(s.adminMan, s.syncMan))
	router.Methods(http.MethodPut).Path("/v1/config/projects/{project}/eventing/schema/{type}").HandlerFunc(handlers.HandleSetEventingSchema(s.adminMan, s.syncMan))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/eventing/schema/{type}").HandlerFunc(handlers.HandleDeleteEventingSchema(s.adminMan, s.syncMan))
	router.Methods(http.MethodPut).Path("/v1/config/projects/{project}/eventing/rules/{type}").HandlerFunc(handlers.HandleAddEventingSecurityRules(s.adminMan, s.syncMan))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/eventing/rules/{type}").HandlerFunc(handlers.HandleDeleteEventingSecurityRules(s.adminMan, s.syncMan))
	// Initialize route for file storage config
	router.Methods(http.MethodPut).Path("/v1/config/projects/{project}/file-storage/config").HandlerFunc(handlers.HandleSetFileStore(s.adminMan, s.syncMan))
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/file-storage/connection-state").HandlerFunc(handlers.HandleGetFileState(s.adminMan, s.syncMan))
	router.Methods(http.MethodPut).Path("/v1/config/projects/{project}/file-storage/rules/{ruleName}").HandlerFunc(handlers.HandleSetFileRule(s.adminMan, s.syncMan))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/file-storage/rules/{ruleName}").HandlerFunc(handlers.HandleDeleteFileRule(s.adminMan, s.syncMan))

	// Initialize route for getting database config
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/database/{dbType}/connection-state").HandlerFunc(handlers.HandleGetConnectionState(s.adminMan, s.crud))
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/database/{dbType}/list-collections").HandlerFunc(handlers.HandleGetCollections(s.adminMan, s.crud, s.syncMan)) // TODO: Check response type
	router.Methods(http.MethodPut).Path("/v1/config/projects/{project}/database/{dbType}/collections/{col}/rules").HandlerFunc(handlers.HandleCollectionRules(s.adminMan, s.syncMan))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/database/{dbType}/collections/{col}").HandlerFunc(handlers.HandleDeleteCollection(s.adminMan, s.crud, s.syncMan))
	router.Methods(http.MethodPut).Path("/v1/config/projects/{project}/database/{dbType}/config").HandlerFunc(handlers.HandleDatabaseConnection(s.adminMan, s.crud, s.syncMan))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/database/{dbType}").HandlerFunc(handlers.HandleRemoveDatabaseConfig(s.adminMan, s.crud, s.syncMan))
	router.Methods(http.MethodPut).Path("/v1/config/projects/{project}/database/{dbType}/modify-schema").HandlerFunc(handlers.HandleModifyAllSchema(s.adminMan, s.schema, s.syncMan))
	router.Methods(http.MethodPut).Path("/v1/config/projects/{project}/database/{dbType}/collections/{col}/modify-schema").HandlerFunc(handlers.HandleModifySchema(s.adminMan, s.schema, s.syncMan))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/database/{dbType}/reload-schema").HandlerFunc(handlers.HandleReloadSchema(s.adminMan, s.schema, s.syncMan))
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/database/{dbType}/collections/{col}/inspect-schema").HandlerFunc(handlers.HandleSchemaInspection(s.adminMan, s.schema, s.syncMan))

	// Initialize route for getting all schemas for all the collections present in config.crud
	router.Methods(http.MethodGet).Path("/v1/config/inspect/{project}/{dbType}").HandlerFunc(handlers.HandleGetCollectionSchemas(s.adminMan, s.schema))

	// Initialize routes for the global modules
	router.Methods(http.MethodPut).Path("/v1/config/projects/{project}/letsencrypt").HandlerFunc(handlers.HandleLetsEncryptWhitelistedDomain(s.adminMan, s.syncMan))

	// Initialize routes for routing module configuration
	router.Methods(http.MethodPut).Path("/v1/config/projects/{project}/routing").HandlerFunc(handlers.HandleRoutingConfigRequest(s.adminMan, s.syncMan))
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/routing").HandlerFunc(handlers.HandleGetRoutingConfig(s.adminMan, s.syncMan))
	router.Methods(http.MethodPut).Path("/v1/config/projects/{project}/routing/{routeId}").HandlerFunc(handlers.HandleSetProjectRoute(s.adminMan, s.syncMan))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/routing/{routeId}").HandlerFunc(handlers.HandleDeleteProjectRoute(s.adminMan, s.syncMan))

	// Initialize route for graphql
	router.Path("/v1/api/{project}/graphql").HandlerFunc(handlers.HandleGraphQLRequest(s.graphql))

	// Initialize the route for websocket
	router.HandleFunc("/v1/api/{project}/socket/json", s.handleWebsocket())

	// Initialize the route for graphql websocket
	router.HandleFunc("/v1/api/{project}/graphql/socket", s.handleGraphqlSocket(s.adminMan))

	// Initialize the routes for services module
	router.Methods(http.MethodPost).Path("/v1/api/{project}/services/{service}/{func}").HandlerFunc(handlers.HandleFunctionCall(s.functions, s.auth))

	// Initialize the routes for realtime service
	router.Methods(http.MethodPost).Path("/v1/api/{project}/realtime/handle").HandlerFunc(handlers.HandleRealtimeEvent(s.auth, s.realtime))
	router.Methods(http.MethodPost).Path("/v1/api/{project}/realtime/process").HandlerFunc(handlers.HandleRealtimeProcessRequest(s.auth, s.realtime))

	// Initialize the routes for eventing service
	router.Methods(http.MethodPost).Path("/v1/api/{project}/event-triggers/queue").HandlerFunc(handlers.HandleQueueEvent(s.eventing))
	router.Methods(http.MethodPost).Path("/v1/api/{project}/eventing/process").HandlerFunc(handlers.HandleProcessEvent(s.adminMan, s.eventing))

	// Initialize the routes for the crud operations
	router.Methods(http.MethodPost).Path("/v1/api/{project}/crud/{dbType}/batch").HandlerFunc(handlers.HandleCrudBatch(s.auth, s.crud, s.realtime))

	crudRouter := router.Methods(http.MethodPost).PathPrefix("/v1/api/{project}/crud/{dbType}/{col}").Subrouter()
	crudRouter.HandleFunc("/create", handlers.HandleCrudCreate(s.auth, s.crud, s.realtime))
	crudRouter.HandleFunc("/read", handlers.HandleCrudRead(s.auth, s.crud))
	crudRouter.HandleFunc("/update", handlers.HandleCrudUpdate(s.auth, s.crud, s.realtime))
	crudRouter.HandleFunc("/delete", handlers.HandleCrudDelete(s.auth, s.crud, s.realtime))
	crudRouter.HandleFunc("/aggr", handlers.HandleCrudAggregate(s.auth, s.crud))

	// Initialize the routes for the user management operations
	userRouter := router.PathPrefix("/v1/api/{project}/auth/{dbType}").Subrouter()
	userRouter.Methods(http.MethodPost).Path("/email/signin").HandlerFunc(handlers.HandleEmailSignIn(s.user))
	userRouter.Methods(http.MethodPost).Path("/email/signup").HandlerFunc(handlers.HandleEmailSignUp(s.user))
	userRouter.Methods(http.MethodGet).Path("/profile/{id}").HandlerFunc(handlers.HandleProfile(s.user))
	userRouter.Methods(http.MethodGet).Path("/profiles").HandlerFunc(handlers.HandleProfiles(s.user))
	userRouter.Methods(http.MethodGet).Path("/edit_profile/{id}").HandlerFunc(handlers.HandleEmailEditProfile(s.user))

	// Initialize the routes for the file management operations
	router.Methods(http.MethodPost).Path("/v1/api/{project}/files").HandlerFunc(handlers.HandleCreateFile(s.file))
	router.Methods(http.MethodGet).PathPrefix("/v1/api/{project}/files").HandlerFunc(handlers.HandleRead(s.file))
	router.Methods(http.MethodDelete).PathPrefix("/v1/api/{project}/files").HandlerFunc(handlers.HandleDelete(s.file))

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

	// Add addresses for runner
	router.PathPrefix("/v1/runner").HandlerFunc(s.syncMan.HandleRunnerRequests(s.adminMan))

	// forward request to artifact store
	router.PathPrefix("/v1/artifact").HandlerFunc(s.syncMan.HandleArtifactRequests(s.adminMan))

	// Add handler for mission control
	router.PathPrefix("/mission-control").HandlerFunc(handlers.HandleMissionControl(staticPath))

	// Add handler for routing module
	router.PathPrefix("/").HandlerFunc(s.routing.HandleRoutes())
	return router
}
