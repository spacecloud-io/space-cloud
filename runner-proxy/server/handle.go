package server

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/spaceuptech/helpers"
	"github.com/spaceuptech/space-cloud/runner-proxy/utils"
)

func (s *Server) handleProxy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Minute)
		defer cancel()

		// Close the body of the request
		defer utils.CloseTheCloser(r.Body)

		// http: Request.RequestURI can't be set in client requests.
		// http://golang.org/src/pkg/net/http/client.go
		r.RequestURI = ""

		// Get the meta data from headers
		project := r.Header.Get("x-og-project")
		service := r.Header.Get("x-og-service")
		ogHost := r.Header.Get("x-og-host")
		ogPort := r.Header.Get("x-og-port")
		ogVersion := r.Header.Get("x-og-version")

		// Delete the headers
		r.Header.Del("x-og-project")
		r.Header.Del("x-og-service")
		r.Header.Del("x-og-host")
		r.Header.Del("x-og-port")
		r.Header.Del("x-og-version")

		// Change the destination with the original host and port
		r.Host = ogHost
		r.URL.Host = fmt.Sprintf("%s:%s", ogHost, ogPort)

		// Set the url scheme to http
		r.URL.Scheme = "http"

		helpers.Logger.LogDebug(helpers.GetRequestID(ctx), fmt.Sprintf("Proxy is making request to host (%s) port (%s)", ogHost, ogPort), nil)

		// get token from header
		token, err := s.auth.CreateToken(ctx, map[string]interface{}{})
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusServiceUnavailable, err)
			return
		}

		// Update the ttl of cached deployment
		id := fmt.Sprintf("%s-%s-%s", project, service, ogVersion)
		exist := s.cache.GetDeployment(id)
		if !exist {

			// makes http request to instruct driver to scale up
			var vPtr interface{}
			url := fmt.Sprintf("/v1/runner/%s/scale-up/%s/%s", project, service, ogVersion)
			if err := utils.MakeHTTPRequest(ctx, "POST", url, token, "", map[string]interface{}{}, &vPtr); err != nil {
				_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusServiceUnavailable, err)
				return
			}

			// Wait for the service to scale up
			url = fmt.Sprintf("/v1/runner/%s/wait/%s/%s", project, service, ogVersion)
			if err := utils.MakeHTTPRequest(ctx, "GET", url, token, "", map[string]interface{}{}, &vPtr); err != nil {
				_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusServiceUnavailable, err)
				return
			}
		}

		//after successfull http request make a new entry in TTLMap with id as key
		s.cache.Put(id)

		var res *http.Response
		for i := 0; i < 5; i++ {
			// Fire the request
			var err error
			res, err = http.DefaultClient.Do(r)
			if err != nil {
				_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusInternalServerError, err)
				return
			}

			// TODO: Make this retry logic better
			if res.StatusCode != http.StatusServiceUnavailable {
				break
			}

			time.Sleep(350 * time.Millisecond)

			// Close the body
			_, _ = io.Copy(ioutil.Discard, res.Body)
			utils.CloseTheCloser(res.Body)
		}

		defer utils.CloseTheCloser(res.Body)

		// Copy headers and status code
		for k, v := range res.Header {
			w.Header()[k] = v
		}

		w.WriteHeader(res.StatusCode)
		_, _ = io.Copy(w, res.Body)
	}
}
