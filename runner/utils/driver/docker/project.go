package docker

import (
	"context"
	"fmt"
	"os"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/runner/model"
)

// DeleteProject deletes the project
func (d *Docker) DeleteProject(ctx context.Context, projectID string) error {
	if err := d.DeleteService(ctx, projectID, "", ""); err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("error deleting project in docker unable to delete services with project id (%s)", projectID), err, nil)
	}
	if err := os.RemoveAll(fmt.Sprintf("%s/%s", d.secretPath, projectID)); err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("error deleting project in docker unable to delete secrets folder at (%s)", d.secretPath), err, nil)
	}
	return nil
}

// CreateProject creates a new project
func (d *Docker) CreateProject(ctx context.Context, project *model.Project) error {
	projectPath := fmt.Sprintf("%s/%s", d.secretPath, project.ID)
	if err := d.createDir(projectPath); err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("error creating secret in docker unable to create directory (%s)", projectPath), err, nil)
	}
	return nil
}
