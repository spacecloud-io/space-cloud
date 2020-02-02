package docker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/runner/model"
)

// CreateSecret create a file for every secret & update the secret if already exists and has same type
func (d *docker) CreateSecret(projectID string, secretObj *model.Secret) error {
	// create folder for project
	projectPath := fmt.Sprintf("%s/%s", d.secretPath, projectID)
	if err := d.createDir(projectPath); err != nil {
		return nil
	}

	// check if file exists
	filePath := fmt.Sprintf("%s/%s.json", projectPath, secretObj.Name)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// create file & set it's content
		return d.writeIntoFile(secretObj, filePath)
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
	if fileContent.Type != secretObj.Type {
		return fmt.Errorf("file already exists but secrets to set have different types")
	}
	// update existing file
	return d.writeIntoFile(secretObj, filePath)
}

func (d *docker) ListSecrets(projectID string) ([]*model.Secret, error) {
	projectPath := fmt.Sprintf("%s/%s", d.secretPath, projectID)
	files, err := ioutil.ReadDir(projectPath)
	if err != nil {
		return nil, err
	}

	secretArr := make([]*model.Secret, len(files))
	for _, file := range files {
		// there will be only files in this directory
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
		for key := range fileContent.Data {
			secrets[key] = ""
		}
		secretArr = append(secretArr, &model.Secret{
			Name:     fileContent.Name,
			Type:     fileContent.Type,
			RootPath: fileContent.RootPath,
			Data:     secrets,
		})
	}
	return secretArr, nil
}

func (d *docker) DeleteSecret(projectID, secretName string) error {
	return os.RemoveAll(fmt.Sprintf("%s/%s/%s.json", d.secretPath, projectID, secretName))
}

func (d *docker) SetKey(projectID, secretName, secretKey string, secretObj *model.SecretValue) error {
	// check if file exists
	filePath := fmt.Sprintf("%s/%s/%s.json", d.secretPath, projectID, secretName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file doesn't exists")
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
	return d.writeIntoFile(fileContent, filePath)
}

func (d *docker) DeleteKey(projectID, secretName, secretKey string) error {
	// check if file exists
	filePath := fmt.Sprintf("%s/%s/%s.json", d.secretPath, projectID, secretName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file doesn't exists")
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
	return d.writeIntoFile(fileContent, filePath)
}

// writeIntoFile writes json data into specified file
func (d *docker) writeIntoFile(secretObj *model.Secret, filePath string) error {
	data, err := json.Marshal(secretObj)
	if err != nil {
		logrus.Error("error creating secret in docker unable to marshal data")
		return err
	}
	// create / update file content
	return ioutil.WriteFile(filePath, data, 0644)
}

func (d *docker) createDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0644)
	}
	return nil
}
