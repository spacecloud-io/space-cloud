package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
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

// createFileIfNotExist creates a file with the provided content if it doesn't already exists
func createFileIfNotExist(path, content string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return ioutil.WriteFile(path, []byte(content), 0755)
	}
	return nil
}

func getSpaceCLIConfigPath() string {
	return fmt.Sprintf("%s/config.json", getSpaceCLIDirectory())
}

// getLatestVersion retrieves the latest Space Cloud version based on the current version
func getLatestVersion() (*cliVersionDoc, error) {
	// Create a db object
	db := api.New("spacecloud", "localhost:4122", false).DB("db")

	// Create a context
	ctx := context.Background()

	var result *spaceApiTypes.Response
	var err error
	result, err = db.Get("cli_version").Sort("-version_cod").Limit(1).Apply(ctx)
	if err != nil {
		return nil, err
	}
	r := new(cliVersionResponse)
	if err := result.Unmarshal(&r); err != nil {
		return nil, err
	}
	doc := new(cliVersionDoc)
	for _, val := range r.Docs {
		if val.VersionCode > doc.VersionCode {
			doc.VersionNo = val.VersionNo
			doc.VersionCode = val.VersionCode
			doc.ID = val.ID
		}
	}
	return doc, nil
}

func readVersionConfig() (*cliVersionDoc, error) {
	file, err := ioutil.ReadFile(fmt.Sprintf("%s/config.json", getSpaceCLIDirectory()))
	if err != nil {
		return nil, err
	}

	doc := new(cliVersionDoc)
	err = json.Unmarshal([]byte(file), doc)
	return doc, err
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

func logError(message string, err error) error {
	// Log with error if provided
	if err != nil {
		logrus.WithField("error", err.Error()).Errorln(message)
	} else {
		logrus.Errorln(message)
	}

	// Return the error message
	return errors.New(message)
}
