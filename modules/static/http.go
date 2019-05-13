package static

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// HandleRequest creates a static request endpoint
func (m *Module) HandleRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Return 404 if the static module is not enabled
		if !m.isEnabled() {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "This feature isn't enabled"})
			return
		}

		staticRouter := mux.NewRouter().StrictSlash(true)
		staticRouter.PathPrefix(m.URLPrefix).Handler(http.FileServer(&SpaFileSystem{http.Dir(m.getDirPath())}))

		// compress if gzip enabled in config
		if m.Gzip {
			handlers.CompressHandler(staticRouter)
		}
	}
}
