package server

func (s *Server) routes() {
	s.router.Methods("POST").Path("/v1/runner/project").HandlerFunc(s.handleCreateProject())
	s.router.Methods("POST").Path("/v1/runner/service").HandlerFunc(s.handleServiceRequest())
	s.router.Methods("POST").Path("/v1/runner/events/deploy-service").HandlerFunc(s.HandleApplyService())

	s.router.HandleFunc("/v1/runner/socket", s.handleWebsocketRequest())

	//secret routes :P
	s.router.Methods("POST").Path("/v1/runner/{projectID}/secrets").HandlerFunc(s.handleUpsertSecret())
	s.router.Methods("GET").Path("/v1/runner/{projectID}/secrets").HandlerFunc(s.handleListSecrets())
	s.router.Methods("DEL").Path("/v1/runner/{projectID}/secrets").HandlerFunc(s.handleDeleteSecrets())
	s.router.Methods("POST").Path("/v1/runner/{projectID}/{secretName}/{secretKey}").HandlerFunc(s.handleSetSecretKey())
	s.router.Methods("DEL").Path("/v1/runner/{projectID}/{secretName}/{secretKey}").HandlerFunc(s.handleDeleteSecretKey())

}
