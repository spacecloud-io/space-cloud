package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func ListAllSources(client *http.Client, baseUrl string) ([]schema.GroupVersionResource, error) {
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

func GetResources(client *http.Client, gvr schema.GroupVersionResource, baseUrl string, pkgName string) (*unstructured.UnstructuredList, error) {
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

func ApplyResources(client *http.Client, gvr schema.GroupVersionResource, baseUrl string, spec *unstructured.Unstructured) error {
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

func DeleteResources(client *http.Client, gvr schema.GroupVersionResource, baseUrl string, pkgName string) error {
	path := fmt.Sprintf("%s/sc/v1/config/%s/%s/%s/%s", baseUrl, gvr.Group, gvr.Version, gvr.Resource, pkgName)

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

func ListAllPlugins(client *http.Client, baseUrl string, token string) ([]v1alpha1.HTTPPlugin, error) {
	var allPlugins []v1alpha1.HTTPPlugin

	path := baseUrl + "/sc/v1/plugins"
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch plugins")
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch plugins")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch plugins")
	}

	json.Unmarshal(body, &allPlugins)
	return allPlugins, nil
}
