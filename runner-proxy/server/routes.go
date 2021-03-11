package server

func (s *Server) routes() {
	s.router.PathPrefix("/").HandlerFunc(s.handleProxy())
}
