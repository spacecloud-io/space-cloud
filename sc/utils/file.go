package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ghodss/yaml"
)

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
