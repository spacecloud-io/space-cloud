package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/spaceuptech/space-cli/model"
)

// AppendConfigToDisk creates a yml file or appends to existing
func AppendConfigToDisk(specObj *model.SpecObject, filename string) error {
	// Marshal spec object to yaml
	data, err := yaml.Marshal(specObj)
	if err != nil {
		return err
	}

	// Check if file exists. We need to ammend the file if it does.
	if FileExists(filename) {
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

// ReadSpecObjectsFromFile returns the spec objects present in the file
func ReadSpecObjectsFromFile(fileName string) ([]*model.SpecObject, error) {
	var specs []*model.SpecObject

	// Read the file first
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	// Split the files into independent objects
	dataStrings := strings.Split(string(data), "---")
	for _, dataString := range dataStrings {

		// Skip if string is too small to be a spec object
		if len(dataString) <= 5 {
			continue
		}

		// Unmarshal spec object
		spec := new(model.SpecObject)
		if err := UnmarshalYAML([]byte(dataString), spec); err != nil {
			return nil, err
		}

		// Append the spec object into the array
		specs = append(specs, spec)
	}

	return specs, nil
}

// CreateDirIfNotExist creates a directory if it doesn't already exists
func CreateDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateFileIfNotExist creates a file with the provided content if it doesn't already exists
func CreateFileIfNotExist(path, content string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return ioutil.WriteFile(path, []byte(content), 0755)
	}
	return nil
}

func generateYamlFile(credential *model.Credential) error {
	d, err := yaml.Marshal(&credential)
	if err != nil {
		return err
	}

	if err := CreateDirIfNotExist(GetSpaceCloudDirectory()); err != nil {
		logrus.Errorf("error in generating yaml file unable to create space cli directory - %v", err)
		return err
	}

	fileName := getAccountConfigPath()
	err = ioutil.WriteFile(fileName, d, 0644)
	if err != nil {
		return err
	}

	return nil
}

// FileExists checks if the file exists
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// UnmarshalYAML converts to map[string]interface{} instead of map[interface{}]interface{}.
func UnmarshalYAML(in []byte, out *model.SpecObject) error {
	if err := yaml.Unmarshal(in, out); err != nil {
		return err
	}
	out.Spec = cleanupMapValue(out.Spec)
	return nil
}

func cleanupInterfaceArray(in []interface{}) []interface{} {
	res := make([]interface{}, len(in))
	for i, v := range in {
		res[i] = cleanupMapValue(v)
	}
	return res
}

func cleanupInterfaceMap(in map[interface{}]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range in {
		res[fmt.Sprintf("%v", k)] = cleanupMapValue(v)
	}
	return res
}

func cleanupMapValue(v interface{}) interface{} {
	switch v := v.(type) {
	case []interface{}:
		return cleanupInterfaceArray(v)
	case map[interface{}]interface{}:
		return cleanupInterfaceMap(v)
	default:
		return v
	}
}
