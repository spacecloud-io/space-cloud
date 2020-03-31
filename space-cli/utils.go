package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"

	api "github.com/spaceuptech/space-api-go"
	spaceApiTypes "github.com/spaceuptech/space-api-go/types"
)

func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func getSpaceCloudDirectory() string {
	return fmt.Sprintf("%s/.space-cloud", getHomeDirectory())
}

func getHomeDirectory() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func getSpaceCLIDirectory() string {
	return fmt.Sprintf("%s/cli", getSpaceCloudDirectory())
}

// CreateFileIfNotExist creates a file with the provided content if it doesn't already exists
func createFileIfNotExist(path, content string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return ioutil.WriteFile(path, []byte(content), 0755)
	}
	return nil
}

func getSpaceCLIConfigPath() string {
	return fmt.Sprintf("%s/config.json", getSpaceCLIDirectory())
}

// GetLatestVersion retrieves the latest Space Cloud version based on the current version
func getLatestVersion() (string, int32, error) {
	// Create a db object
	db := api.New("spacecloud", "localhost:4122", false).DB("db")

	// Create a context
	ctx := context.Background()

	var result *spaceApiTypes.Response
	var err error
	result, err = db.Get("cli_version").Sort("-version_code").Limit(1).Apply(ctx)
	if err != nil {
		return "", 0, err
	}

	r := cliVersionResponse{}
	if err := result.Unmarshal(&r); err != nil {
		return "", 0, err
	}
	newVersion, newVersionCode := "", int32(0)
	for _, val := range r.Docs {
		if val.VersionCode > newVersionCode {
			newVersion = val.VersionNo
			newVersionCode = val.VersionCode
		}
	}
	return newVersion, newVersionCode, nil
}

func readVersionConfig() (cliVersionResponse, error) {
	var data cliVersionResponse
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		return data, err
	}
	err = json.Unmarshal([]byte(file), &data)
	return data, err

}

func downloadFile(url string, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
