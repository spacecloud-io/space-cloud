package docker

import (
	"context"
	"fmt"
	"os"

	"github.com/spaceuptech/space-cloud/runner/model"

	"github.com/sirupsen/logrus"
)

// DeleteProject deletes the project
func (d *Docker) DeleteProject(ctx context.Context, projectID string) error {
	if err := d.DeleteService(ctx, projectID, "", ""); err != nil {
		logrus.Errorf("error deleting project in docker unable to delete services with project id (%s) - %s", projectID, err.Error())
		return err
	}
	if err := os.RemoveAll(fmt.Sprintf("%s/%s", d.secretPath, projectID)); err != nil {
		logrus.Errorf("error deleting project in docker unable to delete secrets folder at (%s) - %s", d.secretPath, err.Error())
		return err
	}
	return nil
}

// CreateProject creates a new project
func (d *Docker) CreateProject(ctx context.Context, project *model.Project) error {
	projectPath := fmt.Sprintf("%s/%s", d.secretPath, project.ID)
	if err := d.createDir(projectPath); err != nil {
		logrus.Errorf("error creating secret in docker unable to create directory (%s) - %s", projectPath, err.Error())
		return err
	}
	return nil
}
