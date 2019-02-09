package main

func (s *server) routes() {

	// Initialize the routes for the crud operations
	crudRouter := s.router.Methods("POST").PathPrefix("/v1/api/{project}/crud/{dbType}/{col}").Subrouter()
	crudRouter.HandleFunc("/create", s.handleCreate())
	crudRouter.HandleFunc("/read", s.handleRead())
	crudRouter.HandleFunc("/update", s.handleUpdate())
	crudRouter.HandleFunc("/delete", s.handleDelete())
}
