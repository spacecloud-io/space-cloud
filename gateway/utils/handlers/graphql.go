package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils/projects"
)

// HandleGraphQLRequest creates the graphql operation endpoint
func HandleGraphQLRequest(p *projects.Projects) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		projectID := vars["project"]

		projectConfig, err := syncMan.GetConfig(projectID)
		if err != nil {
			logrus.Errorf("Error handling graphql query execution unable to get project config of %s - %s", projectID, err.Error())
			utils.SendErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		if projectConfig.ContextTime == 0 {
			projectConfig.ContextTime = 10
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(projectConfig.ContextTime)*time.Second)
		defer cancel()

		// Load the request from the body
		req := model.GraphQLRequest{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		defer utils.CloseTheCloser(r.Body)

		state, err := p.LoadProject(project)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		ch := make(chan struct{}, 1)

		state.Graph.ExecGraphQLQuery(ctx, &req, token, func(op interface{}, err error) {
			defer func() { ch <- struct{}{} }()

			if err != nil {
				errMes := map[string]interface{}{"message": err.Error()}
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(map[string]interface{}{"errors": []interface{}{errMes}})
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": op})
			// return
		})

		select {
		case <-ch:
			return
		case <-time.After(time.Duration(projectConfig.ContextTime) * time.Second):
			errMes := map[string]interface{}{"message": "GraphQL Handler: Request timed out"}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"errors": []interface{}{errMes}})
			return
		}
	}
}
