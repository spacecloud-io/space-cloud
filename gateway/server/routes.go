package server

import (
	"net/http"
	"net/http/pprof"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/gateway/server/handlers"
)

func (s *Server) routes(profiler bool, staticPath string, restrictedHosts []string) http.Handler {
	// Add a '*' to the restricted hosts if length is zero
	if len(restrictedHosts) == 0 {
		restrictedHosts = append(restrictedHosts, "*")
	}

	router := mux.NewRouter()
	router.Methods(http.MethodGet).Path("/v1/config/credentials").HandlerFunc(handlers.HandleGetCredentials(s.managers.Admin()))
	router.Methods(http.MethodGet).Path("/v1/config/permissions").HandlerFunc(handlers.HandleGetPermissions(s.managers.Admin()))

	// Initialize the routes for config management
	router.Methods(http.MethodGet).Path("/v1/config/env").HandlerFunc(handlers.HandleLoadEnv(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodPost).Path("/v1/config/login").HandlerFunc(handlers.HandleAdminLogin(s.managers.Admin()))
	router.Methods(http.MethodGet).Path("/v1/config/refresh-token").HandlerFunc(handlers.HandleRefreshToken(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}").HandlerFunc(handlers.HandleGetProjectConfig(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}").HandlerFunc(handlers.HandleApplyProject(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}").HandlerFunc(handlers.HandleDeleteProjectConfig(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/generate-internal-token").HandlerFunc(handlers.HandleGenerateTokenForMissionControl(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodGet).Path("/v1/config/cluster").HandlerFunc(handlers.HandleGetClusterConfig(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodPost).Path("/v1/config/cluster").HandlerFunc(handlers.HandleSetClusterConfig(s.managers.Admin(), s.managers.Sync()))

	router.Methods(http.MethodGet).Path("/v1/config/integrations").HandlerFunc(handlers.HandleGetIntegrations(s.managers.Admin(), s.managers.Sync()))

	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/remote-service/service").HandlerFunc(handlers.HandleGetService(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/remote-service/service/{id}").HandlerFunc(handlers.HandleAddService(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/remote-service/service/{id}").HandlerFunc(handlers.HandleDeleteService(s.managers.Admin(), s.managers.Sync()))

	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/user-management/provider").HandlerFunc(handlers.HandleGetUserManagement(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/user-management/provider/{id}").HandlerFunc(handlers.HandleSetUserManagement(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/user-management/provider/{id}").HandlerFunc(handlers.HandleDeleteUserManagement(s.managers.Admin(), s.managers.Sync()))

	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/security/function").HandlerFunc(handlers.HandleGetSecurityFunction(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/security/function/{id}").HandlerFunc(handlers.HandleSetSecurityFunction(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/security/function/{id}").HandlerFunc(handlers.HandleDeleteSecurityFunction(s.managers.Admin(), s.managers.Sync()))

	router.Methods(http.MethodGet).Path("/v1/config/caching/config").HandlerFunc(handlers.HandleGetCacheConfig())
	router.Methods(http.MethodPost).Path("/v1/config/caching/config/{id}").HandlerFunc(handlers.HandleSetCacheConfig())
	router.Methods(http.MethodGet).Path("/v1/external/caching/connection-state").HandlerFunc(handlers.HandleGetCacheConnectionState())

	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/eventing/config").HandlerFunc(handlers.HandleGetEventingConfig(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/eventing/config/{id}").HandlerFunc(handlers.HandleSetEventingConfig(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/eventing/triggers").HandlerFunc(handlers.HandleGetEventingTriggers(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/eventing/triggers/{id}").HandlerFunc(handlers.HandleAddEventingTriggerRule(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/eventing/triggers/{id}").HandlerFunc(handlers.HandleDeleteEventingTriggerRule(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/eventing/schema").HandlerFunc(handlers.HandleGetEventingSchema(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/eventing/schema/{id}").HandlerFunc(handlers.HandleSetEventingSchema(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/eventing/schema/{id}").HandlerFunc(handlers.HandleDeleteEventingSchema(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/eventing/rules").HandlerFunc(handlers.HandleGetEventingSecurityRules(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/eventing/rules/{id}").HandlerFunc(handlers.HandleAddEventingSecurityRule(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/eventing/rules/{id}").HandlerFunc(handlers.HandleDeleteEventingSecurityRule(s.managers.Admin(), s.managers.Sync()))

	router.Methods(http.MethodGet).Path("/v1/external/projects/{project}/file-storage/connection-state").HandlerFunc(handlers.HandleGetFileState(s.managers.Admin(), s.modules))
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/file-storage/config").HandlerFunc(handlers.HandleGetFileStore(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/file-storage/config/{id}").HandlerFunc(handlers.HandleSetFileStore(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/file-storage/rules").HandlerFunc(handlers.HandleGetFileRule(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/file-storage/rules/{id}").HandlerFunc(handlers.HandleSetFileRule(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/file-storage/rules/{id}").HandlerFunc(handlers.HandleDeleteFileRule(s.managers.Admin(), s.managers.Sync()))

	router.Methods(http.MethodGet).Path("/v1/external/projects/{project}/database/{dbAlias}/connection-state").HandlerFunc(handlers.HandleGetDatabaseConnectionState(s.managers.Admin(), s.modules))
	router.Methods(http.MethodGet).Path("/v1/external/projects/{project}/database/{dbAlias}/list-collections").HandlerFunc(handlers.HandleGetAllTableNames(s.managers.Admin(), s.modules))
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/database/collections/rules").HandlerFunc(handlers.HandleGetTableRules(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/database/config").HandlerFunc(handlers.HandleGetDatabaseConfig(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/database/collections/schema/mutate").HandlerFunc(handlers.HandleGetSchemas(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/database/{dbAlias}/collections/{col}/rules").HandlerFunc(handlers.HandleSetTableRules(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/database/{dbAlias}/collections/{col}/rules").HandlerFunc(handlers.HandleDeleteTableRules(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/database/{dbAlias}/config/{id}").HandlerFunc(handlers.HandleSetDatabaseConfig(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/database/{dbAlias}/config/{id}").HandlerFunc(handlers.HandleRemoveDatabaseConfig(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/database/prepared-queries").HandlerFunc(handlers.HandleGetPreparedQuery(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/database/{dbAlias}/prepared-queries/{id}").HandlerFunc(handlers.HandleSetPreparedQueries(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/database/{dbAlias}/prepared-queries/{id}").HandlerFunc(handlers.HandleRemovePreparedQueries(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/database/{dbAlias}/collections/{col}").HandlerFunc(handlers.HandleDeleteTable(s.managers.Admin(), s.modules, s.managers.Sync()))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/database/{dbAlias}/schema/mutate").HandlerFunc(handlers.HandleModifyAllSchema(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/database/{dbAlias}/collections/{col}/schema/mutate").HandlerFunc(handlers.HandleModifySchema(s.managers.Admin(), s.modules, s.managers.Sync()))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/database/{dbAlias}/schema/inspect").HandlerFunc(handlers.HandleReloadSchema(s.managers.Admin(), s.modules, s.managers.Sync()))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/database/{dbAlias}/collections/{col}/schema/track").HandlerFunc(handlers.HandleInspectCollectionSchema(s.managers.Admin(), s.modules, s.managers.Sync()))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/database/{dbAlias}/collections/{col}/schema/untrack").HandlerFunc(handlers.HandleUntrackCollectionSchema(s.managers.Admin(), s.modules, s.managers.Sync()))
	router.Methods(http.MethodGet).Path("/v1/external/projects/{project}/database/{dbAlias}/schema/inspect").HandlerFunc(handlers.HandleInspectTrackedCollectionsSchema(s.managers.Admin(), s.modules))

	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/letsencrypt/config").HandlerFunc(handlers.HandleGetEncryptWhitelistedDomain(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/letsencrypt/config/{id}").HandlerFunc(handlers.HandleLetsEncryptWhitelistedDomain(s.managers.Admin(), s.managers.Sync()))

	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/routing/ingress").HandlerFunc(handlers.HandleGetProjectRoute(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/routing/ingress/global").HandlerFunc(handlers.HandleSetGlobalRouteConfig(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/routing/ingress/global").HandlerFunc(handlers.HandleGetGlobalRouteConfig(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/routing/ingress/{id}").HandlerFunc(handlers.HandleSetProjectRoute(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/routing/ingress/{id}").HandlerFunc(handlers.HandleDeleteProjectRoute(s.managers.Admin(), s.managers.Sync()))

	router.Methods(http.MethodPost).Path("/v1/config/batch-apply").HandlerFunc(handlers.HandleBatchApplyConfig(s.managers.Admin()))

	// Health check
	router.Methods(http.MethodGet).Path("/v1/api/health-check").HandlerFunc(handlers.HandleHealthCheck(s.managers.Sync()))

	// Initialize route for graphql
	router.Path("/v1/api/{project}/graphql").HandlerFunc(handlers.HandleGraphQLRequest(s.modules, s.managers.Sync()))

	// Initialize the route for websocket
	router.HandleFunc("/v1/api/{project}/socket/json", handlers.HandleWebsocket(s.modules))

	// Initialize the route for graphql websocket
	router.HandleFunc("/v1/api/{project}/graphql/socket", handlers.HandleGraphqlSocket(s.modules))

	// Initialize the routes for services module
	router.Methods(http.MethodPost).Path("/v1/api/{project}/services/{service}/{func}").HandlerFunc(handlers.HandleFunctionCall(s.modules))

	// Initialize the routes for realtime service
	router.Methods(http.MethodPost).Path("/v1/api/{project}/realtime/handle").HandlerFunc(handlers.HandleRealtimeEvent(s.modules))

	// Initialize the routes for eventing service
	router.Methods(http.MethodPost).Path("/v1/api/{project}/eventing/queue").HandlerFunc(handlers.HandleQueueEvent(s.modules))
	router.Methods(http.MethodPost).Path("/v1/api/{project}/eventing/admin-queue").HandlerFunc(handlers.HandleAdminQueueEvent(s.managers.Admin(), s.modules))

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
	runnerRouter := router.PathPrefix("/v1/runner").HandlerFunc(s.managers.Sync().HandleRunnerRequests(s.managers.Admin())).Subrouter()
	// secret routes
	runnerRouter.Methods(http.MethodPost).Path("/{project}/secrets/{id}").HandlerFunc(s.managers.Sync().HandleRunnerApplySecret(s.managers.Admin()))
	runnerRouter.Methods(http.MethodGet).Path("/{project}/secrets").HandlerFunc(s.managers.Sync().HandleRunnerListSecret(s.managers.Admin()))
	runnerRouter.Methods(http.MethodPost).Path("/{project}/secrets/{id}/root-path").HandlerFunc(s.managers.Sync().HandleRunnerSetFileSecretRootPath(s.managers.Admin()))
	runnerRouter.Methods(http.MethodDelete).Path("/{project}/secrets/{id}").HandlerFunc(s.managers.Sync().HandleRunnerDeleteSecret(s.managers.Admin()))
	runnerRouter.Methods(http.MethodPost).Path("/{project}/secrets/{id}/{key}").HandlerFunc(s.managers.Sync().HandleRunnerSetSecretKey(s.managers.Admin()))
	runnerRouter.Methods(http.MethodDelete).Path("/{project}/secrets/{id}/{key}").HandlerFunc(s.managers.Sync().HandleRunnerDeleteSecretKey(s.managers.Admin()))
	// service
	runnerRouter.Methods(http.MethodPost).Path("/{project}/services/{serviceId}/{version}").HandlerFunc(s.managers.Sync().HandleRunnerApplyService(s.managers.Admin()))
	runnerRouter.Methods(http.MethodGet).Path("/{project}/services").HandlerFunc(s.managers.Sync().HandleRunnerGetServices(s.managers.Admin()))
	runnerRouter.Methods(http.MethodDelete).Path("/{project}/services/{serviceId}/{version}").HandlerFunc(s.managers.Sync().HandleRunnerDeleteService(s.managers.Admin()))
	// service routes
	runnerRouter.Methods(http.MethodPost).Path("/{project}/service-routes/{serviceId}").HandlerFunc(s.managers.Sync().HandleRunnerServiceRoutingRequest(s.managers.Admin()))
	runnerRouter.Methods(http.MethodGet).Path("/{project}/service-routes").HandlerFunc(s.managers.Sync().HandleRunnerGetServiceRoutingRequest(s.managers.Admin()))

	// service role
	runnerRouter.Methods(http.MethodPost).Path("/{project}/service-roles/{serviceId}/{roleId}").HandlerFunc(s.managers.Sync().HandleRunnerSetServiceRole(s.managers.Admin()))
	runnerRouter.Methods(http.MethodGet).Path("/{project}/service-roles").HandlerFunc(s.managers.Sync().HandleRunnerGetServiceRoleRequest(s.managers.Admin()))
	runnerRouter.Methods(http.MethodDelete).Path("/{project}/service-roles/{serviceId}/{roleId}").HandlerFunc(s.managers.Sync().HandleRunnerDeleteServiceRole(s.managers.Admin()))

	runnerRouter.Methods(http.MethodGet).Path("/{project}/services/logs").HandlerFunc(s.managers.Sync().HandleRunnerGetServiceLogs(s.managers.Admin()))
	runnerRouter.Methods(http.MethodGet).Path("/{project}/services/status").HandlerFunc(s.managers.Sync().HandleRunnerGetDeploymentStatus(s.managers.Admin()))

	if staticPath != "" {
		// Add handler for mission control
		router.PathPrefix("/mission-control").HandlerFunc(handlers.HandleMissionControl(staticPath))
	}

	// Add handler for routing module
	router.PathPrefix("/").HandlerFunc(s.modules.Routing().HandleRoutes(s.modules))
	return s.restrictDomainMiddleware(restrictedHosts, router)
}
