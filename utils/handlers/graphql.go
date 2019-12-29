package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils/projects"
)

// HandleGraphQLRequest creates the graphql operation endpoint
func HandleGraphQLRequest(p *projects.Projects) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		project := vars["project"]

		// Get the path parameters
		token := getRequestMetaData(r).token

		// Load the request from the body
		req := model.GraphQLRequest{}
		json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()

		w.Header().Set("Content-Type", "application/json")

		state, err := p.LoadProject(project)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		ch := make(chan struct{}, 1)

		state.Graph.ExecGraphQLQuery(ctx, &req, token, func(op interface{}, err error) {
			defer func() { ch <- struct{}{} }()

			if err != nil {
				errMes := map[string]interface{}{"message": err.Error()}
				json.NewEncoder(w).Encode(map[string]interface{}{"errors": []interface{}{errMes}})
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{"data": op})
			return
		})

		select {
		case <-ch:
			return
		case <-time.After(10 * time.Second):
			errMes := map[string]interface{}{"message": "GraphQL Handler: Request timed out"}
			json.NewEncoder(w).Encode(map[string]interface{}{"errors": []interface{}{errMes}})
			return
		}
	}
}
