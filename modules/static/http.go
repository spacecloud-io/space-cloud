package static

import (
	"net/http"

	"github.com/gorilla/mux"
)

// HandleRequest creates a static request endpoint
func (m *Module) HandleRequest(r *mux.Router) {
	if m.Enabled {
		r.PathPrefix(m.URLPrefix).Handler(http.FileServer(http.Dir(m.getDirPath())))
	}
}
