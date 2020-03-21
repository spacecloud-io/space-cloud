package utils

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cli/model"
)

// GetProjectsNamesFromArray returns the array of projects names
func GetProjectsNamesFromArray(projects []*model.Projects) ([]string, error) {
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

// GetProjectsFromSC returns the projects array from sc
func GetProjectsFromSC() ([]*model.Projects, error) {
	type response struct {
		Projects []*model.Projects `json:"projects"`
	}

	res := new(response)
	if err := Get(http.MethodGet, "/v1/config/projects", map[string]string{}, res); err != nil {
		return nil, err
	}

	return res.Projects, nil
}
