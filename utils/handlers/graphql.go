package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils/graphql"
)

// HandleGraphQLRequest creates the graphql operation endpoint
func HandleGraphQLRequest(graphql *graphql.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		projectID := vars["project"]
		pid := graphql.GetProjectID()

		// Load the request from the body
		req := model.GraphQLRequest{}
		json.NewDecoder(r.Body).Decode(&req)
		defer r.Body.Close()

		if projectID != pid {
			//throw some error
			w.WriteHeader(http.StatusInternalServerError) //http status codee
			json.NewEncoder(w).Encode(map[string]string{"error": "project id doesn't match"})
			return
		}

		// Get the path parameters
		token := getRequestMetaData(r).token

		ch := make(chan struct{}, 1)

		graphql.ExecGraphQLQuery(ctx, &req, token, func(op interface{}, err error) {
			defer func() { ch <- struct{}{} }()
			if err != nil {
				errMes := map[string]interface{}{"message": err.Error()}
				json.NewEncoder(w).Encode(map[string]interface{}{"errors": []interface{}{errMes}})
				return
			}

			w.WriteHeader(http.StatusOK) //http status codee
			json.NewEncoder(w).Encode(map[string]interface{}{"data": op})
			return
		})

		select {
		case <-ch:
			return
		case <-time.After(10 * time.Second):
			log.Println("GraphQL Handler: Request timed out")
			return
		}
	}

}
