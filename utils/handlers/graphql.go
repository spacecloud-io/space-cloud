package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils/graphql"
)

// HandleCrudCreate creates the create operation endpoint
func HandleGraphQLRequest(graphql *graphql.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Create a context of execution
		_, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		vars := mux.Vars(r)
		projectID := vars["project"]
		pid := graphql.GetProjectID()
		if projectID == pid {

			// Get the path parameters
			meta := getRequestMetaData(r)
			// Load the request from the body
			req := model.GraphQLRequest{}

			json.NewDecoder(r.Body).Decode(&req)
			defer r.Body.Close()

			op, err := graphql.ExecGraphQLQuery(&req, meta.token)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError) //http status codee
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			w.WriteHeader(http.StatusOK) //http status codee
			json.NewEncoder(w).Encode(map[string]string{"result": op.(string)})
			return
		}
		//throw some error
		w.WriteHeader(http.StatusInternalServerError) //http status codee
		json.NewEncoder(w).Encode(map[string]string{"error": "project id doesn't match"})
		return
	}
}
