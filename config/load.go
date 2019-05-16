package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

func LoadEnvironmentVariable(p *Project) {
	if strings.HasPrefix(p.Secret, "$") {
		TempString := strings.TrimPrefix(p.Secret, "$")
		TempEnvVar, DoesItExist := os.LookupEnv(TempString)

		if DoesItExist {
			p.Secret = TempEnvVar
		}
	}

	for i := range p.Modules.Crud {
		if strings.HasPrefix(p.Modules.Crud[i].Conn, "$") {
			TempStringC := strings.TrimPrefix(p.Modules.Crud[i].Conn, "$")
			TempEnvVarC, DoesItExistC := os.LookupEnv(TempStringC)

			if DoesItExistC {
				p.Modules.Crud[i].Conn = TempEnvVarC
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

	LoadEnvironmentVariable(conf)
	return conf, nil
}
