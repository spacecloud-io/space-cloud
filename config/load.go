package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

func loadEnvironmentVariable(p *Project) {
	if strings.HasPrefix(p.Secret, "$") {
		tempString := strings.TrimPrefix(p.Secret, "$")
		tempEnvVar, doesItExist := os.LookupEnv(tempString)

		if doesItExist {
			p.Secret = tempEnvVar
		}
	}

	for _, i := range p.Modules.Crud {
		if strings.HasPrefix(i.Conn, "$") {
			tempStringC := strings.TrimPrefix(i.Conn, "$")
			tempEnvVarC, doesItExistC := os.LookupEnv(tempStringC)

			if doesItExistC {
				i = tempEnvVarC
			}
		}
	}
}

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

	loadEnvironmentVariable(conf)
	return conf, nil
}
