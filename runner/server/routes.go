package server

import "net/http"

func (s *Server) routes() {
	s.router.Methods(http.MethodPost).Path("/v1/runner/project").HandlerFunc(s.handleCreateProject())
	s.router.Methods(http.MethodPost).Path("/v1/runner/service").HandlerFunc(s.handleServiceRequest())
	s.router.Methods(http.MethodPost).Path("/v1/runner/events/deploy-service").HandlerFunc(s.handleApplyService())

	s.router.HandleFunc("/v1/runner/socket", s.handleWebsocketRequest())

	//secret routes :P
	s.router.Methods(http.MethodPost).Path("/v1/runner/{project}/secrets").HandlerFunc(s.handleApplySecret())
	s.router.Methods(http.MethodGet).Path("/v1/runner/{project}/secrets").HandlerFunc(s.handleListSecrets())
	s.router.Methods(http.MethodDelete).Path("/v1/runner/{project}/secrets/secrets").HandlerFunc(s.handleDeleteSecret())
	s.router.Methods(http.MethodPost).Path("/v1/runner/{project}/secrets/{name}/{key}").HandlerFunc(s.handleSetSecretKey())
	s.router.Methods(http.MethodDelete).Path("/v1/runner/{project}/secrets/{name}/{key}").HandlerFunc(s.handleDeleteSecretKey())

}
