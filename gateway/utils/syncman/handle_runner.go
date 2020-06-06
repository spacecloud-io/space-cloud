package syncman

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"

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

		streamData := false
		// Copy headers and status code
		for k, v := range response.Header {
			// check if data is available in chunks
			if k == "X-Content-Type-Options" && v[0] == "nosniff" {
				streamData = true
			}
			fmt.Println("header", k, v)
			w.Header().Set(k, v[0])
		}
		if streamData {
			rd := bufio.NewReader(response.Body)

			// get signal when client stops listening
			done := r.Context().Done()
			flusher, ok := w.(http.Flusher)
			if !ok {
				panic("expected http.ResponseWriter to be an http.Flusher")
			}
			w.Header().Set("X-Content-Type-Options", "nosniff")

		loop:
			for {
				select {
				case <-done:
					glog.Infof("Client stopped listening")
					break loop
				default:
					str, _ := rd.ReadString('\n')
					if str != "\n" {
						fmt.Fprintf(w, "%s\n", str)
						flusher.Flush() // Trigger "chunked" encoding and send a chunk...
						time.Sleep(500 * time.Millisecond)
					}
				}
			}
			logrus.Debugf("Successfully received data from upstream server (%s)", r.URL.String())
		} else {
			// Copy the body
			w.WriteHeader(response.StatusCode)
			n, err := io.Copy(w, response.Body)
			if err != nil {
				logrus.Errorf("Failed to copy upstream (%s) response to downstream - %s", r.URL.String(), err.Error())
			}

			logrus.Debugf("Successfully copied %d bytes from upstream server (%s)", n, r.URL.String())
		}
	}
}

// GetRunnerAddr returns runner address
func (s *Manager) GetRunnerAddr() string {
	return s.runnerAddr
}

// GetClusterType returns cluster type
func (s *Manager) GetClusterType(admin *admin.Manager) (string, error) {
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
