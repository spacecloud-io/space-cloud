package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/spaceuptech/space-cli/model"
)

func getHomeDirectory() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func getProjects(projects []*model.Projects) ([]string, error) {
	var projectNames []string
	if len(projects) == 0 {
		logrus.Error("error getting projects no projects founds, create new project from mission control")
		return nil, fmt.Errorf("projects array empty")
	}
	for _, project := range projects {
		projectNames = append(projectNames, project.Name)
	}
	return projectNames, nil
}

func generateYamlFile(credential *model.Credential) error {
	d, err := yaml.Marshal(&credential)
	if err != nil {
		return err
	}

	if err := createDirIfNotExist(getSpaceCloudDirectory()); err != nil {
		logrus.Errorf("error in generating yaml file unable to create space cli directory - %v", err)
		return err
	}

	fileName := getAccountConfigPath()
	err = ioutil.WriteFile(fileName, d, 0644)
	if err != nil {
		return err
	}

	return nil
}

func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func createFileIfNotExist(path, content string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return ioutil.WriteFile(path, []byte(content), 0755)
	}
	return nil
}

// CloseTheCloser closes the closer
func CloseTheCloser(c io.Closer) {
	_ = c.Close()
}
