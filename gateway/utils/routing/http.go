package routing

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/modules/auth"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

type modulesInterface interface {
	Auth() *auth.Module
}

// HandleRoutes handles incoming http requests and routes them according to the configured rules.
func (r *Routing) HandleRoutes(modules modulesInterface) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// Close the body of the request
		defer utils.CloseTheCloser(request.Body)

		// Extract the host and url to select route
		host, url := getHostAndURL(request)

		// Select a route based on host and url
		route, err := r.selectRoute(host, request.Method, url)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(writer).Encode(map[string]string{"error": err.Error()})
			return
		}

		token, auth, status, err := r.modifyRequest(request.Context(), modules, route, request)
		if err != nil {
			writer.WriteHeader(status)
			_ = json.NewEncoder(writer).Encode(map[string]string{"error": err.Error()})
			return
		}

		logrus.Debugf("selected route (%v) for request (%s)", route, request.URL.String())

		// Apply the rewrite url if provided. It is the users responsibility to make sure both url
		// and rewrite url starts with a '/'
		url = rewriteURL(url, route)

		// Proxy the request

		if err := setRequest(request, route, url); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(writer).Encode(map[string]string{"error": err.Error()})
			logrus.Errorf("Failed set request for route (%v) - %s", route, err.Error())
			return
		}

		// TODO: Use http2 client if that was the incoming request protocol
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(writer).Encode(map[string]string{"error": err.Error()})
			logrus.Errorf("Failed to make request for route (%v) - %s", route, err.Error())
			return
		}
		defer utils.CloseTheCloser(response.Body)

		if err := r.modifyResponse(response, route, token, auth); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(writer).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Copy headers and status code
		for k, v := range response.Header {
			writer.Header().Set(k, v[0])
		}
		writer.WriteHeader(response.StatusCode)

		// Copy the body
		n, err := io.Copy(writer, response.Body)
		if err != nil {
			logrus.Errorf("Failed to copy upstream (%s) response to downstream - %s", request.URL.String(), err.Error())
		}

		logrus.Debugf("Successfully copied %d bytes from upstream server (%s)", n, request.URL.String())
	}
}

func getHostAndURL(request *http.Request) (string, string) {
	return strings.Split(request.Host, ":")[0], request.URL.Path
}

func rewriteURL(url string, route *config.Route) string {
	if route.Source.RewriteURL != "" {
		// First strip away the url provided
		url = strings.TrimPrefix(url, route.Source.URL)

		// Apply the rewrite url at the prefix
		url = route.Source.RewriteURL + url
	}
	return url
}

func setRequest(request *http.Request, route *config.Route, url string) error {
	// http: Request.RequestURI can't be set in client requests.
	// http://golang.org/src/pkg/net/http/client.go
	request.RequestURI = ""

	// Change the request with the destination host, port and url
	target, err := route.SelectTarget(-1) // pass a -ve weight to randomly generate
	if err != nil {
		return err
	}

	request.Host = target.Host
	request.URL.Host = fmt.Sprintf("%s:%d", target.Host, target.Port)
	request.URL.Path = url

	// Set the url scheme to http
	if target.Scheme == "" {
		target.Scheme = "http"
	}
	request.URL.Scheme = target.Scheme
	return nil
}
