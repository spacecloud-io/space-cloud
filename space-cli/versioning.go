package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
)

type cliVersionResponse struct {
	Docs []*cliVersionDoc `mapstructure:"result"`
}

type cliVersionDoc struct {
	VersionNo   string `mapstructure:"version_no" json:"versionNo"`
	VersionCode int32  `mapstructure:"version_code" json:"versionCode"`
	ID          string `mapstructure:"id" json:"id"`
}

func getmodule() (*cobra.Command, error) {

	_ = createDirIfNotExist(getSpaceCloudDirectory())
	_ = createDirIfNotExist(getSpaceCLIDirectory())
	_ = createFileIfNotExist(getSpaceCLIConfigPath(), "{}")

	currentVersion, err1 := readVersionConfig()
	latestVersion, err2 := getLatestVersion()

	//Return error if we could not get the current or latest version
	if err1 != nil && err2 != nil {
		return nil, logError("Could not fetch space-cli plugin", err2)
	}
	// Return currentVersion if we could not get the latest version
	if err1 == nil && err2 != nil {
		return getplugin(currentVersion.VersionNo)
	}

	if err2 == nil {
		if err1 == nil {
			if latestVersion.VersionCode <= currentVersion.VersionCode {
				return getplugin(currentVersion.VersionNo)
			}
		}
		url := fmt.Sprintf("http://localhost:5000/cmd_%s.so", latestVersion.VersionNo)
		filepath := fmt.Sprintf("%s/cmd_%s.so", getSpaceCLIDirectory(), latestVersion.VersionNo)
		err := downloadFile(url, filepath)
		if err != nil {
			return nil, err
		}
		docs := &cliVersionDoc{
			VersionNo:   latestVersion.VersionNo,
			VersionCode: latestVersion.VersionCode,
			ID:          latestVersion.ID,
		}
		file, _ := json.Marshal(docs)
		_ = ioutil.WriteFile(fmt.Sprintf("%s/config.json", getSpaceCLIDirectory()), file, 0644)
	}
	return getplugin(latestVersion.VersionNo)
}
