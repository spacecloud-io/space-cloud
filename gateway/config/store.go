package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/ghodss/yaml"

	"github.com/spaceuptech/space-cloud/gateway/utils"
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
		return utils.LogError(fmt.Sprintf("Invalid config file type (%s) provided", path), "config", "StoreConfigToFile", nil)
	}

	// Check if error occured while marshaling
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0644)
}
