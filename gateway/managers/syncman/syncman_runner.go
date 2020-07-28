package syncman

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"github.com/spaceuptech/space-cloud/gateway/managers/admin"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// HandleRunnerRequests handles requests of the runner
func (s *Manager) HandleRunnerRequests(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)
		reqParams, err := admin.IsTokenValid(token, "runner", "modify", nil)
		if err != nil {
			logrus.Errorf("error handling forwarding runner request failed to validate token -%v", err)
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		s.forwardRequestToRunner(ctx, w, r, admin, reqParams)
	}
}

// HandleRunnerApplySecret handles requests of the runner
func (s *Manager) HandleRunnerApplySecret(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]

		reqParams, err := admin.IsTokenValid(token, "secret", "modify", map[string]string{"project": projectID})
		if err != nil {
			logrus.Errorf("error handling forwarding runner request failed to validate token -%v", err)
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		s.forwardRequestToRunner(ctx, w, r, admin, reqParams)
	}
}

// HandleRunnerListSecret handles requests of the runner
func (s *Manager) HandleRunnerListSecret(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]

		reqParams, err := admin.IsTokenValid(token, "secret", "read", map[string]string{"project": projectID})
		if err != nil {
			logrus.Errorf("error handling forwarding runner request failed to validate token -%v", err)
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		s.forwardRequestToRunner(ctx, w, r, admin, reqParams)
	}
}

// HandleRunnerSetFileSecretRootPath handles requests of the runner
func (s *Manager) HandleRunnerSetFileSecretRootPath(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]

		reqParams, err := admin.IsTokenValid(token, "secret", "modify", map[string]string{"project": projectID})
		if err != nil {
			logrus.Errorf("error handling forwarding runner request failed to validate token -%v", err)
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		s.forwardRequestToRunner(ctx, w, r, admin, reqParams)
	}
}

// HandleRunnerDeleteSecret handles requests of the runner
func (s *Manager) HandleRunnerDeleteSecret(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]

		reqParams, err := admin.IsTokenValid(token, "secret", "modify", map[string]string{"project": projectID})
		if err != nil {
			logrus.Errorf("error handling forwarding runner request failed to validate token -%v", err)
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		s.forwardRequestToRunner(ctx, w, r, admin, reqParams)
	}
}

// HandleRunnerSetSecretKey handles requests of the runner
func (s *Manager) HandleRunnerSetSecretKey(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]

		reqParams, err := admin.IsTokenValid(token, "secret", "modify", map[string]string{"project": projectID})
		if err != nil {
			logrus.Errorf("error handling forwarding runner request failed to validate token -%v", err)
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		s.forwardRequestToRunner(ctx, w, r, admin, reqParams)
	}
}

// HandleRunnerDeleteSecretKey handles requests of the runner
func (s *Manager) HandleRunnerDeleteSecretKey(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]

		reqParams, err := admin.IsTokenValid(token, "secret", "modify", map[string]string{"project": projectID})
		if err != nil {
			logrus.Errorf("error handling forwarding runner request failed to validate token -%v", err)
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		s.forwardRequestToRunner(ctx, w, r, admin, reqParams)
	}
}

// HandleRunnerApplyService handles requests of the runner
func (s *Manager) HandleRunnerApplyService(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]

		reqParams, err := admin.IsTokenValid(token, "service", "modify", map[string]string{"project": projectID})
		if err != nil {
			logrus.Errorf("error handling forwarding runner request failed to validate token -%v", err)
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		s.forwardRequestToRunner(ctx, w, r, admin, reqParams)
	}
}

// HandleRunnerGetServices handles requests of the runner
func (s *Manager) HandleRunnerGetServices(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]

		reqParams, err := admin.IsTokenValid(token, "service", "read", map[string]string{"project": projectID})
		if err != nil {
			logrus.Errorf("error handling forwarding runner request failed to validate token -%v", err)
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		s.forwardRequestToRunner(ctx, w, r, admin, reqParams)
	}
}

// HandleRunnerGetDeploymentStatus handles requests of the runner
func (s *Manager) HandleRunnerGetDeploymentStatus(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]
		params, err := admin.IsTokenValid(token, "service", "read", map[string]string{"project": projectID})
		if err != nil {
			logrus.Errorf("error handling forwarding runner request failed to validate token -%v", err)
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		s.forwardRequestToRunner(ctx, w, r, admin, params)
	}
}

// HandleRunnerDeleteService handles requests of the runner
func (s *Manager) HandleRunnerDeleteService(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]

		reqParams, err := admin.IsTokenValid(token, "service", "modify", map[string]string{"project": projectID})
		if err != nil {
			logrus.Errorf("error handling forwarding runner request failed to validate token -%v", err)
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		s.forwardRequestToRunner(ctx, w, r, admin, reqParams)
	}
}

// HandleRunnerServiceRoutingRequest handles requests of the runner
func (s *Manager) HandleRunnerServiceRoutingRequest(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]

		reqParams, err := admin.IsTokenValid(token, "service-route", "modify", map[string]string{"project": projectID})
		if err != nil {
			logrus.Errorf("error handling forwarding runner request failed to validate token -%v", err)
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		s.forwardRequestToRunner(ctx, w, r, admin, reqParams)
	}
}

// HandleRunnerGetServiceRoutingRequest handles requests of the runner
func (s *Manager) HandleRunnerGetServiceRoutingRequest(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]

		reqParams, err := admin.IsTokenValid(token, "service-route", "read", map[string]string{"project": projectID})
		if err != nil {
			logrus.Errorf("error handling forwarding runner request failed to validate token -%v", err)
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		s.forwardRequestToRunner(ctx, w, r, admin, reqParams)
	}
}

// HandleRunnerGetServiceLogs handles requests of the runner
func (s *Manager) HandleRunnerGetServiceLogs(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userToken := utils.GetTokenFromHeader(r)
		defer logrus.Println("Closing handle of gateway for logs")

		vars := mux.Vars(r)
		projectID := vars["project"]
		utils.LogDebug("Forwarding request to runner for getting service logs", "syncman", "HandleRunnerGetServiceLogs", map[string]interface{}{})

		_, err := admin.IsTokenValid(userToken, "service-logs", "read", map[string]string{"project": projectID})
		if err != nil {
			logrus.Errorf("error handling forwarding runner request failed to validate token -%v", err)
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		_, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

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

		streamData := false
		// Copy headers and status code
		for k, v := range response.Header {
			// check if data is available in chunks
			if k == "X-Content-Type-Options" && v[0] == "nosniff" {
				streamData = true
			}
			w.Header().Set(k, v[0])
		}

		if streamData {
			if response.StatusCode != 200 {
				respBody := map[string]interface{}{}
				if err := json.NewDecoder(response.Body).Decode(&respBody); err != nil {
					_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
					return
				}
				_ = utils.SendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("received invalid status code (%d) got error - %v", response.StatusCode, respBody["error"]))
				return
			}

			rd := bufio.NewReader(response.Body)

			// get signal when client stops listening
			done := r.Context().Done()
			flusher, ok := w.(http.Flusher)
			if !ok {
				_ = utils.SendErrorResponse(w, http.StatusInternalServerError, "expected http.ResponseWriter to be an http.Flusher")
				return
			}
			w.Header().Set("X-Content-Type-Options", "nosniff")

			for {
				select {
				case <-done:
					utils.LogDebug("Client stopped listening", "syncman", "HandleRunnerGetServiceLogs", nil)
					return
				default:
					str, err := rd.ReadString('\n')
					if err != nil {
						if err == io.EOF {
							_ = utils.SendOkayResponse(w)
							return
						}
						_ = utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
					}
					if str != "\n" {
						fmt.Fprintf(w, "%s\n", str)
						flusher.Flush() // Trigger "chunked" encoding and send a chunk...
						time.Sleep(500 * time.Millisecond)
					}
				}
			}
		} else {
			_ = utils.SendErrorResponse(w, http.StatusBadRequest, "Missing headers X-Content-Type-Options & nosniff")
		}
	}
}

func (s *Manager) forwardRequestToRunner(ctx context.Context, w http.ResponseWriter, r *http.Request, admin *admin.Manager, params model.RequestParams) {

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
