package server

import "net/http"

func (s *Server) routes() {
	s.router.Methods(http.MethodPost).Path("/v1/runner/project").HandlerFunc(s.handleCreateProject())
	s.router.Methods(http.MethodPost).Path("/v1/runner/service").HandlerFunc(s.handleApplyService())

	s.router.Methods(http.MethodPost).Path("/v1/runner/{projectId}/event-service").HandlerFunc(s.HandleApplyEventingService())
	s.router.Methods(http.MethodGet).Path("/v1/runner/{projectId}/service").HandlerFunc(s.HandleGetServices())
	s.router.Methods(http.MethodDelete).Path("/v1/runner/{projectId}/service/{serviceId}/{version}").HandlerFunc(s.HandleDeleteService())

	s.router.HandleFunc("/v1/runner/socket", s.handleWebsocketRequest())
}
