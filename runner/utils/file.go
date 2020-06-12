package utils

import (
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"

	"github.com/spaceuptech/space-cloud/runner/model"
)

// AppendConfigToDisk creates a yml file or appends to existing
func AppendConfigToDisk(specObj *model.SpecObject, filename string) error {
	// Marshal spec object to yaml
	data, err := yaml.Marshal(specObj)
	if err != nil {
		return err
	}

	// Check if file exists. We need to ammend the file if it does.
	if fileExists(filename) {
		f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			return err
		}

		defer func() {
			_ = f.Close()
		}()

		_, err = f.Write(append([]byte("---\n"), data...))
		return err
	}

	// Create a new file with out specs
	return ioutil.WriteFile(filename, data, 0755)
}
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
