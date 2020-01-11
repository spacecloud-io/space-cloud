package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
	"github.com/spaceuptech/space-cloud/gateway/utils/letsencrypt"
)

func HandleLetsEncryptWhitelistedDomain(adminMan *admin.Manager, l *letsencrypt.LetsEncrypt) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		value := config.LetsEncrypt{}
		_ = json.NewDecoder(r.Body).Decode(&value)
		defer r.Body.Close()

		// Check if the request is authorised
		if err := adminMan.IsTokenValid(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if err := l.SetConfig(value.WhitelistedDomains); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{})
	}
}
