package server

import "net/http"

func (s *Server) routes() {
	s.router.Methods(http.MethodGet).Path("/v1/runner/cluster-type").HandlerFunc(s.handleGetClusterType())
	s.router.Methods(http.MethodGet).Path("/v1/runner/metrics").HandlerFunc(s.handleGetMetrics())

	// project routes
	s.router.Methods(http.MethodPost).Path("/v1/runner/project/{project}").HandlerFunc(s.handleCreateProject())
	s.router.Methods(http.MethodDelete).Path("/v1/runner/{project}").HandlerFunc(s.handleDeleteProject())

	// service routes
	s.router.Methods(http.MethodPost).Path("/v1/runner/{project}/services/{serviceId}/{version}").HandlerFunc(s.handleApplyService())
	s.router.Methods(http.MethodGet).Path("/v1/runner/{project}/services").HandlerFunc(s.HandleGetServices())
	s.router.Methods(http.MethodGet).Path("/v1/runner/{project}/services/status").HandlerFunc(s.HandleGetServicesStatus())

	s.router.Methods(http.MethodDelete).Path("/v1/runner/{project}/services/{serviceId}/{version}").HandlerFunc(s.HandleDeleteService())

	s.router.Methods(http.MethodPost).Path("/v1/runner/{project}/service-routes/{serviceId}").HandlerFunc(s.HandleServiceRoutingRequest())
	s.router.Methods(http.MethodGet).Path("/v1/runner/{project}/service-routes").HandlerFunc(s.HandleGetServiceRoutingRequest())

	s.router.Methods(http.MethodPost).Path("/v1/runner/{project}/service-roles/{serviceId}/{roleId}").HandlerFunc(s.HandleSetServiceRole())
	s.router.Methods(http.MethodGet).Path("/v1/runner/{project}/service-roles").HandlerFunc(s.HandleGetServiceRoleRequest())
	s.router.Methods(http.MethodDelete).Path("/v1/runner/{project}/service-roles/{serviceId}/{roleId}").HandlerFunc(s.HandleDeleteServiceRole())

	s.router.Methods(http.MethodGet).Path("/v1/runner/{project}/services/logs").HandlerFunc(s.handleGetLogs())

	s.router.Methods(http.MethodGet).Path("/v1/runner/cluster-type").HandlerFunc(s.handleGetClusterType())

	// secret routes
	s.router.Methods(http.MethodPost).Path("/v1/runner/{project}/secrets/{id}").HandlerFunc(s.handleApplySecret())
	s.router.Methods(http.MethodGet).Path("/v1/runner/{project}/secrets").HandlerFunc(s.handleListSecrets())
	s.router.Methods(http.MethodPost).Path("/v1/runner/{project}/secrets/{id}/root-path").HandlerFunc(s.handleSetFileSecretRootPath())
	s.router.Methods(http.MethodDelete).Path("/v1/runner/{project}/secrets/{id}").HandlerFunc(s.handleDeleteSecret())
	s.router.Methods(http.MethodPost).Path("/v1/runner/{project}/secrets/{id}/{key}").HandlerFunc(s.handleSetSecretKey())
	s.router.Methods(http.MethodDelete).Path("/v1/runner/{project}/secrets/{id}/{key}").HandlerFunc(s.handleDeleteSecretKey())
}
