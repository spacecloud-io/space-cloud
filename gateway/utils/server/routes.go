package server

import (
	"net/http"
	"net/http/pprof"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/gateway/utils/handlers"
)

func (s *Server) routes(profiler bool, staticPath string, restrictedHosts []string) http.Handler {
	// Add a '*' to the restricted hosts if length is zero
	if len(restrictedHosts) == 0 {
		restrictedHosts = append(restrictedHosts, "*")
	}

	router := mux.NewRouter()
	router.Methods(http.MethodGet).Path("/v1/config/credentials").HandlerFunc(handlers.HandleGetCredentials(s.adminMan))

	// Initialize the routes for config management
	router.Methods(http.MethodGet).Path("/v1/config/env").HandlerFunc(handlers.HandleLoadEnv(s.adminMan, s.syncMan))
	router.Methods(http.MethodPost).Path("/v1/config/login").HandlerFunc(handlers.HandleAdminLogin(s.adminMan, s.syncMan))
	router.Methods(http.MethodGet).Path("/v1/config/refresh-token").HandlerFunc(handlers.HandleRefreshToken(s.adminMan, s.syncMan))
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}").HandlerFunc(handlers.HandleGetProjectConfig(s.adminMan, s.syncMan))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}").HandlerFunc(handlers.HandleApplyProject(s.adminMan, s.syncMan))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}").HandlerFunc(handlers.HandleDeleteProjectConfig(s.adminMan, s.syncMan, s.configFilePath))
	router.Methods(http.MethodGet).Path("/v1/config/projects").HandlerFunc(handlers.HandleLoadProjects(s.adminMan, s.syncMan, s.configFilePath))

	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/remote-service/service").HandlerFunc(handlers.HandleGetService(s.adminMan, s.syncMan))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/remote-service/service/{id}").HandlerFunc(handlers.HandleAddService(s.adminMan, s.syncMan))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/remote-service/service/{id}").HandlerFunc(handlers.HandleDeleteService(s.adminMan, s.syncMan))

	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/user-management/provider").HandlerFunc(handlers.HandleGetUserManagement(s.adminMan, s.syncMan))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/user-management/provider/{id}").HandlerFunc(handlers.HandleSetUserManagement(s.adminMan, s.syncMan))

	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/eventing/config").HandlerFunc(handlers.HandleGetEventingConfig(s.adminMan, s.syncMan))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/eventing/config/{id}").HandlerFunc(handlers.HandleSetEventingConfig(s.adminMan, s.syncMan))
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/eventing/triggers").HandlerFunc(handlers.HandleGetEventingTriggers(s.adminMan, s.syncMan))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/eventing/triggers/{id}").HandlerFunc(handlers.HandleAddEventingTriggerRule(s.adminMan, s.syncMan))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/eventing/triggers/{id}").HandlerFunc(handlers.HandleDeleteEventingTriggerRule(s.adminMan, s.syncMan))
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/eventing/schema").HandlerFunc(handlers.HandleGetEventingSchema(s.adminMan, s.syncMan))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/eventing/schema/{id}").HandlerFunc(handlers.HandleSetEventingSchema(s.adminMan, s.syncMan))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/eventing/schema/{id}").HandlerFunc(handlers.HandleDeleteEventingSchema(s.adminMan, s.syncMan))
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/eventing/rules").HandlerFunc(handlers.HandleGetEventingSecurityRules(s.adminMan, s.syncMan))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/eventing/rules/{id}").HandlerFunc(handlers.HandleAddEventingSecurityRule(s.adminMan, s.syncMan))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/eventing/rules/{id}").HandlerFunc(handlers.HandleDeleteEventingSecurityRule(s.adminMan, s.syncMan))

	router.Methods(http.MethodGet).Path("/v1/external/projects/{project}/file-storage/connection-state").HandlerFunc(handlers.HandleGetFileState(s.adminMan, s.modules))
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/file-storage/config").HandlerFunc(handlers.HandleGetFileStore(s.adminMan, s.syncMan))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/file-storage/config/{id}").HandlerFunc(handlers.HandleSetFileStore(s.adminMan, s.syncMan))
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/file-storage/rules").HandlerFunc(handlers.HandleGetFileRule(s.adminMan, s.syncMan))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/file-storage/rules/{id}").HandlerFunc(handlers.HandleSetFileRule(s.adminMan, s.syncMan))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/file-storage/rules/{id}").HandlerFunc(handlers.HandleDeleteFileRule(s.adminMan, s.syncMan))

	router.Methods(http.MethodGet).Path("/v1/external/projects/{project}/database/{dbAlias}/connection-state").HandlerFunc(handlers.HandleGetDatabaseConnectionState(s.adminMan, s.modules))
	router.Methods(http.MethodGet).Path("/v1/external/projects/{project}/database/{dbAlias}/list-collections").HandlerFunc(handlers.HandleGetAllTableNames(s.adminMan, s.modules))
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/database/collections/rules").HandlerFunc(handlers.HandleGetTableRules(s.adminMan, s.syncMan))
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/database/config").HandlerFunc(handlers.HandleGetDatabaseConfig(s.adminMan, s.syncMan))
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/database/collections/schema/mutate").HandlerFunc(handlers.HandleGetSchemas(s.adminMan, s.syncMan))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/database/{dbAlias}/collections/{col}/rules").HandlerFunc(handlers.HandleSetTableRules(s.adminMan, s.syncMan))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/database/{dbAlias}/config/{id}").HandlerFunc(handlers.HandleSetDatabaseConfig(s.adminMan, s.syncMan))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/database/{dbAlias}/config/{id}").HandlerFunc(handlers.HandleRemoveDatabaseConfig(s.adminMan, s.syncMan))
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/database/prepared-queries").HandlerFunc(handlers.HandleGetPreparedQuery(s.adminMan, s.syncMan))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/database/{dbAlias}/prepared-queries/{id}").HandlerFunc(handlers.HandleSetPreparedQueries(s.adminMan, s.syncMan))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/database/{dbAlias}/prepared-queries/{id}").HandlerFunc(handlers.HandleRemovePreparedQueries(s.adminMan, s.syncMan))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/database/{dbAlias}/collections/{col}").HandlerFunc(handlers.HandleDeleteTable(s.adminMan, s.modules, s.syncMan))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/database/{dbAlias}/schema/mutate").HandlerFunc(handlers.HandleModifyAllSchema(s.adminMan, s.syncMan))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/database/{dbAlias}/collections/{col}/schema/mutate").HandlerFunc(handlers.HandleModifySchema(s.adminMan, s.modules, s.syncMan))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/database/{dbAlias}/schema/inspect").HandlerFunc(handlers.HandleReloadSchema(s.adminMan, s.modules, s.syncMan))
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/database/{dbAlias}/collections/{col}/schema/track").HandlerFunc(handlers.HandleInspectCollectionSchema(s.adminMan, s.modules, s.syncMan))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/database/{dbAlias}/collections/{col}/schema/untrack").HandlerFunc(handlers.HandleUntrackCollectionSchema(s.adminMan, s.modules, s.syncMan))
	router.Methods(http.MethodGet).Path("/v1/external/projects/{project}/database/{dbAlias}/schema/inspect").HandlerFunc(handlers.HandleInspectTrackedCollectionsSchema(s.adminMan, s.modules))

	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/letsencrypt/config").HandlerFunc(handlers.HandleGetEncryptWhitelistedDomain(s.adminMan, s.syncMan))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/letsencrypt/config/{id}").HandlerFunc(handlers.HandleLetsEncryptWhitelistedDomain(s.adminMan, s.syncMan))

	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/routing/ingress").HandlerFunc(handlers.HandleGetProjectRoute(s.adminMan, s.syncMan))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/routing/ingress/{id}").HandlerFunc(handlers.HandleSetProjectRoute(s.adminMan, s.syncMan))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/routing/ingress/{id}").HandlerFunc(handlers.HandleDeleteProjectRoute(s.adminMan, s.syncMan))

	// Endpoints for cluster
	router.Methods(http.MethodGet).Path("/v1/config/clusters").HandlerFunc(handlers.HandleCluster())
	router.Methods(http.MethodPost).Path("/v1/config/clusters/{clusterId}").HandlerFunc(handlers.HandleCluster())
	router.Methods(http.MethodGet).Path("/v1/config/clusters/{projectId}/projects").HandlerFunc(handlers.HandleCluster())
	router.Methods(http.MethodDelete).Path("/v1/config/clusters/{projectId}/projects").HandlerFunc(handlers.HandleCluster())

	// Initialize route for graphql
	router.Path("/v1/api/{project}/graphql").HandlerFunc(handlers.HandleGraphQLRequest(s.modules, s.syncMan))

	// Initialize the route for websocket
	router.HandleFunc("/v1/api/{project}/socket/json", handlers.HandleWebsocket(s.modules))

	// Initialize the route for graphql websocket
	router.HandleFunc("/v1/api/{project}/graphql/socket", handlers.HandleGraphqlSocket(s.modules))

	// Initialize the routes for services module
	router.Methods(http.MethodPost).Path("/v1/api/{project}/services/{service}/{func}").HandlerFunc(handlers.HandleFunctionCall(s.modules))

	// Initialize the routes for realtime service
	router.Methods(http.MethodPost).Path("/v1/api/{project}/realtime/handle").HandlerFunc(handlers.HandleRealtimeEvent(s.modules))
	router.Methods(http.MethodPost).Path("/v1/api/{project}/realtime/process").HandlerFunc(handlers.HandleRealtimeProcessRequest(s.modules))

	// Initialize the routes for eventing service
	router.Methods(http.MethodPost).Path("/v1/api/{project}/eventing/queue").HandlerFunc(handlers.HandleQueueEvent(s.modules))
	router.Methods(http.MethodPost).Path("/v1/api/{project}/eventing/process").HandlerFunc(handlers.HandleProcessEvent(s.adminMan, s.modules))
	router.Methods(http.MethodPost).Path("/v1/api/{project}/eventing/process-event-response").HandlerFunc(handlers.HandleEventResponse(s.modules))

	// Initialize the routes for the crud operations
	router.Methods(http.MethodPost).Path("/v1/api/{project}/crud/{dbAlias}/batch").HandlerFunc(handlers.HandleCrudBatch(s.modules))
	router.Methods(http.MethodPost).Path("/v1/api/{project}/crud/{dbAlias}/prepared-queries/{id}").HandlerFunc(handlers.HandleCrudPreparedQuery(s.modules))
	crudRouter := router.Methods(http.MethodPost).PathPrefix("/v1/api/{project}/crud/{dbAlias}/{col}").Subrouter()
	crudRouter.HandleFunc("/create", handlers.HandleCrudCreate(s.modules))
	crudRouter.HandleFunc("/read", handlers.HandleCrudRead(s.modules))
	crudRouter.HandleFunc("/update", handlers.HandleCrudUpdate(s.modules))
	crudRouter.HandleFunc("/delete", handlers.HandleCrudDelete(s.modules))
	crudRouter.HandleFunc("/aggr", handlers.HandleCrudAggregate(s.modules))

	// Initialize the routes for the user management operations
	userRouter := router.PathPrefix("/v1/api/{project}/auth/{dbAlias}").Subrouter()
	userRouter.Methods(http.MethodPost).Path("/email/signin").HandlerFunc(handlers.HandleEmailSignIn(s.modules))
	userRouter.Methods(http.MethodPost).Path("/email/signup").HandlerFunc(handlers.HandleEmailSignUp(s.modules))
	userRouter.Methods(http.MethodGet).Path("/profile/{id}").HandlerFunc(handlers.HandleProfile(s.modules))
	userRouter.Methods(http.MethodGet).Path("/profiles").HandlerFunc(handlers.HandleProfiles(s.modules))
	userRouter.Methods(http.MethodPost).Path("/edit_profile/{id}").HandlerFunc(handlers.HandleEmailEditProfile(s.modules))

	// Initialize the routes for the file management operations
	router.Methods(http.MethodPost).Path("/v1/api/{project}/files").HandlerFunc(handlers.HandleCreateFile(s.modules))
	router.Methods(http.MethodGet).PathPrefix("/v1/api/{project}/files").HandlerFunc(handlers.HandleRead(s.modules))
	router.Methods(http.MethodDelete).PathPrefix("/v1/api/{project}/files").HandlerFunc(handlers.HandleDelete(s.modules))

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

	// forward request for project mutation, websocket, getting cluster type
	runnerRouter := router.PathPrefix("/v1/runner").HandlerFunc(s.syncMan.HandleRunnerRequests(s.adminMan)).Subrouter()
	// secret routes
	runnerRouter.Methods(http.MethodPost).Path("/{project}/secrets/{id}").HandlerFunc(s.syncMan.HandleRunnerApplySecret(s.adminMan))
	runnerRouter.Methods(http.MethodGet).Path("/{project}/secrets").HandlerFunc(s.syncMan.HandleRunnerListSecret(s.adminMan))
	runnerRouter.Methods(http.MethodPost).Path("/{project}/secrets/{id}/root-path").HandlerFunc(s.syncMan.HandleRunnerSetFileSecretRootPath(s.adminMan))
	runnerRouter.Methods(http.MethodDelete).Path("/{project}/secrets/{id}").HandlerFunc(s.syncMan.HandleRunnerDeleteSecret(s.adminMan))
	runnerRouter.Methods(http.MethodPost).Path("/{project}/secrets/{id}/{key}").HandlerFunc(s.syncMan.HandleRunnerSetSecretKey(s.adminMan))
	runnerRouter.Methods(http.MethodDelete).Path("/{project}/secrets/{id}/{key}").HandlerFunc(s.syncMan.HandleRunnerDeleteSecretKey(s.adminMan))
	// service routes
	runnerRouter.Methods(http.MethodPost).Path("/{project}/services/{serviceId}/{version}").HandlerFunc(s.syncMan.HandleRunnerApplyService(s.adminMan))
	runnerRouter.Methods(http.MethodPost).Path("/{project}/event-service").HandlerFunc(s.syncMan.HandleRunnerApplyEventingService(s.adminMan))
	runnerRouter.Methods(http.MethodGet).Path("/{project}/services").HandlerFunc(s.syncMan.HandleRunnerGetServices(s.adminMan))
	runnerRouter.Methods(http.MethodDelete).Path("/{project}/services/{serviceId}/{version}").HandlerFunc(s.syncMan.HandleRunnerDeleteService(s.adminMan))
	runnerRouter.Methods(http.MethodPost).Path("/{project}/service-routes/{serviceId}").HandlerFunc(s.syncMan.HandleRunnerServiceRoutingRequest(s.adminMan))
	runnerRouter.Methods(http.MethodGet).Path("/{project}/service-routes").HandlerFunc(s.syncMan.HandleRunnerGetServiceRoutingRequest(s.adminMan))

	// Add handler for mission control
	router.PathPrefix("/mission-control").HandlerFunc(handlers.HandleMissionControl(staticPath))

	// Add handler for routing module
	router.PathPrefix("/").HandlerFunc(s.routing.HandleRoutes(s.modules))
	return s.restrictDomainMiddleware(restrictedHosts, router)
}
