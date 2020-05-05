package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/syncman"
)

// HandleGraphQLRequest executes graphql queries
func HandleGraphQLRequest(modules *modules.Modules, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		projectID := vars["project"]

		projectConfig, err := syncMan.GetConfig(projectID)
		if err != nil {
			logrus.Errorf("Error handling graphql query execution unable to get project config of %s - %s", projectID, err.Error())
			_ = utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if projectConfig.ContextTimeGraphQL == 0 {
			projectConfig.ContextTimeGraphQL = 10
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(projectConfig.ContextTimeGraphQL)*time.Second)
		defer cancel()

		graphql := modules.GraphQL()

		// Load the request from the body
		req := model.GraphQLRequest{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		defer utils.CloseTheCloser(r.Body)

		// Get the path parameters
		token := getRequestMetaData(r).token

		ch := make(chan struct{}, 1)

		graphql.ExecGraphQLQuery(ctx, &req, token, func(op interface{}, err error) {
			defer func() { ch <- struct{}{} }()
			if err != nil {
				errMes := map[string]interface{}{"message": err.Error()}
				_ = utils.SendResponse(w, http.StatusOK, map[string]interface{}{"errors": []interface{}{errMes}})
				return
			}
			_ = utils.SendResponse(w, http.StatusOK, map[string]interface{}{"data": op})
		})

		select {
		case <-ch:
			return
		case <-time.After(time.Duration(projectConfig.ContextTimeGraphQL) * time.Second):
			log.Println("GraphQL Handler: Request timed out")
			errMes := map[string]interface{}{"message": "GraphQL Handler: Request timed out"}
			_ = utils.SendResponse(w, http.StatusOK, map[string]interface{}{"errors": []interface{}{errMes}})
			return
		}
	}

}
