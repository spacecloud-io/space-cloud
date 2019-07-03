package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

// StoreConfigToFile stores the config file to disk
func StoreConfigToFile(conf *Config, path string) error {
	var data []byte
	var err error

	if strings.HasSuffix(path, ".yaml") {
		data, err = yaml.Marshal(conf)
	} else if strings.HasSuffix(path, ".json") {
		data, err = json.Marshal(conf)
	} else {
		return errors.New("Invalid config file type")
	}

	// Check if error occured while marshaling
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0644)
}
