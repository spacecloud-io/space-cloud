package server

import (
	"net/http/pprof"

	"github.com/gorilla/mux"
	"github.com/spaceuptech/space-cloud/utils/handlers"
)

// InitRoutes initialises the http routes
func (s Server) InitRoutes(profiler bool, staticPath string) {
	s.routes(s.router, profiler, staticPath)
}

// InitSecureRoutes initialises the http routes
func (s Server) InitSecureRoutes(profiler bool, staticPath string) {
	s.routes(s.routerSecure, profiler, staticPath)
}

func (s *Server) routes(router *mux.Router, profiler bool, staticPath string) {
	// Initialize the routes for config management
	router.Methods("GET").Path("/v1/api/config/env").HandlerFunc(handlers.HandleLoadEnv(s.adminMan))
	router.Methods("POST").Path("/v1/api/config/login").HandlerFunc(handlers.HandleAdminLogin(s.adminMan, s.syncMan))
	router.Methods("GET").Path("/v1/api/config/projects").HandlerFunc(handlers.HandleLoadProjects(s.adminMan, s.syncMan, s.configFilePath))
	router.Methods("POST").Path("/v1/api/config/projects").HandlerFunc(handlers.HandleStoreProjectConfig(s.adminMan, s.syncMan, s.configFilePath))
	router.Methods("GET").Path("/v1/api/config/deploy").HandlerFunc(handlers.HandleLoadDeploymentConfig(s.adminMan, s.syncMan, s.configFilePath))
	router.Methods("POST").Path("/v1/api/config/deploy").HandlerFunc(handlers.HandleStoreDeploymentConfig(s.adminMan, s.syncMan, s.configFilePath))
	router.Methods("GET").Path("/v1/api/config/operation").HandlerFunc(handlers.HandleLoadOperationModeConfig(s.adminMan, s.syncMan, s.configFilePath))
	router.Methods("POST").Path("/v1/api/config/operation").HandlerFunc(handlers.HandleStoreOperationModeConfig(s.adminMan, s.syncMan, s.configFilePath))
	router.Methods("GET").Path("/v1/api/config/static").HandlerFunc(handlers.HandleLoadStaticConfig(s.adminMan, s.syncMan))
	router.Methods("POST").Path("/v1/api/config/static").HandlerFunc(handlers.HandleStoreStaticConfig(s.adminMan, s.syncMan))
	router.Methods("DELETE").Path("/v1/api/config/{project}").HandlerFunc(handlers.HandleDeleteProjectConfig(s.adminMan, s.syncMan, s.configFilePath))

	// Initialize the route for websocket
	router.HandleFunc("/v1/api/socket/json", s.handleWebsocket())

	// Initialize the routes for functions service
	router.Methods("POST").Path("/v1/api/{project}/functions/{service}/{func}").HandlerFunc(handlers.HandleFunctionCall(s.functions, s.auth))

	// Initialize the routes for the crud operations
	router.Methods("POST").Path("/v1/api/{project}/crud/{dbType}/batch").HandlerFunc(handlers.HandleCrudBatch(s.auth, s.crud, s.realtime))

	crudRouter := router.Methods("POST").PathPrefix("/v1/api/{project}/crud/{dbType}/{col}").Subrouter()
	crudRouter.HandleFunc("/create", handlers.HandleCrudCreate(s.auth, s.crud, s.realtime))
	crudRouter.HandleFunc("/read", handlers.HandleCrudRead(s.auth, s.crud))
	crudRouter.HandleFunc("/update", handlers.HandleCrudUpdate(s.auth, s.crud, s.realtime))
	crudRouter.HandleFunc("/delete", handlers.HandleCrudDelete(s.auth, s.crud, s.realtime))
	crudRouter.HandleFunc("/aggr", handlers.HandleCrudAggregate(s.auth, s.crud))

	// Initialize the routes for the user management operations
	userRouter := router.PathPrefix("/v1/api/{project}/auth/{dbType}").Subrouter()
	userRouter.Methods("POST").Path("/email/signin").HandlerFunc(handlers.HandleEmailSignIn(s.user))
	userRouter.Methods("POST").Path("/email/signup").HandlerFunc(handlers.HandleEmailSignUp(s.user))
	userRouter.Methods("GET").Path("/profile/{id}").HandlerFunc(handlers.HandleProfile(s.user))
	userRouter.Methods("GET").Path("/profiles").HandlerFunc(handlers.HandleProfiles(s.user))
	userRouter.Methods("GET").Path("/edit_profile/{id}").HandlerFunc(handlers.HandleEmailEditProfile(s.user))

	// Initialize the routes for the file management operations
	router.Methods("POST").Path("/v1/api/{project}/files").HandlerFunc(handlers.HandleCreateFile(s.auth, s.file))
	router.Methods("GET").PathPrefix("/v1/api/{project}/files").HandlerFunc(handlers.HandleRead(s.auth, s.file))
	router.Methods("DELETE").PathPrefix("/v1/api/{project}/files").HandlerFunc(handlers.HandleDelete(s.auth, s.file))

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
	router.PathPrefix("/").HandlerFunc(handlers.HandleStaticRequest(s.static))
}
