package utils

import (
	"context"

	api "github.com/spaceuptech/space-api-go"
	spaceApiTypes "github.com/spaceuptech/space-api-go/types"
)

type scVersionDoc struct {
	VersionNo         string `mapstructure:"version_no" json:"versionNo"`
	VersionCode       int32  `mapstructure:"version_code" json:"versionCode"`
	CompatibleVersion string `mapstructure:"compatible_version" json:"compatibleVersion"`
}

// GetLatestVersion retrieves the latest Space Cloud version based on the current version
func GetLatestVersion(version string) (string, error) {
	// Create a db object
	db := api.New("spacecloud", "api.spaceuptech.com", true).DB("db")

	// Create a context
	ctx := context.Background()

	var result *spaceApiTypes.Response
	var err error
	if version == "" {
		result, err = db.Get("sc_version").Sort("-version_code").Limit(1).Apply(ctx)
		if err != nil {
			return "", err
		}
	} else {
		result, err = db.Get("sc_version").Where(spaceApiTypes.Cond("compatible_version", "==", version)).Apply(ctx)
		if err != nil {
			return "", err
		}
	}

	docs := []*scVersionDoc{}
	if err := result.Unmarshal(&docs); err != nil {
		return "", err
	}
	newVersion, newVersionCode := "", int32(0)
	for _, val := range docs {
		if val.VersionCode > newVersionCode {
			newVersion = val.VersionNo
			newVersionCode = val.VersionCode
		}
	}
	return newVersion, nil
}
