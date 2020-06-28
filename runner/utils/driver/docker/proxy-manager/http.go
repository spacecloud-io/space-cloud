package manager

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

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
		route, err := m.getRoute(projectID, serviceID, port)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(writer).Encode(map[string]string{"error": err.Error()})
			return
		}

		logrus.Debugf("selected route (%v) for projectID (%s), serviceID (%s) and port (%d)", route, projectID, serviceID, port)

		// Proxy the request
		if err := setRequest(request, route, projectID, serviceID); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(writer).Encode(map[string]string{"error": err.Error()})
			logrus.Errorf("Failed set request for route (%v) - %s", route, err.Error())
			return
		}

		response, err := http.DefaultClient.Do(request)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(writer).Encode(map[string]string{"error": err.Error()})
			logrus.Errorf("Failed to make request for route (%v) - %s", route, err.Error())
			return
		}
		defer utils.CloseTheCloser(response.Body)

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
