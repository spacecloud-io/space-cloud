package config

import (
	"encoding/json"
	"io/ioutil"
)

// LoadConfigFromFile loads the config from the provided file path
func LoadConfigFromFile(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	conf := new(Config)
	err = json.Unmarshal(data, conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
