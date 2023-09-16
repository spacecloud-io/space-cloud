package utils

import "github.com/spf13/viper"

// GetInstanceID returns the instance id that has been configured via viper
func GetInstanceID() string {
	return viper.GetString("id")
}
