package pkg

import (
	"log"
	"os"

	"github.com/ghodss/yaml"
)

type Config struct {
	Name        string            `json:"name"`
	Labels      map[string]string `json:"labels,omitempty"`
	Output      Output            `json:"output"`
	ResourceDir string            `json:"sourceDir"`
}

type Output struct {
	Language  string `json:"language"`
	OutputDir string `json:"outputDir"`
}

func CreateConfig(config Config) {
	yamlBytes, err := yaml.Marshal(config)
	if err != nil {
		log.Fatalf("Failed to marshal YAML: %v", err)
	}

	err = os.WriteFile("sc/package.yaml", yamlBytes, 0777)
	if err != nil {
		log.Fatalf("Failed to write YAML file: %v", err)
	}
}

func ReadConfig() Config {
	yamlFile, err := os.ReadFile("sc/package.yaml")
	if err != nil {
		log.Fatalf("Failed to read YAML file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Failed to unmarshal YAML: %v", err)
	}
	return config
}
