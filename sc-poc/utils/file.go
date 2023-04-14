package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/ghodss/yaml"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// ReadSpecObjectsFromFile returns the spec objects present in the file
func ReadSpecObjectsFromFile(fileName string) ([]*unstructured.Unstructured, error) {
	var specs []*unstructured.Unstructured

	var data []byte
	var err error

	// Read the file first
	data, err = os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, nil
	}

	// Split the files into independent objects
	dataStrings := makeSpecStringArray(string(data))
	for _, dataString := range dataStrings {

		// Skip if string is too small to be a spec object
		if len(dataString) <= 5 {
			continue
		}

		// Unmarshal spec object
		spec := unstructured.Unstructured{}
		if err := yaml.Unmarshal([]byte(dataString), &spec.Object); err != nil {
			return nil, err
		}

		// Append the spec object into the array
		specs = append(specs, &spec)
	}

	return specs, nil
}

// GetBytesFromSpec converts spec of type unstructured.Unstructured into array of bytes
func GetBytesFromSpec(spec *unstructured.Unstructured) ([]byte, error) {
	data, err := yaml.Marshal(spec.Object)
	if err != nil {
		return nil, err
	}

	tempString := string(data) + "---" + "\n"
	return []byte(tempString), nil
}

// AppendToFile appends content to the file
func AppendToFile(filePath string, data []byte) error {
	// Open the file for appending
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	// Write the text to the file
	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("could not write to file: %v", err)
	}

	return nil
}

func makeSpecStringArray(raw string) []string {
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
