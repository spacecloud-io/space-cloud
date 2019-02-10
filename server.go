package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/auth"
	"github.com/spaceuptech/space-cloud/crud"
)

type server struct {
	router *mux.Router
	auth   *auth.Module
	crud   *crud.Module
}

func initServer() *server {
	r := mux.NewRouter()
	c := crud.Init()
	a := auth.Init(c)
	return &server{r, a, c}
}

func (s *server) start(port string) error {
	return http.ListenAndServe(":"+port, s.router)
}
