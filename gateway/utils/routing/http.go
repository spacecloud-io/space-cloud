package routing

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// HandleRoutes handles incoming http requests and routes them according to the configured rules.
func (r *Routing) HandleRoutes() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// Close the body of the request
		defer utils.CloseTheCloser(request.Body)

		// Extract the host and url to select route
		host, url := getHostAndURL(request)

		// Select a route based on host and url
		route, err := r.selectRoute(host, url)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(writer).Encode(map[string]string{"error": err.Error()})
			return
		}

		logrus.Debugf("selected route (%v) for request (%s)", route, request.URL.String())

		// Apply the rewrite url if provided. It is the users responsibility to make sure both url
		// and rewrite url starts with a '/'
		url = rewriteURL(url, route)

		// Proxy the request

		// http: Request.RequestURI can't be set in client requests.
		// http://golang.org/src/pkg/net/http/client.go
		setRequest(request, route, url)

		// TODO: Use http2 client if that was the incoming request protocol
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(writer).Encode(map[string]string{"error": err.Error()})
			return
		}
		defer utils.CloseTheCloser(response.Body)

		// Copy headers and status code
		writer.WriteHeader(response.StatusCode)
		for k, v := range response.Header {
			writer.Header().Set(k, v[0])
		}

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

func setRequest(request *http.Request, route *config.Route, url string) {
	request.RequestURI = ""

	// Change the request with the destination host, port and url
	request.Host = route.Destination.Host
	request.URL.Host = fmt.Sprintf("%s:%s", route.Destination.Host, route.Destination.Port)
	request.URL.Path = url

	// Set the url scheme to http
	request.URL.Scheme = "http"
}
