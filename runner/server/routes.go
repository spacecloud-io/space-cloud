package server

import "net/http"

func (s *Server) routes() {
	s.router.Methods(http.MethodPost).Path("/v1/runner/project").HandlerFunc(s.handleCreateProject())
	s.router.Methods(http.MethodPost).Path("/v1/runner/service").HandlerFunc(s.handleServiceRequest())

	s.router.Methods(http.MethodPost).Path("/v1/runner/{projectId}/service").HandlerFunc(s.HandleApplyService())
	s.router.Methods(http.MethodGet).Path("/v1/runner/{projectId}/service/{serviceId}/{version}").HandlerFunc(s.HandleApplyService())
	s.router.Methods(http.MethodDelete).Path("/v1/runner/{projectId}/service/{serviceId}/{version}").HandlerFunc(s.HandleApplyService())

	s.router.HandleFunc("/v1/runner/socket", s.handleWebsocketRequest())
}
