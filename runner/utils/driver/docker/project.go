package docker

import (
	"context"
	"github.com/spaceuptech/space-cloud/runner/model"

	"github.com/sirupsen/logrus"
)

// DeleteProject deletes the project
func (d *Docker) DeleteProject(ctx context.Context, projectID string) error {
	if err := d.DeleteService(ctx, projectID, "", ""); err != nil {
		logrus.Errorf("error deleting project in docker unable to delete services with project id (%s) - %s", projectID, err.Error())
		return err
	}
	return d.fileSystem.RemoveProjectDirectory(projectID)
}

// CreateProject creates a new project
func (d *Docker) CreateProject(ctx context.Context, project *model.Project) error {
	return d.fileSystem.CreateProjectDirectory(project.ID)
}
