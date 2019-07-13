package server

import (
	"net/http/pprof"

	"github.com/spaceuptech/space-cloud/utils/handlers"
)

// Routes initialises the http routes
func (s *Server) Routes(profiler bool, staticPath string) {
	// Initialize the routes for config management
	s.router.Methods("POST").Path("/v1/api/config/login").HandlerFunc(handlers.HandleAdminLogin(s.adminMan, s.syncMan))
	s.router.Methods("GET").Path("/v1/api/config/projects").HandlerFunc(handlers.HandleLoadProjects(s.adminMan, s.syncMan))
	s.router.Methods("POST").Path("/v1/api/config/projects").HandlerFunc(handlers.HandleStoreProjectConfig(s.adminMan, s.syncMan))
	s.router.Methods("GET").Path("/v1/api/config/deploy").HandlerFunc(handlers.HandleLoadDeploymentConfig(s.adminMan, s.syncMan))
	s.router.Methods("POST").Path("/v1/api/config/deploy").HandlerFunc(handlers.HandleStoreDeploymentConfig(s.adminMan, s.syncMan))
	s.router.Methods("GET").Path("/v1/api/config/operation").HandlerFunc(handlers.HandleLoadOperationModeConfig(s.adminMan, s.syncMan))
	s.router.Methods("POST").Path("/v1/api/config/operation").HandlerFunc(handlers.HandleStoreOperationModeConfig(s.adminMan, s.syncMan))
	s.router.Methods("DELETE").Path("/v1/api/config/{project}").HandlerFunc(handlers.HandleDeleteProjectConfig(s.adminMan, s.syncMan))

	// Initialize routes for the deployment module
	s.router.Methods("POST").Path("/v1/api/deploy").HandlerFunc(handlers.HandleUploadAndDeploy(s.adminMan, s.deploy, s.projects))

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

	// Register pprof handlers if profiler set to true
	if profiler {
		s.router.HandleFunc("/debug/pprof/", pprof.Index)
		s.router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		s.router.HandleFunc("/debug/pprof/profile", pprof.Profile)
		s.router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		s.router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
		s.router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
		s.router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
		s.router.Handle("/debug/pprof/block", pprof.Handler("block"))
	}

	s.router.PathPrefix("/mission-control").HandlerFunc(handlers.HandleMissionControl(staticPath))

	// Initialize the route for handling static files
	s.router.PathPrefix("/").HandlerFunc(handlers.HandleStaticRequest(s.projects))
}
