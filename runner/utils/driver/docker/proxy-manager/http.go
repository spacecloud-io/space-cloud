package manager

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/runner/utils"
)

func (m *Manager) routes(port int32) http.Handler {
	router := mux.NewRouter()
	router.PathPrefix("/").HandlerFunc(m.handleHTTPRequest(port))
	return router
}
func (m *Manager) handleHTTPRequest(port int32) http.HandlerFunc {

	return func(writer http.ResponseWriter, request *http.Request) {

		// Close the request body
		defer utils.CloseTheCloser(request.Body)

		// Extract serviceID and projectID
		projectID, serviceID := getServiceAndProject(request)

		// Select a proper route
		route, err := m.getRoute(request.Context(), projectID, serviceID, port)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(writer).Encode(map[string]string{"error": err.Error()})
			return
		}

		helpers.Logger.LogDebug(helpers.GetRequestID(request.Context()), fmt.Sprintf("selected route (%v) for projectID (%s), serviceID (%s) and port (%d)", route, projectID, serviceID, port), nil)

		// Proxy the request
		if err := setRequest(request.Context(), request, route, projectID, serviceID); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(writer).Encode(map[string]string{"error": err.Error()})
			_ = helpers.Logger.LogError(helpers.GetRequestID(request.Context()), fmt.Sprintf("Failed set request for route (%v)", route), err, nil)
			return
		}

		response, err := http.DefaultClient.Do(request)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(writer).Encode(map[string]string{"error": err.Error()})
			_ = helpers.Logger.LogError(helpers.GetRequestID(request.Context()), fmt.Sprintf("Failed to make request for route (%v)", route), err, nil)
			return
		}
		defer utils.CloseTheCloser(response.Body)

		// Copy headers and status code
		for k, v := range response.Header {
			writer.Header()[k] = v
		}
		writer.WriteHeader(response.StatusCode)

		// Copy the body
		n, err := io.Copy(writer, response.Body)
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(request.Context()), fmt.Sprintf("Failed to copy upstream (%s) response to downstream", request.URL.String()), err, nil)
		}

		helpers.Logger.LogDebug(helpers.GetRequestID(request.Context()), fmt.Sprintf("Successfully copied %d bytes from upstream server (%s)", n, request.URL.String()), nil)
	}
}
