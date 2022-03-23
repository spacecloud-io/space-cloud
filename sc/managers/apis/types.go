package apis

import (
	"net/http"
	"sort"

	"github.com/getkin/kin-openapi/openapi3"
)

type (
	// App returns the paths it intends to expose
	App interface {
		GetRoutes() []*API
		GetHandler(op string) (http.HandlerFunc, error)
	}

	// API describes how to handle a particular path
	API struct {
		Path    string             `json:"path"`
		PathDef *openapi3.PathItem `json:"pathDef"`
		Op      string             `json:"op"`

		app string
	}
)

type (
	apps []app
	app  struct {
		name     string
		priority int
	}
)

func (a apps) sort() {
	sort.SliceStable(a, func(i, j int) bool {
		return a[i].priority > a[j].priority
	})
}
