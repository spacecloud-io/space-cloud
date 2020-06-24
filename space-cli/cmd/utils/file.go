package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/ghodss/yaml"

	"github.com/spaceuptech/space-cli/cmd/model"
	"github.com/spaceuptech/space-cli/cmd/utils/file"
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
		f, err := file.File.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			return err
		}

		defer func() {
			_ = file.File.Close(f)
		}()

		_, err = file.File.Write(f, append([]byte("---\n"), data...))
		return err
	}

	// Create a new file with out specs
	return file.File.WriteFile(filename, data, 0755)
}

// ReadSpecObjectsFromFile returns the spec objects present in the file
func ReadSpecObjectsFromFile(fileName string) ([]*model.SpecObject, error) {
	var specs []*model.SpecObject

	// Read the file first
	data, err := file.File.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		logrus.Infoln("empty file provided")
		return nil, nil
	}

	// Split the files into independent objects
	dataStrings := makeSpecString(string(data))
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

func makeSpecString(raw string) []string {
	lines := strings.Split(strings.Replace(raw, "\r\n", "\n", -1), "\n")
	var finalArray []string
	var tempArray []string
	for _, line := range lines {
		if line == "---" {
			finalArray = append(finalArray, strings.Join(tempArray, "\n"))
			tempArray = make([]string, 0)
			continue
		}
		tempArray = append(tempArray, line)
	}

	if len(tempArray) > 0 {
		finalArray = append(finalArray, strings.Join(tempArray, "\n"))
	}

	return finalArray
}

// CreateDirIfNotExist creates a directory if it doesn't already exists
func CreateDirIfNotExist(dir string) error {
	if _, err := file.File.Stat(dir); file.File.IsNotExist(err) {
		err = file.File.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateFileIfNotExist creates a file with the provided content if it doesn't already exists
func CreateFileIfNotExist(path, content string) error {
	if _, err := file.File.Stat(path); file.File.IsNotExist(err) {
		return file.File.WriteFile(path, []byte(content), 0755)
	}
	return nil
}

// CreateConfigFile create empty config file
func CreateConfigFile(path string) error {
	val := map[string]interface{}{"projects": []struct{}{}, "admin": struct{}{}}
	b, err := yaml.Marshal(val)
	if err != nil {
		return err
	}
	if _, err := file.File.Stat(path); file.File.IsNotExist(err) {
		return file.File.WriteFile(path, b, 0755)
	}
	return nil
}

// GenerateAccountsFile generates the yaml file for accounts
func GenerateAccountsFile(credential *model.Credential) error {
	d, err := yaml.Marshal(&credential)
	if err != nil {
		return err
	}

	if err := CreateDirIfNotExist(GetSpaceCloudDirectory()); err != nil {
		_ = LogError(fmt.Sprintf("error in generating yaml file unable to create space cli directory - %v", err), nil)
		return err
	}

	fileName := getAccountConfigPath()
	err = file.File.WriteFile(fileName, d, 0644)
	if err != nil {
		return err
	}

	return nil
}

// FileExists checks if the file exists
func FileExists(filename string) bool {
	info, err := file.File.Stat(filename)
	if file.File.IsNotExist(err) {
		return false
	}
	return !file.File.IsDir(info)
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
