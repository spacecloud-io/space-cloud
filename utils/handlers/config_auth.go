package handlers

import(
	"encoding/json"
	"net/http"
	"strings"
	"github.com/gorilla/mux"
	"github.com/spaceuptech/space-cloud/utils/admin"
	"github.com/spaceuptech/space-cloud/utils/syncman"
	"github.com/spaceuptech/space-cloud/config"
)
// HandleUserManagement returns the handler to get the project config and validate the user via a REST endpoint
func HandleUserManagement(adminMan *admin.Manager, syncMan *syncman.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		tokens, ok := r.Header["Authorization"]
		if !ok {
			tokens = []string{""}
		}
		token := strings.TrimPrefix(tokens[0], "Bearer ")

		if err:=adminMan.IsTokenValid(token);err!=nil{
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error":err.Error()})
			return
		}

		// Load the body of the request
		value := new(config.AuthStub)
		err := json.NewDecoder(r.Body).Decode(value)
		defer r.Body.Close()
		vars := mux.Vars(r)
		project := vars["project"]
		provider := vars["provider"]

		projectConfig,err := syncMan.GetConfig(project)
		if err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error":err.Error()})
			return
        }

		projectConfig.Modules.Auth[provider]=value
		
		// Sync the config
		if err := syncMan.SetProjectConfig(projectConfig); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Give a positive acknowledgement
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}