package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/auth"
	"github.com/spaceuptech/space-cloud/crud"
	"github.com/spaceuptech/space-cloud/userman"
)

type server struct {
	router *mux.Router
	auth   *auth.Module
	crud   *crud.Module
	user   *userman.Module
}

func initServer() *server {
	r := mux.NewRouter()
	c := crud.Init()
	a := auth.Init(c)
	u := userman.Init(c, a)
	return &server{r, a, c, u}
}

func (s *server) start(port string) error {
	return http.ListenAndServe(":"+port, s.router)
}
