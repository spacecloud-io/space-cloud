package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils/input"
)

// GetProjectsNamesFromArray returns the array of projects names
func GetProjectsNamesFromArray(projects []*model.Projects) ([]string, error) {
	var projectNames []string
	if len(projects) == 0 {
		_ = LogError("error getting projects no projects founds, create new project from mission control", nil)
		return nil, fmt.Errorf("projects array empty")
	}
	for _, project := range projects {
		projectNames = append(projectNames, project.Name)
	}
	return projectNames, nil
}

// GetProjectsFromSC returns the projects array from sc
func GetProjectsFromSC() ([]*model.Projects, error) {
	type response struct {
		Result []*model.Projects `json:"result"`
	}

	res := new(response)
	if err := Get(http.MethodGet, "/v1/config/projects/*", map[string]string{}, res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

// GetProjectID checks if project is specified in flags
func GetProjectID() (string, bool) {
	projectID := viper.GetString("project")
	if projectID == "" {
		creds, err := getSelectedAccount()
		if err != nil || creds.DefaultProject == "" {
			return "", false
		}
		return creds.DefaultProject, true
	}
	return projectID, true
}

// SetDefaultProject sets the default project for the selected account in accounts.yaml file
// if empty value provided it will try to project from server & prompts user to select project from server response
func SetDefaultProject(project string) error {
	// we are disabling logs because getting projects from sc may print unnecessary logs
	logrus.SetOutput(ioutil.Discard)
	objs, _ := GetProjectsFromSC()
	logrus.SetOutput(os.Stdout)

	if project == "" {
		if len(objs) > 0 {
			projects := make([]string, 0)
			for _, obj := range objs {
				projects = append(projects, obj.ID)
			}
			if err := input.Survey.AskOne(&survey.Select{Message: "Select default project for this account", Options: projects}, &project); err != nil {
				return err
			}
		}
	} else {
		// validate the project if it really exists
		if len(objs) > 0 {
			for _, obj := range objs {
				if obj.ID == project {
					break
				}
			}
			return LogError(fmt.Sprintf("Provided project (%s) does not exits in space cloud", project), nil)
		}
	}

	acc, err := getSelectedAccount()
	if err != nil {
		return err
	}
	acc.DefaultProject = project
	return StoreCredentials(acc)
}
