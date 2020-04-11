package main

// type cliVersionDoc struct {
// 	VersionNo   string `mapstructure:"version_no" json:"versionNo"`
// 	VersionCode int32  `mapstructure:"version_code" json:"versionCode"`
// 	ID          string `mapstructure:"id" json:"id"`
// }
//
// func getModule() (*cobra.Command, error) {
//
// 	_ = createDirIfNotExist(getSpaceCloudDirectory())
// 	_ = createDirIfNotExist(getSpaceCLIDirectory())
// 	_ = createFileIfNotExist(getSpaceCLIConfigPath(), "{}")
//
// 	currentVersion, err1 := readVersionConfig()
// 	latestVersion, err2 := getLatestVersion()
//
// 	// Return error if we could not get the current or latest version
// 	if err1 != nil && err2 != nil {
// 		return nil, logError("Could not fetch space-cli plugin", err2)
// 	}
// 	// Return currentVersion if we could not get the latest version
// 	if err1 == nil && err2 != nil {
// 		return getplugin(currentVersion.VersionNo)
// 	}
// 	// Returns latestVersion if latestVersion is availabe or currentVersion not found
// 	if err2 == nil {
// 		if err1 == nil {
// 			if latestVersion.VersionCode <= currentVersion.VersionCode {
// 				return getplugin(currentVersion.VersionNo)
// 			}
// 		}
// 		// Download the latest plugin version
// 		if err := downloadPlugin(latestVersion); err != nil {
// 			return nil, err
// 		}
// 	}
// 	return getplugin(latestVersion.VersionNo)
// }
