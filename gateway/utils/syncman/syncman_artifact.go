package syncman

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/modules/auth"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (s *Manager) HandleArtifactRequests(auth *auth.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// http: Request.RequestURI can't be set in client requests.
		// http://golang.org/src/pkg/net/http/client.go
		r.Host = strings.Split(s.artifactAddr, ":")[0]
		r.URL.Host = s.artifactAddr

		r.RequestURI = ""
		r.URL.Scheme = "http"

		vars := mux.Vars(r)
		project := vars["project"]
		r.URL.Path = fmt.Sprintf("/v1/api/%s/files", project)

		token, err := auth.GetInternalAccessToken()
		if err != nil {
			logrus.Errorf("error handling forwarding artifact request failed to generate internal access token -%v", err)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		// TODO: Use http2 client if that was the incoming request protocol
		response, err := http.DefaultClient.Do(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		defer utils.CloseTheCloser(response.Body)

		// Copy headers and status code
		w.WriteHeader(response.StatusCode)
		for k, v := range response.Header {
			w.Header().Set(k, v[0])
		}

		// Copy the body
		n, err := io.Copy(w, response.Body)
		if err != nil {
			logrus.Errorf("Failed to copy upstream (%s) response to downstream - %s", r.URL.String(), err.Error())
		}

		logrus.Debugf("Successfully copied %d bytes from upstream server (%s)", n, r.URL.String())
	}

}
