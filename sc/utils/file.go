package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/sirupsen/logrus"
	"github.com/spacecloud-io/space-cloud/model"
	"github.com/spacecloud-io/space-cloud/utils/file"
)

// StoreConfigToFile stores the config file to disk
func StoreConfigToFile(conf interface{}, path string) error {
	var data []byte
	var err error

	if strings.HasSuffix(path, ".yaml") {
		data, err = yaml.Marshal(conf)
	} else if strings.HasSuffix(path, ".json") {
		data, err = json.Marshal(conf)
	} else {
		return fmt.Errorf("invalid config file type (%s) provided", path)
	}

	// Check if error occured while marshaling
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0644)
}

// LoadFile loads a yaml or json file
func LoadFile(path string, ptr interface{}) error {
	format := "yaml"
	if strings.HasSuffix(path, "json") {
		format = "json"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	switch format {
	case "yaml":
		return yaml.Unmarshal(data, ptr)
	case "json":
		return json.Unmarshal(data, ptr)
	default:
		return fmt.Errorf("invalid format '%s' provided", format)
	}
}

// StoreFile stores a yaml or json file
func StoreFile(path string, ptr interface{}) error {
	format := "yaml"
	if strings.HasSuffix(path, "json") {
		format = "json"
	}

	var data []byte
	var err error

	switch format {
	case "yaml":
		data, err = yaml.Marshal(ptr)
	case "json":
		data, err = json.Marshal(ptr)
	default:
		return fmt.Errorf("invalid format '%s' provided", format)
	}

	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// ReadSpecObjectsFromFile returns the spec objects present in the file
func ReadSpecObjectsFromFile(fileName string) ([]*model.SpecObject, error) {
	var specs []*model.SpecObject

	var data []byte
	var err error
	if strings.HasPrefix(fileName, "http") {
		valuesFileObj, err := ExtractValuesObj("", fileName)
		if err != nil {
			return nil, err
		}
		data, err = yaml.Marshal(valuesFileObj)
		if err != nil {
			return nil, err
		}
	} else {
		// Read the file first
		data, err = file.File.ReadFile(fileName)
		if err != nil {
			return nil, err
		}
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

// ExtractValuesObj extract chart values from yaml file & cli flags
func ExtractValuesObj(setValuesFlag, valuesYamlFile string) (map[string]interface{}, error) {
	valuesFileObj := map[string]interface{}{}
	if valuesYamlFile != "" {
		var bodyInBytes []byte
		var err error
		if strings.HasPrefix(valuesYamlFile, "http") {
			// download file from the internet
			resp, err := http.Get(valuesYamlFile)
			if err != nil {
				return nil, err
			}
			if resp.StatusCode != http.StatusOK {
				return nil, fmt.Errorf("")
			}
			bodyInBytes, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
		} else {
			// read locally available file
			bodyInBytes, err = ioutil.ReadFile(valuesYamlFile)
			if err != nil {
				return nil, err
			}
		}

		if err := yaml.Unmarshal(bodyInBytes, &valuesFileObj); err != nil {
			return nil, err
		}
	}

	setValuesObj := map[string]interface{}{}
	if setValuesFlag != "" {
		arr := strings.Split(setValuesFlag, ",")
		for _, element := range arr {
			tempArr := strings.Split(element, "=")
			if len(tempArr) != 2 {
				return nil, fmt.Errorf("invalid value (%s) provided for flag --set, it should be in format foo1=bar1,foo2=bar2", tempArr)
			}
			setValuesObj[tempArr[0]] = tempArr[1]
		}
	}

	// override values of yaml file
	for key, value := range setValuesObj {
		valuesFileObj[key] = value
	}

	return valuesFileObj, nil
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
