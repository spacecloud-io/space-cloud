package handlers

import (
	"bufio"
	"net/http"
	"strings"

	"github.com/spaceuptech/space-cloud/utils/projects"
)

// HandleStaticRequest creates a static request endpoint
func HandleStaticRequest(p *projects.Projects) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Path
		host := strings.Split(r.Host, ":")[0]

		completed := p.Iter(func(project string, state *projects.ProjectState) bool {
			static := state.Static
			route, ok := static.SelectRoute(host, url)
			if !ok {
				return true
			}

			path := strings.TrimPrefix(url, route.URLPrefix)
			if !strings.HasPrefix(path, "/") {
				path = "/" + path
			}
			path = route.Path + path

			// Its a proxy request
			if route.Proxy != "" {
				addr := route.Proxy + path
				req, err := http.NewRequest(r.Method, addr, r.Body)
				if err != nil {
					http.Error(w, err.Error(), http.StatusNotFound)
					return false
				}

				// Set the http headers
				req.Header = make(http.Header)
				if contentType, p := r.Header["Content-Type"]; p {
					req.Header["Content-Type"] = contentType
				}

				// Make the http client request
				res, err := http.DefaultClient.Do(req)
				if err != nil {
					http.Error(w, err.Error(), http.StatusNotFound)
					return false
				}
				defer res.Body.Close()

				reader := bufio.NewReader(res.Body)

				w.Header().Set("Content-Type", res.Header.Get("Content-Type"))
				w.WriteHeader(res.StatusCode)
				reader.WriteTo(w)
				return false
			}

			http.ServeFile(w, r, path)
			return false
		})

		if !completed {
			http.Error(w, "Path not found", http.StatusNotFound)
		}
	}
}
