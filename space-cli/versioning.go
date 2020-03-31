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
	//Id          TypeID     `mapstructure:"id" json:"id"`
}

//const TypeID string ="ID"

func getmodule() (*cobra.Command, error) {

	_ = createDirIfNotExist(getSpaceCloudDirectory())
	_ = createDirIfNotExist(getSpaceCLIDirectory())
	_ = createFileIfNotExist(getSpaceCLIConfigPath(), "{}")

	data, err1 := readVersionConfig()
	currentVersion := data.Docs[0].VersionNo
	currentVersionCode := data.Docs[0].VersionCode

	latestVersion, latestVersionCode, err2 := getLatestVersion()

	if err1 != nil && err2 != nil {
		fmt.Println("There is an error please try again")
		return nil, err1
	} else if err1 == nil && err2 != nil {
		return getplugin(currentVersion)
	} else if err1 != nil && err2 == nil {
		url := fmt.Sprintf("http://localhost:5000/Download/cli/%s", latestVersion)
		filepath := fmt.Sprintf("%s/cmd_%s.so", getSpaceCLIDirectory(), latestVersion)
		_ = downloadFile(url, filepath)
		data := &cliVersionDoc{
			VersionNo:   latestVersion,
			VersionCode: latestVersionCode,
		}
		file, _ := json.Marshal(data)
		_ = ioutil.WriteFile("config.json", file, 0644)
		return getplugin(latestVersion)
	}

	if latestVersionCode > currentVersionCode {
		url := fmt.Sprintf("http://localhost:5000/Download/cli/%s", latestVersion)
		filepath := fmt.Sprintf("%s/cmd_%s.so", getSpaceCLIDirectory(), latestVersion)
		_ = downloadFile(url, filepath)
		data := &cliVersionDoc{
			VersionNo:   latestVersion,
			VersionCode: latestVersionCode,
		}
		file, _ := json.Marshal(data)
		_ = ioutil.WriteFile("config.json", file, 0644)
		return getplugin(latestVersion)
	}

	return getplugin(latestVersion)
}
