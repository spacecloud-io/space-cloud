package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ghodss/yaml"
)

func loadEnvironmentVariable(c *Config) {
	for _, p := range c.Projects {
		for i, secret := range p.Secrets {
			if strings.HasPrefix(secret.Secret, "$") {
				tempString := strings.TrimPrefix(secret.Secret, "$")
				tempEnvVar, present := os.LookupEnv(tempString)

				if present {
					p.Secrets[i].Secret = tempEnvVar
				}
			}
		}
		for _, value := range p.Modules.Crud {
			if strings.HasPrefix(value.Conn, "$") {
				tempStringC := strings.TrimPrefix(value.Conn, "$")
				tempEnvVarC, presentC := os.LookupEnv(tempStringC)

				if presentC {
					value.Conn = tempEnvVarC
				}
			}
		}
	}
}

// LoadConfigFromFile loads the config from the provided file path
func LoadConfigFromFile(path string) (*Config, error) {
	// Load the file in memory
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Marshal the configuration
	conf := new(Config)
	if strings.HasSuffix(path, "json") {
		err = json.Unmarshal(dat, conf)
	} else {
		err = yaml.Unmarshal(dat, conf)
	}
	if err != nil {
		return nil, err
	}

	loadEnvironmentVariable(conf)
	return conf, nil
}
