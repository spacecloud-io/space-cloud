package config

import (
	"github.com/spf13/viper"
)

// LoadConfigFromFile loads the config from the provided file path
func LoadConfigFromFile(path string) (*Project, error) {
	viper.SetConfigFile(path)

	conf := new(Project)
	viper.Unmarshal(conf)

	viper.ReadInConfig()
	err := viper.Unmarshal(&conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
