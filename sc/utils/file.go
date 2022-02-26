package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ghodss/yaml"
)

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
