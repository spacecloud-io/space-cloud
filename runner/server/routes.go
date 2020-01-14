package server

func (s *Server) routes() {
	s.router.Methods("POST").Path("/v1/galaxy/project").HandlerFunc(s.handleCreateProject())
	s.router.Methods("POST").Path("/v1/galaxy/service").HandlerFunc(s.handleServiceRequest())
	s.router.HandleFunc("/v1/galaxy/socket", s.handleWebsocketRequest())
}
