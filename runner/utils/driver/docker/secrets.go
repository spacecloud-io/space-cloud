package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/runner/model"
)

// CreateSecret create a file for every secret & update the secret if already exists and has same type
func (d *Docker) CreateSecret(ctx context.Context, projectID string, secretObj *model.Secret) error {
	// create folder for project
	projectPath := fmt.Sprintf("%s/%s", d.secretPath, projectID)

	// Try to create a secrets directory if it doesn't already exist
	if err := d.createDir(projectPath); err != nil {
		return err
	}

	// check if file exists
	filePath := fmt.Sprintf("%s/%s.json", projectPath, secretObj.ID)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// create file & set it's content
		return d.writeIntoFile(ctx, secretObj, filePath)
	} else if err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("error creating secret in docker unable to check if file exists (%s)", projectPath), err, nil)
	}

	// file already exists read it's content
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("error creating secret in docker unable to read file (%s)", filePath), err, nil)
	}
	fileContent := new(model.Secret)
	if err := json.Unmarshal(data, fileContent); err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "error creating secret in docker unable to unmarshal data", err, nil)
	}
	if fileContent.Type != secretObj.Type {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "file already exists but secrets to set have different types", nil, nil)
	}
	// update existing file
	return d.writeIntoFile(ctx, secretObj, filePath)
}

// ListSecrets list all the secrets of the project
func (d *Docker) ListSecrets(ctx context.Context, projectID string) ([]*model.Secret, error) {
	projectPath := fmt.Sprintf("%s/%s", d.secretPath, projectID)

	// Try to create a secrets directory if it doesn't already exist
	if err := d.createDir(projectPath); err != nil {
		return nil, err
	}

	// Read all the files
	files, err := ioutil.ReadDir(projectPath)
	if err != nil {
		return nil, err
	}

	secretArr := make([]*model.Secret, len(files))
	for index, file := range files {
		if !file.IsDir() {
			filePath := fmt.Sprintf("%s/%s", projectPath, file.Name())
			data, err := ioutil.ReadFile(filePath)
			if err != nil {
				return nil, err
			}
			fileContent := new(model.Secret)
			if err := json.Unmarshal(data, fileContent); err != nil {
				return nil, err
			}

			// remove all value of secrets
			secrets := map[string]string{}
			for key, data := range fileContent.Data {
				secrets[key] = data
			}
			secretArr[index] = &model.Secret{
				ID:       fileContent.ID,
				Type:     fileContent.Type,
				RootPath: fileContent.RootPath,
				Data:     secrets,
			}
		}
	}
	return secretArr, nil
}

// DeleteSecret deletes the docker secret
func (d *Docker) DeleteSecret(ctx context.Context, projectID, secretName string) error {
	return os.RemoveAll(fmt.Sprintf("%s/%s/%s.json", d.secretPath, projectID, secretName))
}

// SetFileSecretRootPath sets the file secret root path
func (d *Docker) SetFileSecretRootPath(ctx context.Context, projectID string, secretName, rootPath string) error {
	if secretName == "" || rootPath == "" {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("key/value not provided; got (%s,%s)", secretName, rootPath), nil, nil)
	}
	// check if file exists
	filePath := fmt.Sprintf("%s/%s/%s.json", d.secretPath, projectID, secretName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "file doesn't exists", err, nil)
	}

	// file already exists read it's content
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	fileContent := new(model.Secret)
	if err := json.Unmarshal(data, fileContent); err != nil {
		return err
	}

	fileContent.RootPath = rootPath
	return d.writeIntoFile(ctx, fileContent, filePath)
}

// SetKey sets the secret key for docker file
func (d *Docker) SetKey(ctx context.Context, projectID, secretName, secretKey string, secretObj *model.SecretValue) error {
	// check if file exists
	filePath := fmt.Sprintf("%s/%s/%s.json", d.secretPath, projectID, secretName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "file doesn't exists", err, nil)
	}

	// file already exists read it's content
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	fileContent := new(model.Secret)
	if err := json.Unmarshal(data, fileContent); err != nil {
		return err
	}

	fileContent.Data[secretKey] = secretObj.Value
	return d.writeIntoFile(ctx, fileContent, filePath)
}

// DeleteKey deletes the secret key for docker file
func (d *Docker) DeleteKey(ctx context.Context, projectID, secretName, secretKey string) error {
	// check if file exists
	filePath := fmt.Sprintf("%s/%s/%s.json", d.secretPath, projectID, secretName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "file doesn't exists", err, nil)
	}

	// file already exists read it's content
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	fileContent := new(model.Secret)
	if err := json.Unmarshal(data, fileContent); err != nil {
		return err
	}

	delete(fileContent.Data, secretKey)
	return d.writeIntoFile(ctx, fileContent, filePath)
}

// writeIntoFile writes json data into specified file
func (d *Docker) writeIntoFile(ctx context.Context, secretObj *model.Secret, filePath string) error {
	data, err := json.Marshal(secretObj)
	if err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("error writing data in file (%s) unable to marshal data", filePath), err, nil)
	}
	// create / update file content
	return ioutil.WriteFile(filePath, data, 0777)
}

func (d *Docker) createDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0777)
	}
	return nil
}
