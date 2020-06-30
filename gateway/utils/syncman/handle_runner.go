package syncman

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
)

// HandleRunnerRequests handles requests of the runner
func (s *Manager) HandleRunnerRequests(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)
		if err := admin.IsTokenValid(token, "runner", "modify", nil); err != nil {
			logrus.Errorf("error handling forwarding runner request failed to validate token -%v", err)
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		s.forwardRequestToRunner(w, r, admin)
	}
}

// HandleRunnerApplySecret handles requests of the runner
func (s *Manager) HandleRunnerApplySecret(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]

		if err := admin.IsTokenValid(token, "secret", "modify", map[string]string{"project": projectID}); err != nil {
			logrus.Errorf("error handling forwarding runner request failed to validate token -%v", err)
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		s.forwardRequestToRunner(w, r, admin)
	}
}

// HandleRunnerListSecret handles requests of the runner
func (s *Manager) HandleRunnerListSecret(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]

		if err := admin.IsTokenValid(token, "secret", "read", map[string]string{"project": projectID}); err != nil {
			logrus.Errorf("error handling forwarding runner request failed to validate token -%v", err)
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		s.forwardRequestToRunner(w, r, admin)
	}
}

// HandleRunnerSetFileSecretRootPath handles requests of the runner
func (s *Manager) HandleRunnerSetFileSecretRootPath(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]

		if err := admin.IsTokenValid(token, "secret", "modify", map[string]string{"project": projectID}); err != nil {
			logrus.Errorf("error handling forwarding runner request failed to validate token -%v", err)
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		s.forwardRequestToRunner(w, r, admin)
	}
}

// HandleRunnerDeleteSecret handles requests of the runner
func (s *Manager) HandleRunnerDeleteSecret(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]

		if err := admin.IsTokenValid(token, "secret", "modify", map[string]string{"project": projectID}); err != nil {
			logrus.Errorf("error handling forwarding runner request failed to validate token -%v", err)
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		s.forwardRequestToRunner(w, r, admin)
	}
}

// HandleRunnerSetSecretKey handles requests of the runner
func (s *Manager) HandleRunnerSetSecretKey(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]

		if err := admin.IsTokenValid(token, "secret", "modify", map[string]string{"project": projectID}); err != nil {
			logrus.Errorf("error handling forwarding runner request failed to validate token -%v", err)
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		s.forwardRequestToRunner(w, r, admin)
	}
}

// HandleRunnerDeleteSecretKey handles requests of the runner
func (s *Manager) HandleRunnerDeleteSecretKey(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]

		if err := admin.IsTokenValid(token, "secret", "modify", map[string]string{"project": projectID}); err != nil {
			logrus.Errorf("error handling forwarding runner request failed to validate token -%v", err)
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		s.forwardRequestToRunner(w, r, admin)
	}
}

// HandleRunnerApplyService handles requests of the runner
func (s *Manager) HandleRunnerApplyService(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]

		if err := admin.IsTokenValid(token, "service", "modify", map[string]string{"project": projectID}); err != nil {
			logrus.Errorf("error handling forwarding runner request failed to validate token -%v", err)
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		s.forwardRequestToRunner(w, r, admin)
	}
}

// HandleRunnerGetServices handles requests of the runner
func (s *Manager) HandleRunnerGetServices(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]

		if err := admin.IsTokenValid(token, "service", "read", map[string]string{"project": projectID}); err != nil {
			logrus.Errorf("error handling forwarding runner request failed to validate token -%v", err)
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		s.forwardRequestToRunner(w, r, admin)
	}
}

// HandleRunnerDeleteService handles requests of the runner
func (s *Manager) HandleRunnerDeleteService(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]

		if err := admin.IsTokenValid(token, "service", "modify", map[string]string{"project": projectID}); err != nil {
			logrus.Errorf("error handling forwarding runner request failed to validate token -%v", err)
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		s.forwardRequestToRunner(w, r, admin)
	}
}

// HandleRunnerServiceRoutingRequest handles requests of the runner
func (s *Manager) HandleRunnerServiceRoutingRequest(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]

		if err := admin.IsTokenValid(token, "service-route", "modify", map[string]string{"project": projectID}); err != nil {
			logrus.Errorf("error handling forwarding runner request failed to validate token -%v", err)
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		s.forwardRequestToRunner(w, r, admin)
	}
}

// HandleRunnerGetServiceRoutingRequest handles requests of the runner
func (s *Manager) HandleRunnerGetServiceRoutingRequest(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]

		if err := admin.IsTokenValid(token, "service-route", "read", map[string]string{"project": projectID}); err != nil {
			logrus.Errorf("error handling forwarding runner request failed to validate token -%v", err)
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		s.forwardRequestToRunner(w, r, admin)
	}
}

func (s *Manager) forwardRequestToRunner(w http.ResponseWriter, r *http.Request, admin *admin.Manager) {

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

// GetRunnerAddr returns runner address
func (s *Manager) GetRunnerAddr() string {
	return s.runnerAddr
}

// GetClusterType returns cluster type
func (s *Manager) GetClusterType(admin AdminSyncmanInterface) (string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.runnerAddr == "" {
		return "none", nil
	}

	token, err := admin.GetInternalAccessToken()
	if err != nil {
		logrus.Errorf("GetClusterType failed to generate internal access token -%v", err)
		return "", err
	}

	data := new(model.Response)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = s.MakeHTTPRequest(ctx, http.MethodGet, fmt.Sprintf("http://%s/v1/runner/cluster-type", s.runnerAddr), token, "", map[string]interface{}{}, data)
	if err != nil {
		return "", err
	}

	return data.Result.(string), err
}
