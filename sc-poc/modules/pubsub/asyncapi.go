package pubsub

import (
	"net/http"

	"github.com/spacecloud-io/space-cloud/managers/apis"
)

func (a *App) generateASyncAPIDoc() *apis.API {
	return &apis.API{
		Name: "asyncapi",
		Path: "/v1/api/asyncapi.json",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello World"))
		}),
	}
}
