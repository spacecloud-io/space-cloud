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
	"github.com/spaceuptech/helpers"

	"golang.org/x/net/context"

	"github.com/spaceuptech/space-cloud/gateway/managers/admin"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// HandleRunnerRequests handles requests of the runner
func (s *Manager) HandleRunnerRequests(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		reqParams, err := admin.IsTokenValid(ctx, token, "runner", "modify", nil)
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to forward runner request failed to validate token -%v", err), err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		reqParams = utils.ExtractRequestParams(r, reqParams, nil)

		s.forwardRequestToRunner(ctx, w, r, admin, reqParams)
	}
}

// HandleRunnerApplySecret handles requests of the runner
func (s *Manager) HandleRunnerApplySecret(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]

		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		reqParams, err := admin.IsTokenValid(ctx, token, "secret", "modify", map[string]string{"project": projectID, "id": vars["id"]})
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to forward  runner request failed to validate token -%v", err), err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		reqParams = utils.ExtractRequestParams(r, reqParams, nil)

		s.forwardRequestToRunner(ctx, w, r, admin, reqParams)
	}
}

// HandleRunnerListSecret handles requests of the runner
func (s *Manager) HandleRunnerListSecret(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]
		id := r.URL.Query().Get("id")
		if id == "" {
			id = "*"
		}

		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		reqParams, err := admin.IsTokenValid(ctx, token, "secret", "read", map[string]string{"project": projectID, "id": id})
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to forward  runner request failed to validate token -%v", err), err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		reqParams = utils.ExtractRequestParams(r, reqParams, nil)

		s.forwardRequestToRunner(ctx, w, r, admin, reqParams)
	}
}

// HandleRunnerSetFileSecretRootPath handles requests of the runner
func (s *Manager) HandleRunnerSetFileSecretRootPath(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]

		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		reqParams, err := admin.IsTokenValid(ctx, token, "secret", "modify", map[string]string{"project": projectID, "id": vars["id"]})
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to forward  runner request failed to validate token -%v", err), err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		reqParams = utils.ExtractRequestParams(r, reqParams, nil)

		s.forwardRequestToRunner(ctx, w, r, admin, reqParams)
	}
}

// HandleRunnerDeleteSecret handles requests of the runner
func (s *Manager) HandleRunnerDeleteSecret(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]

		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		reqParams, err := admin.IsTokenValid(ctx, token, "secret", "modify", map[string]string{"project": projectID, "id": vars["id"]})
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to forward  runner request failed to validate token -%v", err), err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		reqParams = utils.ExtractRequestParams(r, reqParams, nil)

		s.forwardRequestToRunner(ctx, w, r, admin, reqParams)
	}
}

// HandleRunnerSetSecretKey handles requests of the runner
func (s *Manager) HandleRunnerSetSecretKey(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]

		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		reqParams, err := admin.IsTokenValid(ctx, token, "secret", "modify", map[string]string{"project": projectID, "id": vars["id"], "key": vars["key"]})
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to forward  runner request failed to validate token -%v", err), err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		reqParams = utils.ExtractRequestParams(r, reqParams, nil)

		s.forwardRequestToRunner(ctx, w, r, admin, reqParams)
	}
}

// HandleRunnerDeleteSecretKey handles requests of the runner
func (s *Manager) HandleRunnerDeleteSecretKey(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]

		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		reqParams, err := admin.IsTokenValid(ctx, token, "secret", "modify", map[string]string{"project": projectID, "id": vars["id"], "key": vars["key"]})
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to forward  runner request failed to validate token -%v", err), err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		reqParams = utils.ExtractRequestParams(r, reqParams, nil)

		s.forwardRequestToRunner(ctx, w, r, admin, reqParams)
	}
}

// HandleRunnerApplyService handles requests of the runner
func (s *Manager) HandleRunnerApplyService(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]

		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		reqParams, err := admin.IsTokenValid(ctx, token, "service", "modify", map[string]string{"project": projectID, "id": vars["serviceId"], "version": vars["version"]})
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to forward  runner request failed to validate token -%v", err), err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		reqParams = utils.ExtractRequestParams(r, reqParams, nil)

		s.forwardRequestToRunner(ctx, w, r, admin, reqParams)
	}
}

// HandleRunnerGetServices handles requests of the runner
func (s *Manager) HandleRunnerGetServices(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]
		id := r.URL.Query().Get("serviceId")
		if id == "" {
			id = "*"
		}
		version := r.URL.Query().Get("version")
		if id == "" {
			version = "*"
		}

		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		reqParams, err := admin.IsTokenValid(ctx, token, "service", "read", map[string]string{"project": projectID, "id": id, "version": version})
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to forward  runner request failed to validate token -%v", err), err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		reqParams = utils.ExtractRequestParams(r, reqParams, nil)

		s.forwardRequestToRunner(ctx, w, r, admin, reqParams)
	}
}

// HandleRunnerGetDeploymentStatus handles requests of the runner
func (s *Manager) HandleRunnerGetDeploymentStatus(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]

		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		reqParams, err := admin.IsTokenValid(ctx, token, "service", "read", map[string]string{"project": projectID})
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to forward  runner request failed to validate token -%v", err), err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}
		// Create a context of execution
		reqParams = utils.ExtractRequestParams(r, reqParams, nil)

		s.forwardRequestToRunner(ctx, w, r, admin, reqParams)
	}
}

// HandleRunnerDeleteService handles requests of the runner
func (s *Manager) HandleRunnerDeleteService(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]

		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		reqParams, err := admin.IsTokenValid(ctx, token, "service", "modify", map[string]string{"project": projectID, "id": vars["serviceId"], "version": vars["version"]})
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to forward  runner request failed to validate token -%v", err), err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		reqParams = utils.ExtractRequestParams(r, reqParams, nil)

		s.forwardRequestToRunner(ctx, w, r, admin, reqParams)
	}
}

// HandleRunnerServiceRoutingRequest handles requests of the runner
func (s *Manager) HandleRunnerServiceRoutingRequest(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]

		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		reqParams, err := admin.IsTokenValid(ctx, token, "service-route", "modify", map[string]string{"project": projectID, "id": vars["serviceId"]})
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to forward  runner request failed to validate token -%v", err), err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		reqParams = utils.ExtractRequestParams(r, reqParams, nil)

		s.forwardRequestToRunner(ctx, w, r, admin, reqParams)
	}
}

// HandleRunnerGetServiceRoutingRequest handles requests of the runner
func (s *Manager) HandleRunnerGetServiceRoutingRequest(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]

		id := r.URL.Query().Get("id")
		if id == "" {
			id = "*"
		}

		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		reqParams, err := admin.IsTokenValid(ctx, token, "service-route", "read", map[string]string{"project": projectID, "id": id})
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to forward  runner request failed to validate token -%v", err), err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusUnauthorized, err.Error())
			return
		}

		// Create a context of execution
		reqParams = utils.ExtractRequestParams(r, reqParams, nil)

		s.forwardRequestToRunner(ctx, w, r, admin, reqParams)
	}
}

// HandleRunnerGetServiceLogs handles requests of the runner
func (s *Manager) HandleRunnerGetServiceLogs(admin *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userToken := utils.GetTokenFromHeader(r)

		vars := mux.Vars(r)
		projectID := vars["project"]

		_, err := admin.IsTokenValid(r.Context(), userToken, "service", "read", map[string]string{"project": projectID})
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(r.Context()), fmt.Sprintf("Unable to forward  runner request failed to validate token -%v", err), err, nil)
			_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusUnauthorized, err.Error())
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
			_ = helpers.Logger.LogError(helpers.GetRequestID(r.Context()), fmt.Sprintf("Unable to forward  runner request failed to generate internal access token -%v", err), err, nil)
			_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusInternalServerError, err.Error())
			return
		}
		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		// TODO: Use http2 client if that was the incoming request protocol
		response, err := http.DefaultClient.Do(r)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusInternalServerError, err.Error())
			return
		}
		defer utils.CloseTheCloser(response.Body)
		if response.StatusCode != 200 {
			respBody := map[string]interface{}{}
			if err := json.NewDecoder(response.Body).Decode(&respBody); err != nil {
				_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusInternalServerError, err.Error())
				return
			}
			_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusInternalServerError, fmt.Sprintf("received invalid status code (%d) got error - %v", response.StatusCode, respBody["error"]))
			return
		}
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
			rd := bufio.NewReader(response.Body)

			// get signal when client stops listening
			done := r.Context().Done()
			flusher, ok := w.(http.Flusher)
			if !ok {
				_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusInternalServerError, "expected http.ResponseWriter to be an http.Flusher")
				return
			}
			w.Header().Set("X-Content-Type-Options", "nosniff")

			for {
				select {
				case <-done:
					helpers.Logger.LogDebug(helpers.GetRequestID(r.Context()), "Connection got closed from client while reading logs of a service", map[string]interface{}{})
					return
				default:
					str, err := rd.ReadString('\n')
					if err != nil {
						if err == io.EOF {
							w.WriteHeader(http.StatusNoContent)
							return
						}
						_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusInternalServerError, err.Error())
					}
					if str != "\n" {
						fmt.Fprintf(w, "%s", str)
						flusher.Flush() // Trigger "chunked" encoding and send a chunk...
					}
				}
			}
		} else {
			_ = helpers.Response.SendErrorResponse(r.Context(), w, http.StatusBadRequest, "Missing headers X-Content-Type-Options & nosniff")
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
		_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to forward  runner request failed to generate internal access token -%v", err), err, nil)
		_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err.Error())
		return
	}
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// TODO: Use http2 client if that was the incoming request protocol
	response, err := http.DefaultClient.Do(r)
	if err != nil {
		_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err.Error())
		return
	}
	defer utils.CloseTheCloser(response.Body)

	// Copy the body
	w.WriteHeader(response.StatusCode)
	n, err := io.Copy(w, response.Body)
	if err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Failed to copy upstream (%s) response to downstream", r.URL.String()), err, nil)
	}

	helpers.Logger.LogDebug(helpers.GetRequestID(ctx), fmt.Sprintf("Successfully copied %d bytes from upstream server (%s)", n, r.URL.String()), nil)
}

// GetRunnerAddr returns runner address
func (s *Manager) GetRunnerAddr() string {
	return s.runnerAddr
}

// GetClusterType returns cluster type
func (s *Manager) GetClusterType(ctx context.Context, admin AdminSyncmanInterface) (string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.runnerAddr == "" {
		return "none", nil
	}

	token, err := admin.GetInternalAccessToken()
	if err != nil {
		return "", helpers.Logger.LogError(helpers.GetRequestID(ctx), "GetClusterType failed to generate internal access token", err, nil)
	}

	data := new(model.Response)

	err = s.MakeHTTPRequest(ctx, http.MethodGet, fmt.Sprintf("http://%s/v1/runner/cluster-type", s.runnerAddr), token, "", map[string]interface{}{}, data)
	if err != nil {
		return "", helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to fetch cluster type from runner", err, nil)
	}

	return data.Result.(string), err
}
