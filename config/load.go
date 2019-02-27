package config

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// LoadConfigFromFile loads the config from the provided file path
func LoadConfigFromFile(path string) (*Project, error) {
	// Load the file in memory
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Marshal the configuration
	conf := new(Project)
	if strings.HasSuffix(path, "json") {
		err = json.Unmarshal(dat, conf)
	} else {
		err = yaml.Unmarshal(dat, conf)
	}
	if err != nil {
		return nil, err
	}

	return conf, nil
}
