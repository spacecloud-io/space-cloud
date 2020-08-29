package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/runner/model"
	"github.com/spaceuptech/space-cloud/runner/utils"
)

var upgrader = websocket.Upgrader{}

func (s *Server) handleWebsocketRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Could not upgrade to websocket (autoscaler):", err, nil)
			return
		}
		defer utils.CloseTheCloser(c)

		// Check if token is valid
		claims, err := s.auth.VerifyProxyToken(utils.GetToken(r))
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to verify autoscaler socket connection", err, nil)
			return
		}

		// Extract node id, project id and service name
		nodeIDTemp, ok1 := claims["id"]
		projectTemp, ok2 := claims["project"]
		serviceTemp, ok3 := claims["service"]
		versionTemp, ok4 := claims["version"]
		if !ok1 || !ok2 || !ok3 || !ok4 {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to establish autoscaler socket connection - token does not contain valid claims", nil, nil)
			return
		}

		nodeID := nodeIDTemp.(string)
		project := projectTemp.(string)
		service := serviceTemp.(string)
		version := versionTemp.(string)

		for {
			msg := new(model.ProxyMessage)
			if err := c.ReadJSON(msg); err != nil {
				_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Failed to receive message from proxy (%s:%s)", project, service), err, nil)
				return
			}

			// Set crucial meta data
			msg.NodeID = nodeID
			msg.Project = project
			msg.Service = service
			msg.Version = version

			// Append msg to disk
			s.chAppend <- msg
		}
	}
}
