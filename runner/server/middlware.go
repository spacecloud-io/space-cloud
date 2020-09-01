package server

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/segmentio/ksuid"
	"github.com/spaceuptech/helpers"
)

func loggerMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		requestID := r.Header.Get(helpers.HeaderRequestID)
		if requestID == "" {
			// set a new request id header of request
			requestID = ksuid.New().String()
			r.Header.Set(helpers.HeaderRequestID, requestID)
		}

		var reqBody []byte
		if r.Header.Get("Content-Type") == "application/json" {
			reqBody, _ = ioutil.ReadAll(r.Body)
			r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
		}

		helpers.Logger.LogInfo(requestID, "Request", map[string]interface{}{"method": r.Method, "url": r.URL.Path, "queryVars": r.URL.Query(), "body": string(reqBody)})
		next.ServeHTTP(w, r.WithContext(helpers.CreateContext(r)))

	})
}
