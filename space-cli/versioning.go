package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
)

type scVersionResponse struct {
	Docs []*scVersionDoc `mapstructure:"result"`
}

type scVersionDoc struct {
	VersionNo   string `mapstructure:"version_no" json:"versionNo"`
	VersionCode int32  `mapstructure:"version_code" json:"versionCode"`
	//Id          TypeID     `mapstructure:"id" json:"id"`
}

//const TypeID string ="ID"

func getmodule() (*cobra.Command, error) {

	_ = createDirIfNotExist(getSpaceCloudDirectory())
	_ = createDirIfNotExist(getSpaceCliDirectory())
	_ = createFileIfNotExist(getSpaceCliConfigPath(), "{}")

	//currentVersionCode := int32(0)
	file, _ := ioutil.ReadFile("config.json")

	var data scVersionResponse

	_ = json.Unmarshal([]byte(file), &data)
	//currentVersion := data.VersionNo
	currentVersionCode := data.Docs[0].VersionCode

	// result := make(map[string]interface{})
	// if err := utils.Get(http.MethodGet, "/v1/config/env", map[string]string{}, &result); err != nil {
	// 	currentVersion := ""
	// 	//return utils.LogError("Unable to get current Space Cloud version. Is Space Cloud running?", err)
	// } else {
	// 	currentVersion := result["version"].(string)
	// }

	var rootCmd = &cobra.Command{}

	latestVersion, latestVersionCode, err := getLatestVersion()
	if err != nil {
		return rootCmd, err
	}

	if latestVersionCode > currentVersionCode {
		url := fmt.Sprintf("http://localhost:5000/Download/cli/%s", latestVersion)
		filepath := fmt.Sprintf("%s/cmd_%s.so", getSpaceCliDirectory(), latestVersion)
		_ = downloadFile(url, filepath)
		data := &scVersionDoc{
			VersionNo:   latestVersion,
			VersionCode: latestVersionCode,
		}
		file, _ := json.Marshal(data)
		_ = ioutil.WriteFile("config.json", file, 0644)
		return getplugin()
	}

	return rootCmd, nil
}
