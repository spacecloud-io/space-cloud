package server

import (
	"net/http"

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
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/config").HandlerFunc(handlers.HandleGlobalConfig(s.adminMan, s.syncMan))
	router.Methods(http.MethodGet).Path("/v1/config/projects").HandlerFunc(handlers.HandleLoadProjects(s.adminMan, s.syncMan))
	router.Methods(http.MethodPut).Path("/v1/config/projects/{project}").HandlerFunc(handlers.HandleStoreProjectConfig(s.adminMan, s.syncMan))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}").HandlerFunc(handlers.HandleDeleteProjectConfig(s.adminMan, s.syncMan))

	// added endpoints for service
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/services/{service}").HandlerFunc(handlers.HandleAddService(s.adminMan, s.syncMan))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/services/{service}").HandlerFunc(handlers.HandleDeleteService(s.adminMan, s.syncMan))
	// Initialize route for graphql schema inspection
	// Initialize route for user management config
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/user-management/{provider}").HandlerFunc(handlers.HandleUserManagement(s.adminMan, s.syncMan))
	// Initialize route for eventing config
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/eventing/triggers/{triggerName}").HandlerFunc(handlers.HandleAddEventingTriggerRule(s.adminMan, s.syncMan))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/eventing/triggers/{triggerName}").HandlerFunc(handlers.HandleDeleteEventingTriggerRule(s.adminMan, s.syncMan))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/eventing/config").HandlerFunc(handlers.HandleSetEventingConfig(s.adminMan, s.syncMan))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/eventing/schema/{type}").HandlerFunc(handlers.HandleSetEventingSchema(s.adminMan, s.syncMan))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/eventing/schema/{type}").HandlerFunc(handlers.HandleDeleteEventingSchema(s.adminMan, s.syncMan))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/eventing/rules/{type}").HandlerFunc(handlers.HandleAddEventingSecurityRule(s.adminMan, s.syncMan))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/eventing/rules/{type}").HandlerFunc(handlers.HandleDeleteEventingSecurityRule(s.adminMan, s.syncMan))
	// Initialize route for file storage config
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/file-storage/config").HandlerFunc(handlers.HandleSetFileStore(s.adminMan, s.syncMan))
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/file-storage/connection-state").HandlerFunc(handlers.HandleGetFileState(s.adminMan, s.syncMan))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/file-storage/rules/{ruleName}").HandlerFunc(handlers.HandleSetFileRule(s.adminMan, s.syncMan))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/file-storage/rules/{ruleName}").HandlerFunc(handlers.HandleDeleteFileRule(s.adminMan, s.syncMan))

	// Initialize route for getting database config
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/database/{dbAlias}/connection-state").HandlerFunc(handlers.HandleGetConnectionState(s.adminMan, s.projects))
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/database/{dbAlias}/list-collections").HandlerFunc(handlers.HandleGetCollections(s.adminMan, s.projects)) // TODO: Check response type
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/database/{dbAlias}/collections/{col}/rules").HandlerFunc(handlers.HandleCollectionRules(s.adminMan, s.syncMan))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/database/{dbAlias}/collections/{col}").HandlerFunc(handlers.HandleDeleteCollection(s.adminMan, s.projects, s.syncMan))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/database/{dbAlias}/config").HandlerFunc(handlers.HandleDatabaseConnection(s.adminMan, s.syncMan))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/database/{dbAlias}").HandlerFunc(handlers.HandleRemoveDatabaseConfig(s.adminMan, s.projects, s.syncMan))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/database/{dbAlias}/modify-schema").HandlerFunc(handlers.HandleModifyAllSchema(s.adminMan, s.projects, s.syncMan))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/database/{dbAlias}/collections/{col}/modify-schema").HandlerFunc(handlers.HandleModifySchema(s.adminMan, s.projects, s.syncMan))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/database/{dbAlias}/reload-schema").HandlerFunc(handlers.HandleReloadSchema(s.adminMan, s.projects, s.syncMan))
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/database/{dbAlias}/collections/{col}/inspect-schema").HandlerFunc(handlers.HandleSchemaInspection(s.adminMan, s.projects, s.syncMan))

	// Initialize route for getting all schemas for all the collections present in config.crud
	router.Methods(http.MethodGet).Path("/v1/config/inspect/{project}/{dbAlias}").HandlerFunc(handlers.HandleGetCollectionSchemas(s.adminMan, s.projects))

	// Initialize routes for the global modules
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/letsencrypt").HandlerFunc(handlers.HandleLetsEncryptWhitelistedDomain(s.adminMan, s.syncMan))

	// Initialize routes for routing module configuration
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/routing").HandlerFunc(handlers.HandleRoutingConfigRequest(s.adminMan, s.syncMan))
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/routing").HandlerFunc(handlers.HandleGetRoutingConfig(s.adminMan, s.syncMan))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/routing/{routeId}").HandlerFunc(handlers.HandleSetProjectRoute(s.adminMan, s.syncMan))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/routing/{routeId}").HandlerFunc(handlers.HandleDeleteProjectRoute(s.adminMan, s.syncMan))

	// Initialize route for graphql
	router.Path("/v1/api/{project}/graphql").HandlerFunc(handlers.HandleGraphQLRequest(s.projects))

	// Initialize the route for websocket
	router.HandleFunc("/v1/api/{project}/socket/json", s.handleWebsocket())

	// Initialize the route for graphql websocket
	router.HandleFunc("/v1/api/{project}/graphql/socket", s.handleGraphqlSocket())

	// Initialize the routes for services module
	router.Methods(http.MethodPost).Path("/v1/api/{project}/services/{service}/{func}").HandlerFunc(handlers.HandleFunctionCall(s.projects))

	// Initialize the routes for realtime service
	router.Methods(http.MethodPost).Path("/v1/api/{project}/realtime/handle").HandlerFunc(handlers.HandleRealtimeEvent(s.projects))
	router.Methods(http.MethodPost).Path("/v1/api/{project}/realtime/process").HandlerFunc(handlers.HandleRealtimeProcessRequest(s.projects))

	// Initialize the routes for eventing service
	router.Methods(http.MethodPost).Path("/v1/api/{project}/eventing/queue").HandlerFunc(handlers.HandleQueueEvent(s.projects))
	router.Methods(http.MethodPost).Path("/v1/api/{project}/eventing/process").HandlerFunc(handlers.HandleProcessEvent(s.adminMan, s.projects))

	// Initialize the routes for the crud operations
	router.Methods(http.MethodPost).Path("/v1/api/{project}/crud/{dbAlias}/batch").HandlerFunc(handlers.HandleCrudBatch(s.projects))

	crudRouter := router.Methods(http.MethodPost).PathPrefix("/v1/api/{project}/crud/{dbAlias}/{col}").Subrouter()
	crudRouter.HandleFunc("/create", handlers.HandleCrudCreate(s.projects))
	crudRouter.HandleFunc("/read", handlers.HandleCrudRead(s.projects))
	crudRouter.HandleFunc("/update", handlers.HandleCrudUpdate(s.projects))
	crudRouter.HandleFunc("/delete", handlers.HandleCrudDelete(s.projects))
	crudRouter.HandleFunc("/aggr", handlers.HandleCrudAggregate(s.projects))

	// Initialize the routes for the user management operations
	userRouter := router.PathPrefix("/v1/api/{project}/auth/{dbAlias}").Subrouter()
	userRouter.Methods(http.MethodPost).Path("/email/signin").HandlerFunc(handlers.HandleEmailSignIn(s.projects))
	userRouter.Methods(http.MethodPost).Path("/email/signup").HandlerFunc(handlers.HandleEmailSignUp(s.projects))
	userRouter.Methods(http.MethodGet).Path("/profile/{id}").HandlerFunc(handlers.HandleProfile(s.projects))
	userRouter.Methods(http.MethodGet).Path("/profiles").HandlerFunc(handlers.HandleProfiles(s.projects))
	userRouter.Methods(http.MethodGet).Path("/edit_profile/{id}").HandlerFunc(handlers.HandleEmailEditProfile(s.projects))

	// Initialize the routes for the file management operations
	router.Methods(http.MethodPost).Path("/v1/api/{project}/files").HandlerFunc(handlers.HandleCreateFile(s.projects))
	router.Methods(http.MethodGet).PathPrefix("/v1/api/{project}/files").HandlerFunc(handlers.HandleRead(s.projects))
	router.Methods(http.MethodDelete).PathPrefix("/v1/api/{project}/files").HandlerFunc(handlers.HandleDelete(s.projects))

	// // Register pprof handlers if profiler set to true
	// if profiler {
	// 	router.HandleFunc("/debug/pprof/", pprof.Index)
	// 	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	// 	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	// 	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	// 	router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	// 	router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	// 	router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	// 	router.Handle("/debug/pprof/block", pprof.Handler("block"))
	// }

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
