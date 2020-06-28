package server

import (
	"net/http"
	"strings"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (s *Server) restrictDomainMiddleware(restrictedHosts []string, h http.Handler) http.Handler {
	routingHandler := s.routing.HandleRoutes(s.modules)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get the url and host parameters
		url := r.URL.Path
		host := strings.Split(r.Host, ":")[0]

		// Check if host belongs does not belong to restricted host
		if !utils.StringExists(restrictedHosts, "*") && !utils.StringExists(restrictedHosts, host) {
			// We are here means we just got an excluded host. The config, runner and mission-control routes need to be hidden in this case.
			// So we'll forward the request straight to the routing handler.
			if strings.HasPrefix(url, "/v1/config") || strings.HasPrefix(url, "/v1/runner") || strings.HasPrefix(url, "/mission-control") {
				// Forward the request to the routing handler and don't forget to return
				routingHandler(w, r)
				return
			}
			// We are here means that the request is not a config, runner or mission-control path. Hence we can safely use the main handler to serve the request.
			// Since the main handler includes the routing handler, we need not worry about routing requests
		}
		h.ServeHTTP(w, r)
	})
}
