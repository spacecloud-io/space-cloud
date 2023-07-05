package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func GetCredentials() (Credentials, error) {
	var creds Credentials
	homeDir, _ := os.UserHomeDir()
	file := filepath.Join(homeDir, ".space-cloud", "creds.json")
	data, err := os.ReadFile(file)
	if err != nil {
		return creds, err
	}

	if len(data) == 0 {
		return creds, fmt.Errorf("empty file found")
	}

	err = json.Unmarshal(data, &creds)
	if err != nil {
		return creds, err
	}
	fmt.Println(creds)
	return creds, nil
}

func Login(client *http.Client, creds Credentials) error {
	b, _ := base64.StdEncoding.DecodeString(creds.Username)
	creds.Username = string(b)
	b, _ = base64.StdEncoding.DecodeString(creds.Password)
	creds.Password = string(b)

	path := fmt.Sprintf("%s/sc/v1/login", creds.BaseUrl)
	payload, _ := json.Marshal(creds)

	req, err := http.NewRequest(http.MethodPost, path, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf(string(body))
	}

	return nil
}

func UpdateSpaceCloudCredsFile(creds Credentials) (string, error) {
	homeDir, _ := os.UserHomeDir()
	dirPath := filepath.Join(homeDir, ".space-cloud")
	_ = os.Mkdir(dirPath, 0777)

	output := map[string]string{
		"username": creds.Username,
		"password": creds.Password,
		"baseUrl":  creds.BaseUrl,
	}

	b, err := json.Marshal(output)
	if err != nil {
		return "", err
	}

	location := filepath.Join(dirPath, "creds.json")
	_ = os.WriteFile(location, b, 0777)
	return location, nil
}
