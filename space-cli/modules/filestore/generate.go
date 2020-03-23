package filestore

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spaceuptech/space-cli/model"
)

func generateFilestoreRule() (*model.SpecObject, error) {
	projectID := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter Project ID"}, &projectID); err != nil {
		return nil, err
	}
	ID := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter Collection Name"}, &ID); err != nil {
		return nil, err
	}
	prefix := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter DB Alias"}, &prefix); err != nil {
		return nil, err
	}

	v := &model.SpecObject{
		API:  "/v1/config/projects/{project}/file-storage/rules/{id}",
		Type: "filestore-rule",
		Meta: map[string]string{
			"id":      ID,
			"project": projectID,
		},
		Spec: map[string]interface{}{
			"prefix": prefix,
			"rule": map[string]interface{}{
				"create": map[string]interface{}{
					"rule": "allow",
				},
				"delete": map[string]interface{}{
					"rule": "allow",
				},
				"read": map[string]interface{}{
					"rule": "allow",
				},
			},
		},
	}

	return v, nil
}

func generateFilestoreConfig() (*model.SpecObject, error) {
	projectID := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter Project ID"}, &projectID); err != nil {
		return nil, err
	}

	storeType := ""
	if err := survey.AskOne(&survey.Select{Message: "Enter Storetype", Options: []string{"Local", "AmazonS3", "GCPStorage"}}, &storeType); err != nil {
		return nil, err
	}
	bucket := ""
	endpoint := ""
	conn := ""
	switch storeType {
	case "Local":
		if err := survey.AskOne(&survey.Input{Message: "Enter connection"}, &conn); err != nil {
			return nil, err
		}
	case "AmazonS3":
		if err := survey.AskOne(&survey.Input{Message: "Enter connection"}, &conn); err != nil {
			return nil, err
		}
		if err := survey.AskOne(&survey.Input{Message: "Enter endpoint"}, &endpoint); err != nil {
			return nil, err
		}
		if err := survey.AskOne(&survey.Input{Message: "Enter bucket"}, &bucket); err != nil {
			return nil, err
		}
	case "GCPStorage":
		if err := survey.AskOne(&survey.Input{Message: "Enter bucket"}, &bucket); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Invalid choice")
	}

	v := &model.SpecObject{
		API:  "/v1/external/projects/{project}/file-storage/config/{id}",
		Type: "filestore-config",
		Meta: map[string]string{
			"project": projectID,
			"id":      "filestore-config",
		},
		Spec: map[string]interface{}{
			"bucket":    bucket,
			"conn":      conn,
			"enabled":   true,
			"endpoint":  endpoint,
			"storeType": storeType,
		},
	}

	return v, nil
}
