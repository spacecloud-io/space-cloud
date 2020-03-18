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
	// Initialize the routes for config management
	router.Methods(http.MethodGet).Path("/v1/config/env").HandlerFunc(handlers.HandleLoadEnv(s.adminMan))
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

	router.Methods(http.MethodGet).Path("/v1/external/projects/{project}/file-storage/connection-state").HandlerFunc(handlers.HandleGetFileState(s.adminMan, s.syncMan, s.modules.File))
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/file-storage/config").HandlerFunc(handlers.HandleGetFileStore(s.adminMan, s.syncMan))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/file-storage/config/{id}").HandlerFunc(handlers.HandleSetFileStore(s.adminMan, s.syncMan))
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/file-storage/rules").HandlerFunc(handlers.HandleGetFileRule(s.adminMan, s.syncMan))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/file-storage/rules/{id}").HandlerFunc(handlers.HandleSetFileRule(s.adminMan, s.syncMan))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/file-storage/rules/{id}").HandlerFunc(handlers.HandleDeleteFileRule(s.adminMan, s.syncMan))

	router.Methods(http.MethodGet).Path("/v1/external/projects/{project}/database/{dbAlias}/connection-state").HandlerFunc(handlers.HandleGetConnectionState(s.adminMan, s.modules.Crud))
	router.Methods(http.MethodGet).Path("/v1/external/projects/{project}/database/{dbAlias}/list-collections").HandlerFunc(handlers.HandleGetCollections(s.adminMan, s.modules.Crud, s.syncMan)) // TODO: Check response type
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/database/collections/rules").HandlerFunc(handlers.HandleGetCollectionRules(s.adminMan, s.syncMan))
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/database/config").HandlerFunc(handlers.HandleGetDatabaseConnection(s.adminMan, s.syncMan))
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/database/collections/schema/mutate").HandlerFunc(handlers.HandleGetSchema(s.adminMan, s.syncMan))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/database/{dbAlias}/collections/{col}/rules").HandlerFunc(handlers.HandleCollectionRules(s.adminMan, s.syncMan))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/database/{dbAlias}/config/{id}").HandlerFunc(handlers.HandleDatabaseConnection(s.adminMan, s.modules.Crud, s.syncMan))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/database/{dbAlias}/config/{id}").HandlerFunc(handlers.HandleRemoveDatabaseConfig(s.adminMan, s.modules.Crud, s.syncMan))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/database/{dbAlias}/collections/{col}").HandlerFunc(handlers.HandleDeleteCollection(s.adminMan, s.modules.Crud, s.syncMan))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/database/{dbAlias}/schema/mutate").HandlerFunc(handlers.HandleModifyAllSchema(s.adminMan, s.modules.Schema, s.syncMan))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/database/{dbAlias}/collections/{col}/schema/mutate").HandlerFunc(handlers.HandleModifySchema(s.adminMan, s.modules.Schema, s.syncMan))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/database/{dbAlias}/schema/inspect").HandlerFunc(handlers.HandleReloadSchema(s.adminMan, s.modules.Schema, s.syncMan))
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/database/{dbAlias}/collections/{col}/schema/inspect").HandlerFunc(handlers.HandleInspectCollectionSchema(s.adminMan, s.modules.Schema, s.syncMan))
	router.Methods(http.MethodGet).Path("/v1/external/projects/{project}/database/{dbAlias}/schema/inspect").HandlerFunc(handlers.HandleInspectTrackedCollectionsSchema(s.adminMan, s.modules.Schema))

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
	router.Path("/v1/api/{project}/graphql").HandlerFunc(handlers.HandleGraphQLRequest(s.modules.Graphql, s.syncMan))

	// Initialize the route for websocket
	router.HandleFunc("/v1/api/{project}/socket/json", s.handleWebsocket())

	// Initialize the route for graphql websocket
	router.HandleFunc("/v1/api/{project}/graphql/socket", s.handleGraphqlSocket(s.adminMan))

	// Initialize the routes for services module
	router.Methods(http.MethodPost).Path("/v1/api/{project}/services/{service}/{func}").HandlerFunc(handlers.HandleFunctionCall(s.modules.Functions, s.modules.Auth))

	// Initialize the routes for realtime service
	router.Methods(http.MethodPost).Path("/v1/api/{project}/realtime/handle").HandlerFunc(handlers.HandleRealtimeEvent(s.modules.Auth, s.modules.Realtime))
	router.Methods(http.MethodPost).Path("/v1/api/{project}/realtime/process").HandlerFunc(handlers.HandleRealtimeProcessRequest(s.modules.Auth, s.modules.Realtime))

	// Initialize the routes for eventing service
	router.Methods(http.MethodPost).Path("/v1/api/{project}/eventing/queue").HandlerFunc(handlers.HandleQueueEvent(s.modules.Eventing))
	router.Methods(http.MethodPost).Path("/v1/api/{project}/eventing/process").HandlerFunc(handlers.HandleProcessEvent(s.adminMan, s.modules.Eventing))
	router.Methods(http.MethodPost).Path("/v1/api/{project}/eventing/process-event-response").HandlerFunc(handlers.HandleEventResponse(s.modules.Auth, s.modules.Eventing))

	// Initialize the routes for the crud operations
	router.Methods(http.MethodPost).Path("/v1/api/{project}/crud/{dbAlias}/batch").HandlerFunc(handlers.HandleCrudBatch(s.modules.Auth, s.modules.Crud, s.modules.Realtime))

	crudRouter := router.Methods(http.MethodPost).PathPrefix("/v1/api/{project}/crud/{dbAlias}/{col}").Subrouter()
	crudRouter.HandleFunc("/create", handlers.HandleCrudCreate(s.modules.Auth, s.modules.Crud, s.modules.Realtime))
	crudRouter.HandleFunc("/read", handlers.HandleCrudRead(s.modules.Auth, s.modules.Crud))
	crudRouter.HandleFunc("/update", handlers.HandleCrudUpdate(s.modules.Auth, s.modules.Crud, s.modules.Realtime))
	crudRouter.HandleFunc("/delete", handlers.HandleCrudDelete(s.modules.Auth, s.modules.Crud, s.modules.Realtime))
	crudRouter.HandleFunc("/aggr", handlers.HandleCrudAggregate(s.modules.Auth, s.modules.Crud))

	// Initialize the routes for the user management operations
	userRouter := router.PathPrefix("/v1/api/{project}/auth/{dbAlias}").Subrouter()
	userRouter.Methods(http.MethodPost).Path("/email/signin").HandlerFunc(handlers.HandleEmailSignIn(s.modules.User))
	userRouter.Methods(http.MethodPost).Path("/email/signup").HandlerFunc(handlers.HandleEmailSignUp(s.modules.User))
	userRouter.Methods(http.MethodGet).Path("/profile/{id}").HandlerFunc(handlers.HandleProfile(s.modules.User))
	userRouter.Methods(http.MethodGet).Path("/profiles").HandlerFunc(handlers.HandleProfiles(s.modules.User))
	userRouter.Methods(http.MethodGet).Path("/edit_profile/{id}").HandlerFunc(handlers.HandleEmailEditProfile(s.modules.User))

	// Initialize the routes for the file management operations
	router.Methods(http.MethodPost).Path("/v1/api/{project}/files").HandlerFunc(handlers.HandleCreateFile(s.modules.File))
	router.Methods(http.MethodGet).PathPrefix("/v1/api/{project}/files").HandlerFunc(handlers.HandleRead(s.modules.File))
	router.Methods(http.MethodDelete).PathPrefix("/v1/api/{project}/files").HandlerFunc(handlers.HandleDelete(s.modules.File))

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
	return s.restrictDomainMiddleware(restrictedHosts, router)
}
