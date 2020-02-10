package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/graphql"
	"github.com/spaceuptech/space-cloud/gateway/utils/syncman"
)

// HandleGraphQLRequest executes graphql queries
func HandleGraphQLRequest(graphql *graphql.Module, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		c := new(config.Project)
		err := json.NewDecoder(r.Body).Decode(c)
		defer utils.CloseTheCloser(r.Body)

		// Throw error if request was of incorrect type
		if err != nil {
			logrus.Errorf("Config was of invalid type - - %s", err.Error())
			utils.SendErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		vars := mux.Vars(r)
		projectID := vars["project"]

		projectConfig, err := syncMan.GetConfig(projectID)
		if err != nil {
			logrus.Errorf("Config was of invalid type - - %s", err.Error())
			utils.SendErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(projectConfig.ContextTime)*time.Second)
		defer cancel()

		pid := graphql.GetProjectID()

		// Load the request from the body
		req := model.GraphQLRequest{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		defer utils.CloseTheCloser(r.Body)

		w.Header().Set("Content-Type", "application/json")

		if projectID != pid {
			// throw some error
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "project id doesn't match"})
			return
		}

		// Get the path parameters
		token := getRequestMetaData(r).token

		ch := make(chan struct{}, 1)

		graphql.ExecGraphQLQuery(ctx, &req, token, func(op interface{}, err error) {
			defer func() { ch <- struct{}{} }()
			if err != nil {
				errMes := map[string]interface{}{"message": err.Error()}
				_ = json.NewEncoder(w).Encode(map[string]interface{}{"errors": []interface{}{errMes}})
				return
			}

			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": op})
			// return
		})

		select {
		case <-ch:
			return
		case <-time.After(time.Duration(projectConfig.ContextTime) * time.Second):
			log.Println("GraphQL Handler: Request timed out")
			errMes := map[string]interface{}{"message": "GraphQL Handler: Request timed out"}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"errors": []interface{}{errMes}})
			return
		}
	}

}
