package server

func (s *Server) routes() {
	s.router.Methods("POST").Path("/v1/runner/project").HandlerFunc(s.handleCreateProject())
	s.router.Methods("POST").Path("/v1/runner/service").HandlerFunc(s.handleServiceRequest())
	s.router.HandleFunc("/v1/runner/socket", s.handleWebsocketRequest())
}
