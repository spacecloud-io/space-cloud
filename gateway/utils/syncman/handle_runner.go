package syncman

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
)

// HandleRunnerRequests handles requests of the runner
func (s *Manager) HandleRunnerRequests(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := admin.IsTokenValid(utils.GetTokenFromHeader(r)); err != nil {
			logrus.Errorf("error handling forwarding runner request failed to validate token -%v", err)
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// http: Request.RequestURI can't be set in client requests.
		// http://golang.org/src/pkg/net/http/client.go
		r.RequestURI = ""

		// Get host from addr
		host := strings.Split(s.runnerAddr, ":")[0]

		// Change the request with the destination host, port and url
		r.Host = host
		r.URL.Host = s.runnerAddr

		// Set the url scheme to http
		r.URL.Scheme = "http"

		token, err := admin.GetInternalAccessToken()
		if err != nil {
			logrus.Errorf("error handling forwarding runner request failed to generate internal access token -%v", err)
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		// TODO: Use http2 client if that was the incoming request protocol
		response, err := http.DefaultClient.Do(r)
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		defer utils.CloseTheCloser(response.Body)

		// Copy headers and status code
		for k, v := range response.Header {
			w.Header().Set(k, v[0])
		}

		// Copy the body
		w.WriteHeader(response.StatusCode)
		n, err := io.Copy(w, response.Body)
		if err != nil {
			logrus.Errorf("Failed to copy upstream (%s) response to downstream - %s", r.URL.String(), err.Error())
		}

		logrus.Debugf("Successfully copied %d bytes from upstream server (%s)", n, r.URL.String())
	}
}

// GetRunnerAddr returns runner address
func (s *Manager) GetRunnerAddr() string {
	return s.runnerAddr
}

// GetRunnerType returns runner type
func (s *Manager) GetRunnerType(admin *admin.Manager) (string, error) {
	if s.runnerAddr == "" {
		return "none", nil
	}

	token, err := admin.GetInternalAccessToken()
	if err != nil {
		logrus.Errorf("error handling forwarding runner request failed to generate internal access token -%v", err)
		return "", err
	}

	httpReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/v1/runner/cluster-type", s.runnerAddr), nil)
	if err != nil {
		logrus.Errorf("error while getting runnerType in handler unable to create http request - %v", err)
		return "", err
	}
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	httpRes, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		logrus.Errorf("error while getting runnerType in handler unable to execute http request - %v", err)
		return "", err
	}

	data := new(model.Response)
	if err = json.NewDecoder(httpRes.Body).Decode(&data); err != nil {
		logrus.Errorf("error while getting runnerType in handler unable to decode response body -%v", err)
		return "", err
	}

	if httpRes.StatusCode != http.StatusOK {
		logrus.Errorf("error while getting runnerType in handler got http request -%v", httpRes.StatusCode)
		return "", fmt.Errorf("error while getting runnerType in handler got http request -%v -%v", httpRes.StatusCode, data.Error)
	}

	return data.Result.(string), err
}
