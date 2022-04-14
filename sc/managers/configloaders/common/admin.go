package common

import (
	"encoding/json"

	"github.com/spf13/viper"
)

func prepareAdminApp() json.RawMessage {
	user := viper.GetString("admin-user")
	pass := viper.GetString("admin-pass")
	secret := viper.GetString("admin-secret")
	isDev := viper.GetBool("dev")

	config := map[string]interface{}{
		"user":   user,
		"pass":   pass,
		"secret": secret,
		"isDev":  isDev,
	}

	data, _ := json.Marshal(config)
	return data
}
