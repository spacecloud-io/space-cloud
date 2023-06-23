package pkg

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func listAllSources(client *http.Client, baseUrl string) ([]schema.GroupVersionResource, error) {
	var sourcesGVR []schema.GroupVersionResource

	path := baseUrl + "/sc/v1/sources"
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sources")
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sources")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sources")
	}

	json.Unmarshal(body, &sourcesGVR)
	return sourcesGVR, nil
}

func getResources(client *http.Client, gvr schema.GroupVersionResource, baseUrl string, pkgName string) (*unstructured.UnstructuredList, error) {
	var unstr *unstructured.UnstructuredList

	path := fmt.Sprintf("%s/sc/v1/config/%s/%s/%s/?package=%s", baseUrl, gvr.Group, gvr.Version, gvr.Resource, pkgName)
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch resources from %s", path)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch resources from %s", path)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch resources from %s", path)
	}

	json.Unmarshal(body, &unstr)
	return unstr, nil
}

func applyResources(client *http.Client, gvr schema.GroupVersionResource, baseUrl string, spec *unstructured.Unstructured) error {
	path := fmt.Sprintf("%s/sc/v1/config/%s/%s/%s/", baseUrl, gvr.Group, gvr.Version, gvr.Resource)

	b, err := json.Marshal(spec)
	if err != nil {
		return fmt.Errorf("failed to apply resource to %s", path)
	}
	req, err := http.NewRequest("PUT", path, bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("failed to apply resource to %s", path)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to apply resource to %s", path)
	}
	defer resp.Body.Close()

	return nil
}

func deleteResources(client *http.Client, path string) error {
	req, err := http.NewRequest("DELETE", path, nil)
	if err != nil {
		return fmt.Errorf("failed to delete resources from %s", path)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete resources from %s", path)
	}
	defer resp.Body.Close()

	return nil
}

func findElement(arr []string, target string) int {
	for i, str := range arr {
		if str == target {
			return i
		}
	}
	return -1
}

func DeleteElement(arr []string, index int) []string {
	// Check if the index is out of range
	if index < 0 || index >= len(arr) {
		return arr
	}

	// Create a new slice with the element at the specified index removed
	return append(arr[:index], arr[index+1:]...)
}

func getCredentials() (map[string]string, error) {
	homeDir, _ := os.UserHomeDir()
	file := filepath.Join(homeDir, ".space-cloud", "creds.json")
	m := make(map[string]string)
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("empty file found")
	}

	err = json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func login(client *http.Client, creds map[string]string) error {
	b, _ := base64.StdEncoding.DecodeString(creds["username"])
	creds["username"] = string(b)
	b, _ = base64.StdEncoding.DecodeString(creds["password"])
	creds["password"] = string(b)

	path := fmt.Sprintf("%s/sc/v1/login", creds["url"])
	payload, _ := json.Marshal(creds)

	req, err := http.NewRequest(http.MethodPost, path, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf(string(body))
	}

	return nil
}
