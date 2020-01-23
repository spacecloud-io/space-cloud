package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
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

func getServiceConfig(filename string) (*model.Service, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	fileName := fmt.Sprintf("/%s/%s", dir, filename)
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	c := new(model.Service)
	err = yaml.Unmarshal(file, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func checkFile(filename string) bool {
	dir, err := os.Getwd()
	if err != nil {
		return false
	}
	fileName := fmt.Sprintf("/%s/%s", dir, filename)
	_, err = ioutil.ReadFile(fileName)
	if err != nil {
		return false
	}
	return true
}

func getProgLang() (string, error) {
	if checkFile("requirement.txt") {
		return "python", nil
	}
	if checkFile("package.json") {
		return "node", nil
	}
	if checkExt(".go") {
		return "go", nil
	}
	return "", fmt.Errorf("No Programming language detected")
}

func getImage(progLang string) (string, error) {
	switch progLang {
	case "python":
		return "spaceuptech/runtime-python", nil
	case "nodejs":
		return "spaceuptech/runtime-node", nil
	case "go":
		return "spaceuptech/runtime-alpine", nil
	default:
		return "", fmt.Errorf("%s is not supported", progLang)
	}
}

func getCmd(progLang string) []string {
	switch progLang {
	case "python":
		return []string{"python3", "app.py"}
	case "nodejs":
		return []string{"npm", "start"}
	case "go":
		return []string{"./app"}
	default:
		return []string{}
	}
}

func getProject(projectID string, projects []*model.Projects) (*model.Projects, error) {
	for _, project := range projects {
		if projectID == project.ID {
			return project, nil
		}
	}
	return nil, fmt.Errorf("Invalid Project Name")
}

func getEnvironment(envID string, environments []*model.Environment) (*model.Environment, error) {
	for _, env := range environments {
		if envID == env.ID {
			return env, nil
		}
	}
	return nil, errors.New("Invalid Project Name")
}

func getProjects(projects []*model.Projects) ([]string, error) {
	var projnames []string
	if len(projects) == 0 {
		return nil, fmt.Errorf("error getting projects no projects founds, create new project from missioin control")
	}
	for _, val := range projects {
		projnames = append(projnames, fmt.Sprintf("%s (%s)", val.ID, val.Name))
	}
	return projnames, nil
}

func getEnvironments(project *model.Projects) []string {
	var envs []string
	for _, val := range project.Environments {
		envs = append(envs, fmt.Sprintf("%s %s", val.ID, val.Name))
	}
	return envs
}

func getClusters(environment *model.Environment) []string {
	var clusternames []string
	for _, val := range environment.Clusters {
		clusternames = append(clusternames, val.ID)
	}
	return clusternames
}

func checkExt(ext string) bool {
	pathS, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	present := false
	filepath.Walk(pathS, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(ext, f.Name())
			if err == nil && r {
				present = true
			}
		}
		return nil
	})
	return present
}

func generateYamlFile(credential *model.Credential) error {
	d, err := yaml.Marshal(&credential)
	if err != nil {
		return err
	}

	if err := createDirIfNotExist(getSpaceCliDirectory()); err != nil {
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
