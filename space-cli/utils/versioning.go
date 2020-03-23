package utils

import (
	"context"

	api "github.com/spaceuptech/space-api-go"
	spaceApiTypes "github.com/spaceuptech/space-api-go/types"
)

type scVersionResponse struct {
	Docs []*scVersionDoc `mapstructure:"result"`
}

type scVersionDoc struct {
	VersionNo         string `mapstructure:"version_no" json:"versionNo"`
	VersionCode       int32  `mapstructure:"version_code" json:"versionCode"`
	CompatibleVersion string `mapstructure:"compatible_version" json:"compatibleVersion"`
}

// GetLatestVersion retrieves the latest Space Cloud version based on the current version
func GetLatestVersion(version string) (string, error) {
	// Create a db object
	db := api.New("spacecloud", "localhost:4122", false).DB("db")

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

	r := scVersionResponse{}
	if err := result.Unmarshal(&r); err != nil {
		return "", err
	}
	newVersion, newVersionCode := "", int32(0)
	for _, val := range r.Docs {
		if val.VersionCode > newVersionCode {
			newVersion = val.VersionNo
			newVersionCode = val.VersionCode
		}
	}
	return newVersion, nil
}
