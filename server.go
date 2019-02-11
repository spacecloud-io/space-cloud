package main

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/auth"
	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/crud"
	"github.com/spaceuptech/space-cloud/userman"
)

type server struct {
	router  *mux.Router
	auth    *auth.Module
	crud    *crud.Module
	user    *userman.Module
	project string
	env     string
}

func initServer(project, env string) *server {
	r := mux.NewRouter()
	c := crud.Init()
	a := auth.Init(c)
	u := userman.Init(c, a)
	return &server{r, a, c, u, project, env}
}

func (s *server) start(port string) error {
	return http.ListenAndServe(":"+port, s.router)
}

func (s *server) loadConfig(config *config.Config) error {
	proj, p := config.Projects[s.project]
	if !p {
		return errors.New("Config doesn't include " + s.project + " project")
	}

	env, p := proj.Env[s.env]
	if !p {
		return errors.New("Config doesn't include " + s.env + " environment")
	}

	// Set the configuration for the auth module
	s.auth.SetConfig(env.Secret, env.Modules.Crud)

	// Set the configuration for the user management module
	s.user.SetConfig(env.Modules.Auth)

	// Set the configuration for the curd module
	return s.crud.SetConfig(env.Modules.Crud)
}
